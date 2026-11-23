package adapters

// P85 F07 (design rev 3, section 3 / test row M6): the proved name chain must be the
// chain the KERNEL walks.
//
// filepath.Dir/Base/Clean fold "..", ".", doubled and trailing separators away; the
// kernel does not. Before the fix the folded form was proven while the raw form was
// executed, so "/a/../b" could be proven as "/b" and run as whatever "/a/.." resolved
// to. The rule is fail-closed: a spelling that is not already canonical is refused,
// and on a canonical spelling Clean/Dir/Base are identities.
//
// Mutation check: delete the `filepath.Clean(path) != path` refusal => every case
// below is accepted and the test goes red.
//
// Companion (already pinned, Linux): TestP85T3LinuxLinkedInterpreterChainIsRefused
// covers the link component of the same rule - a canonical spelling whose chain
// carries a symlink is refused by openat2, not followed.

import (
	"path/filepath"
	"testing"
)

func TestP85F07CanonicalSpellingRefusesFoldedNameChains(t *testing.T) {
	root := t.TempDir()
	separator := string(filepath.Separator)

	cases := map[string]string{
		"parent-component":   filepath.Join(root, "a") + separator + ".." + separator + "b",
		"dot-component":      filepath.Join(root, "a") + separator + "." + separator + "b",
		"doubled-separator":  root + separator + separator + "a",
		"trailing-separator": filepath.Join(root, "a") + separator,
		"relative":           "bin" + separator + "sh",
		"empty":              "",
		"nul-byte":           filepath.Join(root, "a") + "\x00" + "b",
	}
	for name, path := range cases {
		t.Run(name, func(t *testing.T) {
			if err := installSourceRequireCanonicalSpelling(path); err == nil {
				t.Fatalf("%q must be refused: the kernel walks the raw chain, not the cleaned one", path)
			}
		})
	}

	// Control: the canonical form of the SAME components is accepted, so the refusals
	// above are caused by the spelling and not by the path being unusable.
	canonical := filepath.Join(root, "a", "b")
	if err := installSourceRequireCanonicalSpelling(canonical); err != nil {
		t.Fatalf("a canonical absolute spelling must be accepted, got %v", err)
	}
	if cleaned := filepath.Clean(canonical); cleaned != canonical {
		t.Fatalf("the accepted spelling must be Clean's fixed point, got %q -> %q", canonical, cleaned)
	}
	if rebuilt := filepath.Join(filepath.Dir(canonical), filepath.Base(canonical)); rebuilt != canonical {
		t.Fatalf("Dir+Base must be an identity on the accepted spelling, got %q", rebuilt)
	}
}
