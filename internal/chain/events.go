package chain

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type EventStore interface {
	Append(context.Context, evidence.StateEvent) error
	Load(context.Context) ([]evidence.StateEvent, error)
}

type FileEventStore struct {
	path string
}

func NewFileEventStore(path string) (*FileEventStore, error) {
	if path == "" {
		return nil, errors.New("event log path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create event log directory: %w", err)
	}
	return &FileEventStore{path: path}, nil
}

func (store *FileEventStore) Append(ctx context.Context, event evidence.StateEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(store.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open event log: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("append event: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync event log: %w", err)
	}
	return nil
}

func (store *FileEventStore) Load(ctx context.Context) ([]evidence.StateEvent, error) {
	file, err := os.Open(store.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open event log: %w", err)
	}
	defer file.Close()
	var events []evidence.StateEvent
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		event, err := evidence.DecodeStateEvent(scanner.Bytes())
		if err != nil {
			return nil, fmt.Errorf("decode event %d: %w", len(events)+1, err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read event log: %w", err)
	}
	if len(events) > 0 {
		if err := evidence.VerifyEventChain(events); err != nil {
			return nil, err
		}
	}
	return events, nil
}
