//go:build linux

package evidence

// P83 FIX ROUND 4 (F1) - a FIFO artifact must be refused, not waited on.
//
// Before round 4 the bound open used O_RDONLY alone, so opening a FIFO with no
// writer parked the collector inside the kernel forever: no error, no timeout, no
// evidence. Round 4 opens with O_NONBLOCK and then refuses anything that is not a
// regular file on the returned descriptor (boundArtifactRegular), so the FIFO
// becomes an unsafe-artifact refusal.
//
// This test cannot run on the Windows host - it needs mkfifo(3), a Linux
// openat2(2) kernel and the unix half of internal/evidence. Its execution venue is
// the Ubuntu CI leg; on Windows the file is excluded by the build constraint.
//
// Mutation checks, both on Linux:
//   - removing syscall.O_NONBLOCK from openat2Resolve blocks the open, the watcher
//     below fires, and the test fails fast instead of hanging the suite;
//   - removing the boundArtifactRegular call returns a handle for the FIFO and the
//     "no handle and ErrBoundFileNotRegular" assertion turns red.

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// p83Fix4FIFODecision is the message the watched open reports back.
type p83Fix4FIFODecision struct {
	err       error
	gotHandle bool
}

// TestBoundDirRefusesFIFOWithoutBlocking is the F1 counterexample.
func TestBoundDirRefusesFIFOWithoutBlocking(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "D")
	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	const fifoName = "artifact"
	if err := syscall.Mkfifo(filepath.Join(boundPath, fifoName), 0o600); err != nil {
		t.Fatalf("create the FIFO artifact: %v", err)
	}
	// Positive control in the same bound: the refusal below must be about the
	// artifact type, not about the bound directory being unusable.
	if err := os.WriteFile(filepath.Join(boundPath, "regular"), []byte("p83-fix4-payload"), 0o644); err != nil {
		t.Fatalf("write the regular in-bound artifact: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	// Close is deliberately not deferred: on a regression the open goroutine holds
	// dir.mu inside the kernel, and a deferred Close would block on that mutex and
	// turn a fast red into a hung suite.

	decisions := make(chan p83Fix4FIFODecision, 1)
	go func() {
		file, openErr := OpenBoundDirFile(dir, fifoName, 1024)
		if file != nil {
			_ = file.Close()
		}
		decisions <- p83Fix4FIFODecision{err: openErr, gotHandle: file != nil}
	}()

	select {
	case decision := <-decisions:
		if decision.gotHandle {
			_ = dir.Close()
			t.Fatalf("a FIFO artifact must not be returned as a readable handle")
		}
		if decision.err == nil {
			_ = dir.Close()
			t.Fatalf("a FIFO artifact must be refused with an error")
		}
		if !errors.Is(decision.err, ErrBoundFileNotRegular) {
			_ = dir.Close()
			t.Fatalf("a FIFO artifact must be refused as ErrBoundFileNotRegular, got %v", decision.err)
		}
	case <-time.After(10 * time.Second):
		// Deliberately no dir.Close(): see above.
		t.Fatalf("opening a FIFO artifact blocked instead of refusing; the collector would hang here")
	}

	file, err := OpenBoundDirFile(dir, "regular", 1024)
	if err != nil {
		_ = dir.Close()
		t.Fatalf("a regular in-bound artifact must still open after the FIFO refusal: %v", err)
	}
	_ = file.Close()
	if err := dir.Close(); err != nil {
		t.Fatalf("close the bound directory: %v", err)
	}
}
