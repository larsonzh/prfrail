package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// The differential test builds the pinned deterministic CLI from source. It is
// the same artifact the production launcher runs, so the comparison measures the
// production capture/stop/mapping protocol instead of a mock.
const parityStubVersion = "agent-stub-0.1.0"

var (
	parityStubOnce sync.Once
	parityStubPath string
	parityStubErr  error
)

// buildAgentRunnerParityStub builds tools/agent-stub once per test binary. It is
// skipped rather than failed when the toolchain is unavailable, so a host without
// a Go toolchain reports an explicit SKIP instead of a red.
func buildAgentRunnerParityStub(t *testing.T) string {
	t.Helper()
	parityStubOnce.Do(func() {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			parityStubErr = err
			return
		}
		target := filepath.Join(os.TempDir(), "proofrail-agent-stub-parity")
		if err := os.MkdirAll(target, 0o755); err != nil {
			parityStubErr = err
			return
		}
		name := "agent-stub"
		if strings.EqualFold(filepath.Ext(os.Args[0]), ".exe") {
			name += ".exe"
		}
		parityStubPath = filepath.Join(target, name)
		build := exec.Command("go", "build", "-o", parityStubPath, "./tools/agent-stub")
		build.Dir = root
		if output, err := build.CombinedOutput(); err != nil {
			parityStubErr = fmt.Errorf("build agent-stub: %v: %s", err, output)
		}
	})
	if parityStubErr != nil {
		t.Skipf("the differential test needs a buildable agent-stub: %v", parityStubErr)
	}
	return parityStubPath
}

// ---------------------------------------------------------------------------
// Postflight port legs
// ---------------------------------------------------------------------------

// parityRequestRecord builds the request record the postflight port binds to.
// The effect mapping pair is mirrored from unexported constants in
// internal/adapters/agent_runner.go: the record is refused without them and they
// are contract values, not tunables, so a mirror is the only in-package way to
// assemble a legal request.
const (
	parityEffectMappingVersion = "1"
	parityEffectMappingHash    = "sha256:2c9609ee374bfc79bd0ed997aea6cfe4865c82ed0216d1a4a288139e0b0c92d4"
)

func agentRunnerPostflightRequestRecord(t *testing.T, taskID, stepID string, attempt int) AgentRunnerRequestRecord {
	t.Helper()
	record, err := NewAgentRunnerRequestRecord(AgentRunnerRequest{
		RequestID:            "request-one",
		RunID:                "run-one",
		TaskID:               taskID,
		StepID:               stepID,
		Attempt:              attempt,
		CreatedAt:            "2026-01-02T03:04:05.000Z",
		AdapterID:            "adapter-one",
		Mode:                 "create",
		WorkspaceHash:        evidence.Digest("", []byte("workspace")),
		ContextHash:          evidence.Digest("", []byte("context")),
		ParentSnapshotHash:   evidence.Digest("", []byte("parent")),
		AuthorizationHash:    evidence.Digest("", []byte("authorization")),
		BudgetHash:           evidence.Digest("", []byte("budget")),
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"local-process", "workspace-write"},
		EffectMappingVersion: parityEffectMappingVersion,
		EffectMappingHash:    parityEffectMappingHash,
		Evidence:             []string{evidence.Digest("", []byte("request-evidence"))},
	})
	if err != nil {
		t.Fatalf("build postflight request record: %v", err)
	}
	return record
}

type agentRunnerPostflightFixture struct {
	run          *AgentRunnerPinnedCLIRun
	store        *AgentRunnerReplayStore
	workspace    string
	snapshotRoot string
	record       AgentRunnerRequestRecord
	request      AgentRunnerRunRequest
	evidenceDir  string
	requestPath  chain.PostflightRequest
}

