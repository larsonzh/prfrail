package adapters

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

var (
	// A6 pinned CLI launcher errors.
	//
	// ErrAgentRunnerLaunchVersion reports a pinned CLI whose version does not
	// match the pin. Nothing is spawned in that case.
	ErrAgentRunnerLaunchVersion = errors.New("AgentRunner pinned CLI version mismatch")
	// ErrInvalidAgentRunnerLauncher reports an unusable launcher configuration.
	ErrInvalidAgentRunnerLauncher = errors.New("AgentRunner launcher invalid")
	// ErrAgentRunnerLaunchConflict reports a second launch while another one is
	// still live in the same run root. The durable identity mirror is unique per
	// run root, so a second launch would overwrite it and strip the live process of
	// the identity an operator stop needs.
	ErrAgentRunnerLaunchConflict = errors.New("AgentRunner launch conflict")
	// ErrAgentRunnerLaunchUnsettled reports a run whose own terminal
	// reconciliation did not finish inside the settle bound. Its process result is
	// then a placeholder rather than a fact, so the caller must degrade the outcome.
	ErrAgentRunnerLaunchUnsettled = errors.New("AgentRunner launch unsettled")
)

// agentRunnerEvidenceDirEnvVar tells the pinned CLI where its evidence artifacts
// belong. Redirecting them into the evidence directory keeps the isolated
// workspace a clean candidate output tree, so the manifest and diff only cover
// what the run actually produced.
const agentRunnerEvidenceDirEnvVar = "PROOFRAIL_EVIDENCE_DIR"

// agentRunnerLaunchSlotFreshness bounds how long a launch slot is treated as
// owned without proof. Inside the window a claimed slot may belong to a launch
// that is still publishing its identity, so it is never taken over; after it, an
// unexaminable or dead identity mirror proves the owner is gone and the slot may
// be reused.
const agentRunnerLaunchSlotFreshness = 10 * time.Second

// AgentRunnerPinnedCLIConfig pins the offline CLI one launcher may run. A6 runs
// only the deterministic stub: the launcher refuses a request that names another
// executable, so a caller cannot bypass the pin.
type AgentRunnerPinnedCLIConfig struct {
	Executable   string
	Version      string
	VersionArgs  []string
	RunArgs      []string
	EnvAllow     []string
	CheckTimeout time.Duration
	Grace        time.Duration
}

// AgentRunnerPinnedCLILauncher implements chain.AgentRunnerLauncher for the
// offline slice: version precheck, single spawn, log artifacts, durable process
// identity and the stop entry. It never writes chain state or policy, and it
// returns an error only when nothing is left running.
type AgentRunnerPinnedCLILauncher struct {
	config   AgentRunnerPinnedCLIConfig
	store    *AgentRunnerReplayStore
	registry *agentRunnerProcessRegistry
}

var _ chain.AgentRunnerLauncher = (*AgentRunnerPinnedCLILauncher)(nil)

func NewAgentRunnerPinnedCLILauncher(config AgentRunnerPinnedCLIConfig, store *AgentRunnerReplayStore) (*AgentRunnerPinnedCLILauncher, error) {
	if strings.TrimSpace(config.Executable) == "" || strings.TrimSpace(config.Version) == "" {
		return nil, fmt.Errorf("%w: executable and version are required", ErrInvalidAgentRunnerLauncher)
	}
	if store == nil || strings.TrimSpace(store.RunRoot()) == "" {
		return nil, fmt.Errorf("%w: a replay store bound to a run root is required", ErrInvalidAgentRunnerLauncher)
	}
	if config.CheckTimeout <= 0 {
		config.CheckTimeout = 5 * time.Second
	}
	if config.Grace <= 0 {
		config.Grace = 2 * time.Second
	}
	return &AgentRunnerPinnedCLILauncher{config: config, store: store, registry: newAgentRunnerProcessRegistry()}, nil
}

