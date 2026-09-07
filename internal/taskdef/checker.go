package taskdef

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

const changeSetDomain = "proofrail:managed-change-set:1\n"

var (
	ErrInvalidChangeSet = errors.New("invalid managed change-set")
	ErrAssertionFailed  = errors.New("change-set assertion failed")
	ErrTargetChanged    = errors.New("managed target changed")
)

type Assertion struct {
	Phase string `json:"phase"`
	Kind  string `json:"kind"`
	Text  string `json:"text"`
	Count int    `json:"count"`
}

type Marker struct {
	MarkerID    string `json:"markerId"`
	Text        string `json:"text"`
	BeforeCount int    `json:"beforeCount"`
	AfterCount  int    `json:"afterCount"`
}

type Operation struct {
	OperationID     string      `json:"operationId"`
	Sequence        int         `json:"sequence"`
	Kind            string      `json:"kind"`
	Path            string      `json:"path"`
	BeforeHash      *string     `json:"beforeHash"`
	AfterHash       *string     `json:"afterHash"`
	Assertions      []Assertion `json:"assertions"`
	Content         *string     `json:"content,omitempty"`
	LineEnding      *string     `json:"lineEnding,omitempty"`
	Match           *string     `json:"match,omitempty"`
	Replacement     *string     `json:"replacement,omitempty"`
	ExpectedMatches *int        `json:"expectedMatches,omitempty"`
	Marker          *Marker     `json:"marker,omitempty"`
	Anchor          *string     `json:"anchor,omitempty"`
}

type ChangeSetBody struct {
	ChangeSetID        string      `json:"changeSetId"`
	CreatedAt          string      `json:"createdAt"`
	RunID              string      `json:"runId"`
	TaskID             string      `json:"taskId"`
	Attempt            int         `json:"attempt"`
	ParentSnapshotHash string      `json:"parentSnapshotHash"`
	BeforeManifestHash string      `json:"beforeManifestHash"`
	AfterManifestHash  string      `json:"afterManifestHash"`
	Mode               string      `json:"mode"`
	Operations         []Operation `json:"operations"`
}

type ChangeSet struct {
	SchemaVersion string        `json:"schemaVersion"`
	ChangeSet     ChangeSetBody `json:"changeSet"`
	ChangeSetHash string        `json:"changeSetHash"`
}

type FileState struct {
	Exists bool
	Bytes  []byte
	Mode   os.FileMode
}

type PlanEntry struct {
	Path   string
	Before FileState
	After  FileState
}

type Plan struct {
	ChangeSetHash string
	Entries       []PlanEntry
}

type CheckOptions struct {
	Root               string
	ParentSnapshotHash string
	BeforeManifestHash string
	AfterManifestHash  string
	ManagedPath        func(string) bool
}

func DecodeChangeSet(input []byte) (ChangeSet, error) {
	var record ChangeSet
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return ChangeSet{}, fmt.Errorf("%w: %v", ErrInvalidChangeSet, err)
	}
	canonical, err := evidence.EncodeCanonical(record.ChangeSet)
	if err != nil {
		return ChangeSet{}, err
	}
	if record.SchemaVersion != "1.0.0" || record.ChangeSetHash != evidence.Digest(changeSetDomain, canonical) {
		return ChangeSet{}, fmt.Errorf("%w: schema version or hash mismatch", ErrInvalidChangeSet)
	}
	return record, nil
}

func NewChangeSet(body ChangeSetBody) (ChangeSet, error) {
	canonical, err := evidence.EncodeCanonical(body)
	if err != nil {
		return ChangeSet{}, err
	}
	return ChangeSet{SchemaVersion: "1.0.0", ChangeSet: body, ChangeSetHash: evidence.Digest(changeSetDomain, canonical)}, nil
}

