//go:build windows

package guard

// P83 FIX ROUND 6 counterexamples (items T7a and T7b).
//
// Both tests use environmentBlockScratchObserver, the seam that hands the
// builder's own slices to the test. The observer keeps the slice header, so the
// assertion runs after buildEnvironmentBlock returned and only a real clear() can
// have zeroed the backing array.

import (
	"strings"
	"testing"
)

// p83Fix6AssertZeroed fails when any unit of the observed slice is still readable.
func p83Fix6AssertZeroed(t *testing.T, stage string, block []uint16) {
	t.Helper()
	for offset, unit := range block {
		if unit != 0 {
			t.Fatalf("%s still holds %q at offset %d: the scratch block was not cleared", stage, string(rune(unit)), offset)
		}
	}
}

// TestP83Fix6T7aTerminatorBranchDiscardsAZeroedBlock is the T7a counterexample:
// the block whose encoded length equals the initial capacity reaches the
// terminator-headroom branch with len == cap, that branch repoints the slice and
// abandons the full backing array, so the abandoned array must be zeroed.
//
// The entries are sized so the growth branch inside the loop is never taken
// (3*8+1 = 25 = 9+8+8 units), which makes the terminator branch the only
// abandonment the builder performs.
//
// Mutation check: deleting the `clear(block)` in the len==cap branch leaves the
// canary readable and turns this red.
func TestP83Fix6T7aTerminatorBranchDiscardsAZeroedBlock(t *testing.T) {
	var discarded, grown []uint16
	environmentBlockScratchObserver = func(stage string, slice []uint16) {
		switch stage {
		case "discarded-block-terminator":
			discarded = slice
		case "discarded-block":
			grown = slice
		}
	}
	defer func() { environmentBlockScratchObserver = nil }()

	// Encoded lengths: "A=123456" is 8 chars -> 9 units, "B=12345" and "C=12345"
	// are 7 chars -> 8 units each, so the block ends at exactly cap(block) == 25.
	env := []string{"A=123456", "B=12345", "C=12345"}
	block, err := buildEnvironmentBlock(env)
	if err != nil {
		t.Fatalf("build the environment block: %v", err)
	}
	if len(block) != 26 || block[25] != 0 {
		t.Fatalf("the block must carry 25 units plus the terminator, got len=%d", len(block))
	}
	if len(grown) > 0 {
		t.Fatalf("test precondition failed: the in-loop growth branch was taken (%d units), so the terminator branch was not isolated", len(grown))
	}
	if len(discarded) != 25 {
		t.Fatalf("the terminator branch must abandon the exactly-full 25-unit array, got %d units", len(discarded))
	}
	p83Fix6AssertZeroed(t, "discarded-block-terminator", discarded)

	// The returned block is the grown copy and must still carry the secrets, so a
	// clear() that zeroed the wrong array is caught as well.
	decoded := make([]rune, 0, len(block))
	for _, unit := range block[:len(block)-1] {
		decoded = append(decoded, rune(unit))
	}
	// Each entry keeps its own terminating NUL, so the 25 units of the block spell
	// the three entries separated by NUL bytes and ending with one.
	if got := string(decoded); got != "A=123456\x00B=12345\x00C=12345\x00" {
		t.Fatalf("the returned block must still carry the child environment, got %q", got)
	}
}

// TestP83Fix6T7bFailedPartialBlockIsZeroed is the T7b counterexample: an entry
// that cannot be converted leaves a partially built block behind, and that block
// already holds the earlier entries' secrets.
//
// Mutation check: deleting the `clear(block)` on the conversion-error path leaves
// the canary readable in the reported partial block and turns this red.
func TestP83Fix6T7bFailedPartialBlockIsZeroed(t *testing.T) {
	var partial []uint16
	entries := make([][]uint16, 0, 2)
	environmentBlockScratchObserver = func(stage string, slice []uint16) {
		switch stage {
		case "partial-block":
			partial = slice
		case "entry":
			entries = append(entries, slice)
		}
	}
	defer func() { environmentBlockScratchObserver = nil }()

	// The second entry carries a NUL byte, which UTF16FromString refuses, so the
	// failure happens after the first entry is already in the block.
	env := []string{
		"P83_FIX6_CANARY_SECRET=" + strings.Repeat("c", 8),
		"P83_FIX6_BAD=va\x00lue",
	}
	block, err := buildEnvironmentBlock(env)
	if err == nil {
		t.Fatal("an entry containing a NUL byte must fail the conversion")
	}
	if block != nil {
		t.Fatalf("a failed build must return no block, got %d units", len(block))
	}
	if len(partial) == 0 {
		t.Fatal("the partial block must be non-empty, otherwise the observation proves nothing")
	}
	p83Fix6AssertZeroed(t, "partial-block", partial)
	// The per-entry intermediates are cleared by the deferred helper as well.
	if len(entries) == 0 {
		t.Fatal("no per-entry slice was observed, so the deferred clearing is unproven")
	}
	for index, entry := range entries {
		p83Fix6AssertZeroed(t, "entry "+string(rune('a'+index)), entry)
	}
}
