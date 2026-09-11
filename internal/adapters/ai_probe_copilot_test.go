package adapters

import (
	"context"
	"slices"
	"strings"
	"testing"
)

type fakeCopilotCLIProbeExecutor struct {
	called              bool
	invocation          CopilotCLIProbeInvocation
	retainedEnvironment map[string][]byte
	execution           CopilotCLIProbeExecution
}

func (executor *fakeCopilotCLIProbeExecutor) Execute(_ context.Context, invocation CopilotCLIProbeInvocation) CopilotCLIProbeExecution {
	executor.called = true
	executor.retainedEnvironment = invocation.Environment
	executor.invocation = invocation
	executor.invocation.Args = slices.Clone(invocation.Args)
	executor.invocation.Environment = make(map[string][]byte, len(invocation.Environment))
	for name, value := range invocation.Environment {
		executor.invocation.Environment[name] = slices.Clone(value)
	}
	return executor.execution
}

func TestCopilotCLIProbeBuildsIsolatedDeepSeekSmoke(t *testing.T) {
	executor := &fakeCopilotCLIProbeExecutor{execution: CopilotCLIProbeExecution{
		Outcome:               "exited",
		ExitCode:              0,
		ResultSeen:            true,
		ResultExit:            0,
		ResultUsageSeen:       true,
		Reply:                 "PONG",
		RequestsUsed:          1,
		UsageSeen:             true,
		UsageUserRequests:     1,
		UsageProviderRequests: 1,
		Evidence:              []string{queueHashOne},
	}}
	probe := CopilotCLIProbe{
		Executable:       `C:\bin\copilot.exe`,
		WorkingDirectory: `C:\probe`,
		Executor:         executor,
	}
	credential := []byte("not-a-real-secret")
	result, err := probe.Probe(context.Background(), deepSeekProfile(), credential, "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "available" || result.RequestsUsed != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	wantArgs := []string{
		"-p", "Reply with exactly PONG and nothing else. Do not call any tool.",
		"-C", `C:\probe`,
		"--name", "proofrail-availability-probe",
		"--model", "deepseek-flash",
		"--output-format", "json",
		"--stream", "on",
		"--max-ai-credits", "1",
		"--disable-builtin-mcps",
		"--disallow-temp-dir",
		"--no-custom-instructions",
		"--no-ask-user",
		"--no-auto-update",
		"--no-color",
		"--no-remote",
		"--no-remote-export",
		"--secret-env-vars", "COPILOT_PROVIDER_API_KEY",
		"--deny-url=https://*",
		"--deny-url=http://*",
		"--available-tools=powershell",
		"--deny-tool=shell",
	}
	if !slices.Equal(executor.invocation.Args, wantArgs) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", executor.invocation.Args, wantArgs)
	}
	wantEnvironment := map[string]string{
		"COPILOT_PROVIDER_TYPE":     "anthropic",
		"COPILOT_PROVIDER_BASE_URL": "https://api.deepseek.com/anthropic",
		"COPILOT_PROVIDER_API_KEY":  "not-a-real-secret",
		"COPILOT_MODEL":             "deepseek-flash",
	}
	if len(executor.invocation.Environment) != len(wantEnvironment) {
		t.Fatalf("unexpected environment names: %#v", executor.invocation.Environment)
	}
	for name, want := range wantEnvironment {
		if got := string(executor.invocation.Environment[name]); got != want {
			t.Fatalf("unexpected %s value", name)
		}
	}
	for name, value := range executor.retainedEnvironment {
		for _, character := range value {
			if character != 0 {
				t.Fatalf("temporary %s value was not cleared", name)
			}
		}
	}
}

func TestParseCopilotCLIProbeOutputRequiresOneStructuredRequest(t *testing.T) {
	stdout := strings.Join([]string{
		`{"type":"model.call_start","data":{"model":"deepseek-flash"}}`,
		`{"type":"assistant.message","data":{"content":"PONG","toolRequests":[]}}`,
		`{"type":"result","exitCode":0,"usage":{"premiumRequests":0}}`,
	}, "\n")
	usage := `{"totalUserRequests":1,"modelMetrics":{"deepseek-flash":{"requests":{"count":1}}}}`
	execution := parseCopilotCLIProbeOutput([]byte(stdout), []byte(usage))
	if execution.FailureCode != "" || !execution.ResultSeen || execution.ResultExit != 0 || execution.Reply != "PONG" || execution.ToolCalls != 0 || execution.RequestsUsed != 1 || !execution.UsageSeen {
		t.Fatalf("unexpected execution: %+v", execution)
	}
}

func TestCopilotCLIProbeRejectsMissingOrContradictoryUsage(t *testing.T) {
	for name, execution := range map[string]CopilotCLIProbeExecution{
		"missing": {
			Outcome:         "exited",
			ExitCode:        0,
			ResultSeen:      true,
			ResultUsageSeen: true,
			Reply:           "PONG",
			RequestsUsed:    1,
		},
		"zero-user-requests": {
			Outcome:               "exited",
			ExitCode:              0,
			ResultSeen:            true,
			ResultUsageSeen:       true,
			Reply:                 "PONG",
			RequestsUsed:          1,
			UsageSeen:             true,
			UsageUserRequests:     0,
			UsageProviderRequests: 1,
		},
		"two-provider-requests": {
			Outcome:               "exited",
			ExitCode:              0,
			ResultSeen:            true,
			ResultUsageSeen:       true,
			Reply:                 "PONG",
			RequestsUsed:          2,
			UsageSeen:             true,
			UsageUserRequests:     1,
			UsageProviderRequests: 2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			executor := &fakeCopilotCLIProbeExecutor{execution: execution}
			probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
			result, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "agent-runner-cli", 1)
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != "unavailable" || result.Reason != "response_invalid" {
				t.Fatalf("incomplete usage was accepted: %+v", result)
			}
		})
	}
}

