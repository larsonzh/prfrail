//go:build windows

package guard

// P83 F7 regression: the UTF-16 environment block the child receives carries the
// child's secrets, so it must have an owner and that owner must clear it on both
// the success path and every failure path (CONTRACTS.md:216).
//
// The seam is the unexported envClearHook on platformProcess: it runs after the
// deferred clear and therefore observes what the contract promises. Two paths are
// walked - a real spawn, and a spawn that fails after the block was built (an
// embedded NUL byte in the working directory makes the path conversion refuse).
//
// Mutation check: deleting "defer clear(environment)" leaves the canary in the
// observed block and turns both halves red.

import (
	"fmt"
	"slices"
	"strings"
	"syscall"
	"testing"
	"unicode/utf16"
)

const (
	environmentCanaryEntry = "CANARY_SECRET=proofrail-env-block-canary"
	environmentCanaryName  = "CANARY_SECRET"
)

// observeEnvironmentClear runs one start and reports the block the clear seam
// observed. It always releases the process it may have created.
func observeEnvironmentClear(t *testing.T, spec ProcessSpec) ([]uint16, error) {
	t.Helper()
	process, err := preparePlatformProcess()
	if err != nil {
		t.Fatalf("prepare the platform process: %v", err)
	}
	var observed []uint16
	process.envClearHook = func(block []uint16) { observed = slices.Clone(block) }
	command, startErr := process.start(&spec)
	if startErr == nil && command != nil {
		_, _ = process.stop(processIdentityForTest(process), 0)
	}
	process.close()
	return observed, startErr
}

// processIdentityForTest reports the identity of the child the process just
// created, for the stop call above. A start that produced no process yields the
// zero identity, which stop refuses without touching anything.
func processIdentityForTest(process *platformProcess) ProcessIdentity {
	handle, ok := process.processHandle()
	if !ok || handle == 0 {
		return ProcessIdentity{}
	}
	return ProcessIdentity{PID: int(handle), StartToken: "test"}
}

func TestWindowsStartClearsEnvironmentBlockOnSuccessAndFailure(t *testing.T) {
	// The block is the widened entry, its per-entry terminator, and the block
	// terminator: the length proves the observed slice is the environment block
	// that carried the canary, not an empty leftover.
	wantLength := len(utf16.Encode([]rune(environmentCanaryEntry))) + 2
	successSpec := ProcessSpec{
		Command: "cmd.exe",
		Args:    []string{"/c", "exit", "0"},
		Env:     []string{environmentCanaryEntry},
	}
	successBlock, successErr := observeEnvironmentClear(t, successSpec)
	if successErr != nil {
		t.Fatalf("the success path must spawn: %v", successErr)
	}
	assertEnvironmentBlockCleared(t, "success", successBlock, wantLength)

	failureSpec := successSpec
	// An embedded NUL byte makes the working directory conversion refuse AFTER the
	// environment block was built, which is exactly the early-return path that used
	// to leak the block.
	failureSpec.Dir = "C:\\invalid\x00directory"
	failureBlock, failureErr := observeEnvironmentClear(t, failureSpec)
	if failureErr == nil {
		t.Fatal("the induced failure must refuse the spawn")
	}
	assertEnvironmentBlockCleared(t, "failure", failureBlock, wantLength)
}

func assertEnvironmentBlockCleared(t *testing.T, name string, block []uint16, wantLength int) {
	t.Helper()
	if block == nil {
		t.Fatalf("%s: the clear seam never observed the environment block", name)
	}
	if len(block) != wantLength {
		t.Fatalf("%s: the observed block is %d wide but the canary entry is %d wide", name, len(block), wantLength)
	}
	for _, value := range block {
		if value != 0 {
			t.Fatalf("%s: the environment block still holds %q after the spawn returned", name, environmentCanaryName)
		}
	}
}

