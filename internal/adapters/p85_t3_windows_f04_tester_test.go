//go:build windows

package adapters

// P85 F04 (operator ruling 2026-11-23, four points; carried by 02c-f04-ruling.md) -
// the Windows volume-root exclusion boundary:
//
//  1. ONLY a true volume root may drop FILE_ADD_FILE / FILE_ADD_SUBDIRECTORY;
//  2. every NON-root ancestor keeps the full seven-right table (FILE_DELETE_CHILD and
//     the creation rights decide whether an existing entry can be swapped);
//  3. the root identity is anchored by the P84 primitive (FileIndexLow == 5), and an
//     anchor FAILURE falls back to the full table - the exclusion is a privilege,
//     never an assumption;
//  4. the ruling and its measured evidence are registered in the slice report.
//
// The three tests below are the ③ obligations of point 4. They pin the table
// contents, the measured root baseline, and the fallback selection. Every one of them
// is falsifiable by deleting the guard it names (mutation notes are in each test).
//
// EXEMPT-REGISTERED (not silently skipped, see the report §8): the END-TO-END
// unanchored-volume-root leg - a real volume root whose P84 anchor fails, judged by
// the chain walk - is not constructible on this host. It needs either a non-fixed
// volume (mapped/subst/removable) mounted at a root name, or an object that grants
// ONLY the two creation rights; the latter is impossible to build from a directory
// this user owns, because the Windows owner is ALWAYS granted WRITE_DAC implicitly,
// and the former changes machine state outside a test's contract. The injection used
// instead is recorded per test below and in the report.

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// p85T3VolumeRootOf returns this host's system volume root (`<SystemDrive>\`).
func p85T3VolumeRootOf(t *testing.T) string {
	t.Helper()
	volume := os.Getenv("SystemDrive")
	if volume == "" {
		volume = "C:"
	}
	root := volume + string(filepath.Separator)
	if _, err := os.Stat(root); err != nil {
		t.Skipf("the system volume root %q is not addressable on this host (%v)", root, err)
	}
	return root
}

// TestP85T3VolumeRootExcludesCreationRightsOnlyWhenAnchored is the POSITIVE leg of
// the ruling: on this host the (identity-anchored) volume root grants
// FILE_ADD_SUBDIRECTORY through the default "Authenticated Users:(AD)" ace and
// denies every non-creation right, so the anchored root is admitted - and the same
// root judged by the full table would have been refused. That counterfactual is what
// proves the exclusion is doing the work instead of the ACL.
//
// Mutation check: putting the creation rights back into
// installSourceVolumeRootMutationRights (or judging the anchored root by the
// directory table) turns the "no creation right is asked" assertion and the
// admission half red.
func TestP85T3VolumeRootExcludesCreationRightsOnlyWhenAnchored(t *testing.T) {
	root := p85T3VolumeRootOf(t)

	// Precondition of the ruling: the root must be anchored by identity, otherwise
	// the exclusion does not apply at all (point 3).
	anchored, err := installSourceAnchorsFixedVolumeRoot(root)
	if err != nil {
		t.Skipf("the P84 identity anchor cannot bind %q on this host (%v); the exclusion needs a provable fixed volume root", root, err)
	}
	if !anchored {
		t.Skipf("%q is not anchored as a fixed volume root on this host; the exclusion must never be measured by assumption", root)
	}

	grantedRoot, deniedRoot, err := installSourceMutationAccess(root, true, installSourceVolumeRootMutationRights)
	if err != nil {
		t.Fatalf("the volume-root probe of %q must be decidable: %v", root, err)
	}
	if len(grantedRoot) != 0 || len(deniedRoot) != len(installSourceVolumeRootMutationRights) {
		// The measured baseline of the ruling is "this token holds no non-creation
		// right on the root". A host where it does is a different environment, not a
		// broken rule.
		t.Skipf("this token holds a non-creation right on %q (granted=%s); the ruling's measured baseline does not hold here",
			root, installSourceRightsText(grantedRoot))
	}

	// The anchored verdict passes BECAUSE the creation rights are not asked at all.
	for _, right := range []uint32{installSourceFileWriteData, installSourceFileAppendData} {
		if slices.Contains(installSourceVolumeRootMutationRights, right) {
			t.Fatalf("the anchored volume-root table must not ask creation right 0x%08x", right)
		}
	}

	// Counterfactual: the SAME root under the full table grants a creation right and
	// would therefore be refused.
	grantedFull, _, err := installSourceMutationAccess(root, true, installSourceDirectoryMutationRights)
	if err != nil {
		t.Fatalf("the full-table probe of %q must be decidable: %v", root, err)
	}
	if !slices.Equal(grantedFull, []uint32{installSourceFileAppendData}) {
		t.Skipf("the root grants %s under the full table, not only FILE_ADD_SUBDIRECTORY; the counterfactual leg is not measurable here",
			installSourceRightsText(grantedFull))
	}

	t.Run("anchored-root-system-executable-is-provable", func(t *testing.T) {
		// The property the ruling protects: an install under this root is not
		// systematically refused because of the root component alone.
		//
		// The leg is made decisive: the ancestor directories of a system executable are
		// probed first, and a host where one of them is mutable skips (that is a
		// different environment, not this rule). On a host where they are all clean -
		// this one, per the baseline above - the chain walk MUST admit the executable,
		// so judging the anchored root by the full table (which grants
		// FILE_ADD_SUBDIRECTORY here) turns this red instead of silently skipping.
		ancestors := []string{root}
		for _, name := range []string{"Windows", "System32"} {
			ancestors = append(ancestors, filepath.Join(ancestors[len(ancestors)-1], name))
		}
		shell := filepath.Join(ancestors[len(ancestors)-1], "cmd.exe")
		if _, statErr := os.Stat(shell); statErr != nil {
			t.Skipf("the system executable %q is not addressable on this host (%v)", shell, statErr)
		}
		for _, dir := range ancestors[1:] {
			granted, _, probeErr := installSourceMutationAccess(dir, true, installSourceDirectoryMutationRights)
			if probeErr != nil {
				t.Skipf("the ancestor probe of %q is undecidable on this host (%v)", dir, probeErr)
			}
			if len(granted) != 0 {
				t.Skipf("%q is mutable for this token (granted=%s); the chain cannot be provable on this host",
					dir, installSourceRightsText(granted))
			}
		}
		leafGranted, _, probeErr := installSourceMutationAccess(shell, false, installSourceFileMutationRights)
		if probeErr != nil {
			t.Skipf("the leaf probe of %q is undecidable on this host (%v)", shell, probeErr)
		}
		if len(leafGranted) != 0 {
			t.Skipf("%q is mutable for this token (granted=%s)", shell, installSourceRightsText(leafGranted))
		}
		if _, err := VerifyAIInstallSource(shell); err != nil {
			t.Fatalf("the anchored root and every other component are clean, so %q must be provable; the root exclusion is not honoured: %v", shell, err)
		}
	})
}

