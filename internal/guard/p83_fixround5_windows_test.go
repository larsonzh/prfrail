//go:build windows

package guard

// P83 FIX ROUND 5 counterexample (item T7a): growing the environment block
// abandons a backing array that still holds the child's environment, so that
// array must be zeroed before the slice is repointed.
//
// The seam is environmentBlockScratchObserver, which hands the builder's own
// slices to the test; the collected backing array is inspected AFTER the call, so
// only a real clear() can zero it.
//
// Mutation check: deleting `clear(block)` above the "discarded-block" report
// leaves the canary readable and turns this red.

import (
	"strings"
	"testing"
)

func TestP83T7aDiscardedEnvironmentBlockIsZeroed(t *testing.T) {
	var discarded [][]uint16
	environmentBlockScratchObserver = func(stage string, slice []uint16) {
		if stage == "discarded-block" || stage == "discarded-block-terminator" {
			discarded = append(discarded, slice)
		}
	}
	defer func() { environmentBlockScratchObserver = nil }()

	// Three entries: the initial capacity is 3*8+1=25, so the second entry needs 41
	// and forces a grow-and-repoint while the abandoned array already holds the
	// first entry - the branch that actually carries a secret.
	env := []string{
		"P83_FIX5_CANARY_A=" + strings.Repeat("a", 12),
		"P83_FIX5_CANARY_B=" + strings.Repeat("b", 12),
		"P83_FIX5_CANARY_C=" + strings.Repeat("c", 12),
	}
	if _, err := buildEnvironmentBlock(env); err != nil {
		t.Fatalf("build the environment block: %v", err)
	}
	observedData := false
	for index, block := range discarded {
		if len(block) == 0 {
			continue
		}
		observedData = true
		for offset, unit := range block {
			if unit != 0 {
				t.Fatalf("discarded block %d still holds %q at offset %d: the pre-growth copy was not cleared",
					index, string(rune(unit)), offset)
			}
		}
	}
	if !observedData {
		t.Fatal("no discarded block carried data, so the observation proves nothing")
	}
}
