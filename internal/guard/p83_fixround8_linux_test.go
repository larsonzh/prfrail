//go:build linux

package guard

// P83 FIX ROUND 8 counterexample (item E): the verified executable descriptor is
// inherited by the intended child alone. It is handed over through ExtraFiles,
// which marks it inheritable for that one spawn; a raw syscall.Dup would clear
// close-on-exec in the PARENT as well, so every later child of the process would
// carry a descriptor to the verified image.
//
// The Ubuntu leg of the CI executes this file; the Windows host cannot run it.
//
// Mutation check: replacing ExtraFiles with syscall.Dup on approved.file makes the
// sibling child below list the image in /proc/self/fd and turns the assertion red.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestP83Fix8EVerifiedDescriptorInheritedOnlyByItsChild(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "agent.sh")
	content := []byte("#!/bin/sh\nreadlink /proc/self/fd/3\nsleep 0.4\n")
	if err := os.WriteFile(script, content, 0o755); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	process, err := preparePlatformProcess()
	if err != nil {
		t.Fatal(err)
	}
	command, err := process.start(&ProcessSpec{
		Command:                script,
		Dir:                    dir,
		Env:                    []string{"PATH=/usr/bin:/bin"},
		Stdout:                 &output,
		Stderr:                 &output,
		ExecutableDigest:       evidence.Digest(p83Fix7ExecutableDomain, content),
		ExecutableDigestDomain: p83Fix7ExecutableDomain,
	})
	if err != nil {
		t.Fatalf("start the approved executable: %v", err)
	}

	// A sibling child started while the approved child is alive must not see the
	// verified descriptor. It lists every descriptor it holds; a readlink on a
	// descriptor that closes mid-loop must not fail the inspection itself.
	sibling := exec.Command("/bin/sh", "-c", "for fd in /proc/self/fd/*; do readlink \"$fd\" 2>/dev/null || true; done")
	siblingOutput, siblingErr := sibling.CombinedOutput()
	if siblingErr != nil {
		t.Fatalf("inspect the sibling descriptors: %v (%s)", siblingErr, siblingOutput)
	}
	if strings.Contains(string(siblingOutput), script) {
		t.Fatalf("a concurrent child inherited the verified descriptor: %s", siblingOutput)
	}

	if err := command.wait(); err != nil {
		t.Fatalf("the approved child must run to completion, got %v (%s)", err, output.String())
	}
	// The intended child does see it: the child is created from that descriptor.
	if !strings.Contains(output.String(), script) {
		t.Fatalf("the approved child must reach its own descriptor at fd 3, got %q", output.String())
	}
}