// gappedEnvironmentEntries builds an environment that is far larger than the
// pre-sized block (len(env)*8+1 wide), so buildEnvironmentBlock has to grow the
// block at least once. The canary stays first so the growth branch is entered
// with the canary entry already copied into the array that gets discarded.
func gappedEnvironmentEntries() ([]string, int) {
	entries := []string{environmentCanaryEntry}
	wantLength := len(utf16.Encode([]rune(environmentCanaryEntry))) + 2
	for index := 0; index < 64; index++ {
		entry := fmt.Sprintf("P83_GROWTH_%02d=%s", index, strings.Repeat("g", 48))
		entries = append(entries, entry)
		wantLength += len(utf16.Encode([]rune(entry))) + 1
	}
	return entries, wantLength
}

// TestWindowsStartClearsGrownEnvironmentBlock is the F7 growth-path half of the
// block contract: a block that had to grow still reaches the clear seam with the
// canary gone, on the success path and on the early-return failure path.
//
// Mutation check: deleting the deferred clear(environment) in start leaves the
// canary in the grown block and turns both halves red.
func TestWindowsStartClearsGrownEnvironmentBlock(t *testing.T) {
	entries, wantLength := gappedEnvironmentEntries()
	successSpec := ProcessSpec{
		Command: "cmd.exe",
		Args:    []string{"/c", "exit", "0"},
		Env:     entries,
	}
	successBlock, successErr := observeEnvironmentClear(t, successSpec)
	if successErr != nil {
		t.Fatalf("the grown-block success path must spawn: %v", successErr)
	}
	assertEnvironmentBlockCleared(t, "grown-success", successBlock, wantLength)

	failureSpec := successSpec
	failureSpec.Dir = "C:\\invalid\x00directory"
	failureBlock, failureErr := observeEnvironmentClear(t, failureSpec)
	if failureErr == nil {
		t.Fatal("the induced failure must refuse the spawn")
	}
	assertEnvironmentBlockCleared(t, "grown-failure", failureBlock, wantLength)
}

// TestBuildEnvironmentBlockGrowsWithoutLosingEntries is the content control for
// the same branch: growth must copy every entry, so the zeroization above cannot
// be satisfied by a block that lost its payload. It also asserts that growth
// actually happened, otherwise this test would silently stop covering the branch.
func TestBuildEnvironmentBlockGrowsWithoutLosingEntries(t *testing.T) {
	entries, _ := gappedEnvironmentEntries()
	block, err := buildEnvironmentBlock(entries)
	if err != nil {
		t.Fatalf("build the grown environment block: %v", err)
	}
	if presized := len(entries)*8 + 1; cap(block) <= presized {
		t.Fatalf("the block was not grown: cap=%d is not above the pre-sized %d", cap(block), presized)
	}
	want := make([]uint16, 0, len(block))
	for _, entry := range entries {
		want = append(want, utf16.Encode([]rune(entry))...)
		want = append(want, 0)
	}
	want = append(want, 0)
	if !slices.Equal(block, want) {
		t.Fatalf("the grown block lost or reordered entries: got %d wide, want %d wide", len(block), len(want))
	}
}

// TestClearEnvironmentScratchZeroesEveryEntry pins the helper the deferred clear
// in buildEnvironmentBlock relies on. It is the only part of the per-entry
// zeroization that is observable from outside the builder (see the gap noted in
// the P83 fix-round-4 tester report: buildEnvironmentBlock exposes no seam for the
// scratch slices it allocates, so their zeroization is pinned through this helper
// rather than through the build call).
//
// Mutation check: deleting clear(slice) inside the loop leaves the canary in the
// first slice and turns this red.
func TestClearEnvironmentScratchZeroesEveryEntry(t *testing.T) {
	first, err := syscall.UTF16FromString(environmentCanaryEntry)
	if err != nil {
		t.Fatalf("encode the canary entry: %v", err)
	}
	second, err := syscall.UTF16FromString("P83_SCRATCH_SECOND=" + strings.Repeat("s", 32))
	if err != nil {
		t.Fatalf("encode the second entry: %v", err)
	}
	scratch := [][]uint16{first, second}
	clearEnvironmentScratch(scratch)
	for index, slice := range scratch {
		if len(slice) == 0 {
			t.Fatalf("scratch slice %d was never populated", index)
		}
		for _, value := range slice {
			if value != 0 {
				t.Fatalf("scratch slice %d still holds %q after clearEnvironmentScratch", index, environmentCanaryName)
			}
		}
	}
}
