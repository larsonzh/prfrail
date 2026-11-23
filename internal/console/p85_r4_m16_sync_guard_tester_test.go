package console

// P85 design rev 4, section 4 / test row M16 (rewrite of M8): the re-entry sync
// guard has to be DISCRIMINATING, so the two implementation facts it protects -
// "the re-entry syncs the carrier directory" and "a sync failure refuses" - can each
// be deleted by exactly one leg turning red.
//
// The old M8 leg 2 withdrew the guard AND cleared the failure injection in the same
// step, so it proved nothing about the guard: dropping the re-entry sync, or the
// guard condition, left it green (r2, F10 counterexample).
//
// The three legs below keep the injection orthogonal to the guard:
//
//	leg-P  positive control - first publication fails the sync; failSync=false with
//	       the guard still on adopts the carrier;
//	leg-R  discriminating red - failSync stays TRUE and the guard stays ON => the
//	       re-entry must refuse with the injected sync marker;
//	leg-G  discriminating green - failSync stays TRUE and ONLY the guard is turned
//	       off => the re-entry must adopt, and the adopted carrier must be the very
//	       bytes the refused publication left behind.
//
// Mutation contracts:
//   - deleting the re-entry's directory sync call in recordOrConfirmAIInstallSourceProof
//     makes leg-R adopt and turns this test red;
//   - deleting the "&& aiInstallSourceReentryRequiresSync" condition makes a failing
//     sync refuse unconditionally, leg-G refuses and this test turns red.

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP85R4M16ReentrySyncGuardIsDiscriminating is M16.
func TestP85R4M16ReentrySyncGuardIsDiscriminating(t *testing.T) {
	// The seam legs measure the RE-ENTRY guard, which is platform-independent, so
	// they must not inherit the Linux-only durability precondition of the native
	// leg. The native durability leg of the same entry lives in
	// p85_r4_linux_native_test.go and is skipped elsewhere.
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

	t.Run("leg-P-positive-control", func(t *testing.T) {
		root := t.TempDir()
		carrierPath := p85T3CarrierPath(root)

		failSync = true
		aiInstallSourceReentryRequiresSync = true
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now); !errors.Is(err, ErrAIInstallSourceUnproven) || !strings.Contains(err.Error(), p85F10SyncFailureMarker) {
			t.Fatalf("the fixture must start from a publication refused by the sync failure, got %v", err)
		}

		failSync = false
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now.Add(time.Hour)); err != nil {
			t.Fatalf("a re-entry whose sync succeeds must adopt the carrier, got %v", err)
		}
	})

	t.Run("leg-R-guard-on-while-sync-still-fails-refuses", func(t *testing.T) {
		root := t.TempDir()
		carrierPath := p85T3CarrierPath(root)

		failSync = true
		aiInstallSourceReentryRequiresSync = true
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now); !errors.Is(err, ErrAIInstallSourceUnproven) || !strings.Contains(err.Error(), p85F10SyncFailureMarker) {
			t.Fatalf("the fixture must start from a publication refused by the sync failure, got %v", err)
		}
		if _, statErr := os.Stat(carrierPath); statErr != nil {
			t.Fatalf("the hard link is published before the sync, so the carrier must be visible: %v", statErr)
		}

		// The injection is untouched and the guard is ON: the sync fails again, so
		// this re-entry must refuse - and it must be the SAME injected failure that
		// refuses it, not the carrier reader.
		_, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now.Add(time.Hour))
		if !errors.Is(err, ErrAIInstallSourceUnproven) || !strings.Contains(err.Error(), p85F10SyncFailureMarker) {
			t.Fatalf("with the guard on and the sync still failing the re-entry must refuse with the sync failure, got %v", err)
		}
	})

	t.Run("leg-G-guard-off-while-sync-still-fails-adopts", func(t *testing.T) {
		root := t.TempDir()
		carrierPath := p85T3CarrierPath(root)

		failSync = true
		aiInstallSourceReentryRequiresSync = true
		if _, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now); !errors.Is(err, ErrAIInstallSourceUnproven) {
			t.Fatalf("the fixture must start from a refused publication, got %v", err)
		}
		published, err := os.ReadFile(carrierPath)
		if err != nil {
			t.Fatalf("the carrier published before the failed sync must be readable: %v", err)
		}

		// ONLY the guard moves: the injection still fails. The re-entry must therefore
		// adopt, and the adopted object must be the bytes the refused publication
		// left on disk.
		aiInstallSourceReentryRequiresSync = false
		adopted, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", now.Add(time.Hour))
		if err != nil {
			t.Fatalf("with only the guard withdrawn the unchanged carrier must be adopted, got %v", err)
		}
		if adopted.ExecutablePath == "" {
			t.Fatal("the adopted facts must describe the carrier")
		}
		after, err := os.ReadFile(carrierPath)
		if err != nil {
			t.Fatalf("the adopted carrier must stay readable: %v", err)
		}
		if string(after) != string(published) {
			t.Fatalf("an adopted carrier must never be rewritten, got %q want %q", after, published)
		}
	})
}