// StartAgentRunnerProcess prechecks the pinned version, spawns exactly one
// process, persists its identity mirror and returns the opaque process id. If the
// identity cannot be published the spawn is stopped again, so the caller never
// receives a result it cannot stop.
func (launcher *AgentRunnerPinnedCLILauncher) StartAgentRunnerProcess(ctx context.Context, launch chain.AgentRunnerLaunchRequest) (chain.AgentRunnerLaunchResult, error) {
	if launcher == nil {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: nil launcher", ErrInvalidAgentRunnerLauncher)
	}
	if err := ctx.Err(); err != nil {
		return chain.AgentRunnerLaunchResult{}, err
	}
	if !evidence.ValidID(launch.RequestID) {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: invalid requestId", ErrInvalidAgentRunnerLauncher)
	}
	if command := strings.TrimSpace(launch.Command); command != "" && command != launcher.config.Executable {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: request names %q but the launcher is pinned to %q", ErrInvalidAgentRunnerLauncher, command, launcher.config.Executable)
	}
	workspaceRoot := strings.TrimSpace(launch.WorkspaceRoot)
	if workspaceRoot == "" {
		workspaceRoot = strings.TrimSpace(launch.Dir)
	}
	if workspaceRoot == "" {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: a workspace root is required", ErrInvalidAgentRunnerLauncher)
	}
	launchID, err := agentRunnerLaunchID(launch.RequestID)
	if err != nil {
		return chain.AgentRunnerLaunchResult{}, err
	}
	// Exclusivity comes first: creating the log artifacts truncates them and writing
	// the mirror replaces the durable identity, so neither may happen while another
	// launch is still live in this run root. The slot is claimed atomically, so two
	// concurrent Starts cannot both pass a check and then both spawn.
	if liveID, live := launcher.registry.live(); live {
		if liveID == launchID {
			return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: launch %s", ErrAgentRunnerProcessDuplicate, launchID)
		}
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: launch %s is still live", ErrAgentRunnerLaunchConflict, liveID)
	}
	// The run's grace comes from the request when it names one, so a caller can
	// tighten the stop bound without a second configuration path.
	grace := launcher.config.Grace
	if launch.Grace > 0 {
		grace = launch.Grace
	}
	if err := launcher.checkPinnedVersion(ctx); err != nil {
		return chain.AgentRunnerLaunchResult{}, err
	}
	// The slot is claimed only now, after the precheck: the claim has to cover the
	// window between taking the run root and publishing the identity, and that window
	// must not include the precheck's own budget, which a caller may configure far
	// beyond the freshness of a claimed slot.
	slotToken, err := launcher.claimLaunchSlot(launchID)
	if err != nil {
		return chain.AgentRunnerLaunchResult{}, err
	}
	// Every failure below must release the slot again, otherwise a refused launch
	// would lock the run root for the rest of the process.
	claimed := true
	defer func() {
		if claimed {
			launcher.releaseLaunchSlot(slotToken)
		}
	}()
	evidenceDir, err := launcher.store.RunEvidenceDir(launch.RequestID)
	if err != nil {
		return chain.AgentRunnerLaunchResult{}, err
	}
	logsDirectory := filepath.Join(evidenceDir, agentRunnerLogsDirectoryName)
	if err := os.MkdirAll(logsDirectory, 0o755); err != nil {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: create %s: %v", ErrInvalidAgentRunnerLauncher, logsDirectory, err)
	}
	stdout, err := os.Create(filepath.Join(logsDirectory, agentRunnerStdoutFileName))
	if err != nil {
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: create stdout log: %v", ErrInvalidAgentRunnerLauncher, err)
	}
	stderr, err := os.Create(filepath.Join(logsDirectory, agentRunnerStderrFileName))
	if err != nil {
		_ = stdout.Close()
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: create stderr log: %v", ErrInvalidAgentRunnerLauncher, err)
	}
	args := launch.Args
	if len(args) == 0 {
		args = launcher.config.RunArgs
	}
	process, err := guard.StartManaged(ctx, guard.ProcessSpec{
		Command: launcher.config.Executable,
		Args:    append([]string(nil), args...),
		Dir:     launch.Dir,
		Env:     append(launcher.childEnv(launch.Env), agentRunnerEvidenceDirEnvVar+"="+evidenceDir),
		Stdout:  stdout,
		Stderr:  stderr,
	})
	if err != nil {
		// guard guarantees nothing is running when it returns an error.
		_ = stdout.Close()
		_ = stderr.Close()
		return chain.AgentRunnerLaunchResult{}, err
	}
	processID, err := agentRunnerProcessID(process.Identity())
	if err != nil {
		proof, stopErr := stopSpawnedProcess(process, grace)
		_ = stdout.Close()
		_ = stderr.Close()
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: %v (%s)", ErrAgentRunnerLaunchIdentityUnproven, err, describeAgentRunnerSpawnStop(proof, stopErr))
	}
	if err := agentRunnerIdentityMirrorWriter(launcher.store.RunRoot(), process.Identity()); err != nil {
		// Without the durable identity an operator stop cannot find the process, so
		// the spawn is stopped again instead of being reported as started.
		proof, stopErr := stopSpawnedProcess(process, grace)
		_ = stdout.Close()
		_ = stderr.Close()
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: %v (%s)", ErrAgentRunnerLaunchIdentityUnproven, err, describeAgentRunnerSpawnStop(proof, stopErr))
	}
	startedAt := time.Now().UTC()
	handle := &agentRunnerProcessHandle{
		launchID:  launchID,
		requestID: launch.RequestID,
		process:   process,
		phase:     agentRunnerProcessRunning,
		startedAt: startedAt,
		grace:     grace,
		slotToken: slotToken,
		logs:      []*os.File{stdout, stderr},
	}
	if err := launcher.registry.register(handle); err != nil {
		proof, stopErr := stopSpawnedProcess(process, grace)
		handle.closeLogs()
		return chain.AgentRunnerLaunchResult{}, fmt.Errorf("%w: %v (%s)", err, stopErr, describeAgentRunnerSpawnStop(proof, stopErr))
	}
	// The launch owns the slot until it is retired, so the deferred release must not
	// fire for a successful start.
	claimed = false
	return chain.AgentRunnerLaunchResult{LaunchID: launchID, ProcessID: processID, StartedAt: startedAt}, nil
}

