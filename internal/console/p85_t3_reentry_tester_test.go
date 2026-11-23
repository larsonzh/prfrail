package console

// P85 T3 (finding F10, final review) - the re-entry half of the carrier publication.
//
// The publish path hard-links the carrier into place and THEN syncs the directory. If
// the link became visible but that sync failed, the first call refused - and a
// re-entry used to adopt the visible carrier without repeating the sync, silently
// undoing the refusal. F10 requires the re-entry to confirm and re-sync the directory
// BEFORE it may adopt an existing carrier.
//
// The branch under test is the `if found { ... }` half of
// recordOrConfirmAIInstallSourceProof. It has exactly one return that hands back the
// carrier's facts, and that return sits after aiInstallSourceConfirmDirectory and
// aiInstallSourceSyncDirectory - so "adopted" implies "both ran on this entry" by
// construction. The tests below measure the consequence that is observable from
// outside: a carrier that cannot be confirmed is never adopted, and a carrier that can
// be confirmed is adopted from its recorded bytes (never re-derived).
//
// P85 F10 (design rev 3, section 4.5 / test row M8): the exact injection this file
// used to register as EXEMPT is now executable. The directory fsync is a package
// seam (aiInstallSourceCarrierDirectorySync) and the re-entry guard is a package
// flag (aiInstallSourceReentryRequiresSync), so BOTH halves of the counterexample
// run natively:
//
//	A. the first publication's hard link is visible but its directory sync fails =>
//	   the publication refuses, and a re-entry whose sync succeeds adopts;
//	B. with the carrier reading cleanly, withdrawing ONLY the sync guard
//	   (aiInstallSourceReentryRequiresSync=false) makes the same re-entry adopt =>
//	   the earlier refusal is attributable to the sync guard, not to the carrier
//	   reader that the link counterexample would trip first.
//
// The EXEMPT-REGISTERED P85/F10 note that lived here is therefore retired: the
// behaviour is covered by TestP85F10ReentryRefusesToAdoptWhenTheSameHandleSyncFails
// below (carried-by: the same-handle confirmation + sync of the `found` branch).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85T3ReentryProvableExecutable returns an executable whose installation source the
// live check can prove on this host, or skips the caller.
//
// P85 final review (F18) consequence: the interpreter rule now proves the ORIGINAL
// name chain, so a merged-/usr host refuses /bin/sh (a link) instead of proving its
// target. The candidate list therefore tries several system binaries and reports why
// none was provable, instead of assuming /bin/sh.
func p85T3ReentryProvableExecutable(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skipf("the re-entry durability leg is Linux-only: on %s the publish grants no directory durability at all (atomic visibility only)", runtime.GOOS)
	}
	candidates := []string{"/usr/bin/true", "/bin/true", "/usr/bin/busybox", "/bin/sh"}
	var lastErr error
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}
		if _, err := adapters.VerifyAIInstallSource(candidate); err != nil {
			lastErr = err
			continue
		}
		return candidate
	}
	t.Skipf("no provable probe executable on this host (last refusal: %v); this leg needs a native non-root Ubuntu run whose install layout is link-free", lastErr)
	return ""
}

// TestP85T3ReentryConfirmsDirectoryDurabilityBeforeAdoptingCarrier pins the F10
// consequence at the re-entry branch: adoption is only reachable through the
// confirmation (and sync) the re-entry performs, so a carrier the entry cannot confirm
// is refused and left byte-identical, while a confirmable one is adopted from the
// carrier's own facts rather than from a fresh observation.
//
// Mutation checks: deleting the aiInstallSourceConfirmDirectory call from the `found`
// branch turns "a-carrier-behind-a-link-is-never-adopted" red (the entry would adopt a
// directory it never bound); replacing the returned carrier facts with `facts` (the
// fresh observation) turns "the-adopted-facts-are-the-recorded-facts" red.
func TestP85T3ReentryConfirmsDirectoryDurabilityBeforeAdoptingCarrier(t *testing.T) {
	executable := p85T3ReentryProvableExecutable(t)

	t.Run("the-adopted-facts-are-the-recorded-facts", func(t *testing.T) {
		root := t.TempDir()
		carrierPath := p85T3CarrierPath(root)
		now := time.Date(2026, 11, 23, 8, 0, 0, 0, time.UTC)

		published, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now)
		if err != nil {
			t.Fatalf("a provable source must publish the carrier: %v", err)
		}
		wire, err := os.ReadFile(carrierPath)
		if err != nil {
			t.Fatalf("the published carrier must be readable: %v", err)
		}
		recorded, err := adapters.DecodeAIInstallSourceProof(wire)
		if err != nil {
			t.Fatalf("the published carrier must decode: %v", err)
		}

		// The re-entry: the carrier is already visible, so this takes the `found`
		// branch - the one that must confirm the directory before adopting.
		reentered, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now.Add(time.Hour))
		if err != nil {
			t.Fatalf("a confirmable carrier must be adopted: %v", err)
		}
		if !reentered.Matches(recorded) || reentered.ExecutablePath != published.ExecutablePath {
			t.Fatalf("the re-entry must adopt the RECORDED facts, got %+v want %+v", reentered, recorded.Facts())
		}
		again, err := os.ReadFile(carrierPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(again) != string(wire) {
			t.Fatalf("an adopted carrier must never be rewritten, got %q and %q", string(wire), string(again))
		}
	})

	t.Run("a-carrier-behind-a-link-is-never-adopted", func(t *testing.T) {
		root := t.TempDir()
		store := filepath.Join(root, "carrier-store")
		if err := os.MkdirAll(store, 0o755); err != nil {
			t.Fatal(err)
		}
		carrierPath := filepath.Join(store, adapters.AIInstallSourceProofFileName)
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", time.Now().UTC()); err != nil {
			t.Fatalf("a provable source must publish the carrier: %v", err)
		}
		published, err := os.ReadFile(carrierPath)
		if err != nil {
			t.Fatal(err)
		}
		// The carrier is visible and its bytes MATCH the live observation, which is
		// exactly the state the pre-fix re-entry adopted. The directory that must be
		// confirmed is now reachable only through a link, so the link-free binding the
		// re-entry performs before adoption cannot describe it.
		link := filepath.Join(root, "carrier-link")
		if err := os.Symlink(store, link); err != nil {
			t.Skipf("this host cannot create the directory link the leg needs: %v", err)
		}
		linkedPath := filepath.Join(link, adapters.AIInstallSourceProofFileName)

		_, err = recordOrConfirmAIInstallSourceProof(linkedPath, executable, "", time.Now().UTC())
		if err == nil {
			t.Fatal("a carrier whose directory cannot be confirmed must never be adopted")
		}
		if !errors.Is(err, ErrAIInstallSourceUnproven) {
			t.Fatalf("the refusal must be the install-source classification, got %v", err)
		}
		if strings.Contains(err.Error(), executable) {
			t.Fatalf("the refusal must not echo the executable path, got %q", err)
		}
		after, readErr := os.ReadFile(carrierPath)
		if readErr != nil {
			t.Fatalf("the carrier must stay readable after a refusal: %v", readErr)
		}
		if string(after) != string(published) {
			t.Fatal("a refusal must never rewrite the published carrier")
		}
	})
}