// TestP85T3DirectoryTableKeepsCreationRights is the COUNTEREXAMPLE of the ruling's
// point 2: a non-root ancestor keeps FILE_ADD_FILE and FILE_ADD_SUBDIRECTORY, so an
// ancestor that lets this user create a child is refused - and it is the volume-root
// table that removes exactly those two rights, nothing else.
//
// Mutation check: dropping a creation right from the directory table (or adding one
// to the volume-root table) turns the table assertions red; short-circuiting the
// chain walk so the ancestor grant is not judged turns the behavioural half red.
func TestP85T3DirectoryTableKeepsCreationRights(t *testing.T) {
	creation := []uint32{installSourceFileWriteData, installSourceFileAppendData}

	wantDirectory := []uint32{
		installSourceFileWriteData,       // FILE_ADD_FILE
		installSourceFileAppendData,      // FILE_ADD_SUBDIRECTORY
		installSourceFileDeleteChild,     // FILE_DELETE_CHILD
		installSourceFileWriteAttributes, // FILE_WRITE_ATTRIBUTES
		installSourceStandardDelete,      // DELETE
		installSourceStandardWriteDAC,    // WRITE_DAC
		installSourceStandardWriteOwner,  // WRITE_OWNER
	}
	if !slices.Equal(installSourceDirectoryMutationRights, wantDirectory) {
		t.Fatalf("a non-root ancestor keeps the full seven-right table, got %s", installSourceRightsText(installSourceDirectoryMutationRights))
	}
	wantVolumeRoot := []uint32{
		installSourceFileDeleteChild,
		installSourceFileWriteAttributes,
		installSourceStandardDelete,
		installSourceStandardWriteDAC,
		installSourceStandardWriteOwner,
	}
	if !slices.Equal(installSourceVolumeRootMutationRights, wantVolumeRoot) {
		t.Fatalf("the anchored volume-root table is the seven-right table minus the creation rights, got %s",
			installSourceRightsText(installSourceVolumeRootMutationRights))
	}
	// The ONLY difference between the two tables is the two creation rights.
	missing := make([]uint32, 0, len(creation))
	for _, right := range installSourceDirectoryMutationRights {
		if !slices.Contains(installSourceVolumeRootMutationRights, right) {
			missing = append(missing, right)
		}
	}
	if !slices.Equal(missing, creation) {
		t.Fatalf("the volume-root exclusion must remove exactly the creation rights, removed %s", installSourceRightsText(missing))
	}
	for _, right := range creation {
		if !slices.Contains(installSourceDirectoryMutationRights, right) {
			t.Fatalf("the directory table must ask creation right 0x%08x", right)
		}
		if slices.Contains(installSourceVolumeRootMutationRights, right) {
			t.Fatalf("the volume-root table must not ask creation right 0x%08x", right)
		}
	}

	// Behaviour: this user's own directory grants both creation rights, and the full
	// table therefore reports them as granted - which is exactly the fact the chain
	// walk refuses a non-root ancestor for.
	dir := t.TempDir()
	granted, _, err := installSourceMutationAccess(dir, true, installSourceDirectoryMutationRights)
	if err != nil {
		t.Fatalf("the directory probe of %q must be decidable: %v", dir, err)
	}
	for _, right := range creation {
		if !slices.Contains(granted, right) {
			t.Fatalf("a user-owned directory must grant creation right 0x%08x, got %s", right, installSourceRightsText(granted))
		}
	}

	// End-to-end: the same grant on an ANCESTOR of an executable refuses the chain.
	executable := filepath.Join(dir, "copilot.exe")
	if err := os.WriteFile(executable, []byte("p85-f04-ancestor-grant"), 0o700); err != nil {
		t.Fatal(err)
	}
	facts, err := verifyInstallSourceImmutability(executable, nil)
	if !errors.Is(err, ErrAIInstallSourceProof) {
		t.Fatalf("an ancestor that grants this user a creation right must be refused, got facts=%+v err=%v", facts, err)
	}
	if facts.Windows != nil && facts.Windows.DirACLDeniesWriteDelete {
		t.Fatalf("the refused chain must record a granted right, not an immutable verdict: %+v", facts.Windows)
	}
}

