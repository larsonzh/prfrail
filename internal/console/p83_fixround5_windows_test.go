//go:build windows

package console

// P83 FIX ROUND 5 counterexample (item T7b): the terminal secret reader copies its
// buffer into a larger one as the operator types, so the array it leaves behind
// must be zeroed before the slice is repointed.
//
// The seam is terminalSecretBufferObserver, which hands the reader's own buffers
// to the test; they are inspected AFTER the call, so only a real clear() can zero
// them.
//
// Mutation check: deleting `clear(secret)` above the "discarded" report leaves the
// typed prefix readable and turns this red.

import (
	"strings"
	"testing"
)

func TestP83T7bPreGrowthTerminalSecretBufferIsZeroed(t *testing.T) {
	var discarded [][]byte
	terminalSecretBufferObserver = func(stage string, buffer []byte) {
		if stage == "discarded" {
			discarded = append(discarded, buffer)
		}
	}
	defer func() { terminalSecretBufferObserver = nil }()

	// The initial capacity is 64 bytes, so a 200 byte secret crosses the growth
	// boundary several times.
	secret := strings.Repeat("s", 200)
	read, err := readTerminalSecretFrom(strings.NewReader(secret + "\n"))
	if err != nil {
		t.Fatalf("read the terminal secret: %v", err)
	}
	if string(read) != secret {
		t.Fatalf("the reader must return the typed secret, got %d bytes", len(read))
	}
	if len(discarded) == 0 {
		t.Fatal("the growth path must report the abandoned buffer")
	}
	for index, buffer := range discarded {
		if len(buffer) == 0 {
			t.Fatalf("discarded buffer %d is empty, so the observation proves nothing", index)
		}
		for offset, value := range buffer {
			if value != 0 {
				t.Fatalf("discarded buffer %d still holds %q at offset %d: the pre-growth copy was not cleared",
					index, string(rune(value)), offset)
			}
		}
	}
}
