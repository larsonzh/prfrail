package guard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// naturalExitVerifyWindow bounds the tree check that follows a natural exit.
	naturalExitVerifyWindow = 20 * time.Millisecond
	// startupCleanupGrace bounds reaping a process that failed to start cleanly.
	startupCleanupGrace = time.Second
	// identityStopCommandBound is the synchronous window the identity-only stop
	// path gives its external stop command. taskkill start-up latency is not the
	// tree stop grace, so it takes the widest allowance the design permits:
	// docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md
	// §1 row 8 fixes the settle shape at grace+settle with settle <= 500ms and a
	// hard 10s ceiling; a caller-supplied grace can be far smaller than taskkill's
	// own start-up cost under parallel load, so this path uses the ceiling.
	identityStopCommandBound = 10 * time.Second
)

var (
	ErrInvalidProcess       = errors.New("invalid managed process")
	ErrTerminationUncertain = errors.New("managed process termination uncertain")
	ErrProcessTreeRemained  = errors.New("managed process tree remained after parent exit")
	// ErrSpawnOwnership marks a spawn that never reached a running state because
	// the ownership chain could not be completed: the OnSpawn hook failed, the
	// ownership anchor could not be verified, or the process could not be resumed.
	// It always wraps the underlying cause, so a caller can tell an ownership
	// failure apart from an ordinary exec error with errors.Is.
	ErrSpawnOwnership = errors.New("managed process spawn ownership unproven")
	// errStopCommandExpired marks an external stop command that outlived its
	// synchronous bound. The command is never killed at that point: it keeps
	// running detached, so this error means "unknown", not "failed".
	errStopCommandExpired = errors.New("stop command outlived its synchronous bound")
)

type ProcessSpec struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
	// OnSpawn, when set, runs after the process exists with a complete identity
	// and its ownership anchor verified, and before it is allowed to run. A
	// non-nil error means "not started": the guard stops the process again and
	// StartManaged returns an error that satisfies errors.Is(ErrSpawnOwnership).
	// The hook owns whatever durable record it writes and must clean it up when it
	// fails; the guard never touches the caller's namespace. A nil hook keeps the
	// previous behaviour exactly.
	OnSpawn func(ctx context.Context, identity ProcessIdentity) error
}

// managedProcessHandle is the OS process handle guard owns. It exposes the small
// part of *os.Process the package needs - the pid, a direct kill and a wait - so
// the reaper and every stop path are platform neutral.
type managedProcessHandle struct {
	pid    int
	killFn func() error
	waitFn func() (*os.ProcessState, error)
}

// Pid reports the process id of the handle.
func (handle *managedProcessHandle) Pid() int {
	if handle == nil {
		return 0
	}
	return handle.pid
}

// Kill terminates the process directly. It exists for the reviewed stop paths
// and for a test that has to get rid of a leftover child; it is never the
// production stop, which goes through the platform stop.
func (handle *managedProcessHandle) Kill() error {
	if handle == nil || handle.killFn == nil {
		return ErrInvalidProcess
	}
	return handle.killFn()
}

// Wait reaps the process handle. The package always reaps through the single
// reaper; this mirrors *os.Process so a caller can wait for a child it killed.
func (handle *managedProcessHandle) Wait() (*os.ProcessState, error) {
	if handle == nil || handle.waitFn == nil {
		return nil, ErrInvalidProcess
	}
	return handle.waitFn()
}

// managedCommand is the platform spawn result: the process handle, the platform
// wait, and the exit status the single reaper records. The reaper is the only
// writer of exitStatus, so any reader must have observed waitDone first.
type managedCommand struct {
	Process    *managedProcessHandle
	waitFn     func() error
	exitStatus int
}

// wait reaps the process through the platform wait.
func (command *managedCommand) wait() error { return command.waitFn() }

// exitCode reports the exit status of a reaped process.
func (command *managedCommand) exitCode() int { return command.exitStatus }

type ProcessResult struct {
	ExitCode    int
	Started     bool
	StartedAt   time.Time
	FinishedAt  time.Time
	Termination *TerminationEvidence
}