// newAgentRunnerPostflightFixture runs one real offline run and returns the
// recorded evidence plus the postflight request that binds it.
func newAgentRunnerPostflightFixture(t *testing.T, modeArgs []string) *agentRunnerPostflightFixture {
	t.Helper()
	runRoot := t.TempDir()
	workspace := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	snapshotRoot := filepath.Join(runRoot, "snapshot-store")
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		t.Fatal(err)
	}
	capturer, err := NewAgentRunnerSnapshotCapturer(snapshotStore, "run-one", "task-two", 1)
	if err != nil {
		t.Fatal(err)
	}
	config := agentRunnerPinnedCLIConfig("")
	helper := []string{"-test.run=^TestAgentRunnerPinnedCLIHelper$", "--"}
	config.RunArgs = append(append([]string{}, helper...), modeArgs...)
	launcher, err := NewAgentRunnerPinnedCLILauncher(config, store)
	if err != nil {
		t.Fatal(err)
	}
	record := agentRunnerPostflightRequestRecord(t, "task-two", "code-isolated", 1)
	request := AgentRunnerRunRequest{
		Launch:      agentRunnerPinnedCLIRequest(workspace),
		RequestHash: record.RecordHash,
		TaskID:      "task-two",
		StepID:      "code-isolated",
		Attempt:     1,
		Timeout:     15 * time.Second,
		Capturer:    capturer,
		MaxLogBytes: 1 << 20,
		VerifyWait:  1500 * time.Millisecond,
	}
	fixture := &agentRunnerPostflightFixture{
		run:          &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store},
		store:        store,
		workspace:    workspace,
		snapshotRoot: snapshotRoot,
		record:       record,
		request:      request,
	}
	outcome, err := fixture.run.Run(context.Background(), request)
	if err != nil {
		t.Fatalf("the fixture run must complete: %v", err)
	}
	if outcome.Completion.Status != "completed" || outcome.Facts == nil {
		t.Fatalf("the fixture must record a completed run with facts: %+v", outcome.Completion)
	}
	evidenceDir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatal(err)
	}
	fixture.evidenceDir = evidenceDir
	agentRunnerSeedEmptyChildPIDReport(t, evidenceDir)
	fixture.requestPath = chain.PostflightRequest{
		RunID:              "run-one",
		TaskID:             "task-two",
		Attempt:            1,
		ParentSnapshotHash: record.Request.ParentSnapshotHash,
		Workspace:          chain.Workspace{Root: workspace},
		Facts:              outcome.Facts.toChain(),
	}
	return fixture
}

func (fixture *agentRunnerPostflightFixture) port(t *testing.T) *AgentRunnerPostflight {
	t.Helper()
	port, err := NewAgentRunnerPostflight(fixture.store, fixture.record, fixture.snapshotRoot)
	if err != nil {
		t.Fatal(err)
	}
	return port
}

// TestAgentRunnerPostflightPassesOnRecordedEvidence proves the six-leg
// reconciliation accepts a clean, fully recorded run and carries every fact
// digest plus every leg digest.
func TestAgentRunnerPostflightPassesOnRecordedEvidence(t *testing.T) {
	fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
	decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != chain.PostflightPassed {
		t.Fatalf("a clean recorded run must pass: %+v", decision)
	}
	if len(decision.Evidence) != len(agentRunnerFactHashes(fixture.requestPath.Facts))+6 {
		t.Fatalf("a passed verdict must carry the five facts and the six legs: %d", len(decision.Evidence))
	}
	for index, hash := range agentRunnerFactHashes(fixture.requestPath.Facts) {
		if decision.Evidence[index] != hash {
			t.Fatalf("fact digest %d must travel with the pass: %s != %s", index, decision.Evidence[index], hash)
		}
	}
	if len(decision.ErrorEvidence) != 0 {
		t.Fatalf("a passed verdict carries no error evidence: %+v", decision.ErrorEvidence)
	}
}

// TestAgentRunnerPostflightRejectsTamperedLegs isolates each reconciliation leg:
// only one artifact is disturbed per case, so the failing leg is the one the
// setup targets.
func TestAgentRunnerPostflightRejectsTamperedLegs(t *testing.T) {
	t.Run("workspace change is a manifest failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		if err := os.WriteFile(filepath.Join(fixture.workspace, "tampered.txt"), []byte("unaccounted\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed {
			t.Fatalf("an unaccounted workspace write must fail: %+v", decision)
		}
		want := evidence.Digest(agentRunnerPostflightManifestFailDomain, []byte("added tampered.txt"))
		if len(decision.ErrorEvidence) != 1 || decision.ErrorEvidence[0] != want {
			t.Fatalf("the manifest leg must name the changed path: %+v", decision.ErrorEvidence)
		}
	})

	t.Run("tampered diff artifact is a diff failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		forged := []byte(`{"added":[],"removed":[],"changed":[{"path":"x"}]}`)
		if err := os.WriteFile(filepath.Join(fixture.evidenceDir, agentRunnerDiffFileName), forged, 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a tampered diff must fail on the diff leg: %+v", decision)
		}
		// Same strength as "a rewritten diff is a diff failure" below: a diff that
		// cannot bind to the frozen DiffHash must fail on the binding leg and name
		// the stored digest, so this verdict is distinguishable from leg 2 (the
		// re-derived diff comparison).
		want := evidence.Digest(agentRunnerPostflightDiffFailDomain, []byte(evidence.Digest("", forged)))
		if decision.ErrorEvidence[0] != want {
			t.Fatalf("the diff binding must name the stored digest: %+v", decision.ErrorEvidence)
		}
	})

	t.Run("tampered log artifact is a log failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		logPath := filepath.Join(fixture.evidenceDir, agentRunnerLogsDirectoryName, agentRunnerStdoutFileName)
		data, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(logPath, append(data, []byte("tampered\n")...), 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a tampered log must fail on the log leg: %+v", decision)
		}
	})

	t.Run("tampered usage artifact is a usage failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		if err := os.WriteFile(filepath.Join(fixture.evidenceDir, agentRunnerUsageFileName), []byte(`{"calls":9,"tokens":9,"durationMs":9}`), 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a tampered usage artifact must fail on the usage leg: %+v", decision)
		}
	})

	t.Run("a recorded live pid is a process failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		if err := os.WriteFile(filepath.Join(fixture.evidenceDir, "stub.pid"), fmt.Appendf(nil, "%d\n", os.Getpid()), 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a live recorded pid must fail on the process leg: %+v", decision)
		}
	})

	t.Run("a secret marker in the workspace is a secret failure", func(t *testing.T) {
		// The secret file exists before the run, so the frozen post manifest and the
		// fresh capture agree on it and the manifest leg cannot mask the secret leg.
		fixture := newAgentRunnerPostflightFixtureWithContent(t, "leaked.pem", "-----BEGIN PRIVATE KEY-----\n")
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a secret marker must fail on the secret leg: %+v", decision)
		}
		want := evidence.Digest(agentRunnerPostflightSecretFailDomain, []byte(gates.ErrSecretDetected.Error()))
		if decision.ErrorEvidence[0] != want {
			t.Fatalf("the secret leg must name the scanner verdict: %+v", decision.ErrorEvidence)
		}
	})

	t.Run("a missing artifact is uncertain, never a pass", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		if err := os.Remove(filepath.Join(fixture.evidenceDir, agentRunnerDiffFileName)); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightUncertain || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("an unreadable artifact must be uncertain: %+v", decision)
		}
	})
}