// StopAgentRunnerProcess stops one launch through its durable identity, which is
// the same mechanism the operator stop path uses. A process that already exited
// is reported as already-stopped instead of failing on a stale handle, and a live
// one is stopped by the reviewed platform stop. It always retires the launch, so
// the log artifacts are closed by the time the caller reads them.
func (launcher *AgentRunnerPinnedCLILauncher) StopAgentRunnerProcess(ctx context.Context, launchID string) (guard.TerminationEvidence, error) {
	if launcher == nil {
		return guard.TerminationEvidence{}, fmt.Errorf("%w: nil launcher", ErrInvalidAgentRunnerLauncher)
	}
	identity, err := launcher.stopIdentity(launchID)
	if err != nil {
		// The launch is retired even when the identity cannot be resolved: leaving the
		// handle behind would keep the run root claimed with no way to release it.
		launcher.releaseLaunch(launchID)
		return guard.TerminationEvidence{}, err
	}
	if err := launcher.registry.advance(launchID, agentRunnerProcessStopping); err != nil && !errors.Is(err, ErrAgentRunnerProcessMissing) {
		return guard.TerminationEvidence{}, err
	}
	// The launch may have been started with a tighter grace than the configuration,
	// and the operator path must not be more patient than the run it is stopping.
	proof, stopErr := guard.StopProcessIdentity(ctx, identity, launcher.stopGrace(launchID))
	launcher.releaseLaunch(launchID)
	return proof, stopErr
}

