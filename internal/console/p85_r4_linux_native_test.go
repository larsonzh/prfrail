//go:build linux

package console

// P85 design rev 4, section 6: the Linux-NATIVE legs of M13 and M16.
//
// The counterexamples in the platform-independent files reach the same properties
// through the observation seam. What is only executable on Linux is the kernel-level
// fact underneath them, and that is what these two legs pin:
//
//	M13 - a directory replaced at the same name is a DIFFERENT dev:ino object while
//	      the bound descriptor keeps addressing the original one, so "the name still
//	      resolves to the bound object" is a real identity comparison and not a
//	      spelling comparison;
//	M16 - the carrier publication on this platform runs entirely inside the bound
//	      descriptor (create + write + fsync + link + unlink + fsync), so a successful
//	      first publication leaves no temporary name behind and the re-entry adopts
//	      the very bytes that were published.
//
// These legs are only built for linux, so a Windows host reports them as absent
// rather than as a silent pass. They are still type-checked on every host through
// `GOOS=linux go vet ./internal/console`.

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85R4SwapPolicyDirectory renames the policy directory away and puts a complete
// stand-in (same names, same bytes, different directory object) at the path.
//
// P85 R=7 (user ruling of 2026-11-23): this helper moved here from
// p85_r4_m13_chain_tester_test.go. Only this linux-tagged leg calls it, so living
// in the platform-independent file made it look unused in the windows and darwin
// IDE views. See the note on p85R4ReplaceArtifact for why the Windows legs stage
// the artifact instead.
func p85R4SwapPolicyDirectory(t *testing.T, runtime *agentRunnerRuntime) error {
	t.Helper()
	directory := agentRunnerPolicyDirectory(runtime.runRoot)
	retired := directory + ".retired"
	if err := os.Rename(directory, retired); err != nil {
		return err
	}
	if err := os.Mkdir(directory, 0o755); err != nil {
		return err
	}
	for _, path := range []string{
		agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1),
		agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1),
	} {
		payload, err := os.ReadFile(filepath.Join(retired, filepath.Base(path)))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(directory, filepath.Base(path)), payload, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// TestP85R4M13LinuxNativeIdentityIsTheDirectoryObject is the M13 native leg.
func TestP85R4M13LinuxNativeIdentityIsTheDirectoryObject(t *testing.T) {
	runtime := p85R4FrozenCreationRuntime(t)
	directory := agentRunnerPolicyDirectory(runtime.runRoot)

	var before syscall.Stat_t
	if err := syscall.Lstat(directory, &before); err != nil {
		t.Fatalf("lstat the policy directory: %v", err)
	}

	chain, err := bindAgentRunnerPolicyDomain(runtime.runRoot, p83TaskID, p83StepID, 1)
	if err != nil {
		t.Fatalf("bind the policy domain chain: %v", err)
	}
	defer func() { _ = chain.Close() }()

	if err := p85R4SwapPolicyDirectory(t, runtime); err != nil {
		t.Fatalf("swap the policy directory for a same-named stand-in: %v", err)
	}

	var after syscall.Stat_t
	if err := syscall.Lstat(directory, &after); err != nil {
		t.Fatalf("lstat the stand-in directory: %v", err)
	}
	if before.Dev == after.Dev && before.Ino == after.Ino {
		t.Fatalf("the fixture must produce a DIFFERENT dev:ino at the same name, got %#x/%d twice", after.Dev, after.Ino)
	}
	// The bound handle still addresses the original object: identity is a property of
	// the descriptor, which is what makes the terminal comparison meaningful.
	if err := evidence.VerifyBoundDir(chain.dir, chain.dirID); err != nil {
		t.Fatalf("the bound handle must still address the object it bound, got %v", err)
	}
	if err := chain.ReassertIdentities(); err == nil {
		t.Fatal("a directory replaced at the same name must refuse the terminal identity re-assertion")
	}
}

// TestP85R4M16LinuxNativeCarrierPublicationLeavesNoTemporaryName is the M16 native
// leg.
func TestP85R4M16LinuxNativeCarrierPublicationLeavesNoTemporaryName(t *testing.T) {
	executable := p85T3ReentryProvableExecutable(t)
	root := t.TempDir()
	carrierPath := p85T3CarrierPath(root)

	published, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", p83Clock())
	if err != nil {
		t.Fatalf("a provable source must publish the carrier: %v", err)
	}
	if published.ExecutablePath != executable {
		t.Fatalf("the published facts must describe the probed executable, got %q want %q", published.ExecutablePath, executable)
	}
	wire, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatalf("the published carrier must be readable: %v", err)
	}

	// The publication is create + write + fsync + link + unlink inside ONE descriptor,
	// so nothing may be left beside the final name.
	entries, err := os.ReadDir(filepath.Dir(carrierPath))
	if err != nil {
		t.Fatalf("enumerate the carrier directory: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() != filepath.Base(carrierPath) {
			t.Fatalf("the handle-relative publication must leave no second name, found %q", entry.Name())
		}
	}

	adopted, err := recordOrConfirmAIInstallSourceProof(carrierPath, executable, "", p83Clock().Add(time.Hour))
	if err != nil {
		t.Fatalf("the re-entry must adopt the published carrier, got %v", err)
	}
	if !adopted.MatchesFacts(published) {
		t.Fatalf("the re-entry must adopt the recorded facts, got %+v want %+v", adopted, published)
	}
	again, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatalf("the adopted carrier must stay readable: %v", err)
	}
	if string(again) != string(wire) {
		t.Fatalf("an adopted carrier must never be rewritten, got %q want %q", again, wire)
	}
}