// TestAgentRunnerPostflightRejectsArtifactsNotBoundToFrozenFacts proves the
// stored artifacts are anchored to the frozen facts: rewriting manifest.post.json
// or diff.json must fail before any leg reconciles the artifact against itself.
// The rewrites keep the JSON valid, so the stored artifact still decodes and
// still agrees with the live workspace; only the frozen-fact binding catches
// them.
func TestAgentRunnerPostflightRejectsArtifactsNotBoundToFrozenFacts(t *testing.T) {
	t.Run("a rewritten post manifest is a manifest failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		path := filepath.Join(fixture.evidenceDir, agentRunnerManifestPostFileName)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		forged := append(append([]byte{}, raw...), ' ')
		if err := os.WriteFile(path, forged, 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a rewritten post manifest must fail on the manifest leg: %+v", decision)
		}
		want := evidence.Digest(agentRunnerPostflightManifestFailDomain, []byte(evidence.Digest("", forged)))
		if decision.ErrorEvidence[0] != want {
			t.Fatalf("the manifest binding must name the stored digest: %+v", decision.ErrorEvidence)
		}
	})

	t.Run("a rewritten diff is a diff failure", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		path := filepath.Join(fixture.evidenceDir, agentRunnerDiffFileName)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		forged := append(append([]byte{}, raw...), ' ')
		if err := os.WriteFile(path, forged, 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightFailed || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("a rewritten diff must fail on the diff leg: %+v", decision)
		}
		want := evidence.Digest(agentRunnerPostflightDiffFailDomain, []byte(evidence.Digest("", forged)))
		if decision.ErrorEvidence[0] != want {
			t.Fatalf("the diff binding must name the stored digest: %+v", decision.ErrorEvidence)
		}
	})
}

// TestAgentRunnerPostflightIsUncertainWithoutPIDReport proves the process leg
// fails closed when either pid report is missing: both reports are candidate
// written, so a run that left one or both cannot be verified and must never read
// as "no live process".
func TestAgentRunnerPostflightIsUncertainWithoutPIDReport(t *testing.T) {
	cases := []struct {
		name         string
		removeChild  bool
		removeParent bool
	}{
		{name: "both pid reports deleted is uncertain", removeChild: true, removeParent: true},
		{name: "a missing child report is uncertain", removeChild: true},
		{name: "a missing parent report is uncertain", removeParent: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
			remove := func(name string) {
				t.Helper()
				if err := os.Remove(filepath.Join(fixture.evidenceDir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
					t.Fatal(err)
				}
			}
			if testCase.removeChild {
				remove(agentRunnerChildPIDReportFileName)
			}
			if testCase.removeParent {
				remove(agentRunnerParentPIDReportFileName)
			}
			decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
			if err != nil {
				t.Fatal(err)
			}
			if decision.Outcome != chain.PostflightUncertain || len(decision.ErrorEvidence) != 1 {
				t.Fatalf("a run missing a pid report must be uncertain: %+v", decision)
			}
			want := evidence.Digest(agentRunnerPostflightUncertainDomain, []byte(agentRunnerMissingPIDReportDetail))
			if decision.ErrorEvidence[0] != want {
				t.Fatalf("the process leg must name the missing pid report: %+v", decision.ErrorEvidence)
			}
		})
	}
}

