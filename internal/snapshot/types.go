package snapshot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	SchemaVersion          = "1.0.0"
	snapshotManifestDomain = "proofrail:snapshot-manifest:1\n"
)

var (
	ErrInvalidManifest  = errors.New("invalid snapshot manifest")
	ErrInvalidPath      = errors.New("invalid snapshot path")
	ErrReservedPath     = errors.New("windows reserved path")
	ErrCaseCollision    = errors.New("path case collision")
	ErrDirectoryClosure = errors.New("directory closure violation")
	ErrHardlinkMismatch = errors.New("hardlink group mismatch")
	ErrUnsortedEntries  = errors.New("unsorted entries")
	ErrConcurrentChange = errors.New("concurrent modification detected")
	ErrQuotaExceeded    = errors.New("storage quota exceeded")
	ErrObjectNotFound   = errors.New("snapshot object not found")
	ErrObjectCorrupt    = errors.New("snapshot object corrupted")
	ErrTargetNotEmpty   = errors.New("restore target directory is not empty")
	ErrUnsupportedLink  = errors.New("unsupported link or reparse point")
)

var (
	envVarRegex = regexp.MustCompile(`^[^=\x00-\x1f\x7f]+$`)
)

func isRelativePathValid(path string) bool {
	if len(path) == 0 || len(path) > 1024 {
		return false
	}
	if strings.HasPrefix(path, "/") || strings.Contains(path, "\\") || strings.Contains(path, "\x00") {
		return false
	}
	if len(path) >= 2 && ((path[0] >= 'a' && path[0] <= 'z') || (path[0] >= 'A' && path[0] <= 'Z')) && path[1] == ':' {
		return false
	}
	colonIdx := strings.Index(path, ":")
	if colonIdx > 0 && !strings.Contains(path[:colonIdx], "/") {
		first := path[0]
		if (first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') {
			isScheme := true
			for i := 1; i < colonIdx; i++ {
				ch := path[i]
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '+' || ch == '.' || ch == '-') {
					isScheme = false
					break
				}
			}
			if isScheme {
				return false
			}
		}
	}
	segments := strings.Split(path, "/")
	for _, seg := range segments {
		if seg == "" || seg == "." || seg == ".." {
			return false
		}
	}
	return true
}

func isManagedGlobValid(pattern string) bool {
	if len(pattern) == 0 || len(pattern) > 1024 {
		return false
	}
	if strings.HasPrefix(pattern, "!") || strings.HasPrefix(pattern, "/") || strings.Contains(pattern, "\\") || strings.Contains(pattern, "\x00") {
		return false
	}
	if len(pattern) >= 2 && ((pattern[0] >= 'a' && pattern[0] <= 'z') || (pattern[0] >= 'A' && pattern[0] <= 'Z')) && pattern[1] == ':' {
		return false
	}
	colonIdx := strings.Index(pattern, ":")
	if colonIdx > 0 && !strings.Contains(pattern[:colonIdx], "/") {
		first := pattern[0]
		if (first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') {
			isScheme := true
			for i := 1; i < colonIdx; i++ {
				ch := pattern[i]
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '+' || ch == '.' || ch == '-') {
					isScheme = false
					break
				}
			}
			if isScheme {
				return false
			}
		}
	}
	segments := strings.Split(pattern, "/")
	for _, seg := range segments {
		if seg == "." || seg == ".." {
			return false
		}
	}
	return true
}

func isLinkTargetValid(target string) bool {
	if len(target) == 0 || len(target) > 1024 {
		return false
	}
	if strings.HasPrefix(target, "/") || strings.Contains(target, "\\") || strings.Contains(target, "\x00") {
		return false
	}
	if len(target) >= 2 && ((target[0] >= 'a' && target[0] <= 'z') || (target[0] >= 'A' && target[0] <= 'Z')) && target[1] == ':' {
		return false
	}
	colonIdx := strings.Index(target, ":")
	if colonIdx > 0 && !strings.Contains(target[:colonIdx], "/") {
		first := target[0]
		if (first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') {
			isScheme := true
			for i := 1; i < colonIdx; i++ {
				ch := target[i]
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '+' || ch == '.' || ch == '-') {
					isScheme = false
					break
				}
			}
			if isScheme {
				return false
			}
		}
	}
	return true
}

