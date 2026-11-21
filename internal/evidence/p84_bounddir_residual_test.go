//go:build windows

package evidence

// P84 BOUNDDIR-RESIDUAL counterexamples (items W1, W2, W3) on the Windows half.
//
// W1: a case-variant twin directory of the bound object is a SIBLING, so it must
// never be admitted as a descendant. The containment test is Ordinal precisely for
// this shape; case-folding would let a twin whose name differs only by case pass.
//
// W2: a junction in the MIDDLE of a walked name is followed only when its target is
// proven inside the bound object. An out-of-bound target is refused (both for an
// artifact open and for an enumeration), and an in-bound target stays admitted,
// which is the P52 control shape the whole walk has to preserve.
//
// W3: "" and "." name the bound object itself. The enumeration half must accept
// both spellings and return the same entries, and the artifact half must refuse
// both through one shared error.
//
// W4 (root is never opened when no artifact is expected) is already pinned by
// internal/gates/p83_fixround5_tester_test.go:37-63 and is registered only.
//
// Round 3 of the tester adds the counterexamples for the seven-defence blind review
// (P84 F1/F5/F6): W5b walks a suffix component that is itself a reparse point (the
// stored target's prefix is in bounds while the suffix leaves it), W5a adds no test
// of its own and is falsified by removing the Ordinal prefix test in
// boundRelativeSuffix, which turns W2a red, W6 renames the bound object and proves
// the root enumeration still reads the bound handle, and W7 pins the UTF-16 length
// guard of one component.
//
// Mutation contracts are stated on each test.

import (
	"encoding/binary"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"unicode/utf16"
	"unsafe"
)