// TestAgentRunnerPostflightRejectsUnreadablePIDReports proves the process leg is
// fail-closed on the two silent-ignore shapes the stronger criterion removes: an
// unparsable pid line and a liveness answer that cannot be observed.
func TestAgentRunnerPostflightRejectsUnreadablePIDReports(t *testing.T) {
	t.Run("an unparsable pid line is uncertain", func(t *testing.T) {
		fixture := newAgentRunnerPostflightFixture(t, []string{"exit"})
		if err := os.WriteFile(filepath.Join(fixture.evidenceDir, agentRunnerParentPIDReportFileName), []byte("not-a-pid\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		decision, err := fixture.port(t).RunPostflight(context.Background(), fixture.requestPath)
		if err != nil {
			t.Fatal(err)
		}
		if decision.Outcome != chain.PostflightUncertain || len(decision.ErrorEvidence) != 1 {
			t.Fatalf("an unparsable pid line must be uncertain: %+v", decision)
		}
	})

	t.Run("an unobservable liveness answer is an error", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, agentRunnerParentPIDReportFileName), []byte("4242\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, agentRunnerChildPIDReportFileName), nil, 0o600); err != nil {
			t.Fatal(err)
		}
		alive, missing, err := agentRunnerAlivePIDsWith(dir,
			func(pid int) (guard.ProcessIdentity, error) {
				return guard.ProcessIdentity{PID: pid, StartToken: "token"}, nil
			},
			func(guard.ProcessIdentity) (bool, error) {
				return false, errors.New("the liveness probe exploded")
			})
		if !errors.Is(err, ErrInvalidAgentRunnerEvidence) {
			t.Fatalf("an unobservable liveness answer must fail closed: err=%v", err)
		}
		if alive != nil || missing != nil {
			t.Fatalf("a failed observation must not answer a verdict: alive=%v missing=%v", alive, missing)
		}
	})
}

// agentRunnerSeedEmptyChildPIDReport writes the "no child process was spawned"
// pid report. The offline harness stub writes stub-child.pid only when it spawns
// a child, so a run with no descendants leaves only stub.pid; leg 5 requires both
// reports, and an empty child report is the honest "reported no children" form.
func agentRunnerSeedEmptyChildPIDReport(t *testing.T, evidenceDir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName), nil, 0o600); err != nil {
		t.Fatal(err)
	}
}

// newAgentRunnerPostflightFixtureWithContent writes one workspace file before
// the run, so the file is part of both the frozen and the fresh manifest.
func newAgentRunnerPostflightFixtureWithContent(t *testing.T, name, content string) *agentRunnerPostflightFixture {
	t.Helper()
	runRoot := t.TempDir()
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	store := replayStoreForRunRoot(t, runRoot)
	snapshotRoot := filepath.Join(runRoot, "snapshot-store")
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		t.Fatal(err)
	}
	capturer, err := NewAgentRunnerSnapshotCapturer(snapshotStore, "run-one", "task-two", 1)
	if err != nil {
		t.Fatal(err)
	}
	config := agentRunnerPinnedCLIConfigForMode("exit")
	launcher, err := NewAgentRunnerPinnedCLILauncher(config, store)
	if err != nil {
		t.Fatal(err)
	}
	record := agentRunnerPostflightRequestRecord(t, "task-two", "code-isolated", 1)
	request := AgentRunnerRunRequest{
		Launch:      agentRunnerPinnedCLIRequest(workspace),
		RequestHash: record.RecordHash,
		TaskID:      "task-two",
		StepID:      "code-isolated",
		Attempt:     1,
		Timeout:     15 * time.Second,
		Capturer:    capturer,
		MaxLogBytes: 1 << 20,
		VerifyWait:  1500 * time.Millisecond,
	}
	fixture := &agentRunnerPostflightFixture{
		run:          &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store},
		store:        store,
		workspace:    workspace,
		snapshotRoot: snapshotRoot,
		record:       record,
		request:      request,
	}
	outcome, err := fixture.run.Run(context.Background(), request)
	if err != nil {
		t.Fatalf("the fixture run must complete: %v", err)
	}
	if outcome.Facts == nil {
		t.Fatalf("the fixture must record facts: %+v", outcome.Completion)
	}
	evidenceDir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatal(err)
	}
	fixture.evidenceDir = evidenceDir
	agentRunnerSeedEmptyChildPIDReport(t, evidenceDir)
	fixture.requestPath = chain.PostflightRequest{
		RunID:              "run-one",
		TaskID:             "task-two",
		Attempt:            1,
		ParentSnapshotHash: record.Request.ParentSnapshotHash,
		Workspace:          chain.Workspace{Root: workspace},
		Facts:              outcome.Facts.toChain(),
	}
	return fixture
}

// ---------------------------------------------------------------------------
// Differential test: AgentRunnerPinnedCLIRun.Run vs WaitAgentRunnerDispatched
// ---------------------------------------------------------------------------

// fixedParityClock is the one clock both paths inject, so CompletedAt is a
// comparable field rather than a timing artifact.
func fixedParityClock() time.Time {
	return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
}

type agentRunnerParityFixture struct {
	run       *AgentRunnerPinnedCLIRun
	store     *AgentRunnerReplayStore
	launcher  *AgentRunnerPinnedCLILauncher
	workspace string
	capturer  *AgentRunnerSnapshotCapturer
	request   AgentRunnerRunRequest
}

