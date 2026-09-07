package applier

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/taskdef"
)

type fixedFencer struct{ err error }

func (f fixedFencer) Validate(context.Context, string, guard.FencingToken) error { return f.err }

type expiringFencer struct{ calls, failAt int }

func (f *expiringFencer) Validate(context.Context, string, guard.FencingToken) error {
	f.calls++
	if f.calls >= f.failAt {
		return guard.ErrFencingTokenMismatch
	}
	return nil
}

func TestApplyCommitsCreateReplaceDelete(t *testing.T) {
	root, _, plan, options := fixture(t)
	result, err := Apply(context.Background(), plan, options)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "applied" {
		t.Fatalf("unexpected result: %+v", result)
	}
	assertContent(t, filepath.Join(root, "modify.txt"), "after\n")
	assertContent(t, filepath.Join(root, "nested", "create.txt"), "created\n")
	if _, err := os.Stat(filepath.Join(root, "delete.txt")); !os.IsNotExist(err) {
		t.Fatalf("delete target remains: %v", err)
	}
	recovered, err := Recover(context.Background(), options)
	if err != nil || len(recovered) != 1 || recovered[0].State != "applied" {
		t.Fatalf("recover applied: %+v %v", recovered, err)
	}
}

func TestApplyCrashPointsRecoverWholeTransaction(t *testing.T) {
	points := []struct {
		point string
		entry int
	}{{"journal-durable", -1}, {"before-write", 0}, {"after-write", 0}, {"before-write", 1}, {"after-write", 1}, {"after-verify", 3}}
	for _, item := range points {
		t.Run(item.point, func(t *testing.T) {
			root, _, plan, options := fixture(t)
			options.Failpoint = func(point string, entry int) error {
				if point == item.point && (item.entry < 0 || entry == item.entry) {
					return ErrCrashInjected
				}
				return nil
			}
			if _, err := Apply(context.Background(), plan, options); !errors.Is(err, ErrCrashInjected) {
				t.Fatalf("expected crash, got %v", err)
			}
			options.Failpoint = nil
			results, err := Recover(context.Background(), options)
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != 1 || results[0].State != "rolled-back" {
				t.Fatalf("unexpected recovery: %+v", results)
			}
			assertContent(t, filepath.Join(root, "modify.txt"), "before\n")
			assertContent(t, filepath.Join(root, "delete.txt"), "delete\n")
			if _, err := os.Stat(filepath.Join(root, "nested", "create.txt")); !os.IsNotExist(err) {
				t.Fatalf("created file remains: %v", err)
			}
		})
	}
}

func TestApplyRejectsStaleFenceBeforeJournal(t *testing.T) {
	_, journalDir, plan, options := fixture(t)
	options.Fencer = fixedFencer{err: guard.ErrFencingTokenMismatch}
	if _, err := Apply(context.Background(), plan, options); !errors.Is(err, guard.ErrFencingTokenMismatch) {
		t.Fatalf("expected fencing error: %v", err)
	}
	items, _ := os.ReadDir(journalDir)
	if len(items) != 0 {
		t.Fatalf("journal written for stale fence: %v", items)
	}
}

func TestRecoverRefusesUnknownExternalContent(t *testing.T) {
	root, _, plan, options := fixture(t)
	options.Failpoint = func(point string, entry int) error {
		if point == "after-write" && entry == 0 {
			return ErrCrashInjected
		}
		return nil
	}
	if _, err := Apply(context.Background(), plan, options); !errors.Is(err, ErrCrashInjected) {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "modify.txt"), []byte("external\n"), 0644); err != nil {
		t.Fatal(err)
	}
	options.Failpoint = nil
	if _, err := Recover(context.Background(), options); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("expected uncertain recovery, got %v", err)
	}
	assertContent(t, filepath.Join(root, "modify.txt"), "external\n")
}

