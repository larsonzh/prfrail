package guard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
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
}

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
	cmd      *exec.Cmd
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
			process.waitErr = process.cmd.Wait()
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

func StartManaged(ctx context.Context, spec ProcessSpec) (*ManagedProcess, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if spec.Command == "" {
		return nil, fmt.Errorf("%w: command is empty", ErrInvalidProcess)
	}
	cmd := exec.Command(spec.Command, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = spec.Env
	cmd.Stdout = spec.Stdout
	cmd.Stderr = spec.Stderr
	platform, err := preparePlatformProcess(cmd)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		platform.close()
		return nil, fmt.Errorf("start managed process: %w", err)
	}
	process := &ManagedProcess{cmd: cmd, platform: &platform, waitDone: make(chan struct{})}
	if err := platform.attach(cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		return nil, process.abortStartup(fmt.Errorf("attach managed process tree: %w", err))
	}
	identity, err := platformIdentity(cmd.Process.Pid)
	if err != nil {
		_ = platform.terminate()
		return nil, process.abortStartup(fmt.Errorf("capture managed process identity: %w", err))
	}
	process.identity = identity
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
		return ProcessResult{ExitCode: process.cmd.ProcessState.ExitCode(), Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC()}, process.waitErr
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
