package console

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
)

// P51 CONSOLE-CLI-MULTITASK-WIRING -- CLI wiring red baseline.
//
// P50 landed the multi-task serial assembly, but the `run` command still exposes
// only the single-value `--request-record` flag. This file pins the CLI surface
// the wiring slice must deliver:
//
//   - a repeatable `--request-records <path>` flag that reaches the assembly
//     instead of being rejected by the flag parser;
//   - `--request-record` and `--request-records` are mutually exclusive;
//   - the per-task record registry refuses duplicate (taskId, stepId) keys;
//   - a record whose runId disagrees with `--run-id` is refused;
//   - the single-value `--request-record` path stays compatible (green leg).
//
// Every test is deterministic: temporary directories only, no network, no
// external binary, no wall-clock assertion. The green leg must stay green on the
// pre-wiring HEAD; the other four are expected red until the wiring lands.

const (
	p51RunID     = "run-p51-cli-multitask"
	p51TaskOne   = "task-one"
	p51TaskTwo   = "task-two"
	p51StepOne   = "code-one"
	p51StepTwo   = "code-two"
	p51CreatedAt = "2026-11-03T00:00:00.000Z"
	p51DigestDo  = "proofrail:p51-cli-multitask:1\n"
)

// p51Digest derives a deterministic sha256 digest for fixture values.
func p51Digest(parts ...string) string {
	return digest(p51DigestDo, parts...)
}

// p51ChainTask is one task/step pair of the CLI fixture chain.
type p51ChainTask struct {
	taskID string
	stepID string
}

// p51TwoTaskChain is the two-task chain shape: each task carries one autonomous
// code step, which the CLI maps to the isolated-workspace mode.
func p51TwoTaskChain() []p51ChainTask {
	return []p51ChainTask{
		{taskID: p51TaskOne, stepID: p51StepOne},
		{taskID: p51TaskTwo, stepID: p51StepTwo},
	}
}

// p51SingleTaskChain is the single-task control shape.
func p51SingleTaskChain() []p51ChainTask {
	return []p51ChainTask{{taskID: p51TaskOne, stepID: p51StepOne}}
}

// p51WriteChainConfig validates and writes a chain configuration, returning its
// path. The fixture must validate: a green must never come from a bad chain.
func p51WriteChainConfig(t *testing.T, tasks ...p51ChainTask) string {
	t.Helper()
	cfg := DefaultChainConfig("p51-cli-multitask")
	cfg.Tasks = make([]TaskConfig, 0, len(tasks))
	for _, task := range tasks {
		cfg.Tasks = append(cfg.Tasks, TaskConfig{
			ID:     task.taskID,
			Review: "manual",
			Steps:  []StepConfig{{ID: task.stepID, Kind: "code"}},
		})
	}
	if err := ValidateChainConfig(cfg); err != nil {
		t.Fatalf("the P51 fixture chain must validate: %v", err)
	}
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), DefaultChainConfigName)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// p51WriteRequestRecord persists one valid AgentRunner request record and
// returns its path. The record is decode-complete: the CLI must reach the file
// read before any other gate.
func p51WriteRequestRecord(t *testing.T, dir, name, runID, taskID, stepID, requestID string) string {
	t.Helper()
	record, err := adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            requestID,
		RunID:                runID,
		TaskID:               taskID,
		StepID:               stepID,
		Attempt:              1,
		CreatedAt:            p51CreatedAt,
		AdapterID:            "p51-cli-stub",
		Mode:                 "create",
		WorkspaceHash:        p51Digest("workspace", taskID),
		ContextHash:          p51Digest("context", taskID),
		ParentSnapshotHash:   p51Digest("baseline"),
		AuthorizationHash:    p51Digest("authorization", taskID),
		BudgetHash:           p51Digest("budget", taskID),
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"workspace-write", "local-process"},
		EffectMappingVersion: assemblyEffectMappingVersion,
		EffectMappingHash:    assemblyEffectMappingHash,
		Evidence:             []string{p51Digest("request-evidence", taskID)},
	})
	if err != nil {
		t.Fatalf("the P51 request record fixture must be valid: %v", err)
	}
	content, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// p51SevenAgentRunnerArgs returns the seven non-request AgentRunner input flags.
// The paths are deliberately not created: the mutual-exclusion and registry
// checks must fire before any of them is read.
func p51SevenAgentRunnerArgs(t *testing.T) []string {
	t.Helper()
	dir := t.TempDir()
	return []string{
		"--capability-record", filepath.Join(dir, "capability.json"),
		"--enforcement-record", filepath.Join(dir, "enforcement.json"),
		"--availability-record", filepath.Join(dir, "availability.json"),
		"--authorization-ledger", filepath.Join(dir, "authorization.json"),
		"--cost-ledger", filepath.Join(dir, "cost.json"),
		"--agent-cli", filepath.Join(dir, "agent-cli"),
		"--agent-version", "p51-agent-0.1.0",
		"--probe-permit", filepath.Join(t.TempDir(), "probe-permit.marker"),
	}
}

// p51RunArgs assembles one `run` invocation: chain, run id and extra flags.
func p51RunArgs(chainPath string, extra ...string) []string {
	args := []string{"run", "--chain", chainPath, "--run-id", p51RunID}
	return append(args, extra...)
}

