package console

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testResponse struct {
	Command  string          `json:"command"`
	OK       bool            `json:"ok"`
	ExitCode int             `json:"exitCode"`
	Error    string          `json:"error"`
	Data     json.RawMessage `json:"data"`
}

func TestInitValidateAndExplainJSON(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	code, stdout, stderr := runCLI(t, cli, "init", "--workspace", root, "--json")
	if code != 0 {
		t.Fatalf("init code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if stderr != "" {
		t.Fatalf("init stderr not empty: %q", stderr)
	}
	initResult := decodeResponse(t, stdout)
	if !initResult.OK {
		t.Fatalf("init failed: %s", stdout)
	}
	var initData struct {
		ChainPath string `json:"chainPath"`
	}
	if err := json.Unmarshal(initResult.Data, &initData); err != nil {
		t.Fatal(err)
	}
	if initData.ChainPath == "" {
		t.Fatal("init did not return chain path")
	}

	code, stdout, stderr = runCLI(t, cli, "validate", "--chain", initData.ChainPath, "--json")
	if code != 0 {
		t.Fatalf("validate code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	validateResult := decodeResponse(t, stdout)
	if !validateResult.OK {
		t.Fatalf("validate failed: %s", stdout)
	}
	var summary ConfigValidationSummary
	if err := json.Unmarshal(validateResult.Data, &summary); err != nil {
		t.Fatal(err)
	}
	if summary.ChainID == "" || summary.TaskCount != 1 || summary.StepCount != 1 || !summary.RunnableInCLI {
		t.Fatalf("unexpected validate summary: %+v", summary)
	}

	code, stdout, stderr = runCLI(t, cli, "config", "explain", "--chain", initData.ChainPath, "--json")
	if code != 0 {
		t.Fatalf("config explain code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	explainResult := decodeResponse(t, stdout)
	if !explainResult.OK {
		t.Fatalf("config explain failed: %s", stdout)
	}
	var report ConfigExplainReport
	if err := json.Unmarshal(explainResult.Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.DocumentationPolicy != "if-affected" {
		t.Fatalf("unexpected documentation policy: %s", report.DocumentationPolicy)
	}
	foundDefault := false
	for _, item := range report.Resolution {
		if item.EffectivePointer == "/documentation/policy" && item.SourceKind == "builtin-default" {
			foundDefault = true
			break
		}
	}
	if !foundDefault {
		t.Fatalf("missing builtin-default resolution: %+v", report.Resolution)
	}
}

func TestRunAndReportNoopChain(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	code, stdout, stderr := runCLI(t, cli, "init", "--workspace", root, "--json")
	if code != 0 {
		t.Fatalf("init code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var initData struct {
		ChainPath string `json:"chainPath"`
	}
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &initData); err != nil {
		t.Fatal(err)
	}

	runDir := filepath.Join(root, "tmp", "prfrail-runs", "run-test")
	code, stdout, stderr = runCLI(t, cli, "run", "--chain", initData.ChainPath, "--run-id", "run-test", "--run-dir", runDir, "--json")
	if code != 0 {
		t.Fatalf("run code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	runResult := decodeResponse(t, stdout)
	if !runResult.OK {
		t.Fatalf("run failed: %s", stdout)
	}
	var runSummary RunSummary
	if err := json.Unmarshal(runResult.Data, &runSummary); err != nil {
		t.Fatal(err)
	}
	if runSummary.ChainState != "COMPLETED" || runSummary.RunID != "run-test" || runSummary.Sequence == 0 {
		t.Fatalf("unexpected run summary: %+v", runSummary)
	}
	if _, err := os.Stat(runSummary.EventLogPath); err != nil {
		t.Fatalf("event log missing: %v", err)
	}

	code, stdout, stderr = runCLI(t, cli, "report", "--run-dir", runDir, "--json")
	if code != 0 {
		t.Fatalf("report code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	reportResult := decodeResponse(t, stdout)
	if !reportResult.OK {
		t.Fatalf("report failed: %s", stdout)
	}
	var reportSummary RunSummary
	if err := json.Unmarshal(reportResult.Data, &reportSummary); err != nil {
		t.Fatal(err)
	}
	if reportSummary.ChainState != "COMPLETED" || reportSummary.Sequence == 0 {
		t.Fatalf("unexpected report summary: %+v", reportSummary)
	}
}

func TestRunRejectsExecutableSteps(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	cfg := DefaultChainConfig("chain-one")
	cfg.Tasks[0].Steps = []StepConfig{{
		ID:    "build-one",
		Kind:  "build",
		Hooks: []string{"go-build"},
	}}
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "proofrail.chain.json")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, cli, "run", "--chain", path, "--json")
	if code != 1 {
		t.Fatalf("run should fail for executable step, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if result.OK || !strings.Contains(result.Error, "noop") {
		t.Fatalf("unexpected run failure response: %s", stdout)
	}
}

func TestValidateRejectsUnknownField(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	cfg := DefaultChainConfig("chain-one")
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	raw := strings.Replace(string(content), "{\n", "{\n  \"unknown\": true,\n", 1)
	path := filepath.Join(root, "proofrail.chain.json")
	if err := os.WriteFile(path, []byte(raw), 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, _ := runCLI(t, cli, "validate", "--chain", path, "--json")
	if code != 1 {
		t.Fatalf("validate should fail, code=%d stdout=%s", code, stdout)
	}
	result := decodeResponse(t, stdout)
	if result.OK || !strings.Contains(result.Error, "unknown") {
		t.Fatalf("unexpected validate response: %s", stdout)
	}
}

func TestTextOutputHasNoANSI(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	if code, _, _ := runCLI(t, cli, "init", "--workspace", root); code != 0 {
		t.Fatalf("init failed with code %d", code)
	}
	code, stdout, stderr := runCLI(t, cli, "validate")
	if code != 0 {
		t.Fatalf("validate failed: code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if strings.Contains(stdout, "\x1b") || strings.Contains(stderr, "\x1b") {
		t.Fatalf("unexpected ANSI sequence in output stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestPreviewJSONAndTextParityWithoutSideEffects(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	cfg := DefaultChainConfig("preview-one")
	cfg.Tasks[0].Steps = []StepConfig{{
		ID:    "verify-config",
		Kind:  "verify",
		Hooks: []string{"go-test"},
	}}
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	chainPath := filepath.Join(root, "proofrail.chain.json")
	if err := os.WriteFile(chainPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, cli, "preview", "--chain", chainPath, "--json")
	if code != 0 {
		t.Fatalf("preview code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if stderr != "" {
		t.Fatalf("preview stderr not empty: %q", stderr)
	}
	result := decodeResponse(t, stdout)
	if !result.OK {
		t.Fatalf("preview failed: %s", stdout)
	}
	var report PreviewReport
	if err := json.Unmarshal(result.Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.PreviewRecord.Preview.Outcome != "blocked" || report.PreviewRecord.Preview.Mode != "offline-read-only" {
		t.Fatalf("unexpected preview summary: %+v", report.PreviewRecord.Preview)
	}
	if len(report.Unknowns) == 0 {
		t.Fatal("expected unknowns in preview report")
	}
	if report.CallCounters.CommandCalls != 0 || report.CallCounters.NetworkCalls != 0 || report.CallCounters.ModelCalls != 0 || report.CallCounters.CredentialReads != 0 || report.CallCounters.VersionProbes != 0 {
		t.Fatalf("expected zero call counters, got %+v", report.CallCounters)
	}

	code, textStdout, textStderr := runCLI(t, cli, "preview", "--chain", chainPath)
	if code != 0 {
		t.Fatalf("preview text code=%d stdout=%s stderr=%s", code, textStdout, textStderr)
	}
	if !strings.Contains(textStdout, report.PreviewRecord.PreviewHash) || !strings.Contains(textStdout, "outcome: blocked") {
		t.Fatalf("text output does not match json facts: %s", textStdout)
	}

	after, err := os.ReadFile(chainPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(content) {
		t.Fatal("preview modified the chain file")
	}
	if _, err := os.Stat(filepath.Join(root, "tmp", "prfrail-runs")); !os.IsNotExist(err) {
		t.Fatalf("preview should not create run directory: %v", err)
	}
}

func TestPreviewReadyForNoopChain(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)

	content, err := EncodeChainConfig(DefaultChainConfig("preview-ready"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "proofrail.chain.json")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, cli, "preview", "--chain", path, "--json")
	if code != 0 {
		t.Fatalf("preview code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if !result.OK {
		t.Fatalf("preview failed: %s", stdout)
	}
	var report PreviewReport
	if err := json.Unmarshal(result.Data, &report); err != nil {
		t.Fatal(err)
	}
	if !report.RunnableInCLI || report.PreviewRecord.Preview.Outcome != "ready" || len(report.PreviewRecord.Preview.BlockingEvidence) != 0 {
		t.Fatalf("unexpected ready preview report: %+v", report)
	}
}

func TestUnknownCommandUsageExitCode(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	code, _, stderr := runCLI(t, cli, "unknown")
	if code != 2 {
		t.Fatalf("unexpected code for unknown command: %d", code)
	}
	if !strings.Contains(stderr, "usage: prfrail") {
		t.Fatalf("missing usage for unknown command: %s", stderr)
	}
}

func newTestCLI(cwd string) CLI {
	cli := NewCLI("0.1.0-test")
	cli.now = func() time.Time {
		return time.Date(2026, 9, 8, 10, 30, 0, 0, time.UTC)
	}
	cli.getwd = func() (string, error) {
		return cwd, nil
	}
	return cli
}

func runCLI(t *testing.T, cli CLI, args ...string) (int, string, string) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := cli.Execute(context.Background(), args, stdout, stderr)
	return code, stdout.String(), stderr.String()
}

func decodeResponse(t *testing.T, text string) testResponse {
	t.Helper()
	var response testResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		t.Fatalf("failed to decode response %q: %v", text, err)
	}
	return response
}
