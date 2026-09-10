package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const defaultInteractionsLedgerName = "operator-interactions.jsonl"

type InteractionItem struct {
	InteractionID    string   `json:"interactionId"`
	RunID            string   `json:"runId"`
	TaskID           string   `json:"taskId"`
	StepID           string   `json:"stepId"`
	Attempt          int      `json:"attempt"`
	OperatorID       string   `json:"operatorId"`
	RequestedAt      string   `json:"requestedAt"`
	Question         string   `json:"question"`
	Reason           string   `json:"reason"`
	AllowedResponses []string `json:"allowedResponses"`
	Risk             string   `json:"risk"`
	RequiredAction   string   `json:"requiredAction"`
	ContextHash      string   `json:"contextHash"`
	WaitSeconds      int64    `json:"waitSeconds"`
	NextAction       string   `json:"nextAction"`
}

type InteractionsReport struct {
	LedgerPath   string            `json:"ledgerPath"`
	Pending      []InteractionItem `json:"pending"`
	PendingCount int               `json:"pendingCount"`
	Warnings     []string          `json:"warnings"`
}

func (cli CLI) executeInteractions(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: prfrail interactions <list|respond> [options]")
		return exitUsage
	}
	switch args[0] {
	case "list":
		return cli.executeInteractionsList(args[1:], stdout, stderr)
	case "respond":
		return cli.executeInteractionsRespond(args[1:], stdout, stderr)
	case "tui":
		return cli.executeInteractionsTUI(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown interactions subcommand: %s\n", args[0])
		fmt.Fprintln(stderr, "usage: prfrail interactions <list|respond> [options]")
		return exitUsage
	}
}

func (cli CLI) executeInteractionsList(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("interactions list")
	ledgerPath := set.String("ledger", "", "path to operator interaction ledger")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "interactions list", *jsonOutput, parseErr.String(), "usage: prfrail interactions list [--ledger <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "interactions list", *jsonOutput, "interactions list does not accept positional arguments", "usage: prfrail interactions list [--ledger <path>] [--json]")
	}
	path, err := cli.resolveInteractionsLedgerPath(*ledgerPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions list", *jsonOutput, err)
	}
	report, err := BuildInteractionsReport(path, cli.now().UTC())
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions list", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "interactions list", OK: true, ExitCode: exitSuccess, Data: report})
	}
	fmt.Fprintf(stdout, "ledger: %s\n", report.LedgerPath)
	fmt.Fprintf(stdout, "pending: %d\n", report.PendingCount)
	for _, item := range report.Pending {
		fmt.Fprintf(stdout, "interaction: %s task=%s step=%s attempt=%d operator=%s wait=%ds allowed=%s next=%s\n", item.InteractionID, item.TaskID, item.StepID, item.Attempt, item.OperatorID, item.WaitSeconds, strings.Join(item.AllowedResponses, ","), item.NextAction)
		fmt.Fprintf(stdout, "question: %s\nrisk: %s\nrequiredAction: %s\n", item.Question, item.Risk, item.RequiredAction)
	}
	for _, warning := range report.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	return exitSuccess
}

func (cli CLI) executeInteractionsTUI(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("interactions tui")
	ledgerPath := set.String("ledger", "", "path to operator interaction ledger")
	actorID := set.String("actor-id", "local-operator", "responding operator identity")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "interactions tui", false, parseErr.String(), "usage: prfrail interactions tui [--ledger <path>] [--actor-id <id>]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "interactions tui", false, "interactions tui does not accept positional arguments", "usage: prfrail interactions tui [--ledger <path>] [--actor-id <id>]")
	}
	path, err := cli.resolveInteractionsLedgerPath(*ledgerPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions tui", false, err)
	}
	report, err := BuildInteractionsReport(path, cli.now().UTC())
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions tui", false, err)
	}
	fmt.Fprintf(stdout, "ProofRail operator interactions\npending: %d\n", report.PendingCount)
	if report.PendingCount == 0 {
		return exitSuccess
	}
	for index, item := range report.Pending {
		fmt.Fprintf(stdout, "%d. %s task=%s/%s attempt=%d wait=%ds\n", index+1, item.InteractionID, item.TaskID, item.StepID, item.Attempt, item.WaitSeconds)
		fmt.Fprintf(stdout, "   question: %s\n   reason: %s\n   risk: %s\n   required: %s\n   allowed: %s\n", item.Question, item.Reason, item.Risk, item.RequiredAction, strings.Join(item.AllowedResponses, ", "))
	}
	fmt.Fprint(stdout, "select interaction number (q to quit): ")
	reader := bufio.NewReader(cli.stdin)
	choiceText, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(choiceText) == "" {
		return writeCommandError(stdout, stderr, "interactions tui", false, fmt.Errorf("read interaction selection: %w", err))
	}
	choiceText = strings.TrimSpace(choiceText)
	if choiceText == "q" || choiceText == "quit" {
		return exitSuccess
	}
	choice, err := strconv.Atoi(choiceText)
	if err != nil || choice < 1 || choice > len(report.Pending) {
		return writeCommandError(stdout, stderr, "interactions tui", false, fmt.Errorf("invalid interaction selection %q", choiceText))
	}
	item := report.Pending[choice-1]
	fmt.Fprintf(stdout, "response (%s): ", strings.Join(item.AllowedResponses, ", "))
	selection, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(selection) == "" {
		return writeCommandError(stdout, stderr, "interactions tui", false, fmt.Errorf("read operator response: %w", err))
	}
	return cli.executeInteractionsRespond([]string{
		"--ledger", path,
		"--interaction-id", item.InteractionID,
		"--selection", strings.TrimSpace(selection),
		"--attempt", strconv.Itoa(item.Attempt),
		"--context-hash", item.ContextHash,
		"--actor-id", *actorID,
	}, stdout, stderr)
}