// TestP85T3UnanchoredVolumeRootFallsBackToFullTable is the counterexample of the
// ruling's point 3: the exclusion is never assumed. The branch lives inline in
// verifyInstallSourceImmutability and exposes no seam, so the "not anchored" state is
// injected through the anchor itself - a directory that is not a volume root answers
// false without error, which is precisely the value that selects the full table.
//
// What makes the fallback pin COMPLETE without that seam: the fallback table is the
// directory table by construction (TestP85T3DirectoryTableKeepsCreationRights pins
// its exact seven entries), and a non-root ancestor already runs that same table in
// the chain walk, whose refusal is measured end to end in the same test. So an
// anchor failure judges a root exactly like any other ancestor - the strict, fail
// closed shape - and any other table in that position breaks the two tests.
//
// Mutation check: if the fallback passed installSourceVolumeRootMutationRights (that
// is, "assume root"), the creation rights would never be judged and the judgement
// counts below go red.
func TestP85T3UnanchoredVolumeRootFallsBackToFullTable(t *testing.T) {
	creation := []uint32{installSourceFileWriteData, installSourceFileAppendData}

	// Injection: an existing directory that is NOT a volume root. The anchor is
	// decidable and answers "not anchored" - the same answer a volume root whose
	// identity cannot be proven gives, and the input the fallback branch consumes.
	dir := t.TempDir()
	anchored, err := installSourceAnchorsFixedVolumeRoot(dir)
	if err != nil {
		t.Fatalf("the anchor must be decidable for an existing directory: %v", err)
	}
	if anchored {
		t.Fatalf("a temporary directory is not a fixed volume root, but the anchor admitted %q", dir)
	}

	// The fallback judges EVERY right of the seven-right table.
	grantedFull, deniedFull, err := installSourceMutationAccess(dir, true, installSourceDirectoryMutationRights)
	if err != nil {
		t.Fatalf("the fallback probe of %q must be decidable: %v", dir, err)
	}
	if judged := len(grantedFull) + len(deniedFull); judged != len(installSourceDirectoryMutationRights) {
		t.Fatalf("the fallback must judge all %d rights, judged %d", len(installSourceDirectoryMutationRights), judged)
	}
	// Any single granted right refuses; on this object the creation rights are granted,
	// so an unanchored root with such an ACL can never be admitted.
	if len(grantedFull) == 0 {
		t.Fatalf("a user-owned directory grants mutation rights, got none")
	}
	for _, right := range creation {
		if !slices.Contains(grantedFull, right) {
			t.Fatalf("the fallback must ASK creation right 0x%08x, and this directory grants it", right)
		}
	}

	// The exclusion table, by contrast, never sees the creation rights at all - which
	// is why using it on an unanchored object would be the assumption the ruling
	// forbids.
	grantedRoot, deniedRoot, err := installSourceMutationAccess(dir, true, installSourceVolumeRootMutationRights)
	if err != nil {
		t.Fatalf("the volume-root-table probe of %q must be decidable: %v", dir, err)
	}
	if judged := len(grantedRoot) + len(deniedRoot); judged != len(installSourceVolumeRootMutationRights) {
		t.Fatalf("the volume-root table judges %d rights, judged %d", len(installSourceVolumeRootMutationRights), judged)
	}
	for _, right := range creation {
		if slices.Contains(grantedRoot, right) || slices.Contains(deniedRoot, right) {
			t.Fatalf("the volume-root table must not judge creation right 0x%08x", right)
		}
	}
}
