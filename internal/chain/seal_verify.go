package chain

// P43-T027B-STATE-SEAL: recovery verification of the C4 seal barrier.
//
// An accepted terminal claims that its closure is durable. Recovery therefore
// refuses to trust a PASSED task whose seal is missing, unreadable or does not
// carry the digest the accepted event recorded. A gap is reported as uncertainty;
// nothing reconstructs a seal from the event.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// verifyAcceptedSeals checks every PASSED task attempt against its seal record.
// A store that does not implement sealing carries no seal, so there is nothing to
// verify for it (the in-memory store used by unit tests is such a store); a
// sealing store must have a record for every accepted attempt.
func (engine *Engine) verifyAcceptedSeals(ctx context.Context, events []evidence.StateEvent) error {
	sink, ok := engine.options.Events.(SealSink)
	if !ok {
		return nil
	}
	accepted := map[string][]string{}
	for _, record := range events {
		event := record.Event
		if event.Entity.Kind != "task" || event.ToState != "PASSED" {
			continue
		}
		key := taskKey(event.Entity.TaskID, event.Entity.Attempt)
		accepted[key] = append([]string{}, event.InputEvidence...)
	}
	for key := range accepted {
		taskID, attempt, err := splitTaskKey(key)
		if err != nil {
			return errors.Join(ErrRecoveryUncertain, err)
		}
		_, hash, err := sink.LoadTaskSeal(ctx, engine.options.RunID, taskID, attempt)
		if err != nil {
			return errors.Join(ErrRecoveryUncertain, err)
		}
		if !containsHash(accepted[key], hash) {
			return errors.Join(ErrRecoveryUncertain, fmt.Errorf("%w: task %s attempt %d accepted without its seal", ErrSealMismatch, taskID, attempt))
		}
	}
	return nil
}

func splitTaskKey(key string) (string, int, error) {
	index := strings.LastIndex(key, "\x00")
	if index < 0 {
		return "", 0, fmt.Errorf("%w: malformed task key %q", ErrInvalidState, key)
	}
	attempt, err := strconv.Atoi(key[index+1:])
	if err != nil {
		return "", 0, fmt.Errorf("%w: malformed task key %q", ErrInvalidState, key)
	}
	return key[:index], attempt, nil
}

func containsHash(hashes []string, want string) bool {
	for _, hash := range hashes {
		if hash == want {
			return true
		}
	}
	return false
}
