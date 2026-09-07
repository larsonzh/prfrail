package guard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrInvalidProcess       = errors.New("invalid managed process")
	ErrTerminationUncertain = errors.New("managed process termination uncertain")
)

type ProcessSpec struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
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

func (process *ManagedProcess) Identity() ProcessIdentity { return process.identity }

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
	outcome := "stopped"
	if terminateErr != nil {
		outcome = "uncertain"
	}
	result := TerminationEvidence{
		Identity:  process.identity,
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
