package applier

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/taskdef"
)

type Fencer interface {
	Validate(context.Context, string, guard.FencingToken) error
}
type Failpoint func(point string, entry int) error

type Options struct {
	Root       string
	JournalDir string
	ResourceID string
	Token      guard.FencingToken
	Fencer     Fencer
	Failpoint  Failpoint
}

type Result struct {
	JournalID string
	State     string
}

func Apply(ctx context.Context, plan taskdef.Plan, options Options) (Result, error) {
	if options.Fencer == nil {
		return Result{}, errors.New("writer lease fencer is required")
	}
	if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
		return Result{}, err
	}
	store, err := newJournalStore(options.JournalDir)
	if err != nil {
		return Result{}, err
	}
	record, err := store.create(plan.ChangeSetHash, options.ResourceID, options.Token, plan, options.Root, func() error {
		return options.Fencer.Validate(ctx, options.ResourceID, options.Token)
	})
	if err != nil {
		return Result{}, err
	}
	if err := inject(options.Failpoint, "journal-durable", -1); err != nil {
		return injectedFailure(ctx, store, &record, options, err)
	}
	if err := fencedUpdate(ctx, store, &record, "applying", 0, options); err != nil {
		return Result{}, err
	}
	for index, entry := range plan.Entries {
		if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
			return rollbackFailure(ctx, store, &record, options, err)
		}
		if err := verifyState(options.Root, entry.Path, entry.Before); err != nil {
			return rollbackFailure(ctx, store, &record, options, err)
		}
		if err := inject(options.Failpoint, "before-write", index); err != nil {
			return injectedFailure(ctx, store, &record, options, err)
		}
		if err := writeState(options.Root, entry.Path, entry.After); err != nil {
			return rollbackFailure(ctx, store, &record, options, err)
		}
		if err := fencedUpdate(ctx, store, &record, "applying", index+1, options); err != nil {
			return rollbackFailure(ctx, store, &record, options, err)
		}
		if err := inject(options.Failpoint, "after-write", index); err != nil {
			return injectedFailure(ctx, store, &record, options, err)
		}
	}
	for _, entry := range plan.Entries {
		if err := verifyState(options.Root, entry.Path, entry.After); err != nil {
			return rollbackFailure(ctx, store, &record, options, err)
		}
	}
	if err := inject(options.Failpoint, "after-verify", len(plan.Entries)); err != nil {
		return injectedFailure(ctx, store, &record, options, err)
	}
	if err := fencedUpdate(ctx, store, &record, "applied", len(plan.Entries), options); err != nil {
		return rollbackFailure(ctx, store, &record, options, err)
	}
	return Result{JournalID: record.Journal.JournalID, State: "applied"}, nil
}

func injectedFailure(ctx context.Context, store *journalStore, record *journalRecord, options Options, err error) (Result, error) {
	if errors.Is(err, ErrCrashInjected) {
		return Result{JournalID: record.Journal.JournalID, State: record.Journal.State}, err
	}
	return rollbackFailure(ctx, store, record, options, err)
}

func rollbackFailure(ctx context.Context, store *journalStore, record *journalRecord, options Options, cause error) (Result, error) {
	if err := rollback(ctx, store, record, options); err != nil {
		return Result{JournalID: record.Journal.JournalID, State: "uncertain"}, errors.Join(cause, ErrRecoveryUncertain, err)
	}
	return Result{JournalID: record.Journal.JournalID, State: "rolled-back"}, cause
}

func rollback(ctx context.Context, store *journalStore, record *journalRecord, options Options) error {
	if err := fencedUpdate(ctx, store, record, "rolling-back", record.Journal.Completed, options); err != nil {
		return err
	}
	for index := len(record.Journal.Entries) - 1; index >= 0; index-- {
		entry := record.Journal.Entries[index]
		before, err := journalState(store, entry, false)
		if err != nil {
			_ = store.update(record, "uncertain", record.Journal.Completed)
			return err
		}
		after, err := journalState(store, entry, true)
		if err != nil {
			_ = store.update(record, "uncertain", record.Journal.Completed)
			return err
		}
		if verifyState(options.Root, entry.Path, before) == nil {
			continue
		}
		if verifyState(options.Root, entry.Path, after) != nil {
			_ = fencedUpdate(ctx, store, record, "uncertain", record.Journal.Completed, options)
			return fmt.Errorf("%w: %q has unknown external content", ErrRecoveryUncertain, entry.Path)
		}
		if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
			return err
		}
		if err := writeState(options.Root, entry.Path, before); err != nil {
			_ = fencedUpdate(ctx, store, record, "uncertain", record.Journal.Completed, options)
			return err
		}
	}
	for index := len(record.Journal.CreatedDirs) - 1; index >= 0; index-- {
		if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
			return err
		}
		_ = os.Remove(filepath.Join(options.Root, filepath.FromSlash(record.Journal.CreatedDirs[index])))
	}
	return fencedUpdate(ctx, store, record, "rolled-back", 0, options)
}

func fencedUpdate(ctx context.Context, store *journalStore, record *journalRecord, state string, completed int, options Options) error {
	if err := options.Fencer.Validate(ctx, options.ResourceID, options.Token); err != nil {
		return err
	}
	return store.update(record, state, completed)
}

func journalState(store *journalStore, entry journalEntry, after bool) (taskdef.FileState, error) {
	hash, mode := entry.BeforeHash, entry.BeforeMode
	if after {
		hash, mode = entry.AfterHash, entry.AfterMode
	}
	if hash == nil {
		return taskdef.FileState{}, nil
	}
	data, err := store.getBlob(*hash)
	if err != nil {
		return taskdef.FileState{}, err
	}
	return taskdef.FileState{Exists: true, Bytes: data, Mode: os.FileMode(mode)}, nil
}

func verifyState(root, relative string, expected taskdef.FileState) error {
	actual, err := readState(root, relative)
	if err != nil {
		return err
	}
	if actual.Exists != expected.Exists || (actual.Exists && (evidence.Digest("", actual.Bytes) != evidence.Digest("", expected.Bytes) || actual.Mode.Perm() != expected.Mode.Perm())) {
		return fmt.Errorf("%w: %q", taskdef.ErrTargetChanged, relative)
	}
	return nil
}

func readState(root, relative string) (taskdef.FileState, error) {
	path := filepath.Join(root, filepath.FromSlash(relative))
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return taskdef.FileState{}, nil
	}
	if err != nil {
		return taskdef.FileState{}, err
	}
	if !info.Mode().IsRegular() {
		return taskdef.FileState{}, fmt.Errorf("target %q is not regular", relative)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return taskdef.FileState{}, err
	}
	return taskdef.FileState{Exists: true, Bytes: data, Mode: info.Mode().Perm()}, nil
}

func writeState(root, relative string, state taskdef.FileState) error {
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := ensureSafeParents(root, relative); err != nil {
		return err
	}
	if !state.Exists {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return syncDirectory(filepath.Dir(path))
	}
	mode := state.Mode.Perm()
	if mode == 0 {
		mode = 0644
	}
	return durableReplace(path, state.Bytes, mode)
}

func ensureSafeParents(root, relative string) error {
	parent := filepath.Dir(filepath.Join(root, filepath.FromSlash(relative)))
	root = filepath.Clean(root)
	for current := parent; current != root && len(current) >= len(root); current = filepath.Dir(current) {
		if info, err := os.Lstat(current); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink parent %q", current)
		} else if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return os.MkdirAll(parent, 0755)
}

func inject(failpoint Failpoint, point string, entry int) error {
	if failpoint == nil {
		return nil
	}
	return failpoint(point, entry)
}