// p85F10SyncFailureMarker is the stable text the injected sync failure carries, so
// both halves (the first publication's raw propagation and the re-entry's wrapped
// refusal) can be asserted without depending on either wrapper's wording.
const p85F10SyncFailureMarker = "injected same-handle directory sync failure"

// p85F10SeamExecutable names an executable whose install source this host can prove.
// The seam test below measures the RE-ENTRY guard, which is platform-independent, so
// it must not inherit the Linux-only durability precondition of the native leg.
func p85F10SeamExecutable(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "windows" {
		return p85T3ReentryProvableExecutable(t)
	}
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = `C:\Windows`
	}
	shell := filepath.Join(root, "System32", "cmd.exe")
	if _, err := adapters.VerifyAIInstallSource(shell); err != nil {
		t.Skipf("the shell %s is not a provable install source on this host: %v", shell, err)
	}
	return shell
}

// TestP85F10ReentryRefusesToAdoptWhenTheSameHandleSyncFails is M8. See the file
// header for the two legs and the mutation check.
func TestP85F10ReentryRefusesToAdoptWhenTheSameHandleSyncFails(t *testing.T) {
	executable := p85F10SeamExecutable(t)
	now := time.Date(2026, 11, 23, 8, 0, 0, 0, time.UTC)

	previousSync := aiInstallSourceCarrierDirectorySync
	previousGuard := aiInstallSourceReentryRequiresSync
	failSync := true
	aiInstallSourceCarrierDirectorySync = func(bound *evidence.BoundDir) error {
		if failSync {
			return fmt.Errorf("%w: %s", ErrAIInstallSourceUnproven, p85F10SyncFailureMarker)
		}
		return previousSync(bound)
	}
	t.Cleanup(func() {
		aiInstallSourceCarrierDirectorySync = previousSync
		aiInstallSourceReentryRequiresSync = previousGuard
	})

	t.Run("a-failed-sync-refuses-the-first-publication", func(t *testing.T) {
		root := t.TempDir()
		carrierPath := p85T3CarrierPath(root)
		failSync = true
		_, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now)
		if !errors.Is(err, ErrAIInstallSourceUnproven) || !strings.Contains(err.Error(), p85F10SyncFailureMarker) {
			t.Fatalf("a same-handle directory sync failure must refuse the publication, got %v", err)
		}
		if _, statErr := os.Stat(carrierPath); statErr != nil {
			t.Fatalf("the hard link is published before the sync, so the carrier must be visible: %v", statErr)
		}

		// Positive control: the very same carrier is adopted once the handle sync
		// succeeds, so the refusal above is caused by the sync and nothing else.
		failSync = false
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now.Add(time.Hour)); err != nil {
			t.Fatalf("a re-entry whose same-handle sync succeeds must adopt the carrier, got %v", err)
		}
	})

	// The former second leg withdrew the guard AND cleared the failure injection in
	// one step, which made it non-discriminating (P85 review r2, F10): deleting the
	// re-entry sync or the guard condition left it green. It is superseded by the
	// three orthogonal legs of TestP85R4M16ReentrySyncGuardIsDiscriminating
	// (p85_r4_m16_sync_guard_tester_test.go): leg-R keeps the injection failing with
	// the guard ON, leg-G keeps it failing with ONLY the guard OFF.
}
