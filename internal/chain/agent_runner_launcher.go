package chain

import (
	"context"
	"time"
)

// AgentRunnerLauncher owns the process start port on the dispatch path.
// The port contract is strict: an implementation returns an error only when no
// process was spawned, and returns a result only when the process was spawned
// and its identity is populated. Callers may replay Start after an error (the
// process was never started) but must never treat a returned result as
// retryable.
type AgentRunnerLauncher interface {
	StartAgentRunnerProcess(ctx context.Context, launch AgentRunnerLaunchRequest) (AgentRunnerLaunchResult, error)
}

// AgentRunnerLaunchRequest carries the immutable launch identity a launcher may
// rely on. Later slices extend it with argv digests, cwd, environment
// allowlists, and guard parameters without changing replay/receipt semantics.
type AgentRunnerLaunchRequest struct {
	RequestID string
	RunID     string
	AdapterID string
}

// AgentRunnerLaunchResult reports a spawned process. LaunchID must equal the
// launch-intent launchId; ProcessID is an opaque platform process identity.
type AgentRunnerLaunchResult struct {
	LaunchID  string
	ProcessID string
	StartedAt time.Time
}