var reservedWindowsNames = map[string]struct{}{
	"CON": {}, "PRN": {}, "AUX": {}, "NUL": {},
	"COM1": {}, "COM2": {}, "COM3": {}, "COM4": {}, "COM5": {},
	"COM6": {}, "COM7": {}, "COM8": {}, "COM9": {},
	"LPT1": {}, "LPT2": {}, "LPT3": {}, "LPT4": {}, "LPT5": {},
	"LPT6": {}, "LPT7": {}, "LPT8": {}, "LPT9": {},
}

type EntryType string

const (
	EntryRegularFile  EntryType = "regular-file"
	EntryDirectory    EntryType = "directory"
	EntrySymbolicLink EntryType = "symbolic-link"
)

type FileRole string

const (
	RoleRegular         FileRole = "regular"
	RolePackageManifest FileRole = "package-manifest"
	RoleLockfile        FileRole = "lockfile"
	RoleGenerated       FileRole = "generated"
)

type TaskBinding struct {
	TaskID  string `json:"taskId"`
	Attempt int    `json:"attempt"`
}

type Entry struct {
	Path          string
	Type          EntryType
	ContentHash   string
	SizeBytes     int64
	Executable    bool
	ReadOnly      bool
	Role          FileRole
	HardlinkGroup *string
	Target        string
	TargetHash    string
}

type regularFileJSON struct {
	Path          string    `json:"path"`
	Type          EntryType `json:"type"`
	ContentHash   string    `json:"contentHash"`
	SizeBytes     int64     `json:"sizeBytes"`
	Executable    bool      `json:"executable"`
	ReadOnly      bool      `json:"readOnly"`
	Role          FileRole  `json:"role"`
	HardlinkGroup *string   `json:"hardlinkGroup"`
}

type directoryJSON struct {
	Path     string    `json:"path"`
	Type     EntryType `json:"type"`
	ReadOnly bool      `json:"readOnly"`
}

type symbolicLinkJSON struct {
	Path       string    `json:"path"`
	Type       EntryType `json:"type"`
	Target     string    `json:"target"`
	TargetHash string    `json:"targetHash"`
}

