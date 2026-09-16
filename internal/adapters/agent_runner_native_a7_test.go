//go:build a7native

package adapters

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// The A7 native experiments kill a real child process at every publish stage.
// Windows Process.Kill is TerminateProcess: it flushes nothing and closes the
// handles, so what survives is the OS page cache, not a power-loss disk state.
// Everything here is therefore evidence about visibility and arbitration only -
// never about power-loss durability, which stays with B3.

const (
	a7StageEnv  = "PROOFRAIL_A7_STAGE"
	a7RootEnv   = "PROOFRAIL_A7_ROOT"
	a7ReadyEnv  = "PROOFRAIL_A7_READY"
	a7RoundsEnv = "PROOFRAIL_A7_ROUNDS"
)

const a7RequestID = "request-one"

// a7HelperBlockDuration keeps the helper alive until the parent kills it.
const a7HelperBlockDuration = 5 * time.Minute

func a7ExperimentRounds(t *testing.T) int {
	t.Helper()
	value := strings.TrimSpace(os.Getenv(a7RoundsEnv))
	if value == "" {
		return 20
	}
	rounds, err := strconv.Atoi(value)
	if err != nil || rounds <= 0 {
		t.Fatalf("invalid %s=%q", a7RoundsEnv, value)
	}
	return rounds
}

func a7PublishStages() []string {
	return []string{
		replayStorePublishStageTempWritten,
		replayStorePublishStageTempSynced,
		replayStorePublishStageTempClosed,
		replayStorePublishStageRecordLinked,
		replayStorePublishStageParentSynced,
	}
}

// a7StageLeavesRecordVisible reports whether the record is linked into place by
// the time the named stage is reached.
func a7StageLeavesRecordVisible(stage string) bool {
	return stage == replayStorePublishStageRecordLinked || stage == replayStorePublishStageParentSynced
}

func a7StoreAt(t *testing.T, root string) *AgentRunnerReplayStore {
	t.Helper()
	store, err := newAgentRunnerReplayStoreAt(root)
	if err != nil {
		t.Fatal(err)
	}
	// The experiments exercise the write path on a platform whose durability is
	// declared unproven; the override is the same seam the existing suite uses and
	// never lifts the declaration itself.
	store.publicationDurability = PublishDurabilityProven
	return store
}

// a7SpawnAndKillAtStage starts a child that publishes one request and stops
// existing at the named stage, waits until the child reports it reached that
// stage, then kills it.
func a7SpawnAndKillAtStage(t *testing.T, stage string) (root string, output string) {
	t.Helper()
	root = t.TempDir()
	ready := filepath.Join(t.TempDir(), "ready")
	cmd := exec.Command(os.Args[0], "-test.run="+a7RequestHelperFilter)
	cmd.Env = append(os.Environ(),
		a7StageEnv+"="+stage,
		a7RootEnv+"="+root,
		a7ReadyEnv+"="+ready,
	)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	if err := cmd.Start(); err != nil {
		t.Fatalf("start publish helper: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			t.Fatalf("publish helper never reached stage %s: %s", stage, buffer.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill publish helper: %v", err)
	}
	_ = cmd.Wait()
	return root, buffer.String()
}

// a7RunPublishHelper is the child half: it publishes one request and, when the
// named stage is reached, reports readiness and then blocks forever so the parent
// can kill it at that exact instant.
func a7RunPublishHelper(t *testing.T) {
	t.Helper()
	stage := os.Getenv(a7StageEnv)
	root := os.Getenv(a7RootEnv)
	ready := os.Getenv(a7ReadyEnv)
	store := a7StoreAt(t, root)
	replayStorePublishStageHook = func(current string, _ string) {
		if current != stage {
			return
		}
		if err := os.WriteFile(ready, []byte(current), 0o600); err != nil {
			panic(err)
		}
		// The parent kills this process while it waits here. A bare select{} would
		// trip the runtime's deadlock detector and the child would exit before the
		// kill, which is not the crash the experiment means to inject.
		time.Sleep(a7HelperBlockDuration)
	}
	if _, err := store.RecordRequest(replayStoreRequestRecord(t, a7RequestID)); err != nil {
		t.Fatalf("helper publish failed: %v", err)
	}
	panic("helper returned from a publish it was told to stop inside")
}

