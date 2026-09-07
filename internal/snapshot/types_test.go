package snapshot

import (
	"errors"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func validManifestBody() ManifestBody {
	return ManifestBody{
		SnapshotID:         "snapshot-one",
		Kind:               "baseline",
		CapturedAt:         "2026-09-07T12:00:00.000Z",
		RunID:              "run-one",
		ParentSnapshotHash: nil,
		Task:               nil,
		Entries: []Entry{
			{Path: "src", Type: EntryDirectory, ReadOnly: false},
			{
				Path: "src/main.go", Type: EntryRegularFile,
				ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
				SizeBytes:   128, Executable: false, ReadOnly: false, Role: RoleRegular, HardlinkGroup: nil,
			},
		},
		Exclusions: []Exclusion{
			{Pattern: ".git/**", Source: "default", Reason: "version-control-metadata"},
		},
		Environment: Environment{
			OS:                       "windows",
			Arch:                     "amd64",
			CaseSensitive:            false,
			EnvironmentVariableNames: []string{"PATH"},
			ToolchainEvidence:        []string{},
		},
	}
}

func TestNewSnapshotManifestValid(t *testing.T) {
	body := validManifestBody()
	manifest, err := NewSnapshotManifest(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manifest.SchemaVersion != SchemaVersion {
		t.Errorf("expected version %s, got %s", SchemaVersion, manifest.SchemaVersion)
	}
	if !evidence.ValidHash(manifest.ManifestHash) {
		t.Errorf("invalid manifest hash: %s", manifest.ManifestHash)
	}

	// Verify decoding roundtrip
	encoded, err := evidence.EncodeCanonical(manifest)
	if err != nil {
		t.Fatalf("encode canonical: %v", err)
	}
	decoded, err := DecodeSnapshotManifest(encoded)
	if err != nil {
		t.Fatalf("decode snapshot manifest: %v", err)
	}
	if decoded.ManifestHash != manifest.ManifestHash {
		t.Errorf("manifest hash mismatch: %s vs %s", decoded.ManifestHash, manifest.ManifestHash)
	}
}

func TestSnapshotManifestCandidate(t *testing.T) {
	body := validManifestBody()
	body.Kind = "candidate"
	parentHash := "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	body.ParentSnapshotHash = &parentHash
	body.Task = &TaskBinding{TaskID: "task-001", Attempt: 1}

	manifest, err := NewSnapshotManifest(body)
	if err != nil {
		t.Fatalf("candidate manifest: %v", err)
	}
	if err := VerifySnapshotManifest(manifest); err != nil {
		t.Fatalf("verify candidate: %v", err)
	}

	// Candidate without task fails
	bodyBad := body
	bodyBad.Task = nil
	if _, err := NewSnapshotManifest(bodyBad); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected ErrInvalidManifest, got %v", err)
	}

	// Baseline with parent fails
	bodyBaselineBad := validManifestBody()
	bodyBaselineBad.ParentSnapshotHash = &parentHash
	if _, err := NewSnapshotManifest(bodyBaselineBad); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected ErrInvalidManifest for baseline with parent, got %v", err)
	}
}