// refuseLiveMirrorProcess refuses a launch while the durable identity mirror names
// a process that is still running. The mirror is the only identity that survives a
// restart, so this is what keeps one run root to one live launch across processes;
// a mirror naming a dead process is stale and may be replaced.
func (launcher *AgentRunnerPinnedCLILauncher) refuseLiveMirrorProcess() error {
	path := filepath.Join(launcher.store.RunRoot(), managedProcessIdentityFileName)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No launch has ever published an identity here, so nothing can be running.
			return nil
		}
		// A mirror that cannot be examined is not proof that nothing is running, so the
		// launch is refused rather than risking a second process in one run root.
		return fmt.Errorf("%w: unreadable identity mirror %s: %v", ErrAgentRunnerLaunchConflict, path, err)
	}
	identity, err := loadAgentRunnerProcessIdentityMirror(launcher.store.RunRoot())
	if err != nil {
		return fmt.Errorf("%w: unreadable identity mirror %s: %v", ErrAgentRunnerLaunchConflict, path, err)
	}
	alive, err := guard.ProcessAlive(identity)
	if err != nil {
		return fmt.Errorf("%w: cannot verify the mirrored process: %v", ErrAgentRunnerLaunchConflict, err)
	}
	if alive {
		return fmt.Errorf("%w: the mirrored process %d is still running", ErrAgentRunnerLaunchConflict, identity.PID)
	}
	return nil
}

// stopIdentity resolves which process one launch id may stop. Stopping has exactly one
// mechanism - the guard's identity stop - with two entry points: the owner of a live
// handle stops through its own identity, and a launcher that holds nothing (a restart
// or an operator side) stops through the persisted identity. No third path exists.
// A foreign or stale launch id may never reach a stranger's process.
func (launcher *AgentRunnerPinnedCLILauncher) stopIdentity(launchID string) (guard.ProcessIdentity, error) {
	if handle, found := launcher.registry.lookup(launchID); found && handle.process != nil {
		return handle.process.Identity(), nil
	}
	if liveID, live := launcher.registry.live(); live && liveID != launchID {
		return guard.ProcessIdentity{}, fmt.Errorf("%w: this launcher holds %s, not %s", ErrAgentRunnerLaunchConflict, liveID, launchID)
	}
	if err := launcher.refuseForeignSlot(launchID); err != nil {
		return guard.ProcessIdentity{}, err
	}
	return loadAgentRunnerProcessIdentityMirror(launcher.store.RunRoot())
}

// refuseForeignSlot refuses to stop through the mirror unless the request is the one
// the run root currently owns. Without a handle the mirror is the only identity
// available, so ownership has to come from the slot the live launch recorded. A slot
// that cannot be read, or that names nobody, is not proof that nothing runs: a live
// launch always holds its slot, so both fail closed instead of falling through to the
// durable identity.
func (launcher *AgentRunnerPinnedCLILauncher) refuseForeignSlot(launchID string) error {
	slot := filepath.Join(launcher.store.RunRoot(), managedProcessLaunchSlotFileName)
	owner, err := os.ReadFile(slot)
	switch {
	case err == nil:
		fields := strings.Fields(string(owner))
		if len(fields) == 0 {
			return fmt.Errorf("%w: the launch slot %s names no owner", ErrAgentRunnerLaunchConflict, slot)
		}
		if fields[0] != launchID {
			return fmt.Errorf("%w: the run root is owned by %s, not %s", ErrAgentRunnerLaunchConflict, fields[0], launchID)
		}
		return nil
	case errors.Is(err, os.ErrNotExist):
		// No slot at all: the request may only fall back to the mirror when that mirror
		// does not name a live process, because a live launch always holds its slot. An
		// unreadable mirror refuses here too, the same rule the launch path applies.
		return launcher.refuseLiveMirrorProcess()
	default:
		return fmt.Errorf("%w: the launch slot %s cannot be read: %v", ErrAgentRunnerLaunchConflict, slot, err)
	}
}

// releaseLaunch retires one launch: the handle is marked stopped, its log
// artifacts are closed, it is dropped and the launch slot it claimed is released,
// so neither stop path leaves live or claimed state behind.
func (launcher *AgentRunnerPinnedCLILauncher) releaseLaunch(launchID string) {
	if handle, found := launcher.registry.lookup(launchID); found {
		launcher.releaseLaunchSlot(handle.slotToken)
	}
	handle, found := launcher.registry.lookup(launchID)
	if !found {
		return
	}
	_ = launcher.registry.advance(launchID, agentRunnerProcessStopped)
	handle.closeLogs()
	launcher.registry.release(launchID)
}

