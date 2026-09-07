package applier

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
	"github.com/larsonzh/prfrail/internal/taskdef"
)

const journalDomain = "proofrail:change-set-journal:1\n"

var (
	ErrInvalidJournal    = errors.New("invalid change-set journal")
	ErrRecoveryUncertain = errors.New("change-set recovery uncertain")
	ErrCrashInjected     = errors.New("simulated process crash")
)

type journalEntry struct {
	Path       string  `json:"path"`
	BeforeHash *string `json:"beforeHash"`
	AfterHash  *string `json:"afterHash"`
	BeforeMode uint32  `json:"beforeMode"`
	AfterMode  uint32  `json:"afterMode"`
}

type journalBody struct {
	JournalID     string             `json:"journalId"`
	ChangeSetHash string             `json:"changeSetHash"`
	ResourceID    string             `json:"resourceId"`
	Token         guard.FencingToken `json:"token"`
	State         string             `json:"state"`
	Completed     int                `json:"completed"`
	CreatedDirs   []string           `json:"createdDirs"`
	Entries       []journalEntry     `json:"entries"`
	UpdatedAt     string             `json:"updatedAt"`
}

type journalRecord struct {
	SchemaVersion string      `json:"schemaVersion"`
	Journal       journalBody `json:"journal"`
	JournalHash   string      `json:"journalHash"`
}

type journalStore struct{ root string }

func newJournalStore(root string) (*journalStore, error) {
	if root == "" {
		return nil, errors.New("journal root is required")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, fmt.Errorf("create journal root: %w", err)
	}
	return &journalStore{root: root}, nil
}

func (store *journalStore) create(changeSetHash, resourceID string, token guard.FencingToken, plan taskdef.Plan, workspaceRoot string, beforeWrite func() error) (journalRecord, error) {
	id := stringsTrimHash(changeSetHash)
	body := journalBody{JournalID: "journal-" + id[:16], ChangeSetHash: changeSetHash, ResourceID: resourceID, Token: token, State: "prepared", CreatedDirs: missingParentDirs(workspaceRoot, plan), Entries: make([]journalEntry, 0, len(plan.Entries)), UpdatedAt: formatTime(time.Now().UTC())}
	for _, entry := range plan.Entries {
		item := journalEntry{Path: entry.Path, BeforeMode: uint32(entry.Before.Mode.Perm()), AfterMode: uint32(entry.After.Mode.Perm())}
		if entry.Before.Exists {
			hash := evidence.Digest("", entry.Before.Bytes)
			item.BeforeHash = &hash
			if err := beforeWrite(); err != nil {
				return journalRecord{}, err
			}
			if err := store.putBlob(hash, entry.Before.Bytes); err != nil {
				return journalRecord{}, err
			}
		}
		if entry.After.Exists {
			hash := evidence.Digest("", entry.After.Bytes)
			item.AfterHash = &hash
			if err := beforeWrite(); err != nil {
				return journalRecord{}, err
			}
			if err := store.putBlob(hash, entry.After.Bytes); err != nil {
				return journalRecord{}, err
			}
		}
		body.Entries = append(body.Entries, item)
	}
	record, err := makeJournal(body)
	if err != nil {
		return journalRecord{}, err
	}
	data, err := evidence.EncodeCanonical(record)
	if err != nil {
		return journalRecord{}, err
	}
	if err := beforeWrite(); err != nil {
		return journalRecord{}, err
	}
	if err := durableCreate(store.path(body.JournalID), append(data, '\n'), 0600); err != nil {
		if os.IsExist(err) {
			return journalRecord{}, fmt.Errorf("%w: journal already exists", ErrInvalidJournal)
		}
		return journalRecord{}, err
	}
	return record, nil
}

func (store *journalStore) update(record *journalRecord, state string, completed int) error {
	record.Journal.State = state
	record.Journal.Completed = completed
	record.Journal.UpdatedAt = formatTime(time.Now().UTC())
	updated, err := makeJournal(record.Journal)
	if err != nil {
		return err
	}
	*record = updated
	return store.write(updated)
}

func makeJournal(body journalBody) (journalRecord, error) {
	canonical, err := evidence.EncodeCanonical(body)
	if err != nil {
		return journalRecord{}, err
	}
	return journalRecord{SchemaVersion: "1.0.0", Journal: body, JournalHash: evidence.Digest(journalDomain, canonical)}, nil
}

func (store *journalStore) write(record journalRecord) error {
	data, err := evidence.EncodeCanonical(record)
	if err != nil {
		return err
	}
	return durableReplace(store.path(record.Journal.JournalID), append(data, '\n'), 0600)
}

