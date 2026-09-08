package gates

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type fakeExecutor struct {
	executions   []Execution
	calls        int
	capabilities *Capabilities
}

func (fake *fakeExecutor) Capabilities() Capabilities {
	if fake.capabilities != nil {
		return *fake.capabilities
	}
	return allCapabilities()
}
func (fake *fakeExecutor) Execute(context.Context, Prepared, Hook) Execution {
	result := fake.executions[fake.calls]
	fake.calls++
	return result
}

type fakeScanner struct{ err error }

func (fake fakeScanner) Scan(context.Context, string, []byte) (string, error) {
	return "not-required", fake.err
}

type memoryResults struct {
	records []ResultRecord
	err     error
}

type memoryObjects struct{ values map[string][]byte }

func (memory *memoryObjects) PutObject(_ context.Context, data []byte) (string, error) {
	hash := evidence.Digest("", data)
	if memory.values == nil {
		memory.values = map[string][]byte{}
	}
	memory.values[hash] = append([]byte(nil), data...)
	return hash, nil
}

func (memory *memoryResults) Append(_ context.Context, record ResultRecord) error {
	if memory.err != nil {
		return memory.err
	}
	memory.records = append(memory.records, record)
	return nil
}

func successfulExecution() Execution {
	started := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	exitCode := 0
	return Execution{StartedAt: &started, FinishedAt: started.Add(time.Second), Outcome: "exited", ExitCode: &exitCode, Stdout: []byte("ok"), Stderr: []byte{}, RunnerEvidence: []string{evidence.Digest("", []byte("runner"))}}
}

func testRunner(executor Executor, scanner Scanner, store ResultStore) Runner {
	sequence := 0
	return Runner{Executor: executor, Scanner: scanner, Objects: &memoryObjects{}, Store: store, Clock: func() time.Time { return time.Date(2026, 9, 7, 12, 0, sequence, 0, time.UTC) }, IDs: func() string { sequence++; return "result-" + string(rune('a'+sequence-1)) }, Env: func() []string { return nil }}
}

func TestRunnerPersistsPassedResultAndArtifacts(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "report.txt"), []byte("report"), 0600); err != nil {
		t.Fatal(err)
	}
	request := validRequest(root)
	request.Hook.ArtifactGlobs = []string{"report.txt"}
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if len(store.records) != 1 || record.Result.Assessment != "passed" || record.Result.PolicyDisposition != "pass" || len(record.Result.Artifacts) != 1 || !evidence.ValidHash(record.ResultHash) {
		t.Fatalf("invalid passed record: %+v", record)
	}
}

func TestSecretScannerRejectsKnownCredentialMaterial(t *testing.T) {
	if _, err := DefaultSecretScanner().Scan(context.Background(), "stdout", []byte("github_pat_example")); !errors.Is(err, ErrSecretDetected) {
		t.Fatalf("got %v", err)
	}
}

func TestRunnerRetriesAndPersistsEveryAttempt(t *testing.T) {
	failed := successfulExecution()
	exitCode := 2
	failed.ExitCode = &exitCode
	request := validRequest(t.TempDir())
	request.Hook.OnFail = FailurePolicy{Action: "retry", MaxAttempts: 2}
	executor := &fakeExecutor{executions: []Execution{failed, successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil || executor.calls != 2 || len(store.records) != 2 || store.records[0].Result.PolicyDisposition != "retry" || record.Result.ExecutionAttempt != 2 {
		t.Fatalf("retry failed: record=%+v err=%v", record, err)
	}
}

func TestRunnerFailsOnScanAndMissingArtifact(t *testing.T) {
	request := validRequest(t.TempDir())
	store := &memoryResults{}
	_, err := testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{err: errors.New("scanner unavailable")}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) || *store.records[0].Result.FailureKind != "scan-failed" {
		t.Fatalf("scan failure: %v %+v", err, store.records)
	}
	request.Hook.ArtifactGlobs = []string{"missing.txt"}
	store = &memoryResults{}
	_, err = testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) || *store.records[0].Result.FailureKind != "artifact-missing" {
		t.Fatalf("artifact failure: %v %+v", err, store.records)
	}
}

func TestRunnerDoesNotReturnResultBeforePersistence(t *testing.T) {
	storeErr := errors.New("disk full")
	_, err := testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{}, &memoryResults{err: storeErr}).Run(context.Background(), validRequest(t.TempDir()))
	if !errors.Is(err, storeErr) {
		t.Fatalf("got %v", err)
	}
}