// claimLaunchSlot takes the durable launch slot of one run root atomically. An
// exclusive create is the only check-and-take that two concurrent Starts, two
// launcher instances or a restart cannot both win. An existing slot is only
// stolen when it is old enough that no launch could still be publishing its
// identity and the identity mirror proves the previous owner is gone; a live or
// unexaminable mirror is never overridden.
func (launcher *AgentRunnerPinnedCLILauncher) claimLaunchSlot(launchID string) (string, error) {
	slot := filepath.Join(launcher.store.RunRoot(), managedProcessLaunchSlotFileName)
	info, statErr := os.Stat(slot)
	switch {
	case errors.Is(statErr, os.ErrNotExist):
		// A free slot is not proof that nothing is running: the identity mirror is the
		// durable record of the last launch, and a live process in it blocks the slot
		// even when the slot file itself is gone.
		if err := launcher.refuseLiveMirrorProcess(); err != nil {
			return "", err
		}
		return launcher.createLaunchSlot(slot, launchID)
	case statErr != nil:
		return "", fmt.Errorf("%w: cannot examine the launch slot %s: %v", ErrAgentRunnerLaunchConflict, slot, statErr)
	case time.Since(info.ModTime()) < agentRunnerLaunchSlotFreshness:
		// A fresh slot may belong to a launch that is still starting: taking it would
		// let two processes run in one run root.
		return "", fmt.Errorf("%w: the launch slot %s was taken %v ago", ErrAgentRunnerLaunchConflict, slot, time.Since(info.ModTime()).Round(time.Millisecond))
	}
	if err := launcher.refuseLiveMirrorProcess(); err != nil {
		return "", err
	}
	if err := os.Remove(slot); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("%w: cannot clear the stale launch slot %s: %v", ErrAgentRunnerLaunchConflict, slot, err)
	}
	return launcher.createLaunchSlot(slot, launchID)
}

// createLaunchSlot performs the exclusive create that actually takes the slot. The
// slot records the launch and a token unique to this claim, so a later release can
// tell its own slot from one another launch has taken in the meantime - even when
// the same launch id is claimed again.
func (launcher *AgentRunnerPinnedCLILauncher) createLaunchSlot(slot, launchID string) (string, error) {
	file, err := os.OpenFile(slot, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// Another launch won the race between the check and the create.
			return "", fmt.Errorf("%w: another launch claimed %s first", ErrAgentRunnerLaunchConflict, slot)
		}
		return "", fmt.Errorf("%w: cannot claim the launch slot %s: %v", ErrAgentRunnerLaunchConflict, slot, err)
	}
	token := launchID + "\n" + strconv.FormatInt(time.Now().UnixNano(), 36) + "\n"
	if _, err := file.WriteString(token); err != nil {
		_ = file.Close()
		_ = os.Remove(slot)
		return "", fmt.Errorf("%w: cannot record the slot owner in %s: %v", ErrAgentRunnerLaunchConflict, slot, err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(slot)
		return "", fmt.Errorf("%w: cannot close the launch slot %s: %v", ErrAgentRunnerLaunchConflict, slot, err)
	}
	return strings.TrimSpace(token), nil
}

// releaseLaunchSlot releases the slot this launch claimed. The slot is first renamed
// aside, which is atomic: whichever of a release and a reclaim renames it wins, so a
// late release can never delete a slot another launch has just created. Only a slot
// whose token still matches this claim is removed; anything else is put back.
func (launcher *AgentRunnerPinnedCLILauncher) releaseLaunchSlot(token string) {
	if strings.TrimSpace(token) == "" {
		return
	}
	slot := filepath.Join(launcher.store.RunRoot(), managedProcessLaunchSlotFileName)
	staging := slot + ".releasing"
	if err := os.Rename(slot, staging); err != nil {
		// Missing, or already taken over by another release: nothing of ours to remove.
		return
	}
	owner, err := os.ReadFile(staging)
	if err != nil || strings.TrimSpace(string(owner)) != token {
		// Another launch owns the slot now: put it back untouched.
		_ = os.Rename(staging, slot)
		return
	}
	_ = os.Remove(staging)
}

