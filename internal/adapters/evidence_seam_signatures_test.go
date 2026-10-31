package adapters

// P42-T027B-EVIDENCE-BIND ③ test-first (RED): compile-time pin on the contracts
// the slice may not change (design section 2 "Q3 contract compatibility", section
// 4 I5). Everything here is a package-level assertion: the package must fail to
// compile until each symbol exists with exactly this signature.
//
// Why signature assertions instead of behavioural tests: the three load seams and
// RunEvidenceDir are patched or called by existing tests (agent_runner_replay_store_test.go,
// agent_runner_replay_root_unix_test.go, console runtime_agent_e2e_test.go). A
// signature drift would break them silently at implementation time; a compile
// error cannot.

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// Q3 item 3: the three replay-store load seams keep their exact signatures.
	_ func(string, string) (AgentRunnerRequestRecord, bool, error)    = replayStoreLoadRequestRecord
	_ func(string, string) (AgentRunnerCompletionRecord, bool, error) = replayStoreLoadCompletionRecord
	_ func(string) (AgentRunnerReplayOwnershipRecord, bool, error)    = replayStoreLoadOwnershipRecord

	// Q3 item 1: RunEvidenceDir keeps returning the directory string.
	_ func(*AgentRunnerReplayStore, string) (string, error) = (*AgentRunnerReplayStore).RunEvidenceDir

	// Q3 item 2: the evidence collector signature the run/postflight/console
	// assembly depends on.
	_ func(context.Context, AgentRunnerEvidenceSpec) (AgentRunnerEvidence, error) = CollectAgentRunnerEvidence

	// Q2: the new adapter surfaces this slice adds.
	_ func(*AgentRunnerReplayStore, string, string, int64) ([]byte, bool, error) = (*AgentRunnerReplayStore).readBoundEvidence
	_ func(*AgentRunnerReplayStore)                                              = (*AgentRunnerReplayStore).CloseBoundDirs
	_ func(string) (*BoundDir, error)                                            = OpenBoundDir
	_ func(*BoundDir, string, int64) (*evidence.BoundFile, error)                = OpenBoundDirFile

	// Q1: the bound-read primitive in internal/evidence.
	_ func(string, int64) (*evidence.BoundFile, error)          = evidence.OpenBoundFile
	_ func(*evidence.BoundFile) ([]byte, bool, error)           = (*evidence.BoundFile).ReadAll
	_ func(*evidence.BoundFile, string) (string, error)         = (*evidence.BoundFile).Hash
	_ func(*evidence.BoundFile) (int64, time.Time, os.FileMode) = (*evidence.BoundFile).Identity
)

// TestEvidenceBindSeamSignaturesAreFrozen exists so the intent above is visible in
// test output; the assertions themselves are compile-time.
func TestEvidenceBindSeamSignaturesAreFrozen(t *testing.T) {
	var _ = t
}
