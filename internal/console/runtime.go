package console

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const timestampLayout = "2006-01-02T15:04:05.000Z"

type RunSummary struct {
	RunID           string       `json:"runId"`
	ChainState      string       `json:"chainState"`
	Sequence        int          `json:"sequence"`
	LastEventHash   string       `json:"lastEventHash"`
	TaskStateCounts []StateCount `json:"taskStateCounts"`
	StepStateCounts []StateCount `json:"stepStateCounts"`
	EventLogPath    string       `json:"eventLogPath"`
}

type StateCount struct {
	State string `json:"state"`
	Count int    `json:"count"`
}

func ExecuteNoopRun(ctx context.Context, runID, runDir string, definition chain.Definition, workspaceRoot string, now func() time.Time) (RunSummary, error) {
	store, err := chain.NewFileEventStore(runEventLogPath(runDir))
	if err != nil {
		return RunSummary{}, err
	}

	runtime := &localRuntime{runDir: runDir, workspaceRoot: workspaceRoot, now: now}
	options := chain.Options{
		RunID:      runID,
		Definition: definition,
		Events:     store,
		Baselines:  runtime,
		Workspaces: runtime,
		Steps:      runtime,
		Acceptance: runtime,
		Reviewer:   runtime,
		Publisher:  runtime,
		Postflight: runtime,
		Stopper:    runtime,
		Reconciler: runtime,
		Clock:      now,
		IDs:        runtime.nextEventID,
	}
	engine, err := chain.New(ctx, options)
	if err != nil {
		return RunSummary{}, err
	}
	if err := engine.Run(ctx); err != nil {
		return RunSummary{}, err
	}
	return summarizeProjection(engine.Projection(), runEventLogPath(runDir)), nil
}

func LoadRunSummary(ctx context.Context, runDir string) (RunSummary, error) {
	path := runEventLogPath(runDir)
	if _, err := os.Stat(path); err != nil {
		return RunSummary{}, err
	}
	store, err := chain.NewFileEventStore(path)
	if err != nil {
		return RunSummary{}, err
	}
	events, err := store.Load(ctx)
	if err != nil {
		return RunSummary{}, err
	}
	if len(events) == 0 {
		return RunSummary{}, fmt.Errorf("event log %q is empty", path)
	}
	projection, err := chain.Rebuild(events)
	if err != nil {
		return RunSummary{}, err
	}
	return summarizeProjection(projection, path), nil
}

func runEventLogPath(runDir string) string {
	return filepath.Join(runDir, "events", "state-events.jsonl")
}

func summarizeProjection(projection chain.Projection, eventLogPath string) RunSummary {
	return RunSummary{
		RunID:           projection.RunID,
		ChainState:      projection.ChainState,
		Sequence:        projection.Sequence,
		LastEventHash:   projection.LastEventHash,
		TaskStateCounts: countStates(projection.TaskStates),
		StepStateCounts: countStates(projection.StepStates),
		EventLogPath:    eventLogPath,
	}
}

func countStates(states map[string]string) []StateCount {
	if len(states) == 0 {
		return []StateCount{}
	}
	counters := map[string]int{}
	for _, state := range states {
		counters[state]++
	}
	names := make([]string, 0, len(counters))
	for name := range counters {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]StateCount, 0, len(names))
	for _, name := range names {
		result = append(result, StateCount{State: name, Count: counters[name]})
	}
	return result
}

func digest(domain string, values ...string) string {
	return evidence.Digest(domain, []byte(strings.Join(values, "\n")))
}