// stopGrace resolves the grace bound one launch is stopped with inside this
// process: the launch's own grace when it named one, otherwise the configured
// value. It is deliberately not persisted - the operator stop path keeps its own
// bound - so this only tightens the stop this launcher performs.
func (launcher *AgentRunnerPinnedCLILauncher) stopGrace(launchID string) time.Duration {
	if handle, found := launcher.registry.lookup(launchID); found && handle.grace > 0 {
		return handle.grace
	}
	return launcher.config.Grace
}

// WaitAgentRunnerProcess waits for one live launch under a watchdog. The watchdog
// only bounds the wait: stopping is still owned by the caller's context and the
// guard handle, so a timeout becomes a provable stop instead of an abandon.
func (launcher *AgentRunnerPinnedCLILauncher) WaitAgentRunnerProcess(ctx context.Context, launchID, requestID string, timeout time.Duration) (AgentRunnerWatchResult, guard.ProcessResult, error) {
	if launcher == nil {
		return AgentRunnerWatchResult{}, guard.ProcessResult{}, fmt.Errorf("%w: nil launcher", ErrInvalidAgentRunnerLauncher)
	}
	handle, found := launcher.registry.lookup(launchID)
	if !found {
		return AgentRunnerWatchResult{}, guard.ProcessResult{}, fmt.Errorf("%w: launch %s", ErrAgentRunnerProcessMissing, launchID)
	}
	manager := AgentRunnerTimeoutManager{Timeout: timeout}
	grace := handle.grace
	if grace <= 0 {
		grace = launcher.config.Grace
	}
	// The run side reconciles its own terminal state - the stop proof, the tree
	// verification, the exit status - and the watchdog may decide before that
	// finishes. Reading the process result at that moment would read a placeholder,
	// so the launcher waits for the reconciliation, bounded by the settle bound.
	var result guard.ProcessResult
	done := make(chan guard.ProcessResult, 1)
	callerLive := ctx.Err() == nil
	watch, err := manager.RunWithWatchdog(ctx, requestID, func(childCtx context.Context) error {
		res, runErr := handle.process.Run(childCtx, grace, handle.startedAt)
		done <- res
		return runErr
	})
	if err != nil {
		return AgentRunnerWatchResult{}, result, err
	}
	select {
	case result = <-done:
		return watch, result, nil
	default:
	}
	if !callerLive {
		// The caller was already done, so the watchdog never started the run: there
		// is no reconciliation to wait for.
		return watch, result, nil
	}
	select {
	case result = <-done:
		return watch, result, nil
	case <-time.After(agentRunnerSettleWait(grace)):
		// An unsettled run is not a clean exit: the caller must degrade the outcome
		// instead of trusting a zero-valued process result. The placeholder stop
		// proof keeps it out of the natural-exit branch of the status mapping.
		watch.Err = fmt.Errorf("%w: %v", ErrAgentRunnerLaunchUnsettled, watch.Err)
		unsettled := guard.TerminationEvidence{Outcome: "unknown"}
		return watch, guard.ProcessResult{
			ExitCode:    -1,
			Started:     true,
			StartedAt:   handle.startedAt,
			FinishedAt:  time.Now().UTC(),
			Termination: &unsettled,
		}, nil
	}
}

// agentRunnerSettleWait bounds how long the launcher waits for the run side to
// reconcile its terminal state after the watchdog has decided. It is derived from
// the grace period (which already bounds the stop) plus a small margin, and it is
// capped so a pathological grace cannot turn into a long wait.
func agentRunnerSettleWait(grace time.Duration) time.Duration {
	if grace <= 0 {
		grace = 2 * time.Second
	}
	wait := grace + 500*time.Millisecond
	if limit := 10 * time.Second; wait > limit {
		wait = limit
	}
	return wait
}