func TestApplyOrdinaryFailureRollsBackWholeTransaction(t *testing.T) {
	root, _, plan, options := fixture(t)
	ordinary := errors.New("ordinary write failure")
	options.Failpoint = func(point string, entry int) error {
		if point == "after-write" && entry == 1 {
			return ordinary
		}
		return nil
	}
	result, err := Apply(context.Background(), plan, options)
	if !errors.Is(err, ordinary) || result.State != "rolled-back" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	assertContent(t, filepath.Join(root, "modify.txt"), "before\n")
	assertContent(t, filepath.Join(root, "delete.txt"), "delete\n")
	if _, err := os.Stat(filepath.Join(root, "nested", "create.txt")); !os.IsNotExist(err) {
		t.Fatalf("partial create remains: %v", err)
	}
}

func TestRecoverIsIdempotent(t *testing.T) {
	_, _, plan, options := fixture(t)
	options.Failpoint = func(point string, entry int) error {
		if point == "after-write" && entry == 0 {
			return ErrCrashInjected
		}
		return nil
	}
	if _, err := Apply(context.Background(), plan, options); !errors.Is(err, ErrCrashInjected) {
		t.Fatal(err)
	}
	options.Failpoint = nil
	for iteration := 0; iteration < 2; iteration++ {
		results, err := Recover(context.Background(), options)
		if err != nil || len(results) != 1 || results[0].State != "rolled-back" {
			t.Fatalf("iteration %d: %+v %v", iteration, results, err)
		}
	}
}

func TestRecoverRejectsTamperedJournalAndMissingBlob(t *testing.T) {
	for _, testCase := range []string{"journal", "blob"} {
		t.Run(testCase, func(t *testing.T) {
			_, journalDir, plan, options := fixture(t)
			options.Failpoint = func(point string, entry int) error {
				if point == "journal-durable" {
					return ErrCrashInjected
				}
				return nil
			}
			result, err := Apply(context.Background(), plan, options)
			if !errors.Is(err, ErrCrashInjected) {
				t.Fatal(err)
			}
			if testCase == "journal" {
				path := filepath.Join(journalDir, result.JournalID+".journal.json")
				data, readErr := os.ReadFile(path)
				if readErr != nil {
					t.Fatal(readErr)
				}
				data[len(data)-2] ^= 1
				if writeErr := os.WriteFile(path, data, 0600); writeErr != nil {
					t.Fatal(writeErr)
				}
			} else {
				blobs, globErr := filepath.Glob(filepath.Join(journalDir, "*.blob"))
				if globErr != nil || len(blobs) == 0 {
					t.Fatalf("blobs: %v %v", blobs, globErr)
				}
				if removeErr := os.Remove(blobs[0]); removeErr != nil {
					t.Fatal(removeErr)
				}
			}
			options.Failpoint = nil
			if _, recoverErr := Recover(context.Background(), options); recoverErr == nil {
				t.Fatal("accepted damaged recovery evidence")
			}
		})
	}
}