func newAgentRunnerParityFixture(t *testing.T, stub string, modeArgs []string, timeout time.Duration) *agentRunnerParityFixture {
	t.Helper()
	runRoot := t.TempDir()
	workspace := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	snapshotStore, err := snapshot.NewStore(filepath.Join(runRoot, "snapshot-store"), 0)
	if err != nil {
		t.Fatal(err)
	}
	capturer, err := NewAgentRunnerSnapshotCapturer(snapshotStore, "run-one", "task-two", 1)
	if err != nil {
		t.Fatal(err)
	}
	launcher, err := NewAgentRunnerPinnedCLILauncher(AgentRunnerPinnedCLIConfig{
		Executable:   stub,
		Version:      parityStubVersion,
		VersionArgs:  []string{"-version"},
		RunArgs:      modeArgs,
		CheckTimeout: 20 * time.Second,
		Grace:        500 * time.Millisecond,
	}, store)
	if err != nil {
		t.Fatal(err)
	}
	return &agentRunnerParityFixture{
		run:       &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store, Clock: fixedParityClock},
		store:     store,
		launcher:  launcher,
		workspace: workspace,
		capturer:  capturer,
		request: AgentRunnerRunRequest{
			Launch:      agentRunnerPinnedCLIRequest(workspace),
			RequestHash: evidence.Digest("", []byte("parity-request-record")),
			TaskID:      "task-two",
			StepID:      "code-isolated",
			Attempt:     1,
			Timeout:     timeout,
			Capturer:    capturer,
			MaxLogBytes: 1 << 20,
			VerifyWait:  1500 * time.Millisecond,
		},
	}
}

