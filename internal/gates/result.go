package gates

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type ResultStore interface {
	Append(context.Context, ResultRecord) error
}

type FileResultStore struct {
	path string
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewFileResultStore(path string) (*FileResultStore, error) {
	if path == "" {
		return nil, errors.New("hook result log path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create hook result directory: %w", err)
	}
	store := &FileResultStore{path: path, seen: map[string]struct{}{}}
	data, err := os.ReadFile(path)
	if err == nil {
		for _, line := range bytes.Split(data, []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var record ResultRecord
			if err := evidence.DecodeStrictJSON(line, &record); err != nil || validateRecord(record) != nil {
				return nil, fmt.Errorf("invalid existing hook result log")
			}
			key := resultKey(record.Result)
			if _, duplicate := store.seen[key]; duplicate {
				return nil, fmt.Errorf("duplicate existing hook result %s", key)
			}
			store.seen[key] = struct{}{}
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return store, nil
}

func (store *FileResultStore) Append(ctx context.Context, record ResultRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateRecord(record); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	key := resultKey(record.Result)
	if _, duplicate := store.seen[key]; duplicate {
		return fmt.Errorf("duplicate hook result %s", key)
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(store.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(data, '\n')); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	store.seen[key] = struct{}{}
	return nil
}

func makeRecord(result Result) (ResultRecord, error) {
	canonical, err := evidence.EncodeCanonical(result)
	if err != nil {
		return ResultRecord{}, err
	}
	return ResultRecord{SchemaVersion: "1.0.0", Result: result, ResultHash: evidence.Digest("proofrail:hook-result:1\n", canonical)}, nil
}

func validateRecord(record ResultRecord) error {
	if record.SchemaVersion != "1.0.0" || !evidence.ValidID(record.Result.ResultID) || !evidence.ValidID(record.Result.RunID) || !evidence.ValidID(record.Result.TaskID) || !evidence.ValidID(record.Result.StepID) || !evidence.ValidID(record.Result.HookID) || !evidence.ValidHash(record.Result.HookDefinitionHash) || record.Result.Attempt < 1 || record.Result.ExecutionAttempt < 1 || len(record.Result.RunnerEvidence) == 0 {
		return errors.New("invalid hook result")
	}
	if (record.Result.EffectObservationHash == nil) != (record.Result.EffectRecoveryAction == nil) {
		return errors.New("invalid effect observation references")
	}
	if record.Result.EffectObservationHash != nil {
		if !evidence.ValidHash(*record.Result.EffectObservationHash) {
			return errors.New("invalid effect observation hash")
		}
		switch *record.Result.EffectRecoveryAction {
		case "none-required", "reconcile-only", "discard-workspace", "compensation-only":
		default:
			return errors.New("invalid effect recovery action")
		}
	}
	canonical, err := evidence.EncodeCanonical(record.Result)
	if err != nil || evidence.Digest("proofrail:hook-result:1\n", canonical) != record.ResultHash {
		return errors.New("invalid hook result hash")
	}
	if record.Result.Assessment == "passed" {
		if record.Result.ExecutionOutcome != "exited" || record.Result.ExitCode == nil || *record.Result.ExitCode != 0 || record.Result.FailureKind != nil || record.Result.PolicyDisposition != "pass" || record.Result.Stdout == nil || record.Result.Stderr == nil || len(record.Result.ErrorEvidence) != 0 || len(record.Result.TerminationEvidence) != 0 {
			return errors.New("invalid passed hook result")
		}
		return nil
	}
	if record.Result.Assessment != "failed" || record.Result.FailureKind == nil || record.Result.PolicyDisposition == "pass" || len(record.Result.ErrorEvidence) == 0 {
		return errors.New("invalid failed hook result")
	}
	if record.Result.ExecutionOutcome == "timed-out" || record.Result.ExecutionOutcome == "resource-limited" || record.Result.ExecutionOutcome == "termination-uncertain" {
		if record.Result.ExitCode != nil || record.Result.StartedAt == nil || len(record.Result.TerminationEvidence) == 0 {
			return errors.New("invalid terminated hook result")
		}
	}
	return nil
}

func resultKey(result Result) string {
	return fmt.Sprintf("%s/%s/%d/%s/%d", result.RunID, result.TaskID, result.Attempt, result.HookID, result.ExecutionAttempt)
}
