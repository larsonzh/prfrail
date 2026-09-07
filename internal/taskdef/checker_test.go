package taskdef

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestCheckSimulatesOrderedOperationsWithoutWriting(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	original := []byte("alpha\nomega\n")
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}
	match, replacement, lineEnding := "alpha", "beta", "preserve-existing"
	expected := 1
	before := evidence.Digest("", original)
	afterBytes := []byte("beta\nomega\n")
	after := evidence.Digest("", afterBytes)
	body := validBody([]Operation{{OperationID: "operation-one", Sequence: 1, Kind: "replace-exact", Path: "file.txt", BeforeHash: &before, AfterHash: &after, Assertions: []Assertion{{Phase: "before", Kind: "exact-occurrences", Text: "alpha", Count: 1}, {Phase: "after", Kind: "absent", Text: "alpha", Count: 0}}, Match: &match, Replacement: &replacement, ExpectedMatches: &expected, LineEnding: &lineEnding}})
	record, _ := NewChangeSet(body)
	plan, err := Check(record, validOptions(root, body))
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Entries) != 1 || string(plan.Entries[0].After.Bytes) != string(afterBytes) {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	actual, _ := os.ReadFile(path)
	if string(actual) != string(original) {
		t.Fatal("checker wrote the workspace")
	}
}

func TestCheckAllOperationsAndSamePathChain(t *testing.T) {
	root := t.TempDir()
	deleted := []byte("remove\n")
	if err := os.WriteFile(filepath.Join(root, "delete.txt"), deleted, 0644); err != nil {
		t.Fatal(err)
	}
	lf, preserve := "lf", "preserve-existing"
	content, anchor := "anchor\n", "anchor"
	beforeInsert, afterInsert := "head ", " tail"
	match, replacement := "anchor", "center"
	one := 1
	h1 := evidence.Digest("", []byte(content))
	h2 := evidence.Digest("", []byte("head anchor\n"))
	h3 := evidence.Digest("", []byte("head anchor tail\n"))
	h4 := evidence.Digest("", []byte("head center tail\n"))
	deleteHash := evidence.Digest("", deleted)
	body := validBody([]Operation{
		{OperationID: "create-one", Sequence: 1, Kind: "create-file", Path: "chain.txt", AfterHash: &h1, Assertions: assertions("anchor", 0, 1), Content: &content, LineEnding: &lf},
		{OperationID: "insert-before-one", Sequence: 2, Kind: "insert-before", Path: "chain.txt", BeforeHash: &h1, AfterHash: &h2, Assertions: assertions("head ", 0, 1), Content: &beforeInsert, LineEnding: &preserve, ExpectedMatches: &one, Anchor: &anchor, Marker: &Marker{MarkerID: "marker-before", Text: "head ", BeforeCount: 0, AfterCount: 1}},
		{OperationID: "insert-after-one", Sequence: 3, Kind: "insert-after", Path: "chain.txt", BeforeHash: &h2, AfterHash: &h3, Assertions: assertions(" tail", 0, 1), Content: &afterInsert, LineEnding: &preserve, ExpectedMatches: &one, Anchor: &anchor, Marker: &Marker{MarkerID: "marker-after", Text: " tail", BeforeCount: 0, AfterCount: 1}},
		{OperationID: "replace-one", Sequence: 4, Kind: "replace-exact", Path: "chain.txt", BeforeHash: &h3, AfterHash: &h4, Assertions: assertions("center", 0, 1), Match: &match, Replacement: &replacement, ExpectedMatches: &one, LineEnding: &preserve},
		{OperationID: "delete-one", Sequence: 5, Kind: "delete-file", Path: "delete.txt", BeforeHash: &deleteHash, Assertions: []Assertion{{Phase: "before", Kind: "exact-occurrences", Text: "remove", Count: 1}}},
	})
	record, _ := NewChangeSet(body)
	plan, err := Check(record, validOptions(root, body))
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Entries) != 2 || string(plan.Entries[0].After.Bytes) != "head center tail\n" || plan.Entries[1].After.Exists {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	if _, err := os.Stat(filepath.Join(root, "chain.txt")); !os.IsNotExist(err) {
		t.Fatalf("checker created file: %v", err)
	}
}