type ProcessIdentity struct {
	PID        int    `json:"pid"`
	StartToken string `json:"startToken"`
}

type TerminationEvidence struct {
	Identity  ProcessIdentity `json:"identity"`
	Requested string          `json:"requestedAt"`
	Verified  string          `json:"verifiedAt"`
	Outcome   string          `json:"outcome"`
	Actions   []string        `json:"actions"`
	Hash      string          `json:"evidenceHash"`
}

type ManagedProcess struct {
	cmd      *managedCommand
	identity ProcessIdentity
	// platform is a pointer: it owns a synchronised handle, so it must never be
	// copied once it may be used from several goroutines.
	platform *platformProcess
	// verify is a test seam: when set it replaces platform.verifyGone so a test
	// can observe the context each call site passes. Production leaves it nil.
	verify func(context.Context) error
	// waitOnce starts the single reaper goroutine: exactly one cmd.Wait may ever
	// run for a managed process, because concurrent cmd.Wait calls panic.
	waitOnce sync.Once
	// waitDone is closed when the unique reaper finishes; it is created by
	// StartManaged so every caller shares one broadcast channel.
	waitDone chan struct{}
	// waitErr is the reaper result. It is written before waitDone is closed, so
	// it is only safe to read after <-waitDone.
	waitErr error
}

// waitResult returns the single reaper broadcast channel, starting the unique
// reaper on first use. It is idempotent and always returns the same channel, so
// concurrent callers can never race a second cmd.Wait.
func (process *ManagedProcess) waitResult() <-chan struct{} {
	process.waitOnce.Do(func() {
		if process.waitDone == nil {
			process.waitDone = make(chan struct{})
		}
		go func() {
			process.waitErr = process.cmd.wait()
			close(process.waitDone)
		}()
	})
	return process.waitDone
}

// verifyGone reports whether the managed tree is gone, using the bounded
// context its caller supplies. Stop verification must always be bounded: an
// unbounded wait turns an unkillable descendant into a permanent hang.
func (process *ManagedProcess) verifyGone(ctx context.Context) error {
	if process.verify != nil {
		return process.verify(ctx)
	}
	return process.platform.verifyGone(ctx)
}

// StartManaged creates one managed process and drives it through the fixed
// startup order of the spawn-ownership design
// (docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// sections 2 and 3): prepare the platform ownership anchor, create the process, verify
// the anchor, capture the identity, run the OnSpawn hook, and only then let the
// process run. Ownership therefore always precedes execution, and the hook always
// observes a complete identity. Every failure runs through abortStartup, so a
// refused start leaves no unrunnable residue behind.
func StartManaged(ctx context.Context, spec ProcessSpec) (*ManagedProcess, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if spec.Command == "" {
		return nil, fmt.Errorf("%w: command is empty", ErrInvalidProcess)
	}
	platform, err := preparePlatformProcess()
	if err != nil {
		return nil, err
	}
	command, err := platform.start(&spec)
	if err != nil {
		platform.close()
		return nil, fmt.Errorf("start managed process: %w", err)
	}
	pid := command.Process.Pid()
	process := &ManagedProcess{cmd: command, platform: platform, waitDone: make(chan struct{})}
	// attach verifies the ownership anchor the platform created the process with;
	// it never has to establish it, which is what removes the window between
	// creation and ownership.
	if err := platform.attach(pid); err != nil {
		_ = command.Process.Kill()
		return nil, process.abortStartup(fmt.Errorf("%w: attach managed process tree: %w", ErrSpawnOwnership, err))
	}
	identity, err := platformIdentity(pid)
	if err != nil {
		_ = platform.terminate()
		return nil, process.abortStartup(fmt.Errorf("%w: capture managed process identity: %w", ErrSpawnOwnership, err))
	}
	process.identity = identity
	if spec.OnSpawn != nil {
		if hookErr := spec.OnSpawn(ctx, identity); hookErr != nil {
			// The process must not be left suspended and unowned: the tree is
			// terminated first, then reaped through the bounded abort.
			_ = platform.terminate()
			return nil, process.abortStartup(fmt.Errorf("%w: on spawn hook: %w", ErrSpawnOwnership, hookErr))
		}
	}
	if err := platform.resume(); err != nil {
		_ = platform.terminate()
		return nil, process.abortStartup(fmt.Errorf("%w: resume managed process: %w", ErrSpawnOwnership, err))
	}
	return process, nil
}