func (cli CLI) executeInteractionsRespond(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("interactions respond")
	ledgerPath := set.String("ledger", "", "path to operator interaction ledger")
	interactionID := set.String("interaction-id", "", "interaction id")
	selection := set.String("selection", "", "one allowed response id")
	actorID := set.String("actor-id", "local-operator", "responding operator identity")
	attempt := set.Int("attempt", 0, "expected task attempt")
	contextHash := set.String("context-hash", "", "expected context hash")
	jsonOutput := set.Bool("json", false, "output as JSON")
	usage := "usage: prfrail interactions respond --interaction-id <id> --selection <id> --attempt <n> --context-hash <hash> [--actor-id <id>] [--ledger <path>] [--json]"
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "interactions respond", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || *interactionID == "" || *selection == "" || *attempt < 1 || *contextHash == "" {
		return writeUsageError(stdout, stderr, "interactions respond", *jsonOutput, "interaction-id, selection, attempt and context-hash are required", usage)
	}
	path, err := cli.resolveInteractionsLedgerPath(*ledgerPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, err)
	}
	ledger, err := loadOperatorInteractionLedger(path)
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, err)
	}
	requestRecord, ok := findPendingInteraction(ledger, *interactionID)
	if !ok {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, fmt.Errorf("pending interaction %q not found", *interactionID))
	}
	request := requestRecord.Request
	if request.Attempt != *attempt || request.ContextHash != *contextHash {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, fmt.Errorf("stale interaction binding for %q", *interactionID))
	}
	response := chain.OperatorInteractionResponse{
		InteractionID: request.InteractionID, RespondedAt: cli.now().UTC().Format(timestampLayout),
		RespondedBy: evidence.Actor{Type: "operator", ID: *actorID},
		RunID:       request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID,
		ContextHash: request.ContextHash, RequestHash: requestRecord.RecordHash,
		Selection: *selection, Evidence: []string{requestRecord.RecordHash},
	}
	record, err := chain.NewOperatorInteractionResponseRecord(requestRecord, response, cli.now, nil)
	if err != nil {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, err)
	}
	if err := ledger.Append(record); err != nil {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, err)
	}
	appendRecord := cli.appendInteraction
	if appendRecord == nil {
		appendRecord = appendOperatorInteractionRecord
	}
	if err := appendRecord(path, record); err != nil {
		return writeCommandError(stdout, stderr, "interactions respond", *jsonOutput, err)
	}
	data := map[string]any{"interactionId": request.InteractionID, "selection": *selection, "recordHash": record.RecordHash, "response": record.Response}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "interactions respond", OK: true, ExitCode: exitSuccess, Message: "operator response recorded", Data: data})
	}
	fmt.Fprintf(stdout, "responded: %s selection=%s hash=%s\n", request.InteractionID, *selection, record.RecordHash)
	return exitSuccess
}

func (cli CLI) resolveInteractionsLedgerPath(value string) (string, error) {
	cwd, err := cli.getwd()
	if err != nil {
		return "", err
	}
	if value == "" {
		value = defaultInteractionsLedgerName
	}
	return resolvePathFromCWD(cwd, value)
}

func BuildInteractionsReport(path string, now time.Time) (InteractionsReport, error) {
	report := InteractionsReport{LedgerPath: path, Pending: []InteractionItem{}}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		report.Warnings = []string{"ledger not found, showing empty inbox"}
		return report, nil
	}
	ledger, err := loadOperatorInteractionLedger(path)
	if err != nil {
		return InteractionsReport{}, err
	}
	for _, record := range ledger.Pending() {
		request := record.Request
		requestedAt, _ := time.Parse(timestampLayout, request.RequestedAt)
		waitSeconds := int64(0)
		if now.After(requestedAt) {
			waitSeconds = int64(now.Sub(requestedAt).Seconds())
		}
		report.Pending = append(report.Pending, InteractionItem{
			InteractionID: request.InteractionID, RunID: request.RunID, TaskID: request.TaskID,
			StepID: request.StepID, Attempt: request.Attempt, OperatorID: request.Operator.ID,
			RequestedAt: request.RequestedAt, Question: request.Question, Reason: request.Reason,
			AllowedResponses: append([]string{}, request.AllowedResponses...), Risk: request.Risk,
			RequiredAction: request.RequiredAction, ContextHash: request.ContextHash, WaitSeconds: waitSeconds, NextAction: "respond",
		})
	}
	sort.Slice(report.Pending, func(i, j int) bool {
		return report.Pending[i].RequestedAt < report.Pending[j].RequestedAt
	})
	report.PendingCount = len(report.Pending)
	return report, nil
}

func loadOperatorInteractionLedger(path string) (chain.OperatorInteractionLedger, error) {
	ledger := chain.OperatorInteractionLedger{}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return ledger, nil
	}
	if err != nil {
		return ledger, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		record, err := chain.DecodeOperatorInteractionRecord([]byte(text))
		if err != nil {
			return ledger, fmt.Errorf("operator interaction ledger line %d: %w", line, err)
		}
		if err := ledger.Append(record); err != nil {
			return ledger, fmt.Errorf("operator interaction ledger line %d: %w", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return ledger, err
	}
	return ledger, nil
}

func findPendingInteraction(ledger chain.OperatorInteractionLedger, interactionID string) (chain.OperatorInteractionRecord, bool) {
	for _, record := range ledger.Pending() {
		if record.Request.InteractionID == interactionID {
			return record, true
		}
	}
	return chain.OperatorInteractionRecord{}, false
}

func appendOperatorInteractionRecord(path string, record chain.OperatorInteractionRecord) error {
	content, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(content, '\n')); err != nil {
		return err
	}
	return file.Sync()
}