func TestDecodeChangeSetRejectsUnknownFieldAndTamperedHash(t *testing.T) {
	body := validBody([]Operation{})
	record, _ := NewChangeSet(body)
	data, _ := evidence.EncodeCanonical(record)
	unknown := strings.Replace(string(data), `"schemaVersion"`, `"unknown":true,"schemaVersion"`, 1)
	if _, err := DecodeChangeSet([]byte(unknown)); err == nil {
		t.Fatal("accepted unknown field")
	}
	record.ChangeSetHash = evidence.Digest("", []byte("tampered"))
	data, _ = evidence.EncodeCanonical(record)
	if _, err := DecodeChangeSet(data); !errors.Is(err, ErrInvalidChangeSet) {
		t.Fatalf("accepted tampered hash: %v", err)
	}
}

func TestCheckRejectsInvalidInputsWithoutWriting(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*ChangeSetBody)
	}{
		{"sequence-gap", func(body *ChangeSetBody) { body.Operations[0].Sequence = 2 }},
		{"unmanaged", func(body *ChangeSetBody) { body.Operations[0].Path = "other.txt" }},
		{"bad-before-hash", func(body *ChangeSetBody) {
			bad := evidence.Digest("", []byte("bad"))
			body.Operations[0].BeforeHash = &bad
		}},
		{"bad-marker", func(body *ChangeSetBody) {
			body.Operations[0].Marker = &Marker{MarkerID: "marker-one", Text: "beta", BeforeCount: 1, AfterCount: 1}
		}},
		{"bad-assertion", func(body *ChangeSetBody) { body.Operations[0].Assertions[0].Count = 2 }},
		{"nul-payload", func(body *ChangeSetBody) { value := "be\x00ta"; body.Operations[0].Replacement = &value }},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "file.txt")
			original := []byte("alpha\n")
			if err := os.WriteFile(path, original, 0644); err != nil {
				t.Fatal(err)
			}
			match, replacement, preserve, one := "alpha", "beta", "preserve-existing", 1
			before, after := evidence.Digest("", original), evidence.Digest("", []byte("beta\n"))
			body := validBody([]Operation{{OperationID: "operation-one", Sequence: 1, Kind: "replace-exact", Path: "file.txt", BeforeHash: &before, AfterHash: &after, Assertions: assertions("beta", 0, 1), Match: &match, Replacement: &replacement, ExpectedMatches: &one, LineEnding: &preserve}})
			item.mutate(&body)
			record, _ := NewChangeSet(body)
			options := validOptions(root, body)
			if item.name == "unmanaged" {
				options.ManagedPath = func(path string) bool { return path == "file.txt" }
			}
			if _, err := Check(record, options); err == nil {
				t.Fatal("accepted invalid change-set")
			}
			data, _ := os.ReadFile(path)
			if string(data) != string(original) {
				t.Fatal("failed checker wrote workspace")
			}
		})
	}
}

func assertions(text string, before, after int) []Assertion {
	return []Assertion{{Phase: "before", Kind: "exact-occurrences", Text: text, Count: before}, {Phase: "after", Kind: "exact-occurrences", Text: text, Count: after}}
}

func TestCheckRejectsFailureBeforeWriting(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "file.txt"), []byte("alpha\n"), 0644)
	match, replacement, lineEnding := "missing", "beta", "preserve-existing"
	expected := 1
	before := evidence.Digest("", []byte("alpha\n"))
	after := evidence.Digest("", []byte("beta\n"))
	body := validBody([]Operation{{OperationID: "operation-one", Sequence: 1, Kind: "replace-exact", Path: "file.txt", BeforeHash: &before, AfterHash: &after, Assertions: []Assertion{{Phase: "before", Kind: "exact-occurrences", Text: "alpha", Count: 1}}, Match: &match, Replacement: &replacement, ExpectedMatches: &expected, LineEnding: &lineEnding}})
	record, _ := NewChangeSet(body)
	if _, err := Check(record, validOptions(root, body)); err == nil {
		t.Fatal("expected failed exact match")
	}
	actual, _ := os.ReadFile(filepath.Join(root, "file.txt"))
	if string(actual) != "alpha\n" {
		t.Fatal("failed checker wrote the workspace")
	}
}

func validBody(operations []Operation) ChangeSetBody {
	hash := evidence.Digest("", []byte("manifest"))
	return ChangeSetBody{ChangeSetID: "change-one", CreatedAt: "2026-09-07T12:00:00.000Z", RunID: "run-one", TaskID: "task-one", Attempt: 1, ParentSnapshotHash: hash, BeforeManifestHash: hash, AfterManifestHash: hash, Mode: "managed-change-set", Operations: operations}
}

func validOptions(root string, body ChangeSetBody) CheckOptions {
	return CheckOptions{Root: root, ParentSnapshotHash: body.ParentSnapshotHash, BeforeManifestHash: body.BeforeManifestHash, AfterManifestHash: body.AfterManifestHash, ManagedPath: func(string) bool { return true }}
}