// abortStartup stops a process that could not be started cleanly; it reaps through
// the single reaper with a short bound. When the tree cannot be proven exited, the
// platform handle stays owned (a later Reclaim can finish the job) and the error
// reports ErrTerminationUncertain.
func (process *ManagedProcess) abortStartup(cause error) error {
	waitCtx, cancel := boundedTerminationContext(context.Background(), startupCleanupGrace)
	defer cancel()
	select {
	case <-process.waitResult():
		process.platform.close()
		return cause
	case <-waitCtx.Done():
		// Ownership stays with the handle: the tree may still be running, so the
		// boundary rule forbids releasing it here. This process never escapes
		// StartManaged, so nobody could ever hand it to Reclaim - a background
		// reaper waits for the single waiter instead and releases the handle once
		// the tree is finally reaped. The goroutine only reads waitDone, so the
		// bounded return above is unaffected. platform.close is idempotent, so a
		// later Reclaim remains safe if a caller ever obtains the process.
		go func() {
			<-process.waitResult()
			process.platform.close()
		}()
		return errors.Join(cause, fmt.Errorf("%w: process did not exit within %s", ErrTerminationUncertain, startupCleanupGrace))
	}
}

func RunManaged(ctx context.Context, spec ProcessSpec, grace time.Duration) (ProcessResult, error) {
	startedAt := time.Now().UTC()
	process, err := StartManaged(ctx, spec)
	if err != nil {
		return ProcessResult{ExitCode: -1, StartedAt: startedAt, FinishedAt: time.Now().UTC()}, err
	}
	return process.runManaged(ctx, grace, startedAt)
}

// Run waits for a process that was started separately and reconciles its terminal
// state exactly as RunManaged does: a natural exit is reported only after the tree
// is verified gone, and a cancelled or expired context triggers a bounded stop.
// It exists so a caller that must publish the process identity before waiting can
// still reuse one reconciliation path.
func (process *ManagedProcess) Run(ctx context.Context, grace time.Duration, startedAt time.Time) (ProcessResult, error) {
	if process == nil || process.cmd == nil {
		return ProcessResult{ExitCode: -1, StartedAt: startedAt, FinishedAt: time.Now().UTC()}, ErrInvalidProcess
	}
	return process.runManaged(ctx, grace, startedAt)
}

// runManaged waits for the managed process and reconciles its terminal state: a
// natural exit is reported only after the tree is verified gone, and a cancelled
// or expired context triggers a bounded stop followed by bounded verification.
func (process *ManagedProcess) runManaged(ctx context.Context, grace time.Duration, startedAt time.Time) (ProcessResult, error) {
	waitDone := process.waitResult()
	select {
	case <-waitDone:
		verifyCtx, cancel := context.WithTimeout(context.Background(), naturalExitVerifyWindow)
		verifyErr := process.verifyGone(verifyCtx)
		cancel()
		if verifyErr != nil {
			requested := time.Now().UTC()
			actions, stopErr := process.platform.stop(process.identity, grace)
			if stopErr == nil {
				finalCtx, finalCancel := context.WithTimeout(context.Background(), grace)
				stopErr = process.verifyGone(finalCtx)
				finalCancel()
			}
			if stopErr == nil {
				process.platform.close()
			}
			proof, proofErr := buildTerminationEvidence(process.identity, requested, actions, stopErr)
			result := ProcessResult{ExitCode: -1, Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC(), Termination: &proof}
			if proofErr != nil {
				return result, proofErr
			}
			return result, ErrProcessTreeRemained
		}
		process.platform.close()
		return ProcessResult{ExitCode: process.cmd.exitCode(), Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC()}, process.waitErr
	case <-ctx.Done():
		requested := time.Now().UTC()
		actions, stopErr := process.platform.stop(process.identity, grace)
		// The reaper is the only owner of cmd.Wait, so waiting for it is bounded:
		// the caller's context is already cancelled here, so the wait needs its
		// own bound instead of inheriting (or outliving) the cancellation.
		waitCtx, waitCancel := stopVerificationContext(grace)
		select {
		case <-waitDone:
			waitCancel()
		case <-waitCtx.Done():
			waitCancel()
			stopErr = errors.Join(stopErr, ErrTerminationUncertain)
		}
		if stopErr == nil {
			// The caller's context is already cancelled on this path, so the
			// verification needs its own bounded context: it must not inherit the
			// cancellation whose handler is waiting on it, and it must not be
			// unbounded either.
			verifyCtx, verifyCancel := stopVerificationContext(grace)
			stopErr = process.verifyGone(verifyCtx)
			verifyCancel()
		}
		if stopErr == nil {
			process.platform.close()
		}
		proof, proofErr := buildTerminationEvidence(process.identity, requested, actions, stopErr)
		result := ProcessResult{ExitCode: -1, Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC(), Termination: &proof}
		if proofErr != nil {
			return result, proofErr
		}
		return result, ctx.Err()
	}
}