func (e Entry) MarshalJSON() ([]byte, error) {
	switch e.Type {
	case EntryRegularFile:
		return json.Marshal(regularFileJSON{
			Path: e.Path, Type: e.Type, ContentHash: e.ContentHash,
			SizeBytes: e.SizeBytes, Executable: e.Executable, ReadOnly: e.ReadOnly,
			Role: e.Role, HardlinkGroup: e.HardlinkGroup,
		})
	case EntryDirectory:
		return json.Marshal(directoryJSON{
			Path: e.Path, Type: e.Type, ReadOnly: e.ReadOnly,
		})
	case EntrySymbolicLink:
		return json.Marshal(symbolicLinkJSON{
			Path: e.Path, Type: e.Type, Target: e.Target, TargetHash: e.TargetHash,
		})
	default:
		return nil, fmt.Errorf("%w: unknown entry type %q", ErrInvalidManifest, e.Type)
	}
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type EntryType `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	switch probe.Type {
	case EntryRegularFile:
		var rf regularFileJSON
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&rf); err != nil {
			return err
		}
		*e = Entry{
			Path: rf.Path, Type: rf.Type, ContentHash: rf.ContentHash,
			SizeBytes: rf.SizeBytes, Executable: rf.Executable, ReadOnly: rf.ReadOnly,
			Role: rf.Role, HardlinkGroup: rf.HardlinkGroup,
		}
	case EntryDirectory:
		var dj directoryJSON
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&dj); err != nil {
			return err
		}
		*e = Entry{
			Path: dj.Path, Type: dj.Type, ReadOnly: dj.ReadOnly,
		}
	case EntrySymbolicLink:
		var sl symbolicLinkJSON
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&sl); err != nil {
			return err
		}
		*e = Entry{
			Path: sl.Path, Type: sl.Type, Target: sl.Target, TargetHash: sl.TargetHash,
		}
	default:
		return fmt.Errorf("%w: unknown entry type %q", ErrInvalidManifest, probe.Type)
	}
	return nil
}

type Exclusion struct {
	Pattern string `json:"pattern"`
	Source  string `json:"source"`
	Reason  string `json:"reason"`
}

type Environment struct {
	OS                       string   `json:"os"`
	Arch                     string   `json:"arch"`
	CaseSensitive            bool     `json:"caseSensitive"`
	EnvironmentVariableNames []string `json:"environmentVariableNames"`
	ToolchainEvidence        []string `json:"toolchainEvidence"`
}

type ManifestBody struct {
	SnapshotID         string       `json:"snapshotId"`
	Kind               string       `json:"kind"`
	CapturedAt         string       `json:"capturedAt"`
	RunID              string       `json:"runId"`
	ParentSnapshotHash *string      `json:"parentSnapshotHash"`
	Task               *TaskBinding `json:"task"`
	Entries            []Entry      `json:"entries"`
	Exclusions         []Exclusion  `json:"exclusions"`
	Environment        Environment  `json:"environment"`
}

type SnapshotManifest struct {
	SchemaVersion string       `json:"schemaVersion"`
	Manifest      ManifestBody `json:"manifest"`
	ManifestHash  string       `json:"manifestHash"`
}

func NewSnapshotManifest(body ManifestBody) (SnapshotManifest, error) {
	if err := ValidateManifestBody(body); err != nil {
		return SnapshotManifest{}, err
	}
	canonical, err := evidence.EncodeCanonical(body)
	if err != nil {
		return SnapshotManifest{}, err
	}
	return SnapshotManifest{
		SchemaVersion: SchemaVersion,
		Manifest:      body,
		ManifestHash:  evidence.Digest(snapshotManifestDomain, canonical),
	}, nil
}

func DecodeSnapshotManifest(input []byte) (SnapshotManifest, error) {
	var manifest SnapshotManifest
	if err := evidence.DecodeStrictJSON(input, &manifest); err != nil {
		return SnapshotManifest{}, err
	}
	if err := VerifySnapshotManifest(manifest); err != nil {
		return SnapshotManifest{}, err
	}
	return manifest, nil
}

func VerifySnapshotManifest(manifest SnapshotManifest) error {
	if manifest.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: %q", evidence.ErrUnsupportedVersion, manifest.SchemaVersion)
	}
	if err := ValidateManifestBody(manifest.Manifest); err != nil {
		return err
	}
	canonical, err := evidence.EncodeCanonical(manifest.Manifest)
	if err != nil {
		return err
	}
	expectedHash := evidence.Digest(snapshotManifestDomain, canonical)
	if manifest.ManifestHash != expectedHash {
		return fmt.Errorf("%w: snapshot manifest hash mismatch", evidence.ErrInvalidRecord)
	}
	return nil
}

func ValidateManifestBody(body ManifestBody) error {
	if !evidence.ValidID(body.SnapshotID) || !evidence.ValidID(body.RunID) {
		return fmt.Errorf("%w: invalid snapshot identity", ErrInvalidManifest)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", body.CapturedAt); err != nil {
		return fmt.Errorf("%w: invalid capturedAt timestamp", ErrInvalidManifest)
	}
	switch body.Kind {
	case "baseline":
		if body.ParentSnapshotHash != nil || body.Task != nil {
			return fmt.Errorf("%w: baseline snapshot must not have parent or task", ErrInvalidManifest)
		}
	case "candidate":
		if body.ParentSnapshotHash == nil || !evidence.ValidHash(*body.ParentSnapshotHash) {
			return fmt.Errorf("%w: candidate snapshot must have valid parent hash", ErrInvalidManifest)
		}
		if body.Task == nil || !evidence.ValidID(body.Task.TaskID) || body.Task.Attempt < 1 {
			return fmt.Errorf("%w: candidate snapshot must have valid task binding", ErrInvalidManifest)
		}
	default:
		return fmt.Errorf("%w: invalid snapshot kind %q", ErrInvalidManifest, body.Kind)
	}
	if len(body.Exclusions) == 0 {
		return fmt.Errorf("%w: exclusions must not be empty", ErrInvalidManifest)
	}
	for i, ex := range body.Exclusions {
		if !isManagedGlobValid(ex.Pattern) {
			return fmt.Errorf("%w: exclusion %d has invalid pattern", ErrInvalidManifest, i)
		}
		if ex.Source != "default" && ex.Source != "user" && ex.Source != "secret-policy" {
			return fmt.Errorf("%w: exclusion %d has invalid source", ErrInvalidManifest, i)
		}
		if !evidence.ValidID(ex.Reason) {
			return fmt.Errorf("%w: exclusion %d has invalid reason", ErrInvalidManifest, i)
		}
	}
	if err := validateEnvironment(body.Environment); err != nil {
		return err
	}
	return ValidateEntries(body.Entries)
}

func validateEnvironment(env Environment) error {
	if !evidence.ValidID(env.OS) || !evidence.ValidID(env.Arch) {
		return fmt.Errorf("%w: invalid environment OS or Arch", ErrInvalidManifest)
	}
	seenVars := make(map[string]struct{}, len(env.EnvironmentVariableNames))
	for _, v := range env.EnvironmentVariableNames {
		if !envVarRegex.MatchString(v) || len(v) > 256 {
			return fmt.Errorf("%w: invalid environment variable name %q", ErrInvalidManifest, v)
		}
		if _, exists := seenVars[v]; exists {
			return fmt.Errorf("%w: duplicate environment variable name %q", ErrInvalidManifest, v)
		}
		seenVars[v] = struct{}{}
	}
	seenTools := make(map[string]struct{}, len(env.ToolchainEvidence))
	for _, th := range env.ToolchainEvidence {
		if !evidence.ValidHash(th) {
			return fmt.Errorf("%w: invalid toolchain evidence hash %q", ErrInvalidManifest, th)
		}
		if _, exists := seenTools[th]; exists {
			return fmt.Errorf("%w: duplicate toolchain evidence hash %q", ErrInvalidManifest, th)
		}
		seenTools[th] = struct{}{}
	}
	return nil
}

func ValidateEntries(entries []Entry) error {
	dirSet := make(map[string]struct{}, len(entries))
	caseSet := make(map[string]string, len(entries))
	hardlinkGroups := make(map[string]Entry, len(entries))

	for i := 0; i < len(entries); i++ {
		e := entries[i]
		if i > 0 && entries[i-1].Path >= e.Path {
			return fmt.Errorf("%w: entry %q is not strictly ascending after %q", ErrUnsortedEntries, e.Path, entries[i-1].Path)
		}
		if err := ValidatePath(e.Path); err != nil {
			return err
		}
		lower := strings.ToLower(e.Path)
		if prev, exists := caseSet[lower]; exists {
			return fmt.Errorf("%w: %q conflicts with %q", ErrCaseCollision, e.Path, prev)
		}
		caseSet[lower] = e.Path

		switch e.Type {
		case EntryDirectory:
			dirSet[e.Path] = struct{}{}
		case EntryRegularFile:
			if !evidence.ValidHash(e.ContentHash) {
				return fmt.Errorf("%w: invalid contentHash in regular file %q", ErrInvalidManifest, e.Path)
			}
			if e.SizeBytes < 0 {
				return fmt.Errorf("%w: negative sizeBytes in %q", ErrInvalidManifest, e.Path)
			}
			if e.Role != RoleRegular && e.Role != RolePackageManifest && e.Role != RoleLockfile && e.Role != RoleGenerated {
				return fmt.Errorf("%w: invalid role in %q", ErrInvalidManifest, e.Path)
			}
			if e.HardlinkGroup != nil {
				if !evidence.ValidID(*e.HardlinkGroup) {
					return fmt.Errorf("%w: invalid hardlinkGroup ID in %q", ErrInvalidManifest, e.Path)
				}
				if primary, exists := hardlinkGroups[*e.HardlinkGroup]; exists {
					if primary.ContentHash != e.ContentHash || primary.SizeBytes != e.SizeBytes {
						return fmt.Errorf("%w: hardlinkGroup %q mismatch between %q and %q", ErrHardlinkMismatch, *e.HardlinkGroup, primary.Path, e.Path)
					}
				} else {
					hardlinkGroups[*e.HardlinkGroup] = e
				}
			}
		case EntrySymbolicLink:
			if !isLinkTargetValid(e.Target) {
				return fmt.Errorf("%w: invalid symlink target in %q", ErrInvalidManifest, e.Path)
			}
			if err := validateLinkTargetEscape(e.Path, e.Target); err != nil {
				return err
			}
			expectedTargetHash := evidence.Digest("", []byte(e.Target))
			if e.TargetHash != expectedTargetHash {
				return fmt.Errorf("%w: symlink targetHash mismatch in %q", ErrInvalidManifest, e.Path)
			}
		default:
			return fmt.Errorf("%w: unknown entry type %q in %q", ErrInvalidManifest, e.Type, e.Path)
		}
	}

	// Verify directory closure: every entry's parent directories must exist in dirSet
	for _, e := range entries {
		segments := strings.Split(e.Path, "/")
		for j := 1; j < len(segments); j++ {
			parent := strings.Join(segments[:j], "/")
			if _, ok := dirSet[parent]; !ok {
				return fmt.Errorf("%w: parent directory %q missing for %q", ErrDirectoryClosure, parent, e.Path)
			}
		}
	}
	return nil
}

func ValidatePath(path string) error {
	if !isRelativePathValid(path) {
		return fmt.Errorf("%w: %q does not match relative path rules", ErrInvalidPath, path)
	}
	segments := strings.Split(path, "/")
	for _, seg := range segments {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("%w: invalid path segment %q in %q", ErrInvalidPath, seg, path)
		}
		// Windows reserved name check: check base name without extension
		stem := seg
		if dotIdx := strings.Index(seg, "."); dotIdx != -1 {
			stem = seg[:dotIdx]
		}
		upperStem := strings.ToUpper(stem)
		if _, reserved := reservedWindowsNames[upperStem]; reserved {
			return fmt.Errorf("%w: segment %q is reserved in Windows", ErrReservedPath, seg)
		}
		// Trailing dot or space causes aliasing on Windows
		if strings.HasSuffix(seg, ".") || strings.HasSuffix(seg, " ") {
			return fmt.Errorf("%w: segment %q ends with dot or space", ErrInvalidPath, seg)
		}
	}
	return nil
}

func validateLinkTargetEscape(sourcePath, target string) error {
	// Normalizes target and verifies that it does not escape above root
	sourceDir := ""
	if idx := strings.LastIndex(sourcePath, "/"); idx != -1 {
		sourceDir = sourcePath[:idx]
	}
	parts := []string{}
	if sourceDir != "" {
		parts = strings.Split(sourceDir, "/")
	}
	for _, seg := range strings.Split(target, "/") {
		if seg == "" || seg == "." {
			continue
		}
		if seg == ".." {
			if len(parts) == 0 {
				return fmt.Errorf("%w: symlink target %q escapes root from %q", ErrInvalidPath, target, sourcePath)
			}
			parts = parts[:len(parts)-1]
		} else {
			parts = append(parts, seg)
		}
	}
	return nil
}