func TestCopilotCLIProbeRejectsResultWithoutUsage(t *testing.T) {
	executor := &fakeCopilotCLIProbeExecutor{execution: CopilotCLIProbeExecution{
		Outcome:               "exited",
		ExitCode:              0,
		ResultSeen:            true,
		Reply:                 "PONG",
		RequestsUsed:          1,
		UsageSeen:             true,
		UsageUserRequests:     1,
		UsageProviderRequests: 1,
	}}
	probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
	result, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "unavailable" || result.Reason != "response_invalid" {
		t.Fatalf("missing result usage was accepted: %+v", result)
	}
}

func TestParseCopilotCLIProbeOutputClassifiesOnlyStructuredFailure(t *testing.T) {
	stdout := `{"type":"error","data":{"error":{"code":"invalid_api_key","message":"not-a-real-secret"}}}`
	execution := parseCopilotCLIProbeOutput([]byte(stdout), nil)
	if execution.FailureCode != "credential_rejected" {
		t.Fatalf("unexpected failure classification: %+v", execution)
	}
}

func TestManagedCopilotCLIProbeExecutorFailsClosedWhenExecutableIsMissing(t *testing.T) {
	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: t.TempDir()}
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       t.TempDir() + string([]byte{92}) + "missing-copilot.exe",
		WorkingDirectory: t.TempDir(),
	})
	if execution.Outcome != "start-failed" || len(execution.Evidence) == 0 {
		t.Fatalf("unexpected execution: %+v", execution)
	}
}

func TestCopilotCLIProbeMapsStructuredFailuresWithoutDisclosure(t *testing.T) {
	for _, reason := range []string{
		"credential_rejected", "account_unavailable", "quota_exhausted", "network_unreachable",
		"tls_or_proxy_failure", "provider_unavailable", "model_unavailable", "policy_blocked",
	} {
		t.Run(reason, func(t *testing.T) {
			executor := &fakeCopilotCLIProbeExecutor{execution: CopilotCLIProbeExecution{
				Outcome:      "exited",
				ExitCode:     1,
				FailureCode:  reason,
				RequestsUsed: 1,
			}}
			probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
			secret := []byte("not-a-real-secret")
			result, err := probe.Probe(context.Background(), deepSeekProfile(), secret, "agent-runner-cli", 1)
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != "unavailable" || result.Reason != reason || strings.Contains(result.Reason+strings.Join(result.Evidence, ""), string(secret)) {
				t.Fatalf("unexpected result: %+v", result)
			}
		})
	}
}

func TestCopilotCLIProbeRejectsIncompleteResponseAndDoesNotFallback(t *testing.T) {
	invalidExecutions := map[string]CopilotCLIProbeExecution{
		"missing-result": {Outcome: "exited", ExitCode: 0, Reply: "PONG", RequestsUsed: 1},
		"wrong-reply":    {Outcome: "exited", ExitCode: 0, ResultSeen: true, Reply: "pong", RequestsUsed: 1},
		"tool-called":    {Outcome: "exited", ExitCode: 0, ResultSeen: true, Reply: "PONG", ToolCalls: 1, RequestsUsed: 1},
		"no-request":     {Outcome: "exited", ExitCode: 0, ResultSeen: true, Reply: "PONG"},
		"two-requests":   {Outcome: "exited", ExitCode: 0, ResultSeen: true, Reply: "PONG", RequestsUsed: 2},
	}
	for name, execution := range invalidExecutions {
		t.Run(name, func(t *testing.T) {
			executor := &fakeCopilotCLIProbeExecutor{execution: execution}
			probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
			result, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "agent-runner-cli", 1)
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != "unavailable" || result.Reason != "response_invalid" {
				t.Fatalf("unexpected result: %+v", result)
			}
		})
	}

	executor := &fakeCopilotCLIProbeExecutor{}
	probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
	if _, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "sessionbridge-silent", 1); err == nil {
		t.Fatal("expected unsupported channel error")
	}
	if executor.called {
		t.Fatal("probe fell back to the CLI for another channel")
	}
}

func TestCopilotCLIProbeMapsTimeoutWithoutRawError(t *testing.T) {
	executor := &fakeCopilotCLIProbeExecutor{execution: CopilotCLIProbeExecution{
		Outcome:      "timed-out",
		RequestsUsed: 1,
	}}
	probe := CopilotCLIProbe{Executable: "copilot.exe", WorkingDirectory: `C:\probe`, Executor: executor}
	result, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "unknown" || result.Reason != "timeout" {
		t.Fatalf("unexpected result: %+v", result)
	}
}