// TestAgentRunnerPostflightDispatchedParityMatchesRun is the drift guard for the
// (A′) entry: root A runs the self-spawning Run, root B starts the process with
// the dispatcher's public start port and then waits through
// WaitAgentRunnerDispatched. Every field that a real system can make identical
// must be identical; the three physically per-run fields are closed with
// self-consistency and structural equivalence instead.
func TestAgentRunnerPostflightDispatchedParityMatchesRun(t *testing.T) {
	stub := buildAgentRunnerParityStub(t)

	t.Run("completed leg", func(t *testing.T) {
		args := []string{"-mode", "mutate-workspace", "-sleep", "100ms"}
		a := newAgentRunnerParityFixture(t, stub, args, 15*time.Second)
		b := newAgentRunnerParityFixture(t, stub, args, 15*time.Second)

		outcomeA, err := a.run.Run(context.Background(), a.request)
		if err != nil {
			t.Fatalf("root A: %v", err)
		}
		before, err := b.capturer.CaptureWorkspaceManifest(context.Background(), b.workspace)
		if err != nil {
			t.Fatalf("root B pre-capture: %v", err)
		}
		if _, err := b.launcher.StartAgentRunnerProcess(context.Background(), b.request.Launch); err != nil {
			t.Fatalf("root B dispatch: %v", err)
		}
		outcomeB, err := b.run.WaitAgentRunnerDispatched(context.Background(), b.request, before)
		if err != nil {
			t.Fatalf("root B: %v", err)
		}

		if outcomeA.Completion.Status != "completed" || outcomeB.Completion.Status != "completed" {
			t.Fatalf("both paths must complete: A=%s B=%s", outcomeA.Completion.Status, outcomeB.Completion.Status)
		}
		if mismatches := agentRunnerCompletionMismatches(outcomeA.Completion, outcomeB.Completion, false); len(mismatches) > 0 {
			t.Fatalf("completion drift between Run and WaitAgentRunnerDispatched:\n%s", strings.Join(mismatches, "\n"))
		}
		if !reflect.DeepEqual(outcomeA.Settlement, outcomeB.Settlement) {
			t.Fatalf("settlement drift: A=%+v B=%+v", outcomeA.Settlement, outcomeB.Settlement)
		}
		if outcomeA.Facts == nil || outcomeB.Facts == nil {
			t.Fatalf("both completed outcomes must carry facts: A=%v B=%v", outcomeA.Facts, outcomeB.Facts)
		}
		if outcomeA.Facts.DiffHash != outcomeB.Facts.DiffHash ||
			outcomeA.Facts.LogHash != outcomeB.Facts.LogHash ||
			outcomeA.Facts.UsageHash != outcomeB.Facts.UsageHash {
			t.Fatalf("the deterministic facts must match: A=%+v B=%+v", *outcomeA.Facts, *outcomeB.Facts)
		}

		// OutputManifestHash / Facts.ManifestHash: the post manifest bytes carry the
		// capture timestamp, which the snapshot package cannot inject. Closed by
		// self-consistency (each hash digests its own recorded artifact) plus entry
		// equality.
		manifestA := agentRunnerReadRecordedManifest(t, a.store, agentRunnerManifestPostFileName)
		manifestB := agentRunnerReadRecordedManifest(t, b.store, agentRunnerManifestPostFileName)
		if !reflect.DeepEqual(manifestA.Manifest.Entries, manifestB.Manifest.Entries) {
			t.Fatalf("the recorded post entries must describe the same workspace: A=%+v B=%+v", manifestA.Manifest.Entries, manifestB.Manifest.Entries)
		}
		checkManifestSelfConsistency(t, a, outcomeA.Completion, *outcomeA.Facts)
		checkManifestSelfConsistency(t, b, outcomeB.Completion, *outcomeB.Facts)
		if *outcomeA.Facts == *outcomeB.Facts {
			t.Fatal("the two roots must not produce byte-identical facts: identity and time differ")
		}

		// Facts.ProcessStopEvidenceHash: the guard stop evidence canonical carries the
		// process identity and the requested/verified timestamps, so it cannot be
		// rebuilt here. Closed by structural equivalence: both recorded pids are gone
		// and both hashes are valid and distinct.
		if !evidence.ValidHash(outcomeA.Facts.ProcessStopEvidenceHash) || !evidence.ValidHash(outcomeB.Facts.ProcessStopEvidenceHash) {
			t.Fatalf("both stop-evidence hashes must be valid hashes")
		}
		if outcomeA.Facts.ProcessStopEvidenceHash == outcomeB.Facts.ProcessStopEvidenceHash {
			t.Fatal("the two roots must not produce the same stop evidence: the process identity differs")
		}
		agentRunnerRequireRecordedPIDsGone(t, a.store)
		agentRunnerRequireRecordedPIDsGone(t, b.store)
	})

	t.Run("guard stop leg", func(t *testing.T) {
		args := []string{"-mode", "hang"}
		a := newAgentRunnerParityFixture(t, stub, args, 300*time.Millisecond)
		b := newAgentRunnerParityFixture(t, stub, args, 300*time.Millisecond)

		outcomeA, err := a.run.Run(context.Background(), a.request)
		if err != nil {
			t.Fatalf("root A: %v", err)
		}
		before, err := b.capturer.CaptureWorkspaceManifest(context.Background(), b.workspace)
		if err != nil {
			t.Fatalf("root B pre-capture: %v", err)
		}
		if _, err := b.launcher.StartAgentRunnerProcess(context.Background(), b.request.Launch); err != nil {
			t.Fatalf("root B dispatch: %v", err)
		}
		outcomeB, err := b.run.WaitAgentRunnerDispatched(context.Background(), b.request, before)
		if err != nil {
			t.Fatalf("root B: %v", err)
		}

		if outcomeA.Completion.Status != "uncertain" || outcomeB.Completion.Status != "uncertain" {
			t.Fatalf("a run the guard had to stop must be uncertain: A=%s B=%s", outcomeA.Completion.Status, outcomeB.Completion.Status)
		}
		if outcomeA.Completion.ProcessTreeStatus != outcomeB.Completion.ProcessTreeStatus {
			t.Fatalf("tree status drift: A=%s B=%s", outcomeA.Completion.ProcessTreeStatus, outcomeB.Completion.ProcessTreeStatus)
		}
		// Both paths must have retired the launch through the same branch: the slot is
		// free again, which is exactly the B4 mirror divergence (a) the (A′) entry
		// removes.
		agentRunnerRequireLaunchSlotReleased(t, a, a.request.Launch)
		agentRunnerRequireLaunchSlotReleased(t, b, b.request.Launch)
	})

}

// agentRunnerFacts is a local alias so the map literal above stays readable.
type AgentRunnerFacts = AgentRunnerFrozenFacts

// agentRunnerCompletionMismatches lists the fields that must be identical
// between the two paths. OutputManifestHash is excluded from the strict set and
// closed by checkManifestSelfConsistency instead.
func agentRunnerCompletionMismatches(a, b AgentRunnerCompletion, includeManifest bool) []string {
	var mismatches []string
	record := func(name string, left, right any) {
		if !reflect.DeepEqual(left, right) {
			mismatches = append(mismatches, fmt.Sprintf("%s: A=%v B=%v", name, left, right))
		}
	}
	record("Kind", a.Kind, b.Kind)
	record("CompletionID", a.CompletionID, b.CompletionID)
	record("RequestID", a.RequestID, b.RequestID)
	record("RequestHash", a.RequestHash, b.RequestHash)
	record("RunID", a.RunID, b.RunID)
	record("TaskID", a.TaskID, b.TaskID)
	record("StepID", a.StepID, b.StepID)
	record("Attempt", a.Attempt, b.Attempt)
	record("AdapterID", a.AdapterID, b.AdapterID)
	record("SessionID", a.SessionID, b.SessionID)
	record("CompletedAt", a.CompletedAt, b.CompletedAt)
	record("Status", a.Status, b.Status)
	record("ProcessTreeStatus", a.ProcessTreeStatus, b.ProcessTreeStatus)
	record("LogsComplete", a.LogsComplete, b.LogsComplete)
	record("UsageComplete", a.UsageComplete, b.UsageComplete)
	record("Evidence", a.Evidence, b.Evidence)
	record("ErrorEvidence", a.ErrorEvidence, b.ErrorEvidence)
	record("ExitCode", agentRunnerIntValue(a.ExitCode), agentRunnerIntValue(b.ExitCode))
	if includeManifest {
		record("OutputManifestHash", agentRunnerStringValue(a.OutputManifestHash), agentRunnerStringValue(b.OutputManifestHash))
	}
	return mismatches
}

