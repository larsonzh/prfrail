package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// check is one asserted expectation of a face run. Kind is "assertion" for a
// mechanism expectation and "deviation" for a measurement that contradicts the
// ① design: a deviation is recorded, never hidden, and never counted as a hidden
// failure.
type check struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"`
}

// verdict is the machine-readable result of one harness run. Every field is a
// measurement or an explicitly recorded absence; nothing is inferred.
type verdict struct {
	Face      string `json:"face"`
	Case      string `json:"case"`
	Phase     string `json:"phase"`
	RunID     string `json:"runId"`
	RequestID string `json:"requestId"`
	RunRoot   string `json:"runRoot"`
	Workspace string `json:"workspace"`
	StartedAt string `json:"startedAt"`
	EndedAt   string `json:"endedAt"`

	VerdictClass  string `json:"verdictClass"`
	ChainState    string `json:"chainState"`
	TaskState     string `json:"taskState"`
	StepState     string `json:"stepState"`
	EventLogHash  string `json:"eventLogHash"`
	EventSequence int    `json:"eventSequence"`

	Load struct {
		Args    []string `json:"args"`
		Timeout string   `json:"timeout"`
		Grace   string   `json:"grace"`
	} `json:"load"`

	Mapping   *mappingView      `json:"mapping,omitempty"`
	Legs      []legView         `json:"legs"`
	Journal   []journalEntry    `json:"journal"`
	Checks    []check           `json:"checks"`
	Diff      *diffView         `json:"diff,omitempty"`
	Artifacts map[string]string `json:"artifacts"`

	Observations map[string]string `json:"observations"`
	Errors       []string          `json:"errors"`
	Notes        []string          `json:"notes"`
}

type mappingView struct {
	Status         string `json:"status"`
	TreeStatus     string `json:"treeStatus"`
	LogsComplete   bool   `json:"logsComplete"`
	UsageComplete  bool   `json:"usageComplete"`
	EventsComplete bool   `json:"eventsComplete"`
	Manifests      bool   `json:"manifestsComplete"`
	FactCount      int    `json:"factCount"`
	ErrorEvidence  int    `json:"errorEvidenceCount"`
	ElapsedMs      int64  `json:"elapsedMs"`
}

type legView struct {
	Leg         string `json:"leg"`
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
	Alternative string `json:"alternativeEvidence,omitempty"`
}

type diffView struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []string `json:"changed"`
}

func newVerdict() *verdict {
	return &verdict{Observations: map[string]string{}, Artifacts: map[string]string{}, Legs: []legView{}, Checks: []check{}, Errors: []string{}, Notes: []string{}}
}

func (v *verdict) addCheck(id string, pass bool, detail string) {
	v.Checks = append(v.Checks, check{ID: id, Kind: "assertion", Pass: pass, Detail: detail})
}

// addDeviation records a measurement that contradicts the ① design expectation.
func (v *verdict) addDeviation(id string, pass bool, detail string) {
	v.Checks = append(v.Checks, check{ID: id, Kind: "deviation", Pass: pass, Detail: detail})
}

func (v *verdict) addLeg(leg, status, errText, alternative string) {
	v.Legs = append(v.Legs, legView{Leg: leg, Status: status, Error: errText, Alternative: alternative})
}

func (v *verdict) addNote(note string) { v.Notes = append(v.Notes, note) }

func (v *verdict) addError(err error) {
	if err != nil {
		v.Errors = append(v.Errors, err.Error())
	}
}

// write publishes the verdict as JSON without a BOM.
func (v *verdict) write(path string) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "verdict written: %s (class=%s)\n", path, v.VerdictClass)
	return nil
}

func stamp(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

// TaskStateIsRepair reports whether the recorded task state names REPAIR_PENDING.
func (v *verdict) TaskStateIsRepair() bool {
	return strings.Contains(v.TaskState, "REPAIR_PENDING")
}