func TestSnapshotManifestValidationRejections(t *testing.T) {
	// Unsorted entries
	body := validManifestBody()
	body.Entries = []Entry{
		{Path: "z-dir", Type: EntryDirectory},
		{Path: "a-dir", Type: EntryDirectory},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrUnsortedEntries) {
		t.Fatalf("expected ErrUnsortedEntries, got %v", err)
	}

	// Directory closure violation (parent "sub" missing)
	body = validManifestBody()
	body.Entries = []Entry{
		{
			Path: "sub/file.txt", Type: EntryRegularFile,
			ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
			SizeBytes:   10, Executable: false, ReadOnly: false, Role: RoleRegular,
		},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrDirectoryClosure) {
		t.Fatalf("expected ErrDirectoryClosure, got %v", err)
	}

	// Case collision
	body = validManifestBody()
	body.Entries = []Entry{
		{Path: "FILE.TXT", Type: EntryRegularFile, ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", SizeBytes: 10, Role: RoleRegular},
		{Path: "file.txt", Type: EntryRegularFile, ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", SizeBytes: 10, Role: RoleRegular},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrCaseCollision) {
		t.Fatalf("expected ErrCaseCollision, got %v", err)
	}

	// Reserved Windows name
	body = validManifestBody()
	body.Entries = []Entry{
		{Path: "con.txt", Type: EntryRegularFile, ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", SizeBytes: 10, Role: RoleRegular},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrReservedPath) {
		t.Fatalf("expected ErrReservedPath, got %v", err)
	}

	// Trailing dot
	body = validManifestBody()
	body.Entries = []Entry{
		{Path: "foo.", Type: EntryRegularFile, ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", SizeBytes: 10, Role: RoleRegular},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("expected ErrInvalidPath, got %v", err)
	}

	// Hardlink group mismatch
	group := "hlg-001"
	body = validManifestBody()
	body.Entries = []Entry{
		{Path: "a.txt", Type: EntryRegularFile, ContentHash: "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", SizeBytes: 10, Role: RoleRegular, HardlinkGroup: &group},
		{Path: "b.txt", Type: EntryRegularFile, ContentHash: "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", SizeBytes: 20, Role: RoleRegular, HardlinkGroup: &group},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrHardlinkMismatch) {
		t.Fatalf("expected ErrHardlinkMismatch, got %v", err)
	}

	// Invalid schemaVersion
	manifest, _ := NewSnapshotManifest(validManifestBody())
	manifest.SchemaVersion = "2.0.0"
	if err := VerifySnapshotManifest(manifest); !errors.Is(err, evidence.ErrUnsupportedVersion) {
		t.Fatalf("expected ErrUnsupportedVersion, got %v", err)
	}

	// Invalid exclusion source
	bodyBadEx := validManifestBody()
	bodyBadEx.Exclusions[0].Source = "invalid-source"
	if _, err := NewSnapshotManifest(bodyBadEx); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected ErrInvalidManifest for bad exclusion source, got %v", err)
	}

	// Invalid environment
	bodyBadEnv := validManifestBody()
	bodyBadEnv.Environment.OS = "INVALID_OS_ID"
	if _, err := NewSnapshotManifest(bodyBadEnv); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("expected ErrInvalidManifest for bad env OS, got %v", err)
	}
}

func TestValidatePathLongPath(t *testing.T) {
	// Exactly 1024 chars is allowed
	seg := strings.Repeat("a", 100)
	longPath := seg + "/" + seg + "/" + seg + "/" + seg + "/" + seg + "/" + seg + "/" + seg + "/" + seg + "/" + seg + "/" + strings.Repeat("a", 115)
	if len(longPath) != 1024 {
		t.Fatalf("length is %d, want 1024", len(longPath))
	}
	if err := ValidatePath(longPath); err != nil {
		t.Fatalf("expected 1024-char path to pass, got: %v", err)
	}

	// > 1024 chars is rejected
	tooLong := longPath + "x"
	if err := ValidatePath(tooLong); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("expected ErrInvalidPath for >1024, got: %v", err)
	}
}

func TestSnapshotManifestRejectsUnicodeEquivalentPaths(t *testing.T) {
	body := validManifestBody()
	body.Entries = []Entry{
		{Path: "cafe\u0301.txt", Type: EntryRegularFile, SizeBytes: 1, ContentHash: evidence.Digest("", []byte("a")), Role: RoleRegular},
		{Path: "café.txt", Type: EntryRegularFile, SizeBytes: 1, ContentHash: evidence.Digest("", []byte("b")), Role: RoleRegular},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrUnicodeCollision) {
		t.Fatalf("expected ErrUnicodeCollision, got %v", err)
	}
}

func TestSymlinkEscapeRejection(t *testing.T) {
	target := "../../secret"
	body := validManifestBody()
	body.Entries = []Entry{
		{Path: "sub", Type: EntryDirectory},
		{Path: "sub/link", Type: EntrySymbolicLink, Target: target, TargetHash: evidence.Digest("", []byte(target))},
	}
	if _, err := NewSnapshotManifest(body); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("expected ErrInvalidPath for escaping symlink, got: %v", err)
	}
}