func (store *journalStore) read(path string) (journalRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return journalRecord{}, err
	}
	var record journalRecord
	if err := evidence.DecodeStrictJSON(data, &record); err != nil {
		return journalRecord{}, fmt.Errorf("%w: %v", ErrInvalidJournal, err)
	}
	canonical, err := evidence.EncodeCanonical(record.Journal)
	if err != nil || record.SchemaVersion != "1.0.0" || record.JournalHash != evidence.Digest(journalDomain, canonical) {
		return journalRecord{}, fmt.Errorf("%w: hash or version mismatch", ErrInvalidJournal)
	}
	if !evidence.ValidID(record.Journal.JournalID) || !evidence.ValidHash(record.Journal.ChangeSetHash) || !evidence.ValidID(record.Journal.ResourceID) || record.Journal.Token.ClaimID == "" || record.Journal.Token.Generation == 0 || record.Journal.Completed < 0 || record.Journal.Completed > len(record.Journal.Entries) {
		return journalRecord{}, fmt.Errorf("%w: invalid fields", ErrInvalidJournal)
	}
	if filepath.Base(path) != record.Journal.JournalID+".journal.json" || !validJournalState(record.Journal.State) {
		return journalRecord{}, fmt.Errorf("%w: identity or state mismatch", ErrInvalidJournal)
	}
	seenPaths := make(map[string]struct{}, len(record.Journal.Entries))
	for _, entry := range record.Journal.Entries {
		if snapshot.ValidatePath(entry.Path) != nil || (entry.BeforeHash == nil && entry.AfterHash == nil) {
			return journalRecord{}, fmt.Errorf("%w: invalid entry", ErrInvalidJournal)
		}
		if _, exists := seenPaths[entry.Path]; exists {
			return journalRecord{}, fmt.Errorf("%w: duplicate path", ErrInvalidJournal)
		}
		seenPaths[entry.Path] = struct{}{}
		if (entry.BeforeHash != nil && !evidence.ValidHash(*entry.BeforeHash)) || (entry.AfterHash != nil && !evidence.ValidHash(*entry.AfterHash)) {
			return journalRecord{}, fmt.Errorf("%w: invalid entry hash", ErrInvalidJournal)
		}
	}
	seenDirs := make(map[string]struct{}, len(record.Journal.CreatedDirs))
	for _, dir := range record.Journal.CreatedDirs {
		if snapshot.ValidatePath(dir) != nil {
			return journalRecord{}, fmt.Errorf("%w: invalid created directory", ErrInvalidJournal)
		}
		if _, exists := seenDirs[dir]; exists {
			return journalRecord{}, fmt.Errorf("%w: duplicate created directory", ErrInvalidJournal)
		}
		seenDirs[dir] = struct{}{}
	}
	return record, nil
}

func validJournalState(state string) bool {
	switch state {
	case "prepared", "applying", "rolling-back", "uncertain", "applied", "rolled-back":
		return true
	}
	return false
}

func (store *journalStore) list() ([]string, error) {
	items, err := os.ReadDir(store.root)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, item := range items {
		if !item.IsDir() && stringsHasSuffix(item.Name(), ".journal.json") {
			paths = append(paths, filepath.Join(store.root, item.Name()))
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (store *journalStore) putBlob(hash string, data []byte) error {
	path := filepath.Join(store.root, stringsTrimHash(hash)+".blob")
	if existing, err := os.ReadFile(path); err == nil {
		if evidence.Digest("", existing) == hash {
			return nil
		}
		return fmt.Errorf("%w: corrupt journal blob", ErrInvalidJournal)
	}
	return durableCreate(path, data, 0600)
}

func (store *journalStore) getBlob(hash string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(store.root, stringsTrimHash(hash)+".blob"))
	if err != nil || evidence.Digest("", data) != hash {
		return nil, fmt.Errorf("%w: missing or corrupt journal blob %s", ErrInvalidJournal, hash)
	}
	return data, nil
}

func (store *journalStore) path(id string) string {
	return filepath.Join(store.root, id+".journal.json")
}
func stringsTrimHash(hash string) string {
	if len(hash) == 71 && hash[:7] == "sha256:" {
		return hash[7:]
	}
	return hash
}
func stringsHasSuffix(value, suffix string) bool {
	return len(value) >= len(suffix) && value[len(value)-len(suffix):] == suffix
}
func formatTime(value time.Time) string { return value.Format("2006-01-02T15:04:05.000Z") }

func missingParentDirs(root string, plan taskdef.Plan) []string {
	set := make(map[string]struct{})
	for _, entry := range plan.Entries {
		parent := filepath.Dir(filepath.FromSlash(entry.Path))
		for parent != "." {
			if _, err := os.Stat(filepath.Join(root, parent)); err == nil {
				break
			}
			set[filepath.ToSlash(parent)] = struct{}{}
			parent = filepath.Dir(parent)
		}
	}
	dirs := make([]string, 0, len(set))
	for dir := range set {
		dirs = append(dirs, dir)
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) < len(dirs[j]) })
	return dirs
}