func Check(record ChangeSet, options CheckOptions) (Plan, error) {
	body := record.ChangeSet
	canonical, err := evidence.EncodeCanonical(body)
	if err != nil || record.SchemaVersion != "1.0.0" || record.ChangeSetHash != evidence.Digest(changeSetDomain, canonical) {
		return Plan{}, fmt.Errorf("%w: hash mismatch", ErrInvalidChangeSet)
	}
	if options.Root == "" || body.Mode != "managed-change-set" || !evidence.ValidID(body.ChangeSetID) || !evidence.ValidID(body.RunID) || !evidence.ValidID(body.TaskID) || body.Attempt < 1 || len(body.Operations) < 1 || len(body.Operations) > 10000 {
		return Plan{}, fmt.Errorf("%w: invalid identity, mode, root, attempt, or operation count", ErrInvalidChangeSet)
	}
	if !evidence.ValidHash(body.ParentSnapshotHash) || !evidence.ValidHash(body.BeforeManifestHash) || !evidence.ValidHash(body.AfterManifestHash) || body.ParentSnapshotHash != options.ParentSnapshotHash || body.BeforeManifestHash != options.BeforeManifestHash || body.AfterManifestHash != options.AfterManifestHash {
		return Plan{}, fmt.Errorf("%w: manifest binding mismatch", ErrInvalidChangeSet)
	}

	states := make(map[string]FileState)
	initial := make(map[string]FileState)
	seenOperations := make(map[string]struct{})
	seenMarkers := make(map[string]struct{})
	for index, operation := range body.Operations {
		if operation.Sequence != index+1 || !evidence.ValidID(operation.OperationID) {
			return Plan{}, fmt.Errorf("%w: invalid operation sequence or ID", ErrInvalidChangeSet)
		}
		if _, exists := seenOperations[operation.OperationID]; exists {
			return Plan{}, fmt.Errorf("%w: duplicate operation ID", ErrInvalidChangeSet)
		}
		seenOperations[operation.OperationID] = struct{}{}
		if err := snapshot.ValidatePath(operation.Path); err != nil {
			return Plan{}, err
		}
		if options.ManagedPath == nil || !options.ManagedPath(operation.Path) {
			return Plan{}, fmt.Errorf("%w: unmanaged path %q", ErrInvalidChangeSet, operation.Path)
		}
		state, exists := states[operation.Path]
		if !exists {
			state, err = readFileState(options.Root, operation.Path)
			if err != nil {
				return Plan{}, err
			}
			states[operation.Path] = state
			initial[operation.Path] = cloneState(state)
		}
		next, err := simulate(operation, state, seenMarkers)
		if err != nil {
			return Plan{}, fmt.Errorf("operation %q: %w", operation.OperationID, err)
		}
		states[operation.Path] = next
	}
	paths := make([]string, 0, len(states))
	for path := range states {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	plan := Plan{ChangeSetHash: record.ChangeSetHash, Entries: make([]PlanEntry, 0, len(paths))}
	for _, path := range paths {
		plan.Entries = append(plan.Entries, PlanEntry{Path: path, Before: initial[path], After: states[path]})
	}
	return plan, nil
}

func simulate(operation Operation, current FileState, markers map[string]struct{}) (FileState, error) {
	if len(operation.Assertions) == 0 || len(operation.Assertions) > 64 {
		return FileState{}, fmt.Errorf("%w: assertions required", ErrInvalidChangeSet)
	}
	if err := verifyHash(operation.BeforeHash, current); err != nil {
		return FileState{}, err
	}
	if err := checkAssertions(operation.Assertions, "before", current.Bytes); err != nil {
		return FileState{}, err
	}
	if current.Exists && (!utf8.Valid(current.Bytes) || bytes.IndexByte(current.Bytes, 0) >= 0) && operation.Kind != "delete-file" {
		return FileState{}, fmt.Errorf("%w: target is not valid NUL-free UTF-8", ErrInvalidChangeSet)
	}
	next := cloneState(current)
	switch operation.Kind {
	case "create-file":
		if current.Exists || operation.BeforeHash != nil || operation.Content == nil || operation.LineEnding == nil || (*operation.LineEnding != "lf" && *operation.LineEnding != "crlf") {
			return FileState{}, fmt.Errorf("%w: invalid create", ErrInvalidChangeSet)
		}
		if err := validPayload(*operation.Content, true); err != nil {
			return FileState{}, err
		}
		content := []byte(*operation.Content)
		if *operation.LineEnding == "crlf" {
			content = []byte(strings.ReplaceAll(*operation.Content, "\n", "\r\n"))
		}
		next = FileState{Exists: true, Bytes: content, Mode: defaultFileMode}
	case "delete-file":
		if !current.Exists || operation.AfterHash != nil {
			return FileState{}, fmt.Errorf("%w: invalid delete", ErrInvalidChangeSet)
		}
		next = FileState{}
	case "replace-exact":
		if !current.Exists || operation.Match == nil || operation.Replacement == nil || operation.ExpectedMatches == nil || operation.LineEnding == nil || *operation.LineEnding != "preserve-existing" {
			return FileState{}, fmt.Errorf("%w: invalid replace", ErrInvalidChangeSet)
		}
		if err := validPayload(*operation.Match, false); err != nil {
			return FileState{}, err
		}
		if err := validPayload(*operation.Replacement, true); err != nil {
			return FileState{}, err
		}
		count := bytes.Count(current.Bytes, []byte(*operation.Match))
		if count != *operation.ExpectedMatches {
			return FileState{}, fmt.Errorf("%w: match count %d, want %d", ErrAssertionFailed, count, *operation.ExpectedMatches)
		}
		next.Bytes = bytes.ReplaceAll(current.Bytes, []byte(*operation.Match), []byte(*operation.Replacement))
	case "insert-before", "insert-after":
		if !current.Exists || operation.Anchor == nil || operation.Content == nil || operation.ExpectedMatches == nil || *operation.ExpectedMatches != 1 || operation.LineEnding == nil || *operation.LineEnding != "preserve-existing" || operation.Marker == nil {
			return FileState{}, fmt.Errorf("%w: invalid insert", ErrInvalidChangeSet)
		}
		if err := validPayload(*operation.Anchor, false); err != nil {
			return FileState{}, err
		}
		if err := validPayload(*operation.Content, false); err != nil {
			return FileState{}, err
		}
		if bytes.Count(current.Bytes, []byte(*operation.Anchor)) != 1 {
			return FileState{}, fmt.Errorf("%w: insert anchor count", ErrAssertionFailed)
		}
		replacement := *operation.Content + *operation.Anchor
		if operation.Kind == "insert-after" {
			replacement = *operation.Anchor + *operation.Content
		}
		next.Bytes = bytes.Replace(current.Bytes, []byte(*operation.Anchor), []byte(replacement), 1)
	default:
		return FileState{}, fmt.Errorf("%w: unknown operation kind %q", ErrInvalidChangeSet, operation.Kind)
	}
	if operation.Marker != nil {
		marker := operation.Marker
		if !evidence.ValidID(marker.MarkerID) || marker.BeforeCount != 0 || marker.AfterCount != 1 || marker.Text == "" {
			return FileState{}, fmt.Errorf("%w: invalid marker", ErrInvalidChangeSet)
		}
		if _, exists := markers[marker.MarkerID]; exists {
			return FileState{}, fmt.Errorf("%w: duplicate marker ID", ErrInvalidChangeSet)
		}
		markers[marker.MarkerID] = struct{}{}
		if bytes.Count(current.Bytes, []byte(marker.Text)) != 0 || bytes.Count(next.Bytes, []byte(marker.Text)) != 1 {
			return FileState{}, fmt.Errorf("%w: marker count", ErrAssertionFailed)
		}
		produced := ""
		if operation.Replacement != nil {
			produced = *operation.Replacement
		}
		if operation.Content != nil {
			produced = *operation.Content
		}
		if !strings.Contains(produced, marker.Text) {
			return FileState{}, fmt.Errorf("%w: marker not produced", ErrInvalidChangeSet)
		}
	}
	if err := verifyHash(operation.AfterHash, next); err != nil {
		return FileState{}, err
	}
	if err := checkAssertions(operation.Assertions, "after", next.Bytes); err != nil {
		return FileState{}, err
	}
	return next, nil
}

func readFileState(root, relative string) (FileState, error) {
	path := filepath.Join(root, filepath.FromSlash(relative))
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return FileState{}, nil
	}
	if err != nil {
		return FileState{}, err
	}
	if !info.Mode().IsRegular() {
		return FileState{}, fmt.Errorf("%w: target %q is not a regular file", ErrInvalidChangeSet, relative)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return FileState{}, err
	}
	return FileState{Exists: true, Bytes: data, Mode: info.Mode().Perm()}, nil
}