// TestA7NativePublishCrashMatrix kills a real process at every publish stage and
// classifies the replay root afterwards. The legal outcome set is closed: either
// no record is visible, or the record is complete and the state is
// dispatched-unknown. The two stages after the link must classify identically,
// which is the empirical form of "a kill cannot decide whether the directory
// durability step ran".
func TestA7NativePublishCrashMatrix(t *testing.T) {
	if os.Getenv(a7StageEnv) != "" {
		a7RunPublishHelper(t)
		return
	}
	rounds := a7ExperimentRounds(t)
	observed := map[string]map[AgentRunnerReplayState]int{}
	for _, stage := range a7PublishStages() {
		observed[stage] = map[AgentRunnerReplayState]int{}
		for round := 0; round < rounds; round++ {
			root, _ := a7SpawnAndKillAtStage(t, stage)
			restarted := a7StoreAt(t, root)
			state, err := restarted.State(a7RequestID)
			if err != nil {
				t.Fatalf("stage %s round %d: crash residue must reconcile, got %v", stage, round, err)
			}
			if state != AgentRunnerReplayStateAbsent && state != AgentRunnerReplayStateDispatchedUnknown {
				t.Fatalf("stage %s round %d: illegal replay state %s", stage, round, state)
			}
			path, err := restarted.requestPath(a7RequestID)
			if err != nil {
				t.Fatal(err)
			}
			_, statErr := os.Stat(path)
			if visible := statErr == nil; visible != a7StageLeavesRecordVisible(stage) {
				t.Fatalf("stage %s round %d: record visible = %v", stage, round, visible)
			}
			observed[stage][state]++
		}
	}
	for _, stage := range a7PublishStages() {
		t.Logf("A7 kill matrix stage=%-14s outcomes=%v rounds=%d", stage, observed[stage], rounds)
	}
	linked := observed[replayStorePublishStageRecordLinked]
	parentSynced := observed[replayStorePublishStageParentSynced]
	// Compared key by key rather than by formatted map: a formatted map is
	// order-dependent, so the equivalence claim would be unstable evidence.
	for _, state := range []AgentRunnerReplayState{AgentRunnerReplayStateAbsent, AgentRunnerReplayStateDispatchedUnknown} {
		if linked[state] != parentSynced[state] {
			t.Fatalf("the durability boundary must be observationally invisible: state %s counted %d vs %d", state, linked[state], parentSynced[state])
		}
	}
}

// TestA7NativeSingleWinnerAndKilledWinnerDoNotUnlockDispatch covers arbitration
// across a crash: concurrent publishers must produce one winner, and killing the
// winner must not hand the first dispatch to a later attempt.
func TestA7NativeSingleWinnerAndKilledWinnerDoNotUnlockDispatch(t *testing.T) {
	if os.Getenv(a7StageEnv) != "" {
		a7RunPublishHelper(t)
		return
	}
	// Concurrent publishers, one record: exactly one first dispatch.
	root := t.TempDir()
	store := a7StoreAt(t, root)
	request := replayStoreRequestRecord(t, a7RequestID)
	const publishers = 8
	var (
		mu       sync.Mutex
		first    int
		unknown  int
		failures []string
		wg       sync.WaitGroup
	)
	for i := 0; i < publishers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decision, err := store.RecordRequest(request)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failures = append(failures, err.Error())
				return
			}
			switch decision {
			case AgentRunnerReplayDecisionFirstDispatch:
				first++
			case AgentRunnerReplayDecisionUnknownBlock:
				unknown++
			default:
				failures = append(failures, "unexpected decision "+string(decision))
			}
		}()
	}
	wg.Wait()
	if len(failures) > 0 {
		t.Fatalf("concurrent publishers reported %v", failures)
	}
	if first != 1 {
		t.Fatalf("concurrent publishers produced %d first dispatches, want exactly 1", first)
	}
	if unknown != publishers-1 {
		t.Fatalf("concurrent publishers produced %d unknown blocks, want %d", unknown, publishers-1)
	}
	t.Logf("A7 arbitration: publishers=%d first=%d unknown-block=%d", publishers, first, unknown)

	// A killed winner must not unlock the dispatch for anyone else.
	crashRoot, _ := a7SpawnAndKillAtStage(t, replayStorePublishStageRecordLinked)
	restarted := a7StoreAt(t, crashRoot)
	state, err := restarted.State(a7RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("a killed winner leaves an unknown block, got state %s", state)
	}
	decision, err := restarted.RecordRequest(replayStoreRequestRecord(t, a7RequestID))
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionUnknownBlock {
		t.Fatalf("a killed winner must not unlock a second dispatch, decision = %s", decision)
	}
}

// TestA7NativeTempResidueIsInertButAccumulates records the crash residue of the
// pre-visibility stages: it never becomes a record and never blocks a later
// publication, and nothing removes it, so it accumulates until an operator does.
// The finding is reported, not "fixed": retiring residue is a production decision
// this slice is not allowed to make.
func TestA7NativeTempResidueIsInertButAccumulates(t *testing.T) {
	if os.Getenv(a7StageEnv) != "" {
		a7RunPublishHelper(t)
		return
	}
	root, _ := a7SpawnAndKillAtStage(t, replayStorePublishStageTempSynced)
	requests := filepath.Join(root, "requests")
	before, err := os.ReadDir(requests)
	if err != nil {
		t.Fatal(err)
	}
	residue := 0
	for _, entry := range before {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			residue++
		}
	}
	if residue == 0 {
		t.Fatal("a kill before the link must leave the temp file behind")
	}
	restarted := a7StoreAt(t, root)
	if state, err := restarted.State(a7RequestID); err != nil || state != AgentRunnerReplayStateAbsent {
		t.Fatalf("residue must not be a record: state=%s err=%v", state, err)
	}
	if decision, err := restarted.RecordRequest(replayStoreRequestRecord(t, a7RequestID)); err != nil || decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("residue must not block publication: decision=%s err=%v", decision, err)
	}
	after, err := os.ReadDir(requests)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("A7 temp residue: killed-before-link residue=%d files, directory entries after a successful publish=%d", residue, len(after))
	if len(after) == 0 {
		t.Fatal("the successful publish must be visible")
	}
	t.Logf("A7 temp residue finding: no startup cleanup retires %s", filepath.Join("requests", ".*.tmp"))
}
