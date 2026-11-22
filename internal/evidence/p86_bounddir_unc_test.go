//go:build windows

package evidence

// P86 BOUNDDIR-UNC-EQUIVALENCE (N13-UNC) on the Windows half.
//
// The bound directory may itself be a UNC share root. A share root has TWO canonical
// spellings - with and without the trailing separator filepath.Clean keeps - and SMB
// resolves host and share names case-insensitively, so a reparse point that stores
// the other spelling (or an ASCII case variant of it) names the SAME object the
// binding owns. Before P86 that target was refused: the equality step compared the
// bytes of the two spellings, and the strict-descendant step then met a target equal
// to its own prefix.
//
// The equivalence branch is deliberately confined to the volume name of a share root:
// descendants, the drive form and the case-variant twin of a plain directory stay an
// Ordinal comparison (P84 W1 / W10). Each test below states the mutation that turns
// it red, and the branch itself comments the reason in bounddir_windows.go.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// TestP86N13UncShareRootEquivalence pins the positives of the equivalence branch.
//
// Mutation check: deleting the P86 (N13-UNC) branch from boundRelativeSuffix turns
// the subcases whose two spellings differ (trailing-base/bare-target,
// bare-base/trailing-target, ascii-case-folded) red; the identical-spellings and the
// strict-descendant subcases stay green in both worlds, because Clean makes them
// equality-step or strict-descendant cases.
func TestP86N13UncShareRootEquivalence(t *testing.T) {
	const (
		trailing = `\\server\share\`
		bare     = `\\server\share`
	)

	// Fact baseline: the two spellings are canonical, and they name one volume. These
	// are standard-library facts, not this package's behaviour, and the branch only
	// makes sense while they hold.
	if cleaned := filepath.Clean(trailing); cleaned != trailing {
		t.Fatalf("this test assumes filepath.Clean keeps the share root separator, got %q", cleaned)
	}
	if cleaned := filepath.Clean(bare); cleaned != bare {
		t.Fatalf("this test assumes filepath.Clean keeps the bare share root, got %q", cleaned)
	}
	if got, want := filepath.VolumeName(trailing), filepath.VolumeName(bare); got != want {
		t.Fatalf("both share root spellings must carry the same volume name, got %q and %q", got, want)
	}
	if got := uncShareRootName(trailing); got != bare {
		t.Fatalf("uncShareRootName(%q) = %q, want %q", trailing, got, bare)
	}
	if got := uncShareRootName(bare); got != bare {
		t.Fatalf("uncShareRootName(%q) = %q, want %q", bare, got, bare)
	}

	cases := []struct {
		name   string
		base   string
		target string
		want   []string
	}{
		{name: "trailing-base-bare-target", base: trailing, target: bare, want: []string{}},
		{name: "bare-base-trailing-target", base: bare, target: trailing, want: []string{}},
		{name: "ascii-case-folded-trailing-target", base: `\\SERVER\SHARE`, target: trailing, want: []string{}},
		{name: "identical-trailing-spellings", base: trailing, target: trailing, want: []string{}},
		{name: "strict-descendant-regression", base: trailing, target: `\\server\share\sub\f`, want: []string{"sub", "f"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			suffix, ok := boundRelativeSuffix(testCase.base, testCase.target)
			if !ok {
				t.Fatalf("boundRelativeSuffix(%q, %q) must admit the share root equivalence, got suffix %q", testCase.base, testCase.target, suffix)
			}
			if !slices.Equal(suffix, testCase.want) {
				t.Fatalf("boundRelativeSuffix(%q, %q) = %q, %v, want %q, true", testCase.base, testCase.target, suffix, ok, testCase.want)
			}
		})
	}
}

// TestP86N13UncShareRootRefusals pins the negatives: the branch admits exactly the
// volume-name differences above and nothing else.
//
// These subcases are distinguished by the refusal guards that run BEFORE the branch
// (nonCanonicalReparseTarget for the separator, slash and ".." shapes; the ASCII-only
// fold for the non-ASCII pair) and by the guards that run AFTER it (the Ordinal
// strict-descendant prefix for other hosts, other share names, case-variant
// descendants and out-of-bound siblings). They are green with and without the P86
// branch, which is the point: the branch must not reach any of them.
func TestP86N13UncShareRootRefusals(t *testing.T) {
	const base = `\\server\share\`
	cases := []struct {
		name   string
		target string
	}{
		{name: "different-host", target: `\\other\share`},
		{name: "different-host-trailing-separator", target: `\\other\share\`},
		{name: "different-share-name", target: `\\server\other`},
		{name: "case-variant-descendant-stays-ordinal", target: `\\SERVER\SHARE\sub`},
		{name: "repeated-separator", target: `\\server\share\\`},
		{name: "non-root-trailing-separator", target: `\\server\share\x\`},
		{name: "forward-slash", target: `\\server\share/sub`},
		{name: "parent-component", target: `\\server\share\..\share`},
		{name: "outside-sibling-share", target: `\\server\share2\sub`},
		{name: "host-only-volume", target: `\\server`},
		{name: "host-only-volume-trailing-separator", target: `\\server\`},
		{name: "trailing-blank-share-name", target: `\\server\share `},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			suffix, ok := boundRelativeSuffix(base, testCase.target)
			if ok {
				t.Fatalf("boundRelativeSuffix(%q, %q) admitted the target with suffix %q", base, testCase.target, suffix)
			}
			if suffix != nil {
				t.Fatalf("a refused target must report no suffix, got %q", suffix)
			}
		})
	}

	// The classifier must not claim a drive root, a descendant, or an extended/device
	// namespace path: those shapes reach the equivalence branch only if
	// uncShareRootName answers a non-empty volume name for them.
	for _, path := range []string{
		`D:\`,
		`D:\.`,
		`D:`,
		`\\server\share\sub`,
		`\\server`,
		`\\server\`,
		`\\?\UNC\server\share\`,
		`\\?\UNC\server\share`,
		`\\?\C:\`,
		`\\.\C:\`,
	} {
		if got := uncShareRootName(path); got != "" {
			t.Fatalf("uncShareRootName(%q) = %q, want \"\"", path, got)
		}
	}

	// The fold is ASCII-only. "é" and "É" differ in one byte and must not fold, so a
	// non-ASCII host or share can never be admitted by case alone.
	if !uncShareRootEqual(`\\server\share`, `\\SERVER\SHARE`) {
		t.Fatal("ASCII letters of the volume name must fold both ways")
	}
	if uncShareRootEqual(`\\sérver\share`, `\\SÉRVER\share`) {
		t.Fatal("non-ASCII bytes must be compared exactly, so é/É must not fold")
	}
	if uncShareRootEqual(`\\server\share`, `\\server\share2`) {
		t.Fatal("different lengths must not compare equal")
	}
}

// TestP86UncShareRootProductionEntry drives the equivalence branch through the
// production entry: a real share root is bound with OpenBoundDir, and a real
// mount-point reparse point stores the share root in the OTHER canonical spelling
// (with or without the trailing separator the handle-derived name happens to carry).
//
// Every fixture step that this host may legitimately refuse (the share command, the
// directory handle on a UNC root, the hand-built reparse buffer) skips with the
// reason spelled out rather than pretending the case ran. Under CI those failures are
// hard failures instead, so a green CI proves the production entry really ran.
func TestP86UncShareRootProductionEntry(t *testing.T) {
	local := t.TempDir()
	share := fmt.Sprintf("P86UNC%d", time.Now().UnixNano())
	if output, err := exec.Command("net", "share", share+"="+local).CombinedOutput(); err != nil {
		// A CI runner is elevated, so a failure there is a broken fixture rather than a
		// host boundary. Letting it skip would leave the production entry unproven while
		// the suite still reports ok; a red CI is the honest signal. Locally (no
		// elevation) the share cannot be created and the case skips with the reason.
		if os.Getenv("CI") != "" {
			t.Fatalf("net share cannot create the UNC share root fixture, so the production entry would stay unproven: %v (%s)", err, strings.TrimSpace(string(output)))
		}
		t.Skipf("net share cannot create the UNC share root fixture on this host: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	t.Cleanup(func() {
		output, err := exec.Command("net", "share", share, "/delete").CombinedOutput()
		if err != nil {
			t.Logf("net share %s /delete failed: %v (%s)", share, err, strings.TrimSpace(string(output)))
		}
	})

	if err := os.MkdirAll(filepath.Join(local, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	const payload = "p86-unc-share-root-bytes"
	if err := os.WriteFile(filepath.Join(local, "sub", "f"), []byte(payload), 0o600); err != nil {
		t.Fatal(err)
	}

	uncRoot := `\\localhost\` + share
	dir, err := OpenBoundDir(uncRoot)
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("this CI host must be able to open the UNC share root %s as a directory handle: %v", uncRoot, err)
		}
		t.Skipf("this host cannot open the UNC share root %s as a directory handle: %v", uncRoot, err)
	}
	defer func() { _ = dir.Close() }()

	base, err := dir.RootPath()
	if err != nil {
		t.Fatalf("derive the bound share root name from the handle: %v", err)
	}
	// Precondition: the handle-derived name is the share root in one of the two
	// canonical spellings, and the stored target below takes the OTHER one. N13 is
	// exactly that pair: the two spellings differ only by the trailing separator
	// Clean keeps for a share root, so the equality step misses and only the P86
	// branch can admit the target. Which side carries the separator is host
	// dependent (a CI runner derives the bare form), so the fixture mirrors the
	// observed spelling instead of assuming one.
	const sep = `\`
	var storedShare, printed string
	switch base {
	case uncRoot + sep:
		storedShare, printed = share, uncRoot
	case uncRoot:
		storedShare, printed = share+sep, uncRoot+sep
	default:
		if os.Getenv("CI") != "" {
			t.Fatalf("a CI host must derive the share root %q in one of its two canonical spellings, got %q", uncRoot, base)
		}
		t.Skipf("this host derives %q for the supplied share root %q; N13 needs one of the two canonical spellings", base, uncRoot)
	}
	if got := uncShareRootName(base); got != uncRoot {
		t.Fatalf("uncShareRootName(%q) = %q, want %q", base, got, uncRoot)
	}

	// Positive: the reparse point stores the share root in the OTHER canonical spelling.
	p84SetMountPointReparse(t, filepath.Join(local, "link"), `\??\UNC\localhost\`+storedShare, printed)
	file, err := OpenBoundDirFile(dir, `link\sub\f`, 1<<20)
	if err != nil {
		t.Fatalf("the share root spelled with the other canonical separator shape must be admitted and walked: %v", err)
	}
	data, truncated, err := file.ReadAll()
	_ = file.Close()
	if err != nil {
		t.Fatalf("read through the admitted share root: %v", err)
	}
	if truncated {
		t.Fatal("a small read must not be truncated")
	}
	if string(data) != payload {
		t.Fatalf("the admitted share root must read the bound bytes, got %q", data)
	}

	// Negative: the same shape with a different share name stays refused by the
	// boundary, so the admission above is the volume-name equivalence and not a
	// blanket "any UNC root" rule.
	p84SetMountPointReparse(t, filepath.Join(local, "linkout"), `\??\UNC\localhost\`+share+"-OTHER", uncRoot+"-OTHER")
	other, err := OpenBoundDirFile(dir, `linkout\sub\f`, 1<<20)
	if err == nil {
		_ = other.Close()
		t.Fatal("a stored target naming a different share must not be admitted")
	}
	if !strings.Contains(err.Error(), "outside the bound directory") {
		t.Fatalf("the refusal must name the boundary, got %v", err)
	}
}
