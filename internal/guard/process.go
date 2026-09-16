package guard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrInvalidProcess       = errors.New("invalid managed process")
	ErrTerminationUncertain = errors.New("managed process termination uncertain")
	ErrProcessTreeRemained  = errors.New("managed process tree remained after parent exit")
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
	platform platformProcess
	// verify is a test seam: when set it replaces platform.verifyGone so a test
	// can observe the context each call site passes. Production leaves it nil.
	verify func(context.Context) error
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
	if err := platform.attach(cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		platform.close()
		return nil, fmt.Errorf("attach managed process tree: %w", err)
	}
	identity, err := platformIdentity(cmd.Process.Pid)
	if err != nil {
		_ = platform.terminate()
		_, _ = cmd.Process.Wait()
		platform.close()
		return nil, fmt.Errorf("capture managed process identity: %w", err)
	}
	return &ManagedProcess{cmd: cmd, identity: identity, platform: platform}, nil
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
	waited := make(chan error, 1)
	go func() { waited <- process.cmd.Wait() }()
	select {
	case waitErr := <-waited:
		verifyCtx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
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
			process.platform.close()
			proof, proofErr := buildTerminationEvidence(process.identity, requested, actions, stopErr)
			result := ProcessResult{ExitCode: -1, Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC(), Termination: &proof}
			if proofErr != nil {
				return result, proofErr
			}
			return result, ErrProcessTreeRemained
		}
		process.platform.close()
		return ProcessResult{ExitCode: process.cmd.ProcessState.ExitCode(), Started: true, StartedAt: startedAt, FinishedAt: time.Now().UTC()}, waitErr
	case <-ctx.Done():
		requested := time.Now().UTC()
		actions, stopErr := process.platform.stop(process.identity, grace)
		<-waited
		if stopErr == nil {
			// The caller's context is already cancelled on this path, so the
			// verification needs its own bounded context: it must not inherit the
			// cancellation whose handler is waiting on it, and it must not be
			// unbounded either.
			verifyCtx, verifyCancel := stopVerificationContext(grace)
			stopErr = process.verifyGone(verifyCtx)
			verifyCancel()
		}
		process.platform.close()
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
	actions, stopErr := platformStopIdentity(identity, grace)
	if stopErr == nil {
		waitCtx, cancel := boundedTerminationContext(ctx, grace)
		stopErr = waitIdentityGone(waitCtx, identity)
		cancel()
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
	actions, terminateErr := process.platform.stop(process.identity, grace)
	_, _ = process.cmd.Process.Wait()
	verifyCtx, verifyCancel := boundedTerminationContext(ctx, grace)
	defer verifyCancel()
	if terminateErr == nil {
		terminateErr = process.verifyGone(verifyCtx)
	}
	if terminateErr == nil {
		terminateErr = waitIdentityGone(verifyCtx, process.identity)
	}
	process.platform.close()
	return buildTerminationEvidence(process.identity, requested, actions, terminateErr)
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