// checkPinnedVersion runs the precheck through the reviewed process boundary and
// refuses to spawn the run on any mismatch. The precheck spawn is contained by
// the same mechanism as the run itself, so an unanswered precheck cannot leak a
// process either.
func (launcher *AgentRunnerPinnedCLILauncher) checkPinnedVersion(ctx context.Context) error {
	checkCtx, cancel := context.WithTimeout(ctx, launcher.config.CheckTimeout)
	defer cancel()
	args := launcher.config.VersionArgs
	if len(args) == 0 {
		args = []string{"--print-version"}
	}
	var stdout bytes.Buffer
	result, err := guard.RunManaged(checkCtx, guard.ProcessSpec{
		Command: launcher.config.Executable,
		Args:    append([]string(nil), args...),
		Stdout:  &stdout,
	}, launcher.config.Grace)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAgentRunnerLaunchVersion, err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("%w: the version precheck exited with %d", ErrAgentRunnerLaunchVersion, result.ExitCode)
	}
	if reported := strings.TrimSpace(stdout.String()); reported != launcher.config.Version {
		return fmt.Errorf("%w: reported %q, pinned %q", ErrAgentRunnerLaunchVersion, reported, launcher.config.Version)
	}
	return nil
}

// childEnv builds the child environment: the request's environment when it names
// one, otherwise the inherited one, always filtered by the allowlist when set.
// The evidence directory is appended by the caller and is never filtered out,
// because the run cannot publish its artifacts without it.
func (launcher *AgentRunnerPinnedCLILauncher) childEnv(requested []string) []string {
	env := requested
	if len(env) == 0 {
		env = os.Environ()
	}
	if len(launcher.config.EnvAllow) == 0 {
		return env
	}
	allowed := make(map[string]struct{}, len(launcher.config.EnvAllow))
	for _, name := range launcher.config.EnvAllow {
		allowed[strings.ToUpper(strings.TrimSpace(name))] = struct{}{}
	}
	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		name := entry
		if index := strings.IndexByte(entry, '='); index >= 0 {
			name = entry[:index]
		}
		if _, ok := allowed[strings.ToUpper(name)]; ok {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// stopSpawnedProcess stops a process the launcher already spawned, using a
// context that cannot be cancelled out from under the cleanup.
func stopSpawnedProcess(process *guard.ManagedProcess, grace time.Duration) (guard.TerminationEvidence, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(context.Background()), grace+time.Second)
	defer cancel()
	return process.Terminate(ctx, grace)
}

// describeAgentRunnerSpawnStop reports what happened to a spawn the launcher had
// to abandon, because an operator has to be able to tell "stopped" from
// "possibly still running".
func describeAgentRunnerSpawnStop(proof guard.TerminationEvidence, err error) string {
	if err != nil && proof.Outcome == "" {
		return "spawn stop failed: " + err.Error()
	}
	return "spawn stop outcome: " + proof.Outcome
}

// agentRunnerIdentityMirrorWriter publishes the durable process identity. It is a
// seam so a test can make the publication fail on every platform instead of relying
// on a platform-specific "cannot replace this file" fixture.
var agentRunnerIdentityMirrorWriter = writeAgentRunnerProcessIdentityMirror

// loadAgentRunnerProcessIdentityMirror reads the durable identity an operator
// stop would use. It fails closed on anything a stop would have to guess about.
func loadAgentRunnerProcessIdentityMirror(runRoot string) (guard.ProcessIdentity, error) {
	if strings.TrimSpace(runRoot) == "" {
		return guard.ProcessIdentity{}, fmt.Errorf("%w: empty run root", ErrAgentRunnerProcessIdentity)
	}
	data, err := os.ReadFile(filepath.Join(runRoot, managedProcessIdentityFileName))
	if err != nil {
		return guard.ProcessIdentity{}, fmt.Errorf("%w: %v", ErrAgentRunnerProcessIdentity, err)
	}
	return decodeAgentRunnerProcessID(string(data))
}