func (process *ManagedProcess) Identity() ProcessIdentity { return process.identity }

func InspectProcess(pid int) (ProcessIdentity, error) {
	if pid <= 0 {
		return ProcessIdentity{}, ErrInvalidProcess
	}
	return platformIdentity(pid)
}

func ProcessAlive(identity ProcessIdentity) (bool, error) {
	if identity.PID <= 0 || identity.StartToken == "" {
		return false, ErrInvalidProcess
	}
	return platformIdentityAlive(identity)
}

// StopProcessIdentity performs a best-effort controlled stop using a persisted
// process identity. It returns termination evidence for both success and
// uncertainty outcomes.
func StopProcessIdentity(ctx context.Context, identity ProcessIdentity, grace time.Duration) (TerminationEvidence, error) {
	if identity.PID <= 0 || identity.StartToken == "" {
		return TerminationEvidence{}, ErrInvalidProcess
	}
	alive, err := platformIdentityAlive(identity)
	if err != nil {
		return TerminationEvidence{}, err
	}
	requested := time.Now().UTC()
	if !alive {
		return buildTerminationEvidence(identity, requested, []string{"already-stopped"}, nil)
	}
	// taskkill is an external OS command, so its latency is not the tree-stop
	// grace: it gets identityStopCommandBound. The allowance is still a hard bound,
	// but reaching it no longer means the stop failed - the command keeps running
	// detached and the tree is proven gone (or not) by a bounded verify below.
	stopCtx, stopCancel := boundedTerminationContext(ctx, identityStopCommandBound)
	actions, stopErr := platformStopIdentity(stopCtx, identity, grace, identityStopCommandBound)
	stopCancel()
	detached := errors.Is(stopErr, errStopCommandExpired)
	if detached {
		// The command outlived its bound, so its result is unknown either way: the
		// caller must prove the tree itself instead of trusting the command.
		stopErr = nil
	}
	if stopErr == nil {
		waitCtx, cancel := boundedTerminationContext(ctx, grace)
		verifyErr := waitIdentityGone(waitCtx, identity)
		cancel()
		if verifyErr == nil && detached {
			actions = append(actions, "taskkill-settled-verified-gone")
		}
		stopErr = verifyErr
	}
	if detached && stopErr != nil {
		stopErr = errors.Join(errStopCommandExpired, stopErr)
	}
	proof, proofErr := buildTerminationEvidence(identity, requested, actions, stopErr)
	if proofErr != nil {
		return proof, proofErr
	}
	return proof, nil
}

// boundedTerminationContext returns the context a stop verification may use. It is
// derived from the grace period - never from the caller's deadline, which may be far
// longer - and it is immune to the caller's cancellation, so a caller giving up
// cannot cancel the verification whose result it is waiting for.
func boundedTerminationContext(ctx context.Context, grace time.Duration) (context.Context, context.CancelFunc) {
	if grace <= 0 {
		grace = time.Second
	}
	return context.WithTimeout(context.WithoutCancel(ctx), grace)
}

