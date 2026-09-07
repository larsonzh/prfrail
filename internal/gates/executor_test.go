package gates

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

func TestGuardExecutorUsesArgvWithoutShell(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "injected")
	prepared := helperCommand(t, root, "echo", "; touch "+marker)
	execution := (GuardExecutor{}).Execute(context.Background(), prepared, Hook{Timeout: 5 * time.Second, Resources: ResourceLimits{OutputBytes: 1024}})
	if execution.Outcome != "exited" || execution.ExitCode == nil || *execution.ExitCode != 0 || string(execution.Stdout) != "; touch "+marker || fileExists(marker) {
		t.Fatalf("argv was not preserved: %+v", execution)
	}
}

func TestGuardExecutorTimesOutAndLimitsOutput(t *testing.T) {
	root := t.TempDir()
	timed := (GuardExecutor{}).Execute(context.Background(), helperCommand(t, root, "sleep"), Hook{Timeout: 50 * time.Millisecond, Resources: ResourceLimits{OutputBytes: 1024}})
	if timed.Outcome != "timed-out" || len(timed.TerminationEvidence) == 0 {
		t.Fatalf("timeout result: %+v", timed)
	}
	limited := (GuardExecutor{}).Execute(context.Background(), helperCommand(t, root, "output", "0123456789"), Hook{Timeout: 5 * time.Second, Resources: ResourceLimits{OutputBytes: 4}})
	if limited.Outcome != "resource-limited" || len(limited.Stdout) != 4 || !limited.OutputTruncated || len(limited.TerminationEvidence) == 0 {
		t.Fatalf("output result: %+v", limited)
	}
}

func TestGuardExecutorCleansDescendantsAfterParentExit(t *testing.T) {
	root := t.TempDir()
	ready := filepath.Join(root, "ready")
	ack := filepath.Join(root, "ack")
	completed := make(chan Execution, 1)
	go func() {
		completed <- (GuardExecutor{}).Execute(context.Background(), helperCommand(t, root, "parent-exits", ready, ack), Hook{Timeout: 5 * time.Second, Resources: ResourceLimits{OutputBytes: 1024}})
	}()
	var identity guard.ProcessIdentity
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(ready)
		if err == nil {
			pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
			if parseErr == nil {
				identity, err = guard.InspectProcess(pid)
				if err == nil {
					break
				}
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	if identity.PID == 0 {
		t.Fatal("descendant identity was not observed")
	}
	if err := os.WriteFile(ack, []byte("continue"), 0600); err != nil {
		t.Fatal(err)
	}
	execution := <-completed
	if runtime.GOOS == "linux" && execution.Outcome != "resource-limited" {
		t.Fatalf("linux descendant outcome: %+v", execution)
	}
	alive, err := guard.ProcessAlive(identity)
	if err != nil || alive {
		t.Fatalf("descendant survived managed parent: alive=%v err=%v execution=%+v", alive, err, execution)
	}
}

func helperCommand(t *testing.T, root string, arguments ...string) Prepared {
	t.Helper()
	return Prepared{Executable: os.Args[0], Args: append([]string{"-test.run=TestGateProcessHelper", "--"}, arguments...), CWD: root, Env: append(os.Environ(), "PROOFRAIL_GATE_HELPER=1")}
}

func TestGateProcessHelper(t *testing.T) {
	if os.Getenv("PROOFRAIL_GATE_HELPER") == "" {
		return
	}
	separator := 0
	for index, argument := range os.Args {
		if argument == "--" {
			separator = index
			break
		}
	}
	mode := os.Args[separator+1]
	switch mode {
	case "echo":
		_, _ = os.Stdout.WriteString(os.Args[separator+2])
		os.Exit(0)
	case "output":
		_, _ = os.Stdout.WriteString(os.Args[separator+2])
		for {
			time.Sleep(time.Hour)
		}
	case "sleep":
		for {
			time.Sleep(time.Hour)
		}
	case "parent-exits":
		child := exec.Command(os.Args[0], "-test.run=TestGateProcessHelper", "--", "sleep")
		child.Env = os.Environ()
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(os.Args[separator+2], []byte(strconv.Itoa(child.Process.Pid)), 0600); err != nil {
			os.Exit(4)
		}
		deadline := time.Now().Add(2 * time.Second)
		for !fileExists(os.Args[separator+3]) && time.Now().Before(deadline) {
			time.Sleep(5 * time.Millisecond)
		}
		os.Exit(0)
	default:
		os.Exit(2)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
