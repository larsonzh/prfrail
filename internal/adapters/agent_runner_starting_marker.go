package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// agentRunnerStartingMarker is the durable record of a spawn that is in flight.
//
// It is written inside the run root by the OnSpawn hook, after the process exists
// with a complete identity and before guard lets it run, so a host that dies in
// the starting phase leaves proof that a start was in flight - evidence that a
// "never started" run root cannot produce
// (docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// section 2 Q1/Q3, section 4 I3).
type agentRunnerStartingMarker struct {
	LaunchID   string                `json:"launchId"`
	RequestID  string                `json:"requestId"`
	Identity   guard.ProcessIdentity `json:"identity"`
	StartToken string                `json:"startToken"`
	StartedAt  string                `json:"startedAt"`
}

// startingMarkerClaim is the writer side of one durable starting marker: the open
// file that holds the exclusive host lock, plus the path to remove when the start
// either completes or is rolled back. The lock is what lets a recovering scan tell
// a live host from a crashed one.
type startingMarkerClaim struct {
	path string
	file *os.File
}

// startingMarkerPath is the durable starting marker of one run root.
func startingMarkerPath(runRoot string) string {
	return filepath.Join(runRoot, managedProcessStartingFileName)
}

// writeAgentRunnerStartingMarker publishes the starting marker.
//
// The temp file is locked before it becomes visible: the marker is renamed into
// place while the host already holds the exclusive lock, so a concurrent scan can
// never observe an unlocked marker for a live host (design section 5 risk 4). The rename
// keeps the write atomic, so a reader never sees a half-written record.
func writeAgentRunnerStartingMarker(runRoot string, marker agentRunnerStartingMarker) (*startingMarkerClaim, error) {
	if strings.TrimSpace(runRoot) == "" {
		return nil, fmt.Errorf("%w: empty run root", ErrAgentRunnerProcessIdentity)
	}
	encoded, err := json.Marshal(marker)
	if err != nil {
		return nil, fmt.Errorf("%w: encode the starting marker: %v", ErrAgentRunnerProcessIdentity, err)
	}
	if err := os.MkdirAll(runRoot, 0o755); err != nil {
		return nil, fmt.Errorf("%w: create %s: %v", ErrAgentRunnerProcessIdentity, runRoot, err)
	}
	file, err := createLockedStartingMarkerFile(runRoot)
	if err != nil {
		return nil, fmt.Errorf("%w: create the starting marker: %v", ErrAgentRunnerProcessIdentity, err)
	}
	name := file.Name()
	discard := func() {
		unlockStartingMarker(file)
		_ = file.Close()
		_ = os.Remove(name)
	}
	if _, err := file.Write(encoded); err != nil {
		discard()
		return nil, fmt.Errorf("%w: write the starting marker: %v", ErrAgentRunnerProcessIdentity, err)
	}
	if err := file.Sync(); err != nil {
		discard()
		return nil, fmt.Errorf("%w: sync the starting marker: %v", ErrAgentRunnerProcessIdentity, err)
	}
	target := startingMarkerPath(runRoot)
	if err := os.Rename(name, target); err != nil {
		discard()
		return nil, fmt.Errorf("%w: publish the starting marker: %v", ErrAgentRunnerProcessIdentity, err)
	}
	return &startingMarkerClaim{path: target, file: file}, nil
}

// release drops the host lock and removes the marker. It reports whether the
// marker file was actually removed; a claim that was already released, or a
// marker another scan already cleared, reports false.
func (claim *startingMarkerClaim) release() bool {
	if claim == nil || claim.file == nil {
		return false
	}
	unlockStartingMarker(claim.file)
	_ = claim.file.Close()
	claim.file = nil
	return os.Remove(claim.path) == nil
}

// decodeAgentRunnerStartingMarker reads a decoded marker from an open file. It
// fails closed on anything a recovery decision would have to guess about.
func decodeAgentRunnerStartingMarker(file *os.File) (agentRunnerStartingMarker, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return agentRunnerStartingMarker{}, err
	}
	data, err := io.ReadAll(io.LimitReader(file, 64*1024))
	if err != nil {
		return agentRunnerStartingMarker{}, err
	}
	var marker agentRunnerStartingMarker
	if err := evidence.DecodeStrictJSON(data, &marker); err != nil {
		return agentRunnerStartingMarker{}, err
	}
	return marker, nil
}

// newAgentRunnerStartingMarker builds the marker of one spawn.
func newAgentRunnerStartingMarker(launchID, requestID string, identity guard.ProcessIdentity) agentRunnerStartingMarker {
	return agentRunnerStartingMarker{
		LaunchID:   launchID,
		RequestID:  requestID,
		Identity:   identity,
		StartToken: identity.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}
}
