package console

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

func interactionTestHash(seed byte) string {
	return "sha256:" + strings.Repeat(string(seed), 64)
}

func appendInteractionTestRequest(t *testing.T, path string, cli CLI) chain.OperatorInteractionRecord {
	t.Helper()
	request := chain.OperatorInteractionRequest{
		InteractionID: "interaction-one", RequestedAt: "2026-09-08T10:00:00.000Z",
		RequestedBy: evidence.Actor{Type: "adapter", ID: "agent-one"},
		Operator:    evidence.Actor{Type: "operator", ID: "alice"},
		RunID:       "run-one", TaskID: "task-one", StepID: "step-one", Attempt: 2,
		WorkspaceHash: interactionTestHash('a'), ConversationID: "conversation-one",
		RequestID: "request-one", ContextHash: interactionTestHash('b'),
		Question: "Choose a migration target.", Reason: "Two targets are supported.",
		AllowedResponses: []string{"use-java-21", "use-java-25"},
		Risk:             "The target changes compatibility.", RequiredAction: "Select one target.",
		Evidence: []string{interactionTestHash('c')},
	}
	record, err := chain.NewOperatorInteractionRequestRecord(request, cli.now, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	if err := appendOperatorInteractionRecord(path, record); err != nil {
		t.Fatal(err)
	}
	return record
}

func TestInteractionsListRebuildsPendingInboxAfterRestart(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
	request := appendInteractionTestRequest(t, ledgerPath, cli)

	code, stdout, stderr := runCLI(t, newTestCLI(root), "interactions", "list", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report InteractionsReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.PendingCount != 1 || report.Pending[0].InteractionID != request.Request.InteractionID || report.Pending[0].NextAction != "respond" {
		t.Fatalf("unexpected pending inbox: %+v", report)
	}
	if report.Pending[0].Question == "" || len(report.Pending[0].AllowedResponses) != 2 || report.Pending[0].ContextHash != request.Request.ContextHash {
		t.Fatalf("inbox omitted operator facts: %+v", report.Pending[0])
	}
}

func TestInteractionsRespondPersistsAndClosesPendingItem(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
	request := appendInteractionTestRequest(t, ledgerPath, cli)

	code, stdout, stderr := runCLI(t, cli, "interactions", "respond", "--ledger", ledgerPath,
		"--interaction-id", "interaction-one", "--selection", "use-java-25", "--attempt", "2",
		"--context-hash", request.Request.ContextHash, "--actor-id", "alice", "--json")
	if code != 0 {
		t.Fatalf("respond code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !decodeResponse(t, stdout).OK {
		t.Fatalf("respond failed: %s", stdout)
	}

	report, err := BuildInteractionsReport(ledgerPath, cli.now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if report.PendingCount != 0 {
		t.Fatalf("answered interaction remains pending: %+v", report)
	}
	content, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	if lines := strings.Count(strings.TrimSpace(string(content)), "\n") + 1; lines != 2 {
		t.Fatalf("ledger lines=%d, want request plus response", lines)
	}

	code, stdout, _ = runCLI(t, cli, "interactions", "respond", "--ledger", ledgerPath,
		"--interaction-id", "interaction-one", "--selection", "use-java-25", "--attempt", "2",
		"--context-hash", request.Request.ContextHash, "--actor-id", "alice", "--json")
	if code != 1 || !strings.Contains(decodeResponse(t, stdout).Error, "not found") {
		t.Fatalf("duplicate response should fail: code=%d stdout=%s", code, stdout)
	}
}

func TestInteractionsRespondRejectsStaleOrUnstructuredInput(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
	}{
		{name: "wrong actor", arguments: []string{"--actor-id", "bob"}},
		{name: "stale attempt", arguments: []string{"--attempt", "1"}},
		{name: "wrong context", arguments: []string{"--context-hash", interactionTestHash('d')}},
		{name: "free text", arguments: []string{"--selection", "please use the newest version"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			cli := newTestCLI(root)
			ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
			request := appendInteractionTestRequest(t, ledgerPath, cli)
			arguments := []string{"interactions", "respond", "--ledger", ledgerPath, "--interaction-id", "interaction-one", "--selection", "use-java-21", "--attempt", "2", "--context-hash", request.Request.ContextHash, "--actor-id", "alice", "--json"}
			arguments = append(arguments, test.arguments...)
			code, stdout, _ := runCLI(t, cli, arguments...)
			if code != 1 {
				t.Fatalf("invalid response accepted: code=%d stdout=%s", code, stdout)
			}
			report, err := BuildInteractionsReport(ledgerPath, cli.now().UTC())
			if err != nil || report.PendingCount != 1 {
				t.Fatalf("invalid response changed pending inbox: report=%+v err=%v", report, err)
			}
		})
	}
}

func TestInteractionsListMissingLedgerHasNoSideEffects(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "missing.jsonl")
	code, stdout, stderr := runCLI(t, cli, "interactions", "list", "--ledger", ledgerPath, "--json")
	if code != 0 || stderr != "" {
		t.Fatalf("empty list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report InteractionsReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.PendingCount != 0 || len(report.Warnings) != 1 {
		t.Fatalf("unexpected empty report: %+v", report)
	}
	if _, err := os.Stat(ledgerPath); !os.IsNotExist(err) {
		t.Fatal("list must not create a ledger")
	}
}

func TestInteractionsTUIShowsFactsAndRecordsAllowedResponse(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
	appendInteractionTestRequest(t, ledgerPath, cli)
	cli.stdin = strings.NewReader("1\nuse-java-21\n")

	code, stdout, stderr := runCLI(t, cli, "interactions", "tui", "--ledger", ledgerPath, "--actor-id", "alice")
	if code != 0 || stderr != "" {
		t.Fatalf("tui code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, fact := range []string{"question:", "risk:", "allowed:", "wait=1800s", "responded: interaction-one"} {
		if !strings.Contains(stdout, fact) {
			t.Fatalf("tui output missing %q: %s", fact, stdout)
		}
	}
	if strings.Contains(stdout, "\x1b") {
		t.Fatalf("tui fallback must be ANSI-free: %q", stdout)
	}
	report, err := BuildInteractionsReport(ledgerPath, cli.now().UTC())
	if err != nil || report.PendingCount != 0 {
		t.Fatalf("tui response did not close pending item: report=%+v err=%v", report, err)
	}
}

func TestInteractionsTUIInvalidInputLeavesRequestPending(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
	appendInteractionTestRequest(t, ledgerPath, cli)
	cli.stdin = strings.NewReader("1\nfree form authorization\n")

	code, _, _ := runCLI(t, cli, "interactions", "tui", "--ledger", ledgerPath, "--actor-id", "alice")
	if code != 1 {
		t.Fatalf("free-form response should fail, code=%d", code)
	}
	report, err := BuildInteractionsReport(ledgerPath, cli.now().UTC())
	if err != nil || report.PendingCount != 1 {
		t.Fatalf("invalid tui input changed pending item: report=%+v err=%v", report, err)
	}
}

func TestInteractionsRespondPersistenceFailureLeavesRequestPending(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "operator-interactions.jsonl")
	request := appendInteractionTestRequest(t, ledgerPath, cli)
	cli.appendInteraction = func(string, chain.OperatorInteractionRecord) error {
		return errors.New("injected persistence failure")
	}

	code, stdout, _ := runCLI(t, cli, "interactions", "respond", "--ledger", ledgerPath,
		"--interaction-id", "interaction-one", "--selection", "use-java-21", "--attempt", "2",
		"--context-hash", request.Request.ContextHash, "--actor-id", "alice", "--json")
	if code != 1 || !strings.Contains(decodeResponse(t, stdout).Error, "persistence failure") {
		t.Fatalf("persistence failure should be visible: code=%d stdout=%s", code, stdout)
	}
	report, err := BuildInteractionsReport(ledgerPath, cli.now().UTC())
	if err != nil || report.PendingCount != 1 {
		t.Fatalf("persistence failure changed pending item: report=%+v err=%v", report, err)
	}
}
