package release

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProbeCandidateAcceptsCompleteCLIResponses(t *testing.T) {
	root := t.TempDir()
	binaryPath := filepath.Join(root, "candidate")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}
	sourcePath := filepath.Join(root, "main.go")
	writeFile(t, sourcePath, []byte(probeHelperSource))
	command := exec.Command("go", "build", "-o", binaryPath, sourcePath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build probe helper: %v\n%s", err, output)
	}

	noopPath := filepath.Join(root, "noop.json")
	executablePath := filepath.Join(root, "executable.json")
	outputPath := filepath.Join(root, "candidate", "actual.json")
	writeFile(t, noopPath, []byte("{}"))
	writeFile(t, executablePath, []byte("{}"))
	if err := ProbeCandidate(context.Background(), CandidateProbePlan{
		BinaryPath:          binaryPath,
		NoopChainPath:       noopPath,
		ExecutableChainPath: executablePath,
		WorkRoot:            filepath.Join(root, "candidate", "work"),
		OutputPath:          outputPath,
		ExpectedVersion:     "candidate-test",
	}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var document oracleDocument
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Cases) != 4 {
		t.Fatalf("unexpected oracle cases: %+v", document.Cases)
	}
	for _, item := range document.Cases {
		if item.Outcome != "passed" {
			t.Fatalf("candidate case did not pass: %+v", item)
		}
	}

	response := commandResponse{Data: json.RawMessage(`{"chainId":"wrong","taskCount":1,"stepCount":1,"runnableInCli":true,"executableSteps":[],"configPath":"fixture.json","warnings":[]}`)}
	if err := validateNoopResponse(response); err == nil {
		t.Fatal("mismatched chain semantics must fail closed")
	}
	spoofed := probeCommand(context.Background(), binaryPath, root, "spoofed", 0, func(commandResponse) error { return nil }, "wrong-command")
	if spoofed.err == nil {
		t.Fatal("mismatched response command must fail closed")
	}
}

const probeHelperSource = `package main

import (
	"encoding/json"
	"os"
	"strings"
)

type response struct {
	Command string ` + "`json:\"command\"`" + `
	OK bool ` + "`json:\"ok\"`" + `
	ExitCode int ` + "`json:\"exitCode\"`" + `
	Data any ` + "`json:\"data,omitempty\"`" + `
	Message string ` + "`json:\"message,omitempty\"`" + `
	Error string ` + "`json:\"error,omitempty\"`" + `
	Usage string ` + "`json:\"usage,omitempty\"`" + `
}

func main() {
	command := os.Args[1]
	failed := command == "run" && strings.Contains(strings.Join(os.Args[2:], " "), "executable.json")
	item := response{Command: command, OK: !failed, Message: "validated"}
	switch command {
	case "version":
		item.Data = map[string]string{"version": "candidate-test"}
	case "validate":
		item.Data = map[string]any{"configPath": "noop.json", "chainId": "selfhost-noop", "taskCount": 1, "stepCount": 1, "runnableInCli": true, "executableSteps": []any{}, "warnings": []any{}}
	case "run":
		item.Data = map[string]any{"runId": "selfhost-run", "chainState": "COMPLETED", "sequence": 1, "lastEventHash": "sha256:test", "taskStateCounts": []any{}, "stepStateCounts": []any{}, "eventLogPath": "events.jsonl"}
	case "wrong-command":
		item.Command = "version"
	}
	if failed {
		item.ExitCode = 1
		item.Error = "local runtime only supports noop steps; first executable step is selfhost-task/build-candidate (build)"
		item.Usage = "prfrail run"
		item.Data = nil
	}
	_ = json.NewEncoder(os.Stdout).Encode(item)
	if failed {
		os.Exit(1)
	}
}
`