func agentRunnerIntValue(value *int) string {
	if value == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%d", *value)
}

func agentRunnerStringValue(value *string) string {
	if value == nil {
		return "<nil>"
	}
	return *value
}

// checkManifestSelfConsistency closes OutputManifestHash and Facts.ManifestHash
// with the artifact they claim to describe.
func checkManifestSelfConsistency(t *testing.T, fixture *agentRunnerParityFixture, completion AgentRunnerCompletion, facts AgentRunnerFrozenFacts) {
	t.Helper()
	raw := agentRunnerReadArtifact(t, fixture.store, "request-one", agentRunnerManifestPostFileName)
	want := evidence.Digest("", raw)
	if completion.OutputManifestHash == nil {
		t.Fatal("a completed receipt must publish the output manifest digest")
	}
	if *completion.OutputManifestHash != want {
		t.Fatalf("the output manifest digest must describe the recorded artifact: %s != %s", *completion.OutputManifestHash, want)
	}
	if facts.ManifestHash != want {
		t.Fatalf("the manifest fact must describe the recorded artifact: %s != %s", facts.ManifestHash, want)
	}
}

func agentRunnerReadRecordedManifest(t *testing.T, store *AgentRunnerReplayStore, name string) snapshot.SnapshotManifest {
	t.Helper()
	raw := agentRunnerReadArtifact(t, store, "request-one", name)
	manifest, err := snapshot.DecodeSnapshotManifest(raw)
	if err != nil {
		t.Fatalf("decode %s: %v", name, err)
	}
	return manifest
}

func agentRunnerReadArtifact(t *testing.T, store *AgentRunnerReplayStore, requestID, name string) []byte {
	t.Helper()
	dir, err := store.RunEvidenceDir(requestID)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// agentRunnerRequireRecordedPIDsGone proves the structural half of the stop
// evidence: nothing the CLI reported is alive any more.
func agentRunnerRequireRecordedPIDsGone(t *testing.T, store *AgentRunnerReplayStore) {
	t.Helper()
	dir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatal(err)
	}
	alive, _, err := agentRunnerAlivePIDs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(alive) > 0 {
		t.Fatalf("the run left recorded processes alive: %v", alive)
	}
}

// agentRunnerRequireLaunchSlotReleased proves both paths retired the launch: the
// slot is free and the same launcher can start again.
func agentRunnerRequireLaunchSlotReleased(t *testing.T, fixture *agentRunnerParityFixture, previous chain.AgentRunnerLaunchRequest) {
	t.Helper()
	if _, live := fixture.launcher.registry.live(); live {
		t.Fatal("the launch must be retired after the guard-stop branch")
	}
	reuse := previous
	reuse.RequestID = "request-reuse"
	reuse.RunID = "run-one"
	// Override the run args so the reuse launch exits on its own: this leg measures
	// the slot, not the workload.
	reuse.Args = []string{"-mode", "run", "-sleep", "50ms"}
	if _, err := fixture.launcher.StartAgentRunnerProcess(context.Background(), reuse); err != nil {
		t.Fatalf("a released launch slot must accept a new launch: %v", err)
	}
	// Retire it again. A short-lived workload may already have exited, in which case
	// the reviewed stop reports that race; the launch is retired either way, which is
	// what this leg proves. The recorded pid is then waited out so the suite never
	// leaks the probe.
	_, _ = fixture.launcher.StopAgentRunnerProcess(context.Background(), "launch-"+reuse.RequestID)
	if _, live := fixture.launcher.registry.live(); live {
		t.Fatal("the reuse launch must be retired after the slot proof")
	}
	agentRunnerWaitEvidencePIDsGone(t, fixture.store, reuse.RequestID, 5*time.Second)
}

