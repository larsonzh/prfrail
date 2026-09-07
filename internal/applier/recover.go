package applier

import (
	"context"
	"fmt"

	"github.com/larsonzh/prfrail/internal/guard"
)

type RecoveryResult struct {
	JournalID string
	State     string
}

func Recover(ctx context.Context, options Options) ([]RecoveryResult, error) {
	if options.Fencer == nil {
		return nil, fmt.Errorf("writer lease fencer is required")
	}
	if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
		return nil, err
	}
	store, err := newJournalStore(options.JournalDir)
	if err != nil {
		return nil, err
	}
	paths, err := store.list()
	if err != nil {
		return nil, err
	}
	results := make([]RecoveryResult, 0, len(paths))
	for _, path := range paths {
		record, err := store.read(path)
		if err != nil {
			return results, err
		}
		if record.Journal.ResourceID != options.ResourceID || record.Journal.Token != options.Token {
			return results, guard.ErrFencingTokenMismatch
		}
		switch record.Journal.State {
		case "prepared", "applying", "rolling-back", "uncertain":
			if err := rollback(ctx, store, &record, options); err != nil {
				return results, fmt.Errorf("%w: %v", ErrRecoveryUncertain, err)
			}
		case "applied":
			for _, entry := range record.Journal.Entries {
				expected, err := journalState(store, entry, true)
				if err != nil {
					return results, err
				}
				if err := verifyState(options.Root, entry.Path, expected); err != nil {
					return results, fmt.Errorf("%w: applied journal differs: %v", ErrRecoveryUncertain, err)
				}
			}
		case "rolled-back":
		default:
			return results, fmt.Errorf("%w: unknown state %q", ErrInvalidJournal, record.Journal.State)
		}
		results = append(results, RecoveryResult{JournalID: record.Journal.JournalID, State: record.Journal.State})
	}
	return results, nil
}
