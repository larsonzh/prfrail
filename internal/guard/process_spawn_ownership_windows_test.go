//go:build windows

package guard

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestOnSpawnHookRunsBeforeResumeWindows pins I1 on the only platform where the
// ordering is observable: the child is created suspended inside its job, so its
// first user instruction can only run after guard resumes it, and guard may only
// resume after OnSpawn returned.
//
// The kill-on-close half of I1 (host death during the hook kills the suspended
// tree) needs a throw-away host process and is registered as a later-slice
// candidate: this slice pins the ordering that makes the window unreachable.
func TestOnSpawnHookRunsBeforeResumeWindows(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "spawned.marker")
	blocked := beginBlockedSpawn(t, marker)
	if blocked.identity.PID <= 0 || blocked.identity.StartToken == "" {
		t.Fatalf("OnSpawn must run after the identity is captured, got %+v", blocked.identity)
	}
	// The hook is still blocked, so a child that ran would have written the marker.
	time.Sleep(spawnOwnershipBlockedWindow)
	if _, statErr := os.Stat(marker); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the child executed a user instruction while OnSpawn was blocked (stat err=%v)", statErr)
	}
	process := blocked.finish(t)
	if process.Identity() != blocked.identity {
		t.Fatalf("resume changed the identity: got %+v want %+v", process.Identity(), blocked.identity)
	}
	waitForPath(t, marker)
}
