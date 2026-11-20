//go:build linux

package guard

// P83 FIX ROUND 7 counterexample (item T3d) on the Linux half.
//
// The child is created from the duplicated descriptor that was hashed
// (/proc/self/fd/<n>), so the spawn never resolves the name again. The test drives
// the real platformProcess.start while a writer repeatedly replaces the file at the
// original path: the only admissible outcomes are "the approved image ran" or "the
// binding refused", and running the replacement is a failure.
//
// This test is tagged linux and runs on the Ubuntu leg; the Windows host never
// compiles it, so its RED evidence is recorded from that leg.
//
// Mutation check: restoring `cmd.Path = spec.Command` (path re-resolution instead of
// descriptor exec) makes the writer's replacement runnable in the window between the
// digest and the spawn, and turns this red.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	p83Fix7OriginalScript    = "#!/bin/sh\nprintf 'APPROVED-IMAGE'\n"
	p83Fix7ReplacementScript = "#!/bin/sh\nprintf 'REPLACEMENT-IMAGE'\n"
)

// p83Fix7WriteScript writes one script body at path with the executable bit set.
func p83Fix7WriteScript(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o700); err != nil {
		t.Fatal(err)
	}
}

func TestP83Fix7T3dDescriptorExecNeverResolvesTheReplacedPath(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "agent.sh")
	p83Fix7WriteScript(t, script, p83Fix7OriginalScript)
	digest := executableDigestOf(t, script)

	// The replacer flips the file at the original name for the whole run. It writes
	// through a temporary name and renames, so the path always carries a complete
	// script and the digest either matches the approval or refuses.
	stop := make(chan struct{})
	done := make(chan struct{})
	var flips atomic.Int64
	go func() {
		defer close(done)
		temporary := script + ".swap"
		for {
			select {
			case <-stop:
				return
			default:
			}
			if err := os.WriteFile(temporary, []byte(p83Fix7ReplacementScript), 0o700); err != nil {
				return
			}
			if err := os.Rename(temporary, script); err != nil {
				return
			}
			if err := os.WriteFile(temporary, []byte(p83Fix7OriginalScript), 0o700); err != nil {
				return
			}
			if err := os.Rename(temporary, script); err != nil {
				return
			}
			flips.Add(1)
		}
	}()
	defer func() {
		close(stop)
		<-done
	}()

	// Warm the flipper so the replacement is in the window when the first spawn runs.
	time.Sleep(20 * time.Millisecond)

	runs, refusals := 0, 0
	for iteration := 0; iteration < 40; iteration++ {
		process, err := preparePlatformProcess()
		if err != nil {
			t.Fatal(err)
		}
		var stdout strings.Builder
		command, startErr := process.start(&ProcessSpec{
			Command:                script,
			Dir:                    dir,
			Stdout:                 &stdout,
			ExecutableDigest:       digest,
			ExecutableDigestDomain: p83Fix7ExecutableDomain,
		})
		if startErr != nil {
			if !errors.Is(startErr, ErrExecutableBinding) {
				t.Fatalf("an approved spawn may only refuse with the binding error, got %v", startErr)
			}
			refusals++
			process.close()
			continue
		}
		if err := command.waitFn(); err != nil {
			t.Fatalf("wait for the approved child: %v", err)
		}
		process.close()
		runs++
		if got := stdout.String(); got != "APPROVED-IMAGE" {
			t.Fatalf("the child must be created from the verified descriptor; it printed %q", got)
		}
	}
	// The flipper must actually have replaced the name, otherwise the test proves
	// nothing about the descriptor.
	if flips.Load() == 0 {
		t.Skip("the replacer never completed a swap on this host; the descriptor property is unproven here")
	}
	t.Logf("approved runs=%d binding refusals=%d path swaps=%d", runs, refusals, flips.Load())
}

// executableDigestOf computes the digest of the file as the binding does:
// evidence.BoundFile.Hash is the domain-separated digest of the file bytes.
func executableDigestOf(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return evidence.Digest(p83Fix7ExecutableDomain, content)
}
