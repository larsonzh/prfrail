package guard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestManagedProcessTerminatesDescendantTree(t *testing.T) {
	if os.Getenv("PROOFRAIL_PROCESS_HELPER") != "" {
		runProcessHelper()
		return
	}
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	managed, err := StartManaged(context.Background(), ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=TestManagedProcessTerminatesDescendantTree", "--", "parent", pidFile},
		Env:     append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=parent"),
	})
	if err != nil {
		t.Fatal(err)
	}
	var childIdentity ProcessIdentity
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		data, readErr := os.ReadFile(pidFile)
		if readErr == nil {
			pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
			if parseErr == nil {
				childIdentity, err = platformIdentity(pid)
				if err == nil {
					break
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if childIdentity.PID == 0 {
		t.Fatal("child process identity was not observed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	proof, err := managed.Terminate(ctx, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("terminate tree: %v (%+v)", err, proof)
	}
	if proof.Outcome != "stopped" || !evidence.ValidHash(proof.Hash) {
		t.Fatalf("invalid termination evidence: %+v", proof)
	}
	alive, err := platformIdentityAlive(childIdentity)
	if err != nil || alive {
		t.Fatalf("descendant remained alive: alive=%v err=%v", alive, err)
	}
}

func runProcessHelper() {
	separator := 0
	for index, arg := range os.Args {
		if arg == "--" {
			separator = index
			break
		}
	}
	if separator == 0 || len(os.Args) <= separator+2 {
		os.Exit(2)
	}
	mode, pidFile := os.Args[separator+1], os.Args[separator+2]
	if mode == "parent" {
		child := exec.Command(os.Args[0], "-test.run=TestManagedProcessTerminatesDescendantTree", "--", "child", pidFile)
		child.Env = append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=child")
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(pidFile, fmt.Appendf(nil, "%d\n", child.Process.Pid), 0600); err != nil {
			os.Exit(4)
		}
	}
	for {
		time.Sleep(time.Hour)
	}
}