// agentRunnerWaitEvidencePIDsGone bounds the wait for the recorded pids of one
// request to disappear, so the differential test never leaks a probe process.
func agentRunnerWaitEvidencePIDsGone(t *testing.T, store *AgentRunnerReplayStore, requestID string, bound time.Duration) {
	t.Helper()
	dir, err := store.RunEvidenceDir(requestID)
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(bound)
	for {
		alive, _, err := agentRunnerAlivePIDs(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(alive) == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("the reuse launch left processes alive: %v", alive)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// ---------------------------------------------------------------------------
// The differential-test clock is also used to prove the mapped receipt validates.
// ---------------------------------------------------------------------------

// TestAgentRunnerPostflightDispatchedParityReceiptsValidate proves both parity
// outcomes are legal receipts, which is what a caller hands to the publisher.
func TestAgentRunnerPostflightDispatchedParityReceiptsValidate(t *testing.T) {
	stub := buildAgentRunnerParityStub(t)
	fixture := newAgentRunnerParityFixture(t, stub, []string{"-mode", "mutate-workspace", "-sleep", "100ms"}, 15*time.Second)
	outcome, err := fixture.run.Run(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewAgentRunnerCompletionRecord(outcome.Completion); err != nil {
		t.Fatalf("the Run outcome must be a legal receipt: %v", err)
	}

	other := newAgentRunnerParityFixture(t, stub, []string{"-mode", "mutate-workspace", "-sleep", "100ms"}, 15*time.Second)
	before, err := other.capturer.CaptureWorkspaceManifest(context.Background(), other.workspace)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := other.launcher.StartAgentRunnerProcess(context.Background(), other.request.Launch); err != nil {
		t.Fatal(err)
	}
	dispatched, err := other.run.WaitAgentRunnerDispatched(context.Background(), other.request, before)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewAgentRunnerCompletionRecord(dispatched.Completion); err != nil {
		t.Fatalf("the (A′) outcome must be a legal receipt: %v", err)
	}
	// JSON round-trip: the two receipts must be structurally identical in the shape
	// the replay store persists, with only the three physical fields differing.
	encoder := func(completion AgentRunnerCompletion) string {
		data, err := json.Marshal(completion)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	if encoder(outcome.Completion) == encoder(dispatched.Completion) {
		t.Fatal("the two receipts must not be byte-identical: the output manifest timestamp differs")
	}
}

// TestAgentRunnerReadPIDsMatrix pins the exact pid-parsing contract of
// agentRunnerReadPIDs.
//
// COUPLING: this is a hand-synchronised twin of TestReadAgentRunnerEvidencePIDsMatrix
// in internal/console/runtime_agent_e2e_test.go. Both sides read the same report
// format and must answer the same verdict, so the matrix and its expectations are
// kept equal by hand, exactly like agentRunnerMissingPIDReportDetail and the pid
// report file names. The parse is pure (t.TempDir only), so the case runs
// identically on Windows and Linux.
func TestAgentRunnerReadPIDsMatrix(t *testing.T) {
	cases := []struct {
		name        string
		content     string
		missing     bool
		wantPIDs    []int
		wantPresent bool
		wantInvalid bool
	}{
		{name: "plain 7", content: "7", wantPIDs: []int{7}, wantPresent: true},
		{name: "signed plus 7", content: "+7", wantPIDs: []int{7}, wantPresent: true},
		{name: "padded 7", content: " 7 ", wantPIDs: []int{7}, wantPresent: true},
		{name: "zero padded 007", content: "007", wantPIDs: []int{7}, wantPresent: true},
		{name: "CRLF terminated 7", content: "7\r\n", wantPIDs: []int{7}, wantPresent: true},
		{name: "two pids with an interior blank line", content: "7\n\n8", wantPIDs: []int{7, 8}, wantPresent: true},
		{name: "trailing garbage 12abc", content: "12abc", wantInvalid: true},
		{name: "float 7.0", content: "7.0", wantInvalid: true},
		{name: "underscore 1_0", content: "1_0", wantInvalid: true},
		{name: "hex 0x10", content: "0x10", wantInvalid: true},
		{name: "negative zero -0", content: "-0", wantInvalid: true},
		{name: "negative -5", content: "-5", wantInvalid: true},
		{name: "int64 overflow", content: "9223372036854775808", wantInvalid: true},
		{name: "missing file", missing: true},
		{name: "empty file", content: "", wantPIDs: []int{}, wantPresent: true},
		{name: "whitespace only", content: "   \n\t\n  \n", wantPIDs: []int{}, wantPresent: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), agentRunnerParentPIDReportFileName)
			if !tc.missing {
				if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			pids, present, err := agentRunnerReadPIDs(path)
			if tc.missing {
				if err != nil || present || pids != nil {
					t.Fatalf("a missing report must answer present=false without error: pids=%v present=%v err=%v", pids, present, err)
				}
				return
			}
			if tc.wantInvalid {
				if !errors.Is(err, ErrInvalidAgentRunnerEvidence) || present {
					t.Fatalf("an unparsable report must fail closed: pids=%v present=%v err=%v", pids, present, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("a parsable report must not error: %v", err)
			}
			if present != tc.wantPresent || !agentRunnerPIDsEqual(pids, tc.wantPIDs) {
				t.Fatalf("pids=%v present=%v, want pids=%v present=%v", pids, present, tc.wantPIDs, tc.wantPresent)
			}
		})
	}
}

// agentRunnerPIDsEqual compares pid slices element by element, so an empty result
// compares equal whether it is nil or a zero-length slice.
func agentRunnerPIDsEqual(got, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
