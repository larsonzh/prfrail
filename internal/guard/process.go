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
	waited := make(chan error, 1)
	go func() { waited <- process.cmd.Wait() }()
	select {
	case waitErr := <-waited:
		verifyCtx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		verifyErr := process.platform.verifyGone(verifyCtx)
		cancel()
		if verifyErr != nil {
			requested := time.Now().UTC()
			actions, stopErr := process.platform.stop(process.identity, grace)
			if stopErr == nil {
				finalCtx, finalCancel := context.WithTimeout(context.Background(), grace)
				stopErr = process.platform.verifyGone(finalCtx)
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
			stopErr = process.platform.verifyGone(context.Background())
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
		stopErr = waitIdentityGone(ctx, identity)
	}
	proof, proofErr := buildTerminationEvidence(identity, requested, actions, stopErr)
	if proofErr != nil {
		return proof, proofErr
	}
	return proof, nil
}

func (process *ManagedProcess) Terminate(ctx context.Context, grace time.Duration) (TerminationEvidence, error) {
	requested := time.Now().UTC()
	actions, terminateErr := process.platform.stop(process.identity, grace)
	_, _ = process.cmd.Process.Wait()
	if terminateErr == nil {
		terminateErr = process.platform.verifyGone(ctx)
	}
	if terminateErr == nil {
		terminateErr = waitIdentityGone(ctx, process.identity)
	}
	process.platform.close()
	return buildTerminationEvidence(process.identity, requested, actions, terminateErr)
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
