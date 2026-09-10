package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

type probeSpec struct {
	Command        string   `json:"command"`
	Args           []string `json:"args"`
	Dir            string   `json:"dir"`
	StdoutPath     string   `json:"stdoutPath"`
	StderrPath     string   `json:"stderrPath"`
	CancelMarker   string   `json:"cancelMarker"`
	CancelAfterMS  int      `json:"cancelAfterMs"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
}

type probeResult struct {
	Started               bool                       `json:"started"`
	ExitCode              int                        `json:"exitCode"`
	MarkerSeen            bool                       `json:"markerSeen"`
	CancellationRequested bool                       `json:"cancellationRequested"`
	TimedOut              bool                       `json:"timedOut"`
	Error                 string                     `json:"error,omitempty"`
	Termination           *guard.TerminationEvidence `json:"termination,omitempty"`
}

type markerWriter struct {
	mu       sync.Mutex
	w        io.Writer
	marker   []byte
	tail     []byte
	onMarker func()
	seen     bool
}

func (writer *markerWriter) Write(data []byte) (int, error) {
	written, err := writer.w.Write(data)
	writer.mu.Lock()
	defer writer.mu.Unlock()
	if writer.seen || len(writer.marker) == 0 {
		return written, err
	}
	window := append(append([]byte(nil), writer.tail...), data...)
	if bytes.Contains(window, writer.marker) {
		writer.seen = true
		writer.onMarker()
	}
	keep := len(writer.marker) - 1
	if keep > len(window) {
		keep = len(window)
	}
	writer.tail = append(writer.tail[:0], window[len(window)-keep:]...)
	return written, err
}

func (writer *markerWriter) Seen() bool {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	return writer.seen
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: agent-probe-supervisor <spec.json> <result.json>")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(specPath, resultPath string) error {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return err
	}
	var spec probeSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return err
	}
	if spec.Command == "" || spec.StdoutPath == "" || spec.StderrPath == "" || spec.CancelMarker == "" || spec.CancelAfterMS < 0 || spec.TimeoutSeconds <= 0 {
		return errors.New("invalid probe supervisor spec")
	}
	stdout, err := os.OpenFile(spec.StdoutPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer stdout.Close()
	stderr, err := os.OpenFile(spec.StderrPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer stderr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(spec.TimeoutSeconds)*time.Second)
	defer cancel()
	cancellationRequested := make(chan struct{})
	var requestOnce sync.Once
	writer := &markerWriter{
		w:      stdout,
		marker: []byte(spec.CancelMarker),
		onMarker: func() {
			time.AfterFunc(time.Duration(spec.CancelAfterMS)*time.Millisecond, func() {
				requestOnce.Do(func() {
					close(cancellationRequested)
					cancel()
				})
			})
		},
	}
	processResult, runErr := guard.RunManaged(ctx, guard.ProcessSpec{
		Command: spec.Command,
		Args:    spec.Args,
		Dir:     spec.Dir,
		Env:     os.Environ(),
		Stdout:  writer,
		Stderr:  stderr,
	}, 5*time.Second)
	result := probeResult{
		Started:     processResult.Started,
		ExitCode:    processResult.ExitCode,
		MarkerSeen:  writer.Seen(),
		TimedOut:    errors.Is(runErr, context.DeadlineExceeded),
		Termination: processResult.Termination,
	}
	select {
	case <-cancellationRequested:
		result.CancellationRequested = true
	default:
	}
	if runErr != nil {
		result.Error = runErr.Error()
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(resultPath, encoded, 0600); err != nil {
		return err
	}
	if result.Termination == nil || result.Termination.Outcome != "stopped" || !result.CancellationRequested {
		return errors.New("managed cancellation did not produce stopped termination evidence")
	}
	return nil
}
