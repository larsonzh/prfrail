package console

import (
	"context"
	"fmt"
	"io"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// progressPlaceholder stands in for an empty field so every progress line keeps
// the same column shape.
const progressPlaceholder = "-"

// progressStore decorates the durable event store with a one-line progress trace.
//
// P70-S1: --progress streams one line per state event to the writer resolved from
// the CLI flag. The decorator is a pure observer. The event is appended first and
// the line is printed only after that append succeeded, so the durable record and
// the returned error are exactly what the undecorated store produced. A failing
// writer is ignored on purpose: progress must never change the outcome of a run.
type progressStore struct {
	inner *chain.FileEventStore
	out   io.Writer
}

var (
	_ chain.EventStore = (*progressStore)(nil)
	_ chain.SealSink   = (*progressStore)(nil)
)

func newProgressStore(inner *chain.FileEventStore, out io.Writer) *progressStore {
	return &progressStore{inner: inner, out: out}
}

// Append persists the event first and only then streams its progress line, so an
// event that never reached the log is never reported as progress. The event is
// forwarded untouched: the decorator does not re-marshal, reorder or cache it.
func (s *progressStore) Append(ctx context.Context, event evidence.StateEvent) error {
	if err := s.inner.Append(ctx, event); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(s.out, formatProgressLine(event))
	return nil
}

// Load forwards to the durable store unchanged.
func (s *progressStore) Load(ctx context.Context) ([]evidence.StateEvent, error) {
	return s.inner.Load(ctx)
}

// SealTaskEvidence forwards to the durable store unchanged. The decorator must
// keep the SealSink capability: internal/chain type-asserts Options.Events for the
// seal barrier, and a wrapper without it would silently disable the seal.
func (s *progressStore) SealTaskEvidence(ctx context.Context, record chain.SealRecord) (string, error) {
	return s.inner.SealTaskEvidence(ctx, record)
}

// LoadTaskSeal forwards to the durable store unchanged.
func (s *progressStore) LoadTaskSeal(ctx context.Context, runID, taskID string, attempt int) (chain.SealRecord, string, error) {
	return s.inner.LoadTaskSeal(ctx, runID, taskID, attempt)
}

// formatProgressLine renders one state event as a single fixed-shape line:
//
//	progress: <occurredAt> <entity> <fromState>-><toState> <reasonCode>
func formatProgressLine(event evidence.StateEvent) string {
	body := event.Event
	return fmt.Sprintf("progress: %s %s %s->%s %s",
		progressField(body.OccurredAt),
		progressEntity(body.Entity),
		progressField(body.FromState),
		progressField(body.ToState),
		progressField(body.Reason.Code),
	)
}

// progressEntity names the entity a transition belongs to. A chain event names
// itself, a task event carries its task id and a step event carries both ids.
func progressEntity(entity evidence.Entity) string {
	switch entity.Kind {
	case "chain":
		return "chain"
	case "task":
		return "task " + progressField(entity.TaskID)
	case "step":
		return "step " + progressField(entity.TaskID) + "/" + progressField(entity.StepID)
	default:
		return progressField(entity.Kind)
	}
}

func progressField(value string) string {
	if value == "" {
		return progressPlaceholder
	}
	return value
}
