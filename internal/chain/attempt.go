package chain

// P43-T027B-STATE-SEAL: the C5/N3 attempt source and the immutable context digest.
//
// The chain used to write attempt 1 into every request, receipt and entity key,
// which is only correct while a run has exactly one attempt. N3 shows the failure
// mode: with attempt 2 in flight next to an attempt 1 workspace the hard-coded
// value selects the wrong reopen identity. The resolver below is the single place
// that names the attempt, and the context key below is derived from immutable
// inputs only, so a repeated collection recomputes the same value instead of
// rewriting it.

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const contextKeyDomain = "proofrail:context-key:1\n"

// TaskAttemptResolver names the attempt a task is executing. An assembly that
// knows the durable retry state implements it; the engine fallback derives the
// attempt from the event projection, which yields 1 for a fresh task.
type TaskAttemptResolver interface {
	ResolveTaskAttempt(ctx context.Context, runID, taskID string) (int, error)
}

// resolveTaskAttempt returns the attempt this task attempt must carry.
func (engine *Engine) resolveTaskAttempt(ctx context.Context, taskID string) (int, error) {
	if engine.options.Attempts != nil {
		attempt, err := engine.options.Attempts.ResolveTaskAttempt(ctx, engine.options.RunID, taskID)
		if err != nil {
			return 0, err
		}
		if attempt < 1 {
			return 0, fmt.Errorf("%w: resolved attempt %d for task %q", ErrInvalidState, attempt, taskID)
		}
		return attempt, nil
	}
	return projectionAttempt(engine.state.projection, taskID), nil
}

// projectionAttempt reads the highest attempt the event log ever recorded for the
// task. A task that never started resolves to 1, which keeps every existing run
// byte-identical.
func projectionAttempt(projection Projection, taskID string) int {
	attempt := 0
	prefix := taskID + "\x00"
	for key := range projection.TaskStates {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		value, err := strconv.Atoi(strings.TrimPrefix(key, prefix))
		if err != nil {
			continue
		}
		if value > attempt {
			attempt = value
		}
	}
	if attempt < 1 {
		return 1
	}
	return attempt
}

// ContextKey derives the immutable context digest K. Every input is fixed by the
// time the step leaves PENDING: identity, resolved attempt, parent snapshot,
// workspace root and the dispatch identity. Recomputing K therefore yields the same
// value on re-entry, and an attempt can never inherit another attempt's context.
func ContextKey(runID, taskID, stepID string, attempt int, parentHash, workspaceRoot, executionTarget, requestID, requestHash string) string {
	fields := []string{
		runID,
		taskID,
		stepID,
		strconv.Itoa(attempt),
		parentHash,
		workspaceRoot,
		executionTarget,
		requestID,
		requestHash,
	}
	return evidence.Digest(contextKeyDomain, []byte(strings.Join(fields, "\x00")))
}

// factsRequestID and factsContextHash expose the dispatch identity K covers.
func factsRequestID(facts *AgentRunnerImmutableFacts) string {
	if facts == nil {
		return ""
	}
	return facts.RequestID
}

func factsContextHash(facts *AgentRunnerImmutableFacts) string {
	if facts == nil {
		return ""
	}
	return facts.ContextHash
}

// attemptsFor returns every attempt the projection knows for one task.
func (engine *Engine) taskAttempts(taskID string) []int {
	seen := map[int]bool{}
	prefix := taskID + "\x00"
	for key := range engine.state.projection.TaskStates {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		value, err := strconv.Atoi(strings.TrimPrefix(key, prefix))
		if err != nil {
			continue
		}
		seen[value] = true
	}
	if len(seen) == 0 {
		return []int{1}
	}
	attempts := make([]int, 0, len(seen))
	for attempt := range seen {
		attempts = append(attempts, attempt)
	}
	sort.Ints(attempts)
	return attempts
}

// TaskContextKey is the task-level form of K: it binds the run, the task, the
// resolved attempt, the accepted parent and the workspace root, so a seal can
// only be replayed by the attempt and context that produced it.
func TaskContextKey(runID, taskID string, attempt int, parentHash, workspaceRoot string) string {
	return ContextKey(runID, taskID, "", attempt, parentHash, workspaceRoot, "", "", "")
}
