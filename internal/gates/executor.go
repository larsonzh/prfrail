package gates

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

type Executor interface {
	Capabilities() Capabilities
	Execute(context.Context, Prepared, Hook) Execution
}

type GuardExecutor struct{}

func (GuardExecutor) Capabilities() Capabilities {
	return Capabilities{Process: true, Timeout: true, OutputLimit: true}
}

func (GuardExecutor) Execute(ctx context.Context, prepared Prepared, hook Hook) Execution {
	ctx, cancel := context.WithTimeout(ctx, hook.Timeout)
	defer cancel()
	output := newBoundedOutput(hook.Resources.OutputBytes, cancel)
	result, err := guard.RunManaged(ctx, guard.ProcessSpec{Command: prepared.Executable, Args: prepared.Args, Dir: prepared.CWD, Env: prepared.Env, Stdout: output.stdout, Stderr: output.stderr}, 100*time.Millisecond)
	execution := Execution{FinishedAt: result.FinishedAt, Stdout: output.stdout.Bytes(), Stderr: output.stderr.Bytes(), OutputTruncated: output.exceeded(), RunnerEvidence: []string{runnerEvidence(prepared, hook, (GuardExecutor{}).Capabilities())}}
	if result.Started {
		execution.StartedAt = &result.StartedAt
	}
	if result.Termination != nil {
		execution.TerminationEvidence = []string{result.Termination.Hash}
	}
	switch {
	case !result.Started:
		execution.Outcome = "start-failed"
	case result.Termination != nil && errors.Is(err, guard.ErrTerminationUncertain):
		execution.Outcome = "termination-uncertain"
	case output.exceeded():
		execution.Outcome = "resource-limited"
	case errors.Is(err, guard.ErrProcessTreeRemained):
		execution.Outcome = "resource-limited"
	case errors.Is(err, context.DeadlineExceeded):
		execution.Outcome = "timed-out"
	default:
		execution.Outcome = "exited"
		exitCode := result.ExitCode
		execution.ExitCode = &exitCode
	}
	if err != nil {
		execution.ErrorEvidence = []string{evidence.Digest("proofrail:gate-error:1\n", []byte(err.Error()))}
	}
	return execution
}

type boundedOutput struct {
	stdout    *limitBuffer
	stderr    *limitBuffer
	mu        sync.Mutex
	remaining int64
	overflow  bool
	cancel    context.CancelFunc
}

func newBoundedOutput(limit int64, cancel context.CancelFunc) *boundedOutput {
	output := &boundedOutput{remaining: limit, cancel: cancel}
	output.stdout = &limitBuffer{owner: output}
	output.stderr = &limitBuffer{owner: output}
	return output
}

func (output *boundedOutput) write(buffer *bytes.Buffer, data []byte) (int, error) {
	output.mu.Lock()
	keep := int64(len(data))
	if keep > output.remaining {
		keep = output.remaining
		output.overflow = true
	}
	if keep > 0 {
		_, _ = buffer.Write(data[:keep])
		output.remaining -= keep
	}
	overflow := output.overflow
	output.mu.Unlock()
	if overflow {
		output.cancel()
	}
	return len(data), nil
}

func (output *boundedOutput) exceeded() bool {
	output.mu.Lock()
	defer output.mu.Unlock()
	return output.overflow
}

type limitBuffer struct {
	buffer bytes.Buffer
	owner  *boundedOutput
}

func (buffer *limitBuffer) Write(data []byte) (int, error) {
	return buffer.owner.write(&buffer.buffer, data)
}

func (buffer *limitBuffer) Bytes() []byte { return append([]byte(nil), buffer.buffer.Bytes()...) }

func runnerEvidence(prepared Prepared, hook Hook, capabilities Capabilities) string {
	canonical, _ := evidence.EncodeCanonical(struct {
		Executable   string         `json:"executable"`
		Args         []string       `json:"args"`
		CWD          string         `json:"cwd"`
		EnvAllowlist []string       `json:"envAllowlist"`
		Resources    ResourceLimits `json:"resourceLimits"`
		Network      NetworkPolicy  `json:"networkPolicy"`
		Capabilities Capabilities   `json:"capabilities"`
		OS           string         `json:"os"`
		Arch         string         `json:"arch"`
	}{prepared.Executable, prepared.Args, prepared.CWD, hook.EnvAllowlist, hook.Resources, hook.Network, capabilities, runtime.GOOS, runtime.GOARCH})
	return evidence.Digest("proofrail:gate-runner-evidence:1\n", canonical)
}