func verifyHash(expected *string, state FileState) error {
	if !state.Exists {
		if expected != nil {
			return fmt.Errorf("%w: expected existing content", ErrTargetChanged)
		}
		return nil
	}
	if expected == nil || !evidence.ValidHash(*expected) || *expected != evidence.Digest("", state.Bytes) {
		return fmt.Errorf("%w: content hash mismatch", ErrTargetChanged)
	}
	return nil
}

func checkAssertions(assertions []Assertion, phase string, data []byte) error {
	for _, assertion := range assertions {
		if assertion.Phase != "before" && assertion.Phase != "after" {
			return fmt.Errorf("%w: invalid assertion phase", ErrInvalidChangeSet)
		}
		if assertion.Kind != "exact-occurrences" && assertion.Kind != "absent" {
			return fmt.Errorf("%w: invalid assertion kind", ErrInvalidChangeSet)
		}
		if err := validPayload(assertion.Text, false); err != nil {
			return err
		}
		if assertion.Count < 0 || (assertion.Kind == "absent" && assertion.Count != 0) {
			return fmt.Errorf("%w: invalid assertion count", ErrInvalidChangeSet)
		}
		if assertion.Phase == phase && bytes.Count(data, []byte(assertion.Text)) != assertion.Count {
			return fmt.Errorf("%w: %s assertion", ErrAssertionFailed, phase)
		}
	}
	return nil
}

func validPayload(value string, allowEmpty bool) error {
	if (!allowEmpty && value == "") || len(value) > 1<<20 || strings.ContainsAny(value, "\r\x00") {
		return fmt.Errorf("%w: invalid text payload", ErrInvalidChangeSet)
	}
	return nil
}

func cloneState(state FileState) FileState {
	return FileState{Exists: state.Exists, Bytes: bytes.Clone(state.Bytes), Mode: state.Mode}
}