// TestP84W1CaseVariantTwinIsNotADescendant is the spelling half of W1.
//
// Mutation check: replacing the Ordinal prefix comparison in isUnderDirectory with
// strings.EqualFold turns the first assertion red (the twin would read as a
// descendant) while the genuine-descendant control stays green.
func TestP84W1CaseVariantTwinIsNotADescendant(t *testing.T) {
	cases := []struct {
		name      string
		base      string
		candidate string
		want      bool
	}{
		{name: "case-variant-twin", base: `C:\Bound`, candidate: `C:\BOUND\x`, want: false},
		{name: "case-variant-twin-dir-itself", base: `C:\Bound`, candidate: `C:\BOUND`, want: false},
		{name: "genuine-descendant", base: `C:\Bound`, candidate: `C:\Bound\x`, want: true},
		{name: "genuine-descendant-trailing-separator", base: `C:\Bound\`, candidate: `C:\Bound\sub\x`, want: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := isUnderDirectory(testCase.base, testCase.candidate); got != testCase.want {
				t.Fatalf("isUnderDirectory(%q, %q) = %v, want %v", testCase.base, testCase.candidate, got, testCase.want)
			}
		})
	}
}

// TestP84W1CaseVariantTwinEscapeIsRefused is the end-to-end half of W1. The bound
// object and a twin that differs only by case are made DISTINCT directories by
// enabling per-directory case sensitivity, a junction inside the bound tree points
// at a strict descendant of the twin, and the artifact behind it must be refused.
// The target is a descendant on purpose: the twin directory itself would fail the
// length check under any comparison, so only a descendant makes the Ordinal prefix
// the deciding guard. Without it the twin's handle-derived name reads as a
// descendant of the bound name and the artifact would be admitted.
//
// The whole leg needs fsutil's case-sensitivity attribute (NTFS, may need
// elevation) and junction creation; either being unavailable skips the leg, which
// is registered as a known host boundary for the VM/CI run.
//
// Mutation check: EqualFold in isUnderDirectory turns the refusal assertion red.
func TestP84W1CaseVariantTwinEscapeIsRefused(t *testing.T) {
	root := t.TempDir()
	if output, err := exec.Command("fsutil", "file", "setCaseSensitiveInfo", root, "enable").CombinedOutput(); err != nil {
		t.Skipf("this host cannot enable per-directory case sensitivity: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	bound := filepath.Join(root, "Bound")
	twin := filepath.Join(root, "BOUND")
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(twin, 0o755); err != nil {
		t.Skipf("this host resolves the two spellings to one directory: %v", err)
	}
	// Prove the two names address different objects before relying on it.
	if err := os.WriteFile(filepath.Join(bound, "only-bound.txt"), []byte("bound"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(twin, "only-bound.txt")); err == nil {
		t.Skip("this host resolves the case-variant twin to the bound directory")
	}
	twinSub := filepath.Join(twin, "sub")
	if err := os.Mkdir(twinSub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(twinSub, "x"), []byte("TWIN-OUTSIDE"), 0o600); err != nil {
		t.Fatal(err)
	}
	p83Fix8Junction(t, filepath.Join(bound, "j"), twinSub)

	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := OpenBoundDirFile(dir, `j\x`, 1<<20)
	if err == nil {
		_ = file.Close()
		t.Fatal("a case-variant twin directory outside the bound object must not be reachable")
	}
	if !strings.Contains(err.Error(), "outside the bound directory") {
		t.Fatalf("the refusal must name the boundary, got %v", err)
	}
	// Control: an artifact reached through a plain in-bound name still opens, so the
	// refusal above is the containment proof and not a broken fixture.
	control, err := OpenBoundDirFile(dir, "only-bound.txt", 1<<20)
	if err != nil {
		t.Fatalf("a plain in-bound artifact must still open, got %v", err)
	}
	_ = control.Close()
}

// TestP84W2IntermediateJunctionOutsideIsRefused is the out-of-bound half of W2: an
// intermediate junction that leaves the bound tree refuses the artifact open, the
// enumeration, and a walk that leaves and returns.
//
// The containment proof exists twice on purpose (the walk proves every followed
// target, and confirmOpenedArtifact/belt proves the final name), so a plain escape
// is still refused when either one is removed alone. The excursion leg is the shape
// that separates them: the final name is back inside the bound object, so only the
// intermediate proof can refuse it.
//
// Mutation check: removing the proven-inside step in followIntermediateReparse turns
// the excursion leg red (the walk is admitted) while the plain-escape leg stays
// green; refusing every intermediate reparse point instead turns
// TestP84W2IntermediateJunctionInsideIsAdmitted red.
func TestP84W2IntermediateJunctionOutsideIsRefused(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("SECRET-OUTSIDE-BOUNDARY"), 0o600); err != nil {
		t.Fatal(err)
	}

	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	t.Run("artifact-open", func(t *testing.T) {
		p83Fix8Junction(t, filepath.Join(bound, "j"), outside)
		file, openErr := OpenBoundDirFile(dir, `j\secret`, 1<<20)
		if openErr == nil {
			_ = file.Close()
			t.Fatal("an intermediate junction leaving the bound directory must be refused")
		}
		if !strings.Contains(openErr.Error(), "outside the bound directory") {
			t.Fatalf("the refusal must name the boundary, got %v", openErr)
		}
	})

	t.Run("enumeration", func(t *testing.T) {
		p83Fix8Junction(t, filepath.Join(bound, "k"), outside)
		if names, namesErr := BoundDirNames(dir, "k"); namesErr == nil {
			t.Fatalf("an enumerated directory that is a reparse point must be refused, got %#v", names)
		}
	})

	// The excursion leg is what makes the MIDDLE proof falsifiable on its own: the
	// walk leaves the bound object and comes back before the leaf, so the final
	// containment comparison sees an in-bound name and only the intermediate proof
	// can refuse the excursion.
	t.Run("excursion-out-and-back", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(bound, "visible.txt"), []byte("bound"), 0o600); err != nil {
			t.Fatal(err)
		}
		p83Fix8Junction(t, filepath.Join(bound, "jback"), outside)
		p83Fix8Junction(t, filepath.Join(outside, "back"), bound)
		file, openErr := OpenBoundDirFile(dir, `jback\back\visible.txt`, 1<<20)
		if openErr == nil {
			_ = file.Close()
			t.Fatal("a walk that leaves the bound object and returns must be refused at the exit")
		}
		if !strings.Contains(openErr.Error(), "outside the bound directory") {
			t.Fatalf("the refusal must name the boundary, got %v", openErr)
		}
	})
}

// TestP84W2IntermediateJunctionInsideIsAdmitted is the in-bound half of W2, the
// P52 control shape: an intermediate junction whose target is a strict descendant
// of the bound object is followed and the artifact behind it reads back unchanged.
//
// Mutation check: refusing every intermediate reparse point turns the admission
// assertion red.
func TestP84W2IntermediateJunctionInsideIsAdmitted(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	inside := filepath.Join(bound, "sub")
	if err := os.MkdirAll(inside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inside, "f"), []byte("inside"), 0o600); err != nil {
		t.Fatal(err)
	}

	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()
	p83Fix8Junction(t, filepath.Join(bound, "j"), inside)

	file, err := OpenBoundDirFile(dir, `j\f`, 1<<20)
	if err != nil {
		t.Fatalf("an intermediate junction staying inside the bound object must be admitted: %v", err)
	}
	defer func() { _ = file.Close() }()
	data, truncated, err := file.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll through the in-bound junction: %v", err)
	}
	if truncated {
		t.Fatal("a small read must not be truncated")
	}
	if string(data) != "inside" {
		t.Fatalf("the in-bound junction must read the bound bytes, got %q", data)
	}
}

// TestP84W3RootSpellingsAgree pins the shared root classifier: "" and "." are the
// same name for the enumeration half and the same refusal for the artifact half.
//
// Mutation check: removing the "." branch of boundRootName makes BoundDirNames(dir,
// ".") fall through to cleanBoundRelativePath, which refuses it, and turns this red.
func TestP84W3RootSpellingsAgree(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.MkdirAll(filepath.Join(bound, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bound, "a.txt"), []byte("a"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	empty, err := BoundDirNames(dir, "")
	if err != nil {
		t.Fatalf(`BoundDirNames(dir, "") must enumerate the bound object, got %v`, err)
	}
	dot, err := BoundDirNames(dir, ".")
	if err != nil {
		t.Fatalf(`BoundDirNames(dir, ".") must enumerate the bound object, got %v`, err)
	}
	if !slices.Equal(empty, dot) {
		t.Fatalf("both root spellings must enumerate the same entries, got %#v and %#v", empty, dot)
	}
	if !slices.Contains(empty, "a.txt") || !slices.Contains(empty, "sub") {
		t.Fatalf("the root enumeration must describe the bound object, got %#v", empty)
	}

	_, emptyErr := OpenBoundDirFile(dir, "", 1<<20)
	_, dotErr := OpenBoundDirFile(dir, ".", 1<<20)
	if emptyErr == nil || dotErr == nil {
		t.Fatalf(`an artifact open naming the bound object must be refused for both spellings, got %v and %v`, emptyErr, dotErr)
	}
	for _, refusal := range []error{emptyErr, dotErr} {
		if !strings.Contains(refusal.Error(), "names the bound directory itself") {
			t.Fatalf("both root spellings must use the shared refusal, got %v", refusal)
		}
	}
}

// TestP84W5bChainedJunctionSuffixIsRefused pins the second half of the F1 walk:
// the relative suffix a middle reparse point contributes is itself walked from the
// bound handle, so a suffix component that is a reparse point is classified again
// instead of being trusted because the stored target's PREFIX was in bounds.
//
// The shape is two junctions deep: D\j stores D\sub (prefix inside D, so the prefix
// test passes and the suffix is "sub"), and D\sub is itself a junction leaving D.
// The artifact name D\j\f therefore passes the prefix admission of the first hop and
// can only be refused when the walk re-classifies the suffix component. The control
// leg uses the same two-hop shape with an in-bound second target: it must be
// admitted, so the refusal is the second junction and not a blanket "chained
// junction" rule.
//
// Mutation check: dropping the reparse branch of openDescendantRelative for the
// components that came from a stored-target suffix removes the re-classification,
// and the first assertion below turns red (the walk either reaches the outside
// directory or fails with a status the boundary message does not cover).
func TestP84W5bChainedJunctionSuffixIsRefused(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "D")
	outside := filepath.Join(root, "outside")
	inside := filepath.Join(bound, "real")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(inside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "f"), []byte("SECRET-OUTSIDE"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inside, "f"), []byte("INSIDE"), 0o600); err != nil {
		t.Fatal(err)
	}

	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// Control first: D\k -> D\hop and D\hop -> D\real. Both stored targets are inside
	// the bound object, so the two-hop suffix walk must be admitted.
	p83Fix8Junction(t, filepath.Join(bound, "hop"), inside)
	p83Fix8Junction(t, filepath.Join(bound, "k"), filepath.Join(bound, "hop"))
	admitted, err := OpenBoundDirFile(dir, `k\f`, 1<<20)
	if err != nil {
		t.Fatalf("an in-bound chained junction must be admitted, got %v", err)
	}
	data, _, err := admitted.ReadAll()
	_ = admitted.Close()
	if err != nil || string(data) != "INSIDE" {
		t.Fatalf("the in-bound chain must read the bound bytes, got %q (err %v)", data, err)
	}

	// The escape: D\sub is a junction leaving the bound object and D\j stores it, so
	// the first hop's stored-target PREFIX is in bounds while the suffix is a reparse.
	p83Fix8Junction(t, filepath.Join(bound, "sub"), outside)
	p83Fix8Junction(t, filepath.Join(bound, "j"), filepath.Join(bound, "sub"))
	file, err := OpenBoundDirFile(dir, `j\f`, 1<<20)
	if err == nil {
		_ = file.Close()
		t.Fatal("a junction reached through another junction's stored-target suffix must not be followed out of the bound object")
	}
	if !strings.Contains(err.Error(), "outside the bound directory") {
		t.Fatalf("the refusal must name the boundary, got %v", err)
	}
}

// TestP84W6RootEnumerationReadsTheBoundHandle pins F5: the root enumeration
// (rel "" or ".") reads the entries through the handle the binding already owns and
// resolves no name at all.
//
// The bound object is renamed away and a NEW directory is created at the old name
// carrying an entry the bound object never had. The enumeration must still describe
// the bound object (the original entry present, the replacement entry absent). The
// rename can fail when an indexer or a scanner holds a component of the tree without
// delete sharing; that is a host boundary and is registered as a skip rather than
// masked, because the assertion can only be made about a rename that happened.
//
// Mutation check: replacing the root branch with any name-based reopen - the pre-F5
// `currentName() + CreateFile` shape, or a reopen of the bound path - turns this red:
// the reopened handle either enumerates the replacement or is refused by the identity
// proof, and both contradict the assertions below. (The exact pre-F5 shape resolves
// the CURRENT name, which follows the rename, so it is red only through the identity
// proof; the measurement is recorded in the round-3 report.)
func TestP84W6RootEnumerationReadsTheBoundHandle(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	moved := filepath.Join(root, "bound-moved")
	if err := os.MkdirAll(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bound, "original.txt"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// The bound handle is share-delete (OpenBoundDir), so the directory stays
	// renamable while it is bound; a host that still refuses is a registered boundary.
	if err := os.Rename(bound, moved); err != nil {
		t.Skipf("this host does not allow renaming a directory with an open bound handle: %v", err)
	}
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatalf("create the replacement at the bound name: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bound, "replacement.txt"), []byte("replacement"), 0o600); err != nil {
		t.Fatal(err)
	}

	names, err := BoundDirNames(dir, "")
	if err != nil {
		t.Fatalf("the root enumeration must read the bound handle, got %v", err)
	}
	if !slices.Contains(names, "original.txt") {
		t.Fatalf("the root enumeration must describe the bound object, got %#v", names)
	}
	if slices.Contains(names, "replacement.txt") {
		t.Fatalf("the root enumeration must not describe the replacement at the old name, got %#v", names)
	}
	dot, err := BoundDirNames(dir, ".")
	if err != nil || !slices.Equal(names, dot) {
		t.Fatalf(`BoundDirNames(dir, ".") must read the same handle, got %#v and %v`, dot, err)
	}
}

// TestP84W7OversizeComponentIsRefusedBeforeOpen pins the F6 length guard twice.
//
// Unit half: utf16ComponentLength reports the byte length of a component that fits a
// UNICODE_STRING and refuses one whose UTF-16 form would overflow the 16-bit Length
// field. The boundary is measured in bytes, so the last admitted component is
// (MaxUint16-2)/2 ASCII characters and the first refused one is a single character
// past it.
//
// End-to-end half: an artifact name carrying such a component is refused by the walk
// BEFORE any open, with the message naming the UNICODE_STRING limit. That assertion is
// what makes deleting the guard call in openRelativeComponent falsifiable: without it
// the kernel answers some status for a name no file could carry, and the message
// changes (a filesystem could even resolve a truncated name, which is the defect).
//
// Mutation check: removing the utf16ComponentLength call from openRelativeComponent
// turns the end-to-end assertion red. The comparison operators themselves are NOT
// where this guard is falsifiable: the measured length is always even (two bytes per
// UTF-16 unit) and the bound MaxUint16-2 is odd, so `>` and `>=` accept exactly the
// same set and flipping one leaves every assertion here green (final review N, round
// 4 measurement). The effective mutation is dropping the two-byte reserve - comparing
// against MaxUint16 instead of MaxUint16-2 - which admits the first over-limit
// component and turns the unit assertion red.
func TestP84W7OversizeComponentIsRefusedBeforeOpen(t *testing.T) {
	const componentLimit = (math.MaxUint16 - 2) / 2
	if length, ok := utf16ComponentLength("plain.txt"); !ok || length != uint16(len("plain.txt")*2) {
		t.Fatalf("a normal component must report its UTF-16 byte length, got (%d, %v)", length, ok)
	}
	atLimit := strings.Repeat("x", componentLimit)
	if length, ok := utf16ComponentLength(atLimit); !ok || int(length) != componentLimit*2 {
		t.Fatalf("the largest representable component must be admitted, got (%d, %v)", length, ok)
	}
	overLimit := strings.Repeat("x", componentLimit+1)
	if length, ok := utf16ComponentLength(overLimit); ok {
		t.Fatalf("a component of %d UTF-16 bytes must be refused, got (%d, %v)", len(overLimit)*2, length, ok)
	}

	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.MkdirAll(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := OpenBoundDirFile(dir, overLimit, 1<<20)
	if err == nil {
		_ = file.Close()
		t.Fatalf("a component of %d UTF-16 bytes must be refused before any open", len(overLimit)*2)
	}
	if !strings.Contains(err.Error(), "UNICODE_STRING name is limited to") {
		t.Fatalf("the refusal must be the length guard, got %v", err)
	}
}

// p84FSCTLSetReparsePoint is FSCTL_SET_REPARSE_POINT: FILE_DEVICE_FILE_SYSTEM with
// function 41 and METHOD_BUFFERED (FILE_ANY_ACCESS on this host: the WRITE_ACCESS and
// READ|WRITE variants answer "incorrect function" / "access denied"). It lets C12
// write a reparse point whose two stored names disagree, which no shell command can
// produce.
const p84FSCTLSetReparsePoint = 0x000900A4

// p84MountPointBuffer builds one IO_REPARSE_TAG_MOUNT_POINT REPARSE_DATA_BUFFER in
// which the substitute name (the name the object manager resolves) and the print name
// (the display string) are supplied independently.
//
// Both names are NUL-terminated in the path buffer while their declared lengths
// exclude the terminator, and the print name's offset skips the substitute's
// terminator: that is the layout NTFS writes for a junction, and a buffer without the
// terminators is refused with ERROR_INVALID_REPARSE_DATA (measured on this host).
func p84MountPointBuffer(substitute, printName string) []byte {
	substituteUnits := utf16.Encode([]rune(substitute))
	printUnits := utf16.Encode([]rune(printName))
	body := make([]byte, 8)
	binary.LittleEndian.PutUint16(body[0:2], 0)
	binary.LittleEndian.PutUint16(body[2:4], uint16(len(substituteUnits)*2))
	binary.LittleEndian.PutUint16(body[4:6], uint16(len(substituteUnits)*2+2))
	binary.LittleEndian.PutUint16(body[6:8], uint16(len(printUnits)*2))
	for _, unit := range substituteUnits {
		body = binary.LittleEndian.AppendUint16(body, unit)
	}
	body = binary.LittleEndian.AppendUint16(body, 0)
	for _, unit := range printUnits {
		body = binary.LittleEndian.AppendUint16(body, unit)
	}
	body = binary.LittleEndian.AppendUint16(body, 0)
	buffer := make([]byte, 8, 8+len(body))
	binary.LittleEndian.PutUint32(buffer[0:4], ioReparseTagMountPoint)
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(len(body)))
	return append(buffer, body...)
}

// p84SetMountPointReparse turns an empty directory into a mount-point reparse point
// whose stored names are exactly the two arguments.
//
// A host that refuses a hand-built buffer (NTFS validating more than the structure,
// or a policy that bans reparse points in the temp volume) skips the case: the parse
// authority can only be measured where such a point can be written, and pretending
// otherwise would assert on a fixture that never existed.
func p84SetMountPointReparse(t *testing.T, link, substitute, printName string) {
	t.Helper()
	if err := os.Mkdir(link, 0o755); err != nil {
		t.Fatalf("create the reparse host directory: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(link) })
	path, err := syscall.UTF16PtrFromString(link)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := syscall.CreateFile(
		path,
		syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		t.Fatalf("open the reparse host directory: %v", err)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	buffer := p84MountPointBuffer(substitute, printName)
	var returned uint32
	result, _, callErr := procDeviceIoControl.Call(
		uintptr(handle),
		uintptr(p84FSCTLSetReparsePoint),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
		0,
		uintptr(unsafe.Pointer(&returned)),
		0,
	)
	if result == 0 {
		t.Skipf("this host refuses a hand-built mount point reparse point, so the two stored names cannot be made to disagree: %v", callErr)
	}
}

// p84StoredTarget opens the reparse point itself (never its target) and reads the
// target reparseStoredTarget reports for it.
func p84StoredTarget(t *testing.T, link string) string {
	t.Helper()
	path, err := syscall.UTF16PtrFromString(link)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := syscall.CreateFile(
		path,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		t.Fatalf("open the reparse point itself: %v", err)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	target, err := reparseStoredTarget(handle)
	if err != nil {
		t.Fatalf("read the stored target: %v", err)
	}
	return target
}

// TestP84C12SubstituteNameIsTheAuthority pins N3: the stored target that decides the
// boundary is the SUBSTITUTE name, the one the object manager resolves. The print
// name is a display string a link may carry a different value for, so deciding on it
// - or falling back to it - lets a reparse point whose two names disagree be judged by
// the friendlier one.
//
// Both directions are built for real: a point whose print name is in bounds while its
// substitute name is out of bounds must report the out-of-bounds target (and the
// containment test must refuse it), and the mirrored point must report the in-bounds
// target (and the containment test must admit it). Either direction alone could be
// satisfied by a constant, so the pair is what makes the authority the deciding
// factor.
//
// Mutation check: replacing the substitute name with the print name in
// reparseStoredTarget (the selector and the normalized value) makes each case report
// the OTHER name and turns both target assertions red.
func TestP84C12SubstituteNameIsTheAuthority(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	inside := filepath.Join(bound, "inside")
	outside := filepath.Join(root, "outside")
	for _, directory := range []string{inside, outside} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	cases := []struct {
		name       string
		substitute string
		printName  string
		inside     bool
	}{
		{
			name:       "print-inside-substitute-outside",
			substitute: `\??\` + outside,
			printName:  inside,
			inside:     false,
		},
		{
			name:       "print-outside-substitute-inside",
			substitute: `\??\` + inside,
			printName:  outside,
			inside:     true,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			link := filepath.Join(root, "link-"+testCase.name)
			p84SetMountPointReparse(t, link, testCase.substitute, testCase.printName)
			target := p84StoredTarget(t, link)
			want := filepath.Clean(strings.TrimPrefix(testCase.substitute, `\??\`))
			if target != want {
				t.Fatalf("the stored target must be the substitute name %q, got %q (print name %q)", want, target, testCase.printName)
			}
			if _, ok := boundRelativeSuffix(filepath.Clean(bound), target); ok != testCase.inside {
				t.Fatalf("the containment decision on %q must report inside=%v (print name %q)", target, testCase.inside, testCase.printName)
			}
		})
	}
}

// TestP84W8TrailingBlankingCodePointIsNotStrippedFromTheStoredTarget pins N7.
//
// normalizeReparseTarget once ran strings.TrimSpace over the stored target. U+00A0
// is unicode space AND a legal NTFS file name byte, so a junction whose stored
// target was "<bound>\u00A0" - a SIBLING of the bound object - normalised onto the
// bound name itself and was admitted as if it addressed the bound object. The
// normalization now strips only the exact object-manager prefixes; every remaining
// byte, blanks included, is left for boundRelativeSuffix to judge by Ordinal
// equality / strict-descendant bytes.
//
// Mutation check: restoring `target = strings.TrimSpace(target)` before the return
// of normalizeReparseTarget turns the byte-preservation, the sibling refusal and
// the end-to-end refusal below red (the trimmed target equals the bound name, so
// the containment test admits it and the walk resolves the junction as the bound
// object itself).
func TestP84W8TrailingBlankingCodePointIsNotStrippedFromTheStoredTarget(t *testing.T) {
	const bound = `D:\Bound`
	const sibling = bound + "\u00a0"

	cases := []struct {
		name   string
		stored string
		want   string
	}{
		{name: "nt-prefix", stored: `\??\` + bound, want: bound},
		{name: "extended-prefix", stored: `\\?\` + bound, want: bound},
		{name: "nt-prefix-blank-tail", stored: `\??\` + sibling, want: sibling},
		{name: "extended-prefix-blank-tail", stored: `\\?\` + sibling, want: sibling},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got := normalizeReparseTarget(testCase.stored)
			if got != testCase.want {
				t.Fatalf("normalizeReparseTarget(%q) = %q, want %q", testCase.stored, got, testCase.want)
			}
			// The bound name must be reached by prefix stripping alone: the trailing
			// U+00A0 of the sibling cases survives byte for byte.
			if strings.HasSuffix(testCase.stored, "\u00a0") && !strings.HasSuffix(got, "\u00a0") {
				t.Fatalf("the trailing U+00A0 must survive normalization byte for byte, got %q", got)
			}
		})
	}

	// The production chain, composed: the stored substitute name is normalized, then
	// judged by boundRelativeSuffix against the bound object's handle-derived name.
	if _, ok := boundRelativeSuffix(bound, normalizeReparseTarget(`\??\`+sibling)); ok {
		t.Error("a target that differs only by a trailing U+00A0 is a SIBLING, not the bound object or a descendant of it")
	}
	// Controls: the exact names are still judged exactly as before.
	if suffix, ok := boundRelativeSuffix(bound, normalizeReparseTarget(`\??\`+bound)); !ok || len(suffix) != 0 {
		t.Fatalf("the bound object itself must stay admitted with an empty suffix, got suffix=%q ok=%v", suffix, ok)
	}
	if suffix, ok := boundRelativeSuffix(bound, normalizeReparseTarget(`\??\`+bound+`\sub\`+"\u00a0")); !ok || len(suffix) != 2 || suffix[0] != "sub" || suffix[1] != "\u00a0" {
		t.Errorf(`a genuine descendant must keep a U+00A0 component byte for byte in its suffix, got suffix=%q ok=%v`, suffix, ok)
	}

	// End to end: a junction inside the bound object whose stored target is the
	// blank-suffixed SIBLING, with the same artifact name on both sides but different
	// bytes, so an admitted escape is visible in the content as well as in the
	// refusal. The sibling itself is the target (not a descendant), which is the
	// shape only the byte-exact comparison can refuse.
	root := t.TempDir()
	boundDir := filepath.Join(root, "Bound")
	siblingDir := boundDir + "\u00a0"
	if err := os.MkdirAll(filepath.Join(boundDir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(siblingDir, 0o755); err != nil {
		t.Skipf("this host cannot create a directory whose name ends in U+00A0: %v", err)
	}
	if err := os.WriteFile(filepath.Join(boundDir, "f"), []byte("BOUND"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(siblingDir, "f"), []byte("SIBLING"), 0o600); err != nil {
		t.Fatal(err)
	}
	p84SetMountPointReparse(t, filepath.Join(boundDir, "j"), `\??\`+siblingDir, siblingDir)

	dir, err := OpenBoundDir(boundDir)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := OpenBoundDirFile(dir, `j\f`, 1<<20)
	if err == nil {
		data, _, readErr := file.ReadAll()
		_ = file.Close()
		t.Fatalf("a junction to the blank-suffixed sibling must be refused, got content %q (read err %v)", data, readErr)
	}
	if !strings.Contains(err.Error(), "outside the bound directory") {
		t.Fatalf("the refusal must name the boundary, got %v", err)
	}
	// Control: the same artifact name opens through the plain in-bound path, so the
	// refusal above is the blank byte and not a broken fixture.
	control, err := OpenBoundDirFile(dir, "f", 1<<20)
	if err != nil {
		t.Fatalf("a plain in-bound artifact must still open, got %v", err)
	}
	defer func() { _ = control.Close() }()
	data, _, err := control.ReadAll()
	if err != nil || string(data) != "BOUND" {
		t.Fatalf("the control artifact must read the bound bytes, got %q err=%v", data, err)
	}
}

// TestP84W9NonCanonicalStoredTargetIsRefused pins N10 (fix round 10): a stored
// target is judged by its OWN bytes, so a spelling that only names an admitted
// object after filepath.Clean folds it is refused instead of being rewritten onto
// that object.
//
// The pre-round-10 shape cleaned the target before comparing it, so
// "<bound>\x\..\outside" collapsed to "<bound>\outside" and read as a genuine
// descendant although the bytes the object manager resolves name a different
// object. The precondition is now explicit and fail-closed: after only the exact
// object-manager prefixes are stripped, the target must already BE the form
// filepath.Clean produces, and the containment decision then compares the stored
// bytes themselves (Ordinal equality / strict descendant).
//
// A junction cannot be made to store a folded target - mklink and the shell
// normalise every spelling they write - so the counterexamples are posed where the
// decision is actually taken: boundRelativeSuffix (the containment decision) and
// reparseTargetRelativeSuffix (the production entry that walks a middle reparse
// point), the latter end to end against a hand-written reparse buffer.
//
// Mutation check: deleting the nonCanonicalReparseTarget call from
// boundRelativeSuffix (restoring Clean-before-compare) admits every refused target
// below - "<bound>\x\..\sub" folds onto the genuine descendant "<bound>\sub", and
// "<bound>" with a repeated or trailing separator is returned as the bound object
// itself - so the parsing-layer assertions go red. Deleting the reason check from
// reparseTargetRelativeSuffix turns the folded end-to-end leg red: the refusal then
// keeps only the boundary wording and loses the lossless-interpretation reason.
func TestP84W9NonCanonicalStoredTargetIsRefused(t *testing.T) {
	const bound = `D:\Bound`

	refused := []struct {
		name   string
		target string
		reason string
	}{
		{name: "dot-dot-fold", target: bound + `\x\..\outside`, reason: `a ".." component`},
		{name: "dot-component", target: bound + `\.\sub`, reason: `a "." component`},
		{name: "repeated-separator", target: `D:\\Bound`, reason: "a repeated separator"},
		{name: "forward-slash", target: bound + `/sub`, reason: "a forward slash"},
		{name: "trailing-separator", target: bound + `\sub\`, reason: "a trailing separator"},
	}
	for _, testCase := range refused {
		t.Run(testCase.name, func(t *testing.T) {
			// The reason is what the operator sees, so it has to name the shape.
			if reason := nonCanonicalReparseTarget(testCase.target); !strings.Contains(reason, testCase.reason) {
				t.Fatalf("the stored target %q must be refused as non-canonical for %q, got reason %q", testCase.target, testCase.reason, reason)
			}
			// The object-manager prefix the production path strips first must not launder
			// a non-canonical spelling into a canonical one.
			if reason := nonCanonicalReparseTarget(normalizeReparseTarget(`\??\` + testCase.target)); !strings.Contains(reason, testCase.reason) {
				t.Fatalf("the \\??\\ prefix must not launder the spelling %q, got reason %q", testCase.target, reason)
			}
			if suffix, ok := boundRelativeSuffix(bound, testCase.target); ok || suffix != nil {
				t.Fatalf("the non-canonical target %q names a different object once folded, so it must be refused, got suffix=%q ok=%v", testCase.target, suffix, ok)
			}
			if suffix, ok := boundRelativeSuffix(bound, normalizeReparseTarget(`\??\`+testCase.target)); ok || suffix != nil {
				t.Fatalf("the non-canonical target %q must stay refused through the production normalization, got suffix=%q ok=%v", testCase.target, suffix, ok)
			}
		})
	}

	// Controls: the canonical spellings P52 and W8 rely on are still admitted exactly
	// as before, so the refusal above is the spelling and not a blanket refusal.
	admitted := []struct {
		name   string
		target string
		want   []string
	}{
		{name: "bound-object-itself", target: bound, want: nil},
		{name: "genuine-descendant", target: bound + `\sub`, want: []string{"sub"}},
		{name: "genuine-descendant-deep", target: bound + `\sub\file.txt`, want: []string{"sub", "file.txt"}},
	}
	for _, testCase := range admitted {
		t.Run(testCase.name, func(t *testing.T) {
			if reason := nonCanonicalReparseTarget(testCase.target); reason != "" {
				t.Fatalf("the canonical target %q must be interpreted losslessly, got reason %q", testCase.target, reason)
			}
			suffix, ok := boundRelativeSuffix(bound, testCase.target)
			if !ok || !slices.Equal(suffix, testCase.want) {
				t.Fatalf("boundRelativeSuffix(%q, %q) = %q, %v, want %q, true", bound, testCase.target, suffix, ok, testCase.want)
			}
		})
	}

	// End to end through the production entry, with a reparse buffer written by hand
	// because no shell command stores a folded target. The canonical control keeps the
	// same entry from passing merely because it refuses everything.
	t.Run("production-chain-canonical", func(t *testing.T) {
		root := t.TempDir()
		boundDir := filepath.Join(root, "Bound")
		target := filepath.Join(boundDir, "sub")
		if err := os.MkdirAll(target, 0o755); err != nil {
			t.Fatal(err)
		}
		p84SetMountPointReparse(t, filepath.Join(boundDir, "mid"), `\??\`+target, target)
		dir, err := OpenBoundDir(boundDir)
		if err != nil {
			t.Fatalf("bind the directory: %v", err)
		}
		defer func() { _ = dir.Close() }()
		suffix, err := reparseTargetRelativeSuffix(dir, dir.handle, "mid", "mid")
		if err != nil {
			t.Fatalf("a canonical in-bound stored target must stay admitted, got %v", err)
		}
		if len(suffix) != 1 || suffix[0] != "sub" {
			t.Fatalf("the walked suffix of %q must be [sub], got %q", target, suffix)
		}
	})

	t.Run("production-chain-folded", func(t *testing.T) {
		root := t.TempDir()
		boundDir := filepath.Join(root, "Bound")
		if err := os.MkdirAll(filepath.Join(boundDir, "sub"), 0o755); err != nil {
			t.Fatal(err)
		}
		// Built by concatenation on purpose: filepath.Join would clean the ".." away and
		// the fixture would stop being the shape under test. "x" does not exist, so an
		// admitted fold would be visible as a different object.
		folded := boundDir + `\x\..\sub`
		p84SetMountPointReparse(t, filepath.Join(boundDir, "mid"), `\??\`+folded, folded)
		dir, err := OpenBoundDir(boundDir)
		if err != nil {
			t.Fatalf("bind the directory: %v", err)
		}
		defer func() { _ = dir.Close() }()
		suffix, err := reparseTargetRelativeSuffix(dir, dir.handle, "mid", "mid")
		if err == nil {
			t.Fatalf("a stored target that only reaches a descendant after folding must be refused, got suffix %q", suffix)
		}
		if !strings.Contains(err.Error(), "cannot be interpreted losslessly") {
			t.Fatalf("the refusal of %q must carry the lossless-interpretation reason, got %v", folded, err)
		}
	})
}

// Round 8 (fix round 11) adds the two leftovers of review r3: N13, the volume root
// that a trimmed equality comparison refused although it is deterministically in
// bounds, and N14, a stored name whose declared range carries a NUL unit before its
// last position.

// TestP84W10VolumeRootKeepsItsSeparatorForEquality pins N13.
//
// filepath.Clean keeps the separator of a volume root ("D:\" stays "D:\") while it
// drops the trailing separator of every other path. Judging equality on a TRIMMED
// base therefore compared "D:" against the target "D:\", the bytes differed, and the
// strict-descendant branch then met a target that IS the prefix - a deterministic
// in-bounds refusal of the root itself. Equality now compares the canonical form
// exactly as Clean returns it, and only the strict-descendant prefix strips the
// separator before appending one, so the root and its descendants are admitted while
// the case-variant twin stays refused.
//
// Mutation check: restoring the round-11-before shape - equality against
// strings.TrimRight(filepath.Clean(base), `\/`) with the descendant prefix built from
// that trimmed value - turns the root-equality positive red. Building the prefix from
// the untrimmed canonical form as well (the literal pre-round-11 arrangement, where
// "D:\\" was the prefix) turns the root-descendant positive red too. The
// case-variant refusal is unchanged by either, which is why the positives are what
// the mutation is measured on.
func TestP84W10VolumeRootKeepsItsSeparatorForEquality(t *testing.T) {
	// Everything below rests on Clean keeping the root separator alone.
	if cleaned := filepath.Clean(`D:\`); cleaned != `D:\` {
		t.Fatalf("this test assumes filepath.Clean keeps the root separator, got %q", cleaned)
	}
	if cleaned := filepath.Clean(`D:\sub`); cleaned != `D:\sub` {
		t.Fatalf("this test assumes filepath.Clean drops a non-root separator, got %q", cleaned)
	}

	cases := []struct {
		name   string
		base   string
		target string
		want   []string
		inside bool
	}{
		{name: "root-itself", base: `D:\`, target: `D:\`, want: []string{}, inside: true},
		{name: "root-descendant", base: `D:\`, target: `D:\sub`, want: []string{"sub"}, inside: true},
		{name: "root-descendant-deep", base: `D:\`, target: `D:\sub\file.txt`, want: []string{"sub", "file.txt"}, inside: true},
		{name: "descendant-file-under-root-child", base: `D:\sub`, target: `D:\sub\file.txt`, want: []string{"file.txt"}, inside: true},
		{name: "case-variant-twin-stays-refused", base: `D:\Bound`, target: `D:\BOUND\x`, want: nil, inside: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			suffix, ok := boundRelativeSuffix(testCase.base, testCase.target)
			if ok != testCase.inside {
				t.Fatalf("boundRelativeSuffix(%q, %q) admitted=%v, want %v (suffix %q)", testCase.base, testCase.target, ok, testCase.inside, suffix)
			}
			if !testCase.inside {
				if suffix != nil {
					t.Fatalf("a refused target must report no suffix, got %q", suffix)
				}
				return
			}
			if !slices.Equal(suffix, testCase.want) {
				t.Fatalf("boundRelativeSuffix(%q, %q) = %q, %v, want %q, true", testCase.base, testCase.target, suffix, ok, testCase.want)
			}
		})
	}
}

// p84WidenSubstituteRange rewrites the mount-point fields so the declared substitute
// range takes in the NUL unit that terminates it: the shape the "NUL in the LAST unit
// of the range" exemption exists for, which NTFS may write when it counts the
// terminator in the length.
func p84WidenSubstituteRange(t *testing.T, buffer []byte) []byte {
	t.Helper()
	if len(buffer) < reparseMountPointPathOffset {
		t.Fatalf("the fixture buffer is %d bytes, shorter than the %d byte path header", len(buffer), reparseMountPointPathOffset)
	}
	// The four name fields sit between the reparse header and the path buffer; the
	// stored offsets are relative to the path buffer's start.
	fields := buffer[reparseDataHeaderSize:reparseMountPointPathOffset]
	if offset := binary.LittleEndian.Uint16(fields[0:2]); offset != 0 {
		t.Fatalf("the fixture expects the substitute name at the start of the path buffer, got offset %d", offset)
	}
	length := binary.LittleEndian.Uint16(fields[2:4])
	binary.LittleEndian.PutUint16(fields[2:4], length+2)
	// The print name's offset is a physical offset into the path buffer, so widening
	// the declared range must not move it.
	return buffer
}

// p84WriteMountPointBuffer writes one hand-built buffer onto a fresh directory and
// RETURNS the host's refusal instead of skipping, so a case that only needs the
// kernel to try can report the reason itself.
func p84WriteMountPointBuffer(link string, buffer []byte) error {
	if err := os.Mkdir(link, 0o755); err != nil {
		return err
	}
	path, err := syscall.UTF16PtrFromString(link)
	if err != nil {
		return err
	}
	handle, err := syscall.CreateFile(
		path,
		syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		return err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	var returned uint32
	result, _, callErr := procDeviceIoControl.Call(
		uintptr(handle),
		uintptr(p84FSCTLSetReparsePoint),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
		0,
		uintptr(unsafe.Pointer(&returned)),
		0,
	)
	if result == 0 {
		return callErr
	}
	return nil
}

// p84StoredTargetOf reports the stored target of the reparse point itself, or the
// parse error, so a refusal can be asserted instead of treated as a fixture failure.
func p84StoredTargetOf(link string) (string, error) {
	path, err := syscall.UTF16PtrFromString(link)
	if err != nil {
		return "", err
	}
	handle, err := syscall.CreateFile(
		path,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		return "", err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	return reparseStoredTarget(handle)
}

// TestP84W11EmbeddedNULUnitInDeclaredNameIsRefused pins N14.
//
// syscall.UTF16ToString stops at the first NUL unit it meets, so a stored name whose
// DECLARED range carries a NUL before its last unit decodes to a prefix of the stored
// name, and the boundary is then judged on that prefix. The fixture is built so the
// truncated prefix is the bound object itself - admitted - while the declared bytes
// continue into an upwards-walking tail: decoding first would admit an escape that
// the stored bytes describe. The read closure therefore scans the declared range
// before decoding and refuses an embedded NUL; a NUL in the last unit is still the
// terminator NTFS may count in the length, so that shape keeps decoding.
//
// Mutation check: deleting the embedded-NUL scan from the read closure in
// reparseNames makes the embedded cases decode a truncated name instead of refusing,
// which turns their refusal assertions red. Both controls (terminator outside the
// declared range, and terminator as the last unit inside it) decode identically
// either way, so they cannot mask the mutation.
func TestP84W11EmbeddedNULUnitInDeclaredNameIsRefused(t *testing.T) {
	const bound = `D:\Bound`
	// Truncated at the NUL this spells the bound object; the full declared range names
	// "D:\Bound\..\..\Windows".
	const embeddedSubstitute = `\??\D:\Bound` + "\x00" + `\..\..\Windows`
	const embeddedPrintName = `D:\Bound` + "\x00" + `\..\..\Windows`

	// The truncation the guard prevents, stated as a fact about the fixture: decoding
	// up to the NUL yields the bound object, which the containment test admits.
	truncated := embeddedSubstitute[:strings.IndexByte(embeddedSubstitute, 0)]
	if suffix, ok := boundRelativeSuffix(bound, normalizeReparseTarget(truncated)); !ok || len(suffix) != 0 {
		t.Fatalf("the fixture assumes the truncated name is the bound object, got suffix=%q ok=%v", suffix, ok)
	}

	t.Run("embedded-nul-in-substitute", func(t *testing.T) {
		buffer := p84MountPointBuffer(embeddedSubstitute, bound)
		substitute, printName, err := reparseNames(buffer, reparseMountPointPathOffset)
		if err == nil {
			t.Fatalf("a substitute name with a NUL unit inside its declared range must be refused, got substitute %q print name %q", substitute, printName)
		}
		if !strings.Contains(err.Error(), "embedded NUL unit") || !strings.Contains(err.Error(), "substitute") {
			t.Fatalf("the refusal must name the embedded NUL unit of the substitute name, got %v", err)
		}
	})

	t.Run("embedded-nul-in-print-name", func(t *testing.T) {
		buffer := p84MountPointBuffer(bound, embeddedPrintName)
		substitute, printName, err := reparseNames(buffer, reparseMountPointPathOffset)
		if err == nil {
			t.Fatalf("a print name with a NUL unit inside its declared range must be refused, got substitute %q print name %q", substitute, printName)
		}
		if !strings.Contains(err.Error(), "embedded NUL unit") || !strings.Contains(err.Error(), "print") {
			t.Fatalf("the refusal must name the embedded NUL unit of the print name, got %v", err)
		}
	})

	t.Run("control-terminator-outside-the-range", func(t *testing.T) {
		buffer := p84MountPointBuffer(bound, bound)
		substitute, printName, err := reparseNames(buffer, reparseMountPointPathOffset)
		if err != nil {
			t.Fatalf("the ordinary terminator layout must still parse, got %v", err)
		}
		if substitute != bound || printName != bound {
			t.Fatalf("the names must decode byte for byte, got substitute %q print name %q", substitute, printName)
		}
	})

	t.Run("control-terminator-is-the-last-unit-in-the-range", func(t *testing.T) {
		buffer := p84WidenSubstituteRange(t, p84MountPointBuffer(bound, bound))
		substitute, printName, err := reparseNames(buffer, reparseMountPointPathOffset)
		if err != nil {
			t.Fatalf("a NUL that is the last unit of the declared range is the terminator, not an embedded unit, got %v", err)
		}
		if substitute != bound || printName != bound {
			t.Fatalf("the names must decode byte for byte, got substitute %q print name %q", substitute, printName)
		}
	})

	// The production entry counts as well: the same buffer goes through
	// FSCTL_SET_REPARSE_POINT so reparseStoredTarget has to refuse it, or the host
	// refused to STORE such a name in the first place and the parser cannot be reached
	// through the kernel (registered as a skip, not asserted away).
	t.Run("production-entry", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "embedded-nul")
		if err := p84WriteMountPointBuffer(link, p84MountPointBuffer(embeddedSubstitute, bound)); err != nil {
			t.Skipf("this host refuses to store a name with an embedded NUL unit, so the parser cannot be reached through FSCTL_SET_REPARSE_POINT: %v", err)
		}
		target, err := p84StoredTargetOf(link)
		if err == nil {
			t.Fatalf("a stored substitute name with an embedded NUL unit must be refused by the production reader, got target %q", target)
		}
		if !strings.Contains(err.Error(), "embedded NUL unit") {
			t.Fatalf("the production refusal must carry the embedded-NUL reason, got %v", err)
		}
	})
}

// p84ShortRootAlias manufactures a real 8.3 short name for one temp directory and
// reports the directory's short spelling. The alias is per-file (fsutil file
// setshortname), never a volume policy change, and it dies with the directory.
// ok=false means the host cannot set one and the caller must skip - a skip is
// never read as a pass (R-1 fixture discipline).
func p84ShortRootAlias(t *testing.T, dir string) (string, bool) {
	t.Helper()
	if output, err := exec.Command("fsutil", "file", "setshortname", dir, "R1PRO~1").CombinedOutput(); err != nil {
		t.Logf("setshortname unavailable on this host: %v (%s)", err, strings.TrimSpace(string(output)))
		return "", false
	}
	shortPath, err := p84ShortPathName(dir)
	if err != nil || shortPath == dir {
		t.Logf("the volume produced no short spelling for %s: %v", dir, err)
		return "", false
	}
	return shortPath, true
}

// p84ShortPathName reports the 8.3 short spelling of one existing path.
func p84ShortPathName(path string) (string, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	buffer := make([]uint16, 4096)
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetShortPathNameW")
	written, _, callErr := proc.Call(uintptr(unsafe.Pointer(ptr)), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	if written == 0 {
		if callErr == nil {
			callErr = syscall.EINVAL
		}
		return "", callErr
	}
	return syscall.UTF16ToString(buffer[:written]), nil
}

// TestP84W12ShortAliasRootPrefixIsAdmitted pins the R-1 fallback through the
// production entry: a junction whose stored target spells the bound root's own
// path through a REAL 8.3 short-name alias (the CI windows runner shape,
// RUNNER~1 vs runneradmin) must be admitted and read back the bound bytes.
//
// The fixture manufactures its own short alias at runtime (fsutil file
// setshortname on the temp root, a per-file alias and never a volume policy
// change). A host where setshortname is unavailable skips - a skip is registered,
// never counted as a pass, and the CI windows runner plus this host both execute
// the real leg.
//
// Mutation check: deleting the boundRelativeSuffixShortAlias fallback from
// reparseTargetRelativeSuffix turns the admission assertion red on the same
// fixture (the byte comparison refuses the short spelling).
func TestP84W12ShortAliasRootPrefixIsAdmitted(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "D")
	inside := filepath.Join(bound, "sub")
	if err := os.MkdirAll(inside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inside, "f"), []byte("inside"), 0o600); err != nil {
		t.Fatal(err)
	}
	shortRoot, ok := p84ShortRootAlias(t, root)
	if !ok {
		t.Skipf("this host cannot set an 8.3 short name for the temp root; the skip is not evidence")
	}

	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	link := filepath.Join(bound, "j")
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, filepath.Join(shortRoot, "D", "sub")).CombinedOutput(); err != nil {
		t.Fatalf("create the short-spelled junction: %v (%s)", err, output)
	}
	t.Cleanup(func() { _ = os.Remove(link) })

	file, err := OpenBoundDirFile(dir, `j\f`, 1<<20)
	if err != nil {
		t.Fatalf("a short-alias-spelled in-bound junction must be admitted, got %v", err)
	}
	data, truncated, err := file.ReadAll()
	_ = file.Close()
	if err != nil || truncated || string(data) != "inside" {
		t.Fatalf("the short-alias-spelled target must read the bound bytes, got %q truncated=%v err=%v", data, truncated, err)
	}
}

// TestP84W12AliasBranchRefusesDifferentRecord pins the alias admission's refusal
// of a short-alias prefix whose divergent component names a DIFFERENT object than
// the bound root's own component (a plain sibling directory name). The shape is
// refused by the admission machinery as a whole: the per-level same-record
// comparison, the same-entry short-name evidence and the final candidate-root
// identity comparison overlap for this shape (a divergent component that fails
// the same-record test also fails the same-entry evidence, which is bound to the
// record), so no single-check mutation discriminates, and the final comparison's
// refusal side is unreachable through a static fixture by construction. The test
// therefore pins the COMBINED refusal; the same-entry evidence's own refusal is
// pinned by the case-variant test, the reparse check's by the reparse test, and
// the final comparison's admission side by the W12 positive test.
//
// Mutation check (executed, see the slice record): removing the divergent-level
// proof block and the final identity comparison together turns this test red
// (the sibling name is admitted).
func TestP84W12AliasBranchRefusesDifferentRecord(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "D")
	other := filepath.Join(root, "E")
	if err := os.MkdirAll(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	shortRoot, ok := p84ShortRootAlias(t, root)
	if !ok {
		t.Skipf("this host cannot set an 8.3 short name for the temp root; the skip is not evidence")
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()
	base, err := dir.currentName()
	if err != nil {
		t.Fatalf("derive the bound name: %v", err)
	}
	target := filepath.Join(shortRoot, "E", "sub")
	if suffix, ok := boundRelativeSuffixShortAlias(dir, base, target); ok || suffix != nil {
		t.Fatalf("a candidate root that is not the bound object must be refused, got suffix=%q ok=%v", suffix, ok)
	}
}

// TestP84W12AliasBranchRefusesReparseComponent pins that a short-alias prefix
// whose ancestor chain crosses a junction is refused - the candidate-root walk
// must never follow a reparse point, and a junction leading to the RIGHT tree
// stays refused because the redirect itself is unverified.
//
// Honest scope note (per the review correction): through the helper's
// (suffix, ok) result the per-level reparse check cannot be isolated from the
// overlapping same-entry evidence and identity checks (under the follow-swap
// mutation alone the fixture stays refused by the evidence mismatch, which is
// registered in the slice record). The test therefore pins the COMBINED refusal
// of this shape; the discriminating mutation is the follow-swap combined with the
// wholesale proof removal (the M2C” shape in the slice record), which admits the
// fixture. The reparse check fires before the other proofs by construction -
// every component opens with FILE_OPEN_REPARSE_POINT and the reparse attribute
// refuses immediately.
//
// Mutation check (executed, see the slice record): the follow-swap plus the
// wholesale proof removal turns this test red (the junctioned candidate root is
// admitted); the follow-swap alone does not, because the evidence still refuses.
func TestP84W12AliasBranchRefusesReparseComponent(t *testing.T) {
	root := t.TempDir()
	real := filepath.Join(root, "real")
	bound := filepath.Join(real, "D")
	if err := os.MkdirAll(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", filepath.Join(root, "lnk"), real).CombinedOutput(); err != nil {
		t.Skipf("this host cannot create a junction: %v (%s)", err, output)
	}
	t.Cleanup(func() { _ = os.Remove(filepath.Join(root, "lnk")) })
	shortRoot, ok := p84ShortRootAlias(t, root)
	if !ok {
		t.Skipf("this host cannot set an 8.3 short name for the temp root; the skip is not evidence")
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()
	base, err := dir.currentName()
	if err != nil {
		t.Fatalf("derive the bound name: %v", err)
	}
	target := filepath.Join(shortRoot, "lnk", "D", "sub")
	if suffix, ok := boundRelativeSuffixShortAlias(dir, base, target); ok || suffix != nil {
		t.Fatalf("a candidate root crossing an unverified reparse point must be refused, got suffix=%q ok=%v", suffix, ok)
	}
}

// TestP84W12AliasBranchRefusesCaseVariant pins that a same-record spelling is not
// enough: a case-variant component shares the file record with the bound root's
// own component (the volume is case-insensitive), but the case-only gate refuses
// it before any short-name evidence, so the divergence never enters short-name
// admission - the W1 case discipline is preserved.
//
// Mutation note: the M3 record (removing the same-entry short-name evidence) is
// HISTORICAL - with the case-only gate the fixture refuses before the evidence,
// so M3 no longer discriminates this case. The evidence gate alone is now pinned
// by TestP84W12AliasBranchRefusesMissingEntryEvidence (the seam-injected case).
func TestP84W12AliasBranchRefusesCaseVariant(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "D")
	if err := os.MkdirAll(filepath.Join(bound, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()
	base, err := dir.currentName()
	if err != nil {
		t.Fatalf("derive the bound name: %v", err)
	}
	if !strings.HasSuffix(base, `\D`) {
		t.Fatalf("the fixture needs the bound root to end in a distinct component, got %q", base)
	}
	target := filepath.Join(root, "d", "sub")
	if suffix, ok := boundRelativeSuffixShortAlias(dir, base, target); ok || suffix != nil {
		t.Fatalf("a case-variant spelling is not an 8.3 alias and must be refused, got suffix=%q ok=%v", suffix, ok)
	}
}

// TestP84W12AliasBranchRefusesMissingEntryEvidence pins the evidence gate alone:
// the same-object preconditions hold for real - a genuine 8.3 alias, a non-case
// divergence, and the candidate-root walk proving the same file record - and ONLY
// the same-entry short-name evidence is missing. The absence is injected through
// the production seam (directoryEntryShortNameFn) because NTFS cannot be made to
// report a missing short name for a real alias. The current implementation must
// refuse; removing only the evidence gate admits the fixture.
//
// Mutation check (executed, see the slice record): deleting the evidence gate
// from boundRelativeSuffixShortAlias turns this test red (the alias is admitted
// on the same-record proof alone).
func TestP84W12AliasBranchRefusesMissingEntryEvidence(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "D")
	if err := os.MkdirAll(filepath.Join(bound, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	shortRoot, ok := p84ShortRootAlias(t, root)
	if !ok {
		t.Skipf("this host cannot set an 8.3 short name for the temp root; the skip is not evidence")
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()
	base, err := dir.currentName()
	if err != nil {
		t.Fatalf("derive the bound name: %v", err)
	}

	original := directoryEntryShortNameFn
	directoryEntryShortNameFn = func(syscall.Handle, string, uint32, uint32) (string, bool) {
		return "", false
	}
	t.Cleanup(func() { directoryEntryShortNameFn = original })

	target := filepath.Join(shortRoot, "D", "sub")
	if suffix, ok := boundRelativeSuffixShortAlias(dir, base, target); ok || suffix != nil {
		t.Fatalf("a divergence without the same-entry short-name evidence must be refused, got suffix=%q ok=%v", suffix, ok)
	}
}

// TestP84W12SeedIsFixedVolumeRoot pins the handle-based volume-root proof: a
// handle to the volume root of a fixed local NTFS volume proves, while a handle
// to an ordinary subdirectory never does (its file index is not the NTFS root
// record).
//
// Mutation check (executed, see the slice record): deleting the root-record
// comparison from seedIsFixedVolumeRoot turns the subdirectory assertion red.
func TestP84W12SeedIsFixedVolumeRoot(t *testing.T) {
	volumeRoot := filepath.VolumeName(t.TempDir()) + string(filepath.Separator)
	if driveTypeOf(volumeRoot) != driveFixed {
		t.Skipf("the temp volume is not a fixed local volume; the fallback refuses by design on this host")
	}
	rootHandle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(volumeRoot),
		accessListOrRead|accessReadAttributes|accessSynchronize,
		fileInformationShareAll, nil, syscall.OPEN_EXISTING,
		fileFlagBackupSemantics|fileFlagOpenReparsePoint, 0)
	if err != nil {
		t.Fatalf("open the volume root: %v", err)
	}
	defer syscall.CloseHandle(rootHandle)
	if !seedIsFixedVolumeRoot(rootHandle) {
		t.Skipf("the temp volume is not an NTFS volume root by this proof; the fallback refuses by design on this host")
	}

	subdir, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(t.TempDir()),
		accessListOrRead|accessReadAttributes|accessSynchronize,
		fileInformationShareAll, nil, syscall.OPEN_EXISTING,
		fileFlagBackupSemantics|fileFlagOpenReparsePoint, 0)
	if err != nil {
		t.Fatalf("open a subdirectory: %v", err)
	}
	defer syscall.CloseHandle(subdir)
	if seedIsFixedVolumeRoot(subdir) {
		t.Fatal("an ordinary subdirectory handle must not prove as the volume root")
	}
}

// TestP84W12NtfsRootRecordOf pins the pure 48-bit root-record judgment (Plan A):
// the MFT segment reference is 16 sequence bits in the HIGH 16 bits of the high
// word plus the 48-bit record number below them, so only
// high&0xFFFF == 0 && low == 5 is the root record. Record 0x1_0000_0005 (record
// bits non-zero in the high word) must refuse; the measured development-volume
// shape (high = 0x50000, sequence bits only, low = 5) must admit.
//
// Mutation check (executed, see the slice record): deleting the high-word mask
// from ntfsRootRecordOf turns the 0x1_0000_0005 assertion red.
func TestP84W12NtfsRootRecordOf(t *testing.T) {
	cases := []struct {
		name string
		high uint32
		low  uint32
		want bool
	}{
		{name: "root-record", high: 0x00000000, low: 5, want: true},
		{name: "measured-development-volume-root", high: 0x00050000, low: 5, want: true},
		{name: "sequence-bit-only", high: 0x00010000, low: 5, want: true},
		// Record number 0x1_0000_0005: the low 16 bits of the high word carry the
		// record bits, not the sequence - only the high 16 bits are sequence.
		{name: "record-bit-in-high-word", high: 0x00000001, low: 5, want: false},
		{name: "sequence-and-record-bits", high: 0x00010001, low: 5, want: false},
		{name: "other-low-word", high: 0x00000000, low: 6, want: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := ntfsRootRecordOf(testCase.high, testCase.low); got != testCase.want {
				t.Fatalf("ntfsRootRecordOf(%#x, %d) = %v, want %v", testCase.high, testCase.low, got, testCase.want)
			}
		})
	}
}

// p84DirEntryBuffer builds one FILE_ID_BOTH_DIR_INFORMATION record in an
// 8-byte-padded byte buffer for the focused scanner assertions. The FileName
// range is written WITHOUT a terminator (FileNameLength is authoritative), and
// the ShortName range is written exactly as many units as declared - variants
// mutate the buffer afterwards.
func p84DirEntryBuffer(name, shortName string, fileID uint64, next uint32) []byte {
	// utf16.Encode returns exactly the code units of the name, no terminator:
	// FileNameLength and ShortNameLength are authoritative, so nothing is
	// stripped here (the terminal-NUL variant writes the NUL explicitly).
	nameUnits := utf16.Encode([]rune(name))
	shortUnits := utf16.Encode([]rune(shortName))
	size := fileIDBothDirInfoFileNameOffset + len(nameUnits)*2
	padded := (size + 7) &^ 7
	buffer := make([]byte, padded)
	binary.LittleEndian.PutUint32(buffer[0:4], next)
	binary.LittleEndian.PutUint32(buffer[60:64], uint32(len(nameUnits)*2))
	if len(shortUnits) > 0 {
		buffer[68] = uint8(len(shortUnits) * 2)
	}
	for i, unit := range shortUnits {
		binary.LittleEndian.PutUint16(buffer[70+i*2:72+i*2], unit)
	}
	binary.LittleEndian.PutUint64(buffer[96:104], fileID)
	for i, unit := range nameUnits {
		binary.LittleEndian.PutUint16(buffer[fileIDBothDirInfoFileNameOffset+i*2:fileIDBothDirInfoFileNameOffset+i*2+2], unit)
	}
	return buffer
}

// TestP84W12DirectoryEntryScannerValidation pins the focused hardening of the
// directory-entry evidence scanner (scanFileIDBothDirInfoForShortName): a
// malformed buffer - an invalid NextEntryOffset, an illegal UTF-16 length, or a
// NUL inside a declared name range - refuses the admission fail-closed instead
// of being parsed on a guess, while the ordinary shapes still evidence.
//
// Mutation checks: deleting the NextEntryOffset validation, the length checks,
// or the embedded-NUL scans each turn the corresponding refusal assertion red.
func TestP84W12DirectoryEntryScannerValidation(t *testing.T) {
	const wantHigh, wantLow uint32 = 0x00000012, 0x00000034
	fileID := uint64(wantHigh)<<32 | uint64(wantLow)

	t.Run("valid-entry-with-short-name", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow)
		if refuse || !evidenced || shortName != "R1PRO~1" {
			t.Fatalf("the matching entry must evidence its short name, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("valid-entry-without-short-name", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "", fileID, 0)
		shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow)
		if refuse || evidenced || shortName != "" {
			t.Fatalf("an entry without a short name is no evidence, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("terminal-nul-in-short-name-range-is-the-terminator", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		// Count the terminator in the declared range, the N14 shape.
		buffer[68] = 16
		binary.LittleEndian.PutUint16(buffer[70+14:72+14], 0)
		shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow)
		if refuse || !evidenced || shortName != "R1PRO~1" {
			t.Fatalf("a terminal NUL unit is the terminator, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("embedded-nul-in-short-name-range-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		binary.LittleEndian.PutUint16(buffer[70+4:72+4], 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("an embedded NUL unit in the declared short-name range must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("embedded-nul-in-file-name-range-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("DD", "R1PRO~1", fileID, 0)
		// The NUL must sit INSIDE the declared two-unit name range.
		binary.LittleEndian.PutUint16(buffer[fileIDBothDirInfoFileNameOffset+2:fileIDBothDirInfoFileNameOffset+4], 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "DD", wantHigh, wantLow); !refuse {
			t.Fatalf("an embedded NUL unit in the declared file-name range must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("odd-file-name-length-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		binary.LittleEndian.PutUint32(buffer[60:64], 3)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("an odd FileNameLength must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("zero-file-name-length-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		binary.LittleEndian.PutUint32(buffer[60:64], 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("a zero FileNameLength must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("record-past-buffer-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		binary.LittleEndian.PutUint32(buffer[60:64], uint32(len(buffer)+16))
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("a declared record past the buffer must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("unaligned-next-entry-offset-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 105)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("an unaligned NextEntryOffset must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("next-entry-offset-shorter-than-a-header-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 8)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("a NextEntryOffset shorter than one header must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("next-entry-offset-past-buffer-refuses", func(t *testing.T) {
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		binary.LittleEndian.PutUint32(buffer[0:4], uint32(len(buffer)+104))
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse {
			t.Fatalf("a NextEntryOffset past the buffer must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("chain-second-entry-evidenced", func(t *testing.T) {
		first := p84DirEntryBuffer("A", "", 0, uint32(len(p84DirEntryBuffer("A", "", 0, 0))))
		second := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		buffer := append(first, second...)
		shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow)
		if refuse || !evidenced || shortName != "R1PRO~1" {
			t.Fatalf("the chain must reach the second entry, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("file-id-low-half-mismatch-refuses", func(t *testing.T) {
		// A matched entry whose record belongs to another object is a contradiction,
		// not a "not found": the scan must refuse instead of looking further.
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID^1, 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("a matching entry bound to another object must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("file-id-high-half-mismatch-refuses", func(t *testing.T) {
		// Only the high 32 bits differ: a comparison that degraded to the low half
		// would admit this buffer.
		buffer := p84DirEntryBuffer("D", "R1PRO~1", fileID^uint64(1)<<32, 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("a high-half file id mismatch must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("empty-buffer-refuses", func(t *testing.T) {
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(nil, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("an empty buffer must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("cross-record-name-range-refuses", func(t *testing.T) {
		// A non-last record is bounded by its own NextEntryOffset: a name range that
		// reaches one unit past that boundary must refuse. The padding is filled with
		// a non-NUL value, the range stops before the next record's zero bytes, and
		// the next record is itself a valid chained record, so ONLY the record bound
		// can refuse this buffer: deleting it lets the scan walk the whole chain and
		// return no evidence instead.
		first := p84DirEntryBuffer("DD", "R1PRO~1", fileID, uint32(len(p84DirEntryBuffer("DD", "R1PRO~1", fileID, 0))))
		for i := fileIDBothDirInfoFileNameOffset + 4; i < len(first); i++ {
			first[i] = 0x11
		}
		second := p84DirEntryBuffer("Z", "", 0, uint32(len(first)))
		third := p84DirEntryBuffer("Q", "", 0, 0)
		buffer := append(append(first, second...), third...)
		binary.LittleEndian.PutUint32(buffer[60:64], 10)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("a name range crossing its record boundary must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("unaligned-but-in-buffer-offset-refuses", func(t *testing.T) {
		// The unaligned NextEntryOffset (113) is inside the buffer, leaves a full
		// header behind it, and still contains this record's declared name range, so
		// only the alignment rule can refuse this buffer: removing it would let the
		// scan reach the planted valid entry at 113 and evidence.
		first := p84DirEntryBuffer("A", "", 0, 113)
		second := p84DirEntryBuffer("D", "R1PRO~1", fileID, 0)
		buffer := make([]byte, 113+len(second))
		copy(buffer, first)
		copy(buffer[113:], second)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("an unaligned NextEntryOffset inside the buffer must refuse, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})

	t.Run("embedded-nul-in-non-matching-entry-refuses", func(t *testing.T) {
		// Every entry's declared ranges are validated, not only the matching one.
		buffer := p84DirEntryBuffer("A", "R1PRO~1", 0, 0)
		binary.LittleEndian.PutUint16(buffer[70+4:72+4], 0)
		if shortName, evidenced, refuse := scanFileIDBothDirInfoForShortName(buffer, "D", wantHigh, wantLow); !refuse || evidenced {
			t.Fatalf("a malformed non-matching entry must refuse the whole batch, got %q evidenced=%v refuse=%v", shortName, evidenced, refuse)
		}
	})
}