// TestRunRequestRecordsPathsAreRead pins the repeatable flag: `--request-records`
// must be parsed and its path read before assembly, never rejected as an unknown
// flag.
func TestRunRequestRecordsPathsAreRead(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51TwoTaskChain()...)
	missing := filepath.Join(t.TempDir(), "does-not-exist.json")

	args := p51RunArgs(chainPath, append(p51SevenAgentRunnerArgs(t), "--request-records", missing)...)
	code, stdout, stderr := runCLI(t, cli, args...)
	if code == 0 {
		t.Fatalf("run must fail for a missing request record, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "request record") {
		t.Fatalf("the failure must name the request record, stderr=%q", stderr)
	}
	if strings.Contains(stderr, "flag provided but not defined") {
		t.Fatalf("--request-records must be a defined flag, stderr=%q", stderr)
	}
}

// TestRunRefusesRequestRecordWithRequestRecords pins the mutual exclusion of the
// single-value and repeatable request-record flags.
func TestRunRefusesRequestRecordWithRequestRecords(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51SingleTaskChain()...)
	dir := t.TempDir()
	single := filepath.Join(dir, "a.json")
	repeat := filepath.Join(dir, "b.json")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-record", single, "--request-records", repeat)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must refuse both request-record forms, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "mutually exclusive") {
		t.Fatalf("the refusal must state the flags are mutually exclusive, stderr=%q", stderr)
	}
}

// TestRunRefusesDuplicateRequestRecordKey pins the registry key collision check:
// two records bound to the same (taskId, stepId) are a duplicate key even when
// their requestIds differ.
func TestRunRefusesDuplicateRequestRecordKey(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51TwoTaskChain()...)
	dir := t.TempDir()
	first := p51WriteRequestRecord(t, dir, "record-one.json", p51RunID, p51TaskOne, p51StepOne, "request-p51-one")
	second := p51WriteRequestRecord(t, dir, "record-two.json", p51RunID, p51TaskOne, p51StepOne, "request-p51-two")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-records", first, "--request-records", second)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must refuse duplicate request record keys, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "duplicate key") {
		t.Fatalf("the refusal must name the duplicate key, stderr=%q", stderr)
	}
}

// TestRunRefusesRequestRecordRunIDMismatch pins the cross-record run binding: a
// record whose runId disagrees with `--run-id` is refused.
func TestRunRefusesRequestRecordRunIDMismatch(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51TwoTaskChain()...)
	dir := t.TempDir()
	bound := p51WriteRequestRecord(t, dir, "record-bound.json", p51RunID, p51TaskOne, p51StepOne, "request-p51-bound")
	stray := p51WriteRequestRecord(t, dir, "record-stray.json", "run-p51-other", p51TaskTwo, p51StepTwo, "request-p51-stray")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-records", bound, "--request-records", stray)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must refuse a record bound to another run id, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "does not match run id") {
		t.Fatalf("the refusal must state the run id mismatch, stderr=%q", stderr)
	}
}

// TestRunRefusesRequestRecordsWithoutExecutableSteps pins the fail-closed refusal
// of the repeatable flag on a noop-only chain. Such a chain never enters the
// AgentRunner assembly, so the registry has no consumer; accepting the flag and
// silently running the noop path would hide an unusable declaration.
func TestRunRefusesRequestRecordsWithoutExecutableSteps(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	cfg := DefaultChainConfig("p51-cli-noop")
	if err := ValidateChainConfig(cfg); err != nil {
		t.Fatalf("the P51 noop fixture chain must validate: %v", err)
	}
	if executable := FindExecutableSteps(cfg); len(executable) != 0 {
		t.Fatalf("the P51 noop fixture chain must have no executable step, got %+v", executable)
	}
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	chainPath := filepath.Join(t.TempDir(), DefaultChainConfigName)
	if err := os.WriteFile(chainPath, content, 0o600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(t.TempDir(), "does-not-exist.json")

	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, "--request-records", missing)...)
	if code == 0 {
		t.Fatalf("run must refuse --request-records on a noop-only chain, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "at least one executable step") {
		t.Fatalf("the refusal must state the executable-step requirement, stderr=%q", stderr)
	}
}

// TestRunRefusesRequestRecordOutsideChainDefinition covers the stray-key leg:
// a record whose (taskId, stepId) pair is absent from the chain definition must
// be refused, because the registry may only bind planned task/step pairs.
func TestRunRefusesRequestRecordOutsideChainDefinition(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51TwoTaskChain()...)
	dir := t.TempDir()
	stray := p51WriteRequestRecord(t, dir, "record-stray.json", p51RunID, "task-other", "code-other", "request-p51-other")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-records", stray)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must refuse a record outside the chain definition, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "does not bind a task/step") {
		t.Fatalf("the refusal must state the definition binding, stderr=%q", stderr)
	}
}

// TestRunRefusesRequestRecordCoverageGap covers the coverage leg: a registry
// that binds only one of two planned task/step pairs leaves the other task
// without a request record and must be refused.
func TestRunRefusesRequestRecordCoverageGap(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51TwoTaskChain()...)
	dir := t.TempDir()
	only := p51WriteRequestRecord(t, dir, "record-one.json", p51RunID, p51TaskOne, p51StepOne, "request-p51-only")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-records", only)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must refuse an incomplete registry, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "request record is missing for task") {
		t.Fatalf("the refusal must name the uncovered task, stderr=%q", stderr)
	}
}

// TestRunSingleRequestRecordStaysCompatible is the green leg: the existing
// single-value flag must keep reaching the record read (the label carries the
// colon), so the wiring cannot regress the single-task path.
func TestRunSingleRequestRecordStaysCompatible(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := p51WriteChainConfig(t, p51SingleTaskChain()...)
	missing := filepath.Join(t.TempDir(), "does-not-exist.json")

	extra := append(p51SevenAgentRunnerArgs(t), "--request-record", missing)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code == 0 {
		t.Fatalf("run must fail for a missing request record, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "request record:") {
		t.Fatalf("the single request record must be read, stderr=%q", stderr)
	}
}
