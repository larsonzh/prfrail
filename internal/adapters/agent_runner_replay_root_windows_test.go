//go:build windows

package adapters

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"unsafe"
)

var replayRootTestKernel32 = syscall.NewLazyDLL("kernel32.dll")

func replayRootTestPathProc(procName, path string) (string, error) {
	proc := replayRootTestKernel32.NewProc(procName)
	input, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	buffer := make([]uint16, syscall.MAX_PATH+1)
	for {
		length, _, callErr := proc.Call(uintptr(unsafe.Pointer(input)), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
		if length == 0 {
			return "", callErr
		}
		if int(length) < len(buffer) {
			return syscall.UTF16ToString(buffer[:length]), nil
		}
		buffer = make([]uint16, length+1)
	}
}

func TestAgentRunnerReplayRootWindowsRejectsJunctionReplayRoot(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	link := filepath.Join(runRoot, "agent-runner-replay")
	command := exec.Command("cmd", "/c", "mklink", "/J", link, t.TempDir())
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("junction creation unavailable: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsAcceptsCaseAliasSpelling(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); err != nil {
		t.Fatal(err)
	}
	upper := strings.ToUpper(runRoot)
	if upper == runRoot {
		t.Skip("run root spelling has no letters to fold")
	}
	if _, err := os.Stat(upper); err != nil {
		t.Skipf("case-insensitive alias unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(upper, "run-one"); err != nil {
		t.Fatalf("case alias spelling must reopen the same replay root: %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsRejectsShortNameAliasOverlap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	shortRoot, err := replayRootTestPathProc("GetShortPathNameW", runRoot)
	if err != nil || strings.EqualFold(shortRoot, runRoot) {
		t.Skipf("short-name alias unavailable: short=%s err=%v", shortRoot, err)
	}
	shortReplayRoot := filepath.Join(shortRoot, "agent-runner-replay")
	if _, err := os.Stat(shortReplayRoot); err != nil {
		t.Skipf("short-name replay root unavailable: %v", err)
	}
	if err := store.RejectWriterRoots(shortReplayRoot); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected conflict for short-name writer root, got %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(shortRoot, "run-one"); err != nil {
		t.Fatalf("short-name run root must reopen the same store: %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsRejectsProtectedShortNameAlias(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	shortRoot, err := replayRootTestPathProc("GetShortPathNameW", runRoot)
	if err != nil || strings.EqualFold(shortRoot, runRoot) {
		t.Skipf("short-name alias unavailable: short=%s err=%v", shortRoot, err)
	}
	shortReplayRoot := filepath.Join(shortRoot, "agent-runner-replay")
	if _, err := os.Stat(shortReplayRoot); err != nil {
		t.Skipf("short-name replay root unavailable: %v", err)
	}
	if err := store.RejectProtectedRoots(shortReplayRoot); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected conflict for short-name protected root, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsAcceptsExtendedPathPrefix(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); err != nil {
		t.Fatal(err)
	}
	extended := `\\?\` + runRoot
	if _, err := os.Stat(extended); err != nil {
		t.Skipf("extended-length alias unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(extended, "run-one"); err != nil {
		t.Fatalf("extended-length alias must reopen the same replay root: %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsRejectsDriveRelativeWriterRoot(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	if err := store.RejectWriterRoots(`C:relative-workspace`); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected conflict for drive-relative writer root, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsRejectsRuntimeSubdirectorySwap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	requests := filepath.Join(runRoot, "agent-runner-replay", "requests")
	if err := os.RemoveAll(requests); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	command := exec.Command("cmd", "/c", "mklink", "/J", requests, external)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("junction creation unavailable: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	record := replayStoreRequestRecord(t, "request-runtime-swap")
	if _, err := store.RecordRequest(record); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error after runtime swap, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("external target must stay empty: entries=%v err=%v", entries, err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsRejectsRuntimeLaunchesSwap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	record := replayStoreRequestRecord(t, "request-runtime-launch-swap")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	launches := filepath.Join(runRoot, "agent-runner-replay", "launches")
	if err := os.RemoveAll(launches); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	command := exec.Command("cmd", "/c", "mklink", "/J", launches, external)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("junction creation unavailable: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error after runtime swap, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("external target must stay empty: entries=%v err=%v", entries, err)
	}
}

func TestAgentRunnerReplayStoreForRunWindowsConcurrentDistinctSpellings(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	upper := strings.ToUpper(runRoot)
	if upper == runRoot {
		t.Skip("run root spelling has no letters to fold")
	}
	if _, err := os.Stat(upper); err != nil {
		t.Skipf("case-insensitive alias unavailable: %v", err)
	}
	spellings := []string{runRoot, upper}
	errs := make([]error, len(spellings))
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for index, spelling := range spellings {
		waitGroup.Add(1)
		go func(index int, spelling string) {
			defer waitGroup.Done()
			<-start
			_, errs[index] = NewAgentRunnerReplayStoreForRun(spelling, "run-one")
		}(index, spelling)
	}
	close(start)
	waitGroup.Wait()
	for index, err := range errs {
		if err != nil {
			t.Fatalf("spelling %s failed: %v", spellings[index], err)
		}
	}
}
