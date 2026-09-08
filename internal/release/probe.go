package release

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CandidateProbePlan struct {
	BinaryPath          string
	NoopChainPath       string
	ExecutableChainPath string
	WorkRoot            string
	OutputPath          string
	ExpectedVersion     string
}

type commandResponse struct {
	Command  string          `json:"command"`
	OK       bool            `json:"ok"`
	ExitCode int             `json:"exitCode"`
	Data     json.RawMessage `json:"data"`
	Message  string          `json:"message,omitempty"`
	Error    string          `json:"error,omitempty"`
	Usage    string          `json:"usage,omitempty"`
}

type probeCase struct {
	oracleCase
	err error
}

func ProbeCandidate(ctx context.Context, plan CandidateProbePlan) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", ErrInvalidSelfHost)
	}
	binaryPath, err := existingRealPath(plan.BinaryPath)
	if err != nil {
		return fmt.Errorf("%w: candidate binary: %v", ErrInvalidSelfHost, err)
	}
	noopPath, err := existingRealPath(plan.NoopChainPath)
	if err != nil {
		return fmt.Errorf("%w: noop chain: %v", ErrInvalidSelfHost, err)
	}
	executablePath, err := existingRealPath(plan.ExecutableChainPath)
	if err != nil {
		return fmt.Errorf("%w: executable chain: %v", ErrInvalidSelfHost, err)
	}
	workRoot, err := filepath.Abs(plan.WorkRoot)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(workRoot, 0755); err != nil {
		return err
	}

	cases := []probeCase{
		probeCommand(ctx, binaryPath, workRoot, "version", 0, validateVersionResponse(plan.ExpectedVersion), "version", "--json"),
		probeCommand(ctx, binaryPath, workRoot, "validate-noop", 0, validateNoopResponse, "validate", "--chain", noopPath, "--json"),
		probeCommand(ctx, binaryPath, workRoot, "run-noop", 0, validateRunResponse, "run", "--chain", noopPath, "--run-id", "selfhost-run", "--run-dir", filepath.Join(workRoot, "run"), "--json"),
		probeCommand(ctx, binaryPath, workRoot, "reject-executable", 1, validateRejectResponse, "run", "--chain", executablePath, "--json"),
	}
	oracleCases := make([]oracleCase, 0, len(cases))
	for _, item := range cases {
		if item.err != nil {
			return fmt.Errorf("%w: candidate case %q failed: %v", ErrInvalidSelfHost, item.ID, item.err)
		}
		oracleCases = append(oracleCases, item.oracleCase)
	}
	document := oracleDocument{SchemaVersion: "1.0.0", Cases: oracleCases}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeAtomic(plan.OutputPath, data)
}

func probeCommand(ctx context.Context, binaryPath, workRoot, id string, expectedExit int, validate func(commandResponse) error, args ...string) probeCase {
	item := probeCase{oracleCase: oracleCase{ID: id, Outcome: "failed"}}
	command := exec.CommandContext(ctx, binaryPath, args...)
	command.Dir = workRoot
	command.Env = []string{"PATH=" + os.Getenv("PATH"), "SYSTEMROOT=" + os.Getenv("SYSTEMROOT"), "WINDIR=" + os.Getenv("WINDIR")}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			item.err = err
			return item
		}
		exitCode = exitError.ExitCode()
	}
	if exitCode != expectedExit || stderr.Len() != 0 {
		item.err = fmt.Errorf("exit=%d expected=%d stderr=%q", exitCode, expectedExit, stderr.String())
		return item
	}
	var response commandResponse
	if err := decodeStrict(stdout.Bytes(), &response); err != nil {
		item.err = fmt.Errorf("decode stdout: %v", err)
		return item
	}
	if response.Command != args[0] || response.ExitCode != expectedExit || response.OK != (expectedExit == 0) {
		item.err = fmt.Errorf("response command=%q exit=%d ok=%t", response.Command, response.ExitCode, response.OK)
		return item
	}
	if err := validate(response); err != nil {
		item.err = err
		return item
	}
	item.Outcome = "passed"
	return item
}

func validateVersionResponse(expected string) func(commandResponse) error {
	return func(response commandResponse) error {
		var data struct {
			Version string `json:"version"`
		}
		if expected == "" || decodeStrict(response.Data, &data) != nil || data.Version != expected {
			return fmt.Errorf("unexpected version response")
		}
		return nil
	}
}

func validateNoopResponse(response commandResponse) error {
	var data struct {
		ChainID         string            `json:"chainId"`
		TaskCount       int               `json:"taskCount"`
		StepCount       int               `json:"stepCount"`
		RunnableInCLI   bool              `json:"runnableInCli"`
		ExecutableSteps []json.RawMessage `json:"executableSteps"`
		ConfigPath      string            `json:"configPath"`
		Warnings        []string          `json:"warnings"`
	}
	if err := decodeStrict(response.Data, &data); err != nil {
		return fmt.Errorf("decode validate data: %v", err)
	}
	if data.ChainID != "selfhost-noop" || data.TaskCount != 1 || data.StepCount != 1 || !data.RunnableInCLI || len(data.ExecutableSteps) != 0 || data.ConfigPath == "" {
		return fmt.Errorf("unexpected validate data")
	}
	return nil
}

func validateRunResponse(response commandResponse) error {
	var data struct {
		RunID           string            `json:"runId"`
		ChainState      string            `json:"chainState"`
		Sequence        int               `json:"sequence"`
		LastEventHash   string            `json:"lastEventHash"`
		TaskStateCounts []json.RawMessage `json:"taskStateCounts"`
		StepStateCounts []json.RawMessage `json:"stepStateCounts"`
		EventLogPath    string            `json:"eventLogPath"`
	}
	if err := decodeStrict(response.Data, &data); err != nil {
		return fmt.Errorf("decode run data: %v", err)
	}
	if data.RunID != "selfhost-run" || data.ChainState != "COMPLETED" || data.Sequence < 1 || data.LastEventHash == "" || data.EventLogPath == "" {
		return fmt.Errorf("unexpected run data")
	}
	return nil
}

func validateRejectResponse(response commandResponse) error {
	if !strings.Contains(response.Error, "local runtime only supports noop steps") || len(response.Data) != 0 {
		return fmt.Errorf("unexpected executable rejection")
	}
	return nil
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(temporary, path); err != nil {
		_ = os.Remove(temporary)
		return err
	}
	return nil
}