func TestFencingExpiresBeforeJournalProgressWrite(t *testing.T) {
	root, _, plan, options := fixture(t)
	fencer := &expiringFencer{failAt: 9}
	options.Fencer = fencer
	if _, err := Apply(context.Background(), plan, options); !errors.Is(err, guard.ErrFencingTokenMismatch) {
		t.Fatalf("expected fencing failure: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "delete.txt")); !os.IsNotExist(err) {
		t.Fatalf("stale writer continued unexpectedly: %v", err)
	}
	options.Fencer = fixedFencer{}
	results, err := Recover(context.Background(), options)
	if err != nil || len(results) != 1 || results[0].State != "rolled-back" {
		t.Fatalf("recovery: %+v %v", results, err)
	}
	assertContent(t, filepath.Join(root, "delete.txt"), "delete\n")
}

func TestRecoverRejectsWrongTokenAndRehashedTraversal(t *testing.T) {
	for _, testCase := range []string{"wrong-token", "traversal"} {
		t.Run(testCase, func(t *testing.T) {
			_, journalDir, plan, options := fixture(t)
			options.Failpoint = func(point string, entry int) error {
				if point == "journal-durable" {
					return ErrCrashInjected
				}
				return nil
			}
			result, err := Apply(context.Background(), plan, options)
			if !errors.Is(err, ErrCrashInjected) {
				t.Fatal(err)
			}
			options.Failpoint = nil
			if testCase == "wrong-token" {
				options.Token.Generation++
				if _, err := Recover(context.Background(), options); !errors.Is(err, guard.ErrFencingTokenMismatch) {
					t.Fatalf("expected token rejection: %v", err)
				}
				return
			}
			store, _ := newJournalStore(journalDir)
			record, readErr := store.read(filepath.Join(journalDir, result.JournalID+".journal.json"))
			if readErr != nil {
				t.Fatal(readErr)
			}
			record.Journal.Entries[0].Path = "../escape.txt"
			rehashed, makeErr := makeJournal(record.Journal)
			if makeErr != nil {
				t.Fatal(makeErr)
			}
			if writeErr := store.write(rehashed); writeErr != nil {
				t.Fatal(writeErr)
			}
			if _, recoverErr := Recover(context.Background(), options); !errors.Is(recoverErr, ErrInvalidJournal) {
				t.Fatalf("accepted traversal journal: %v", recoverErr)
			}
		})
	}
}

func TestApplyRejectsSymlinkParent(t *testing.T) {
	root, journalDir := t.TempDir(), t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "victim.txt"), []byte("before\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	info, _ := os.Stat(filepath.Join(outside, "victim.txt"))
	before, after := []byte("before\n"), []byte("after\n")
	plan := taskdef.Plan{ChangeSetHash: evidence.Digest("", []byte(t.Name())), Entries: []taskdef.PlanEntry{{Path: "link/victim.txt", Before: taskdef.FileState{Exists: true, Bytes: before, Mode: info.Mode().Perm()}, After: taskdef.FileState{Exists: true, Bytes: after, Mode: info.Mode().Perm()}}}}
	options := Options{Root: root, JournalDir: journalDir, ResourceID: "workspace-one", Token: guard.FencingToken{ClaimID: "claim-one", Generation: 1}, Fencer: fixedFencer{}}
	if _, err := Apply(context.Background(), plan, options); err == nil {
		t.Fatal("accepted symlink parent")
	}
	assertContent(t, filepath.Join(outside, "victim.txt"), "before\n")
}

func fixture(t *testing.T) (string, string, taskdef.Plan, Options) {
	t.Helper()
	root := t.TempDir()
	journalDir := t.TempDir()
	beforeModify, afterModify := []byte("before\n"), []byte("after\n")
	beforeDelete, afterCreate := []byte("delete\n"), []byte("created\n")
	if err := os.WriteFile(filepath.Join(root, "modify.txt"), beforeModify, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "delete.txt"), beforeDelete, 0644); err != nil {
		t.Fatal(err)
	}
	modifyInfo, err := os.Stat(filepath.Join(root, "modify.txt"))
	if err != nil {
		t.Fatal(err)
	}
	deleteInfo, err := os.Stat(filepath.Join(root, "delete.txt"))
	if err != nil {
		t.Fatal(err)
	}
	plan := taskdef.Plan{ChangeSetHash: evidence.Digest("", []byte(t.Name())), Entries: []taskdef.PlanEntry{
		{Path: "delete.txt", Before: taskdef.FileState{Exists: true, Bytes: beforeDelete, Mode: deleteInfo.Mode().Perm()}, After: taskdef.FileState{}},
		{Path: "modify.txt", Before: taskdef.FileState{Exists: true, Bytes: beforeModify, Mode: modifyInfo.Mode().Perm()}, After: taskdef.FileState{Exists: true, Bytes: afterModify, Mode: modifyInfo.Mode().Perm()}},
		{Path: "nested/create.txt", Before: taskdef.FileState{}, After: taskdef.FileState{Exists: true, Bytes: afterCreate, Mode: createMode()}},
	}}
	token := guard.FencingToken{ClaimID: "claim-one", Generation: 1}
	options := Options{Root: root, JournalDir: journalDir, ResourceID: "workspace-one", Token: token, Fencer: fixedFencer{}}
	return root, journalDir, plan, options
}

func assertContent(t *testing.T, path, expected string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != expected {
		t.Fatalf("%s = %q, %v", path, data, err)
	}
}

func createMode() os.FileMode {
	probe, err := os.CreateTemp("", "prfrail-mode-")
	if err != nil {
		return 0644
	}
	path := probe.Name()
	_ = probe.Close()
	defer os.Remove(path)
	_ = os.Chmod(path, 0644)
	info, err := os.Stat(path)
	if err != nil {
		return 0644
	}
	return info.Mode().Perm()
}