func (process *ManagedProcess) Terminate(ctx context.Context, grace time.Duration) (TerminationEvidence, error) {
	requested := time.Now().UTC()
	verifyCtx, verifyCancel := boundedTerminationContext(ctx, grace)
	defer verifyCancel()
	waitDone := process.waitResult()
	actions, terminateErr := process.platform.stop(process.identity, grace)
	select {
	case <-waitDone:
		// The tree is reaped; verification below proves it is really gone.
	case <-ctx.Done():
		// Cancellation wins, even when the verification bound expires in the
		// same instant.
		terminateErr = errors.Join(terminateErr, ctx.Err())
	case <-verifyCtx.Done():
		if ctxErr := ctx.Err(); ctxErr != nil {
			terminateErr = errors.Join(terminateErr, ctxErr)
		} else {
			terminateErr = errors.Join(terminateErr, ErrTerminationUncertain)
		}
	}
	if terminateErr == nil {
		terminateErr = process.verifyGone(verifyCtx)
	}
	if terminateErr == nil {
		terminateErr = waitIdentityGone(verifyCtx, process.identity)
	}
	if terminateErr == nil {
		// Ownership rule: an uncertain termination keeps the platform handle, so
		// only a proven stop may release it.
		process.platform.close()
	}
	return buildTerminationEvidence(process.identity, requested, actions, terminateErr)
}

// Reclaim releases the platform handle of a process whose termination could not be
// proven in time. It is idempotent and safe to call from several goroutines: the
// handle is released only once the single reaper has finished and the tree is
// verified gone, so an uncertain Terminate can hand ownership over instead of
// leaking it.
func (process *ManagedProcess) Reclaim(ctx context.Context) error {
	if process == nil || process.cmd == nil {
		return ErrInvalidProcess
	}
	if ctx == nil {
		return fmt.Errorf("%w: reclaim requires a context", ErrInvalidProcess)
	}
	select {
	case <-process.waitResult():
	default:
		return fmt.Errorf("%w: managed process is still running", ErrTerminationUncertain)
	}
	if err := process.verifyGone(ctx); err != nil {
		return err
	}
	process.platform.close()
	return nil
}

// stopVerificationContext returns a fresh bounded context for stop verification
// on the cancellation path, where the caller's context is already cancelled.
func stopVerificationContext(grace time.Duration) (context.Context, context.CancelFunc) {
	if grace <= 0 {
		grace = time.Second
	}
	return context.WithTimeout(context.Background(), grace)
}

func buildTerminationEvidence(identity ProcessIdentity, requested time.Time, actions []string, terminateErr error) (TerminationEvidence, error) {
	outcome := "stopped"
	if terminateErr != nil {
		outcome = "uncertain"
	}
	result := TerminationEvidence{
		Identity:  identity,
		Requested: formatTime(requested),
		Verified:  formatTime(time.Now().UTC()),
		Outcome:   outcome,
		Actions:   actions,
	}
	canonical, err := evidence.EncodeCanonical(struct {
		Identity  ProcessIdentity `json:"identity"`
		Requested string          `json:"requestedAt"`
		Verified  string          `json:"verifiedAt"`
		Outcome   string          `json:"outcome"`
		Actions   []string        `json:"actions"`
	}{result.Identity, result.Requested, result.Verified, result.Outcome, result.Actions})
	if err != nil {
		return result, err
	}
	result.Hash = evidence.Digest("proofrail:termination-evidence:1\n", canonical)
	if terminateErr != nil {
		return result, fmt.Errorf("%w: %v", ErrTerminationUncertain, terminateErr)
	}
	return result, nil
}

func waitIdentityGone(ctx context.Context, identity ProcessIdentity) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		alive, err := platformIdentityAlive(identity)
		if err != nil {
			return err
		}
		if !alive {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func formatTime(value time.Time) string { return value.Format("2006-01-02T15:04:05.000Z") }