func TestRunnerDoesNotExecuteWhenPolicyIsUnavailable(t *testing.T) {
	limited := Capabilities{Process: true, Timeout: true, OutputLimit: true}
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}, capabilities: &limited}
	_, err := testRunner(executor, fakeScanner{}, &memoryResults{}).Run(context.Background(), validRequest(t.TempDir()))
	if !errors.Is(err, ErrPolicyUnavailable) || executor.calls != 0 {
		t.Fatalf("unavailable policy executed: calls=%d err=%v", executor.calls, err)
	}
}

func TestFileResultStoreRejectsDuplicateAndTamperedResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.jsonl")
	store, err := NewFileResultStore(path)
	if err != nil {
		t.Fatal(err)
	}
	record, err := makeRecord(testRunner(&fakeExecutor{}, fakeScanner{}, &memoryResults{}).assessForTest(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Append(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	if err := store.Append(context.Background(), record); err == nil {
		t.Fatal("duplicate result accepted")
	}
	record.ResultHash = evidence.Digest("", []byte("tampered"))
	if err := validateRecord(record); err == nil {
		t.Fatal("tampered result accepted")
	}
}

func effectScopeForTest(effectClass string) *EffectScope {
	return &EffectScope{
		EffectClass:       effectClass,
		ExternalSystems:   []string{},
		RecoveryGuarantee: "",
		PolicyHash:        "",
		AuthorizationHash: effectTestHash('e'),
		OperationKey:      "op-run-go-test",
		Evidence:          []string{effectTestHash('c')},
	}
}

func TestRunnerDeniesExternalWriteEffectWithoutExecuting(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Effect = effectScopeForTest("external-write")
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("external-write must be denied: %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("denied effect must never execute: calls=%d", executor.calls)
	}
	if len(store.records) != 1 {
		t.Fatalf("denial must persist one record: %+v", store.records)
	}
	persisted := store.records[0]
	if persisted.Result.ExecutionOutcome != "start-failed" || persisted.Result.Assessment != "failed" {
		t.Fatalf("denied record must be start-failed/failed: %+v", persisted.Result)
	}
	if persisted.Result.EffectObservationHash == nil || persisted.Result.EffectRecoveryAction == nil {
		t.Fatalf("denied record must bind effect observation: %+v", persisted.Result)
	}
	if *persisted.Result.EffectRecoveryAction != "reconcile-only" {
		t.Fatalf("external-write denial must reconcile, got %q", *persisted.Result.EffectRecoveryAction)
	}
	if !evidence.ValidHash(*persisted.Result.EffectObservationHash) {
		t.Fatalf("invalid observation hash %q", *persisted.Result.EffectObservationHash)
	}
	if err := validateRecord(record); err != nil {
		t.Fatalf("denied record must satisfy store validation: %v", err)
	}
}

func TestRunnerBindsObservationToPassedReadOnlyResult(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Effect = effectScopeForTest("read-only")
	store := &memoryResults{}
	record, err := testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if record.Result.Assessment != "passed" || record.Result.EffectObservationHash == nil || record.Result.EffectRecoveryAction == nil {
		t.Fatalf("passed read-only result must bind observation: %+v", record.Result)
	}
	if *record.Result.EffectRecoveryAction != "none-required" {
		t.Fatalf("completed read-only needs no recovery, got %q", *record.Result.EffectRecoveryAction)
	}
	if len(record.Result.ErrorEvidence) != 0 {
		t.Fatalf("passed result must have no error evidence: %+v", record.Result.ErrorEvidence)
	}
	if err := validateRecord(record); err != nil {
		t.Fatalf("passed record must satisfy store validation: %v", err)
	}
}

func (runner Runner) assessForTest(t *testing.T) Result {
	t.Helper()
	request := validRequest(t.TempDir())
	result, err := runner.assess(context.Background(), request, "result-one", time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC), successfulExecution(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestDefinitionHashUsesSchemaFieldNames(t *testing.T) {
	hash, err := DefinitionHash(validRequest(t.TempDir()).Hook)
	if err != nil || !evidence.ValidHash(hash) {
		t.Fatalf("definition hash: %q %v", hash, err)
	}
	canonical, err := evidence.EncodeCanonical(NetworkPolicy{Mode: "allowlist", Hosts: []string{"example.test"}, Ports: []int{443}})
	if err != nil || !strings.Contains(string(canonical), `"mode":"allowlist"`) || strings.Contains(string(canonical), `"Mode"`) {
		t.Fatalf("network policy JSON: %s %v", canonical, err)
	}
}
