//go:build linux

package adapters

// P85 F07 (design rev 3, section 3 / test row M7): the shebang limit is the KERNEL's
// 256-byte window, counted INCLUDING the "#!" magic.
//
// binfmt_script copies at most BINPRM_BUF_SIZE (256) bytes and needs the terminating
// newline INSIDE that window. The previous rule counted the remainder after "#!" and
// did not check where the newline landed, so a spec of 254 bytes (newline at file
// offset 256) passed a rule that claimed to be equivalent to the kernel's.
//
// Mutation check: restore the "remainder length vs 256" rule (and search for the
// newline in the whole buffer) => the offset-256 case stops being refused by the
// window and this test goes red.

import (
	"strings"
	"testing"
)

const p85F07WindowRefusal = "no newline within the kernel's"

func TestP85F07ShebangWindowIncludesTheMagicBytes(t *testing.T) {
	// newline at file offset 255: spec is 253 bytes, exactly what the kernel keeps.
	insideWindow := "#!" + "/" + strings.Repeat("a", 252) + "\n" + "tail"
	if len(insideWindow) <= installSourceShebangWindow {
		t.Fatalf("the fixture must extend past the window, got %d bytes", len(insideWindow))
	}
	if _, err := installSourceShebangInterpreter([]byte(insideWindow), "/p85-f07-offset-255", 0); err != nil {
		if strings.Contains(err.Error(), p85F07WindowRefusal) {
			t.Fatalf("a newline at offset 255 is inside the kernel window and must pass the window rule, got %v", err)
		}
	}

	// newline at file offset 256: one byte past the window, so the kernel never
	// terminates the line. The old "remainder vs 256" rule admitted it.
	pastWindow := "#!" + "/" + strings.Repeat("a", 253) + "\n" + "tail"
	if _, err := installSourceShebangInterpreter([]byte(pastWindow), "/p85-f07-offset-256", 0); err == nil {
		t.Fatal("a newline past the kernel window must be refused")
	} else if !strings.Contains(err.Error(), p85F07WindowRefusal) {
		t.Fatalf("the refusal must come from the window rule, got %v", err)
	}
}
