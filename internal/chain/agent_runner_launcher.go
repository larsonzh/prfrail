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
// rely on, plus the additive process parameters A6 introduced: what to run
// (Command/Args/Env), where to run it (Dir, and WorkspaceRoot for the manifest
// capture), and how much grace it gets when it has to be stopped (Grace). The run
// bound is not part of a launch: the caller that waits owns the timeout, so a
// launcher cannot silently disagree with it. Extending these fields never changes
// replay/receipt semantics: a launcher must still return an error only when no
// process was spawned.
type AgentRunnerLaunchRequest struct {
	RequestID string
	RunID     string
	AdapterID string
	// Process parameters (additive, A6).
	Command       string
	Args          []string
	Dir           string
	Env           []string
	Grace         time.Duration
	WorkspaceRoot string
}

// AgentRunnerLaunchResult reports a spawned process. LaunchID must equal the
// launch-intent launchId; ProcessID is an opaque platform process identity.
type AgentRunnerLaunchResult struct {
	LaunchID  string
	ProcessID string
	StartedAt time.Time
}
