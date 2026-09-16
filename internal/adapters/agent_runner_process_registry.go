package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

var (
	// ErrAgentRunnerProcessDuplicate reports a second registration of one launch.
	ErrAgentRunnerProcessDuplicate = errors.New("AgentRunner process already registered")
	// ErrAgentRunnerProcessMissing reports an unknown launch identity.
	ErrAgentRunnerProcessMissing = errors.New("AgentRunner process not registered")
	// ErrAgentRunnerProcessIdentity reports a malformed or unencodable process
	// identity. An identity that cannot be proven must never be published.
	ErrAgentRunnerProcessIdentity = errors.New("AgentRunner process identity invalid")
)

// agentRunnerProcessPhase is the in-memory lifecycle phase of one launch. A
// phase is not durable state: after a restart the persisted launch receipts are
// the only authority, and this registry only serves the live process.
type agentRunnerProcessPhase string

const (
	agentRunnerProcessStarting agentRunnerProcessPhase = "starting"
	agentRunnerProcessRunning  agentRunnerProcessPhase = "running"
	agentRunnerProcessStopping agentRunnerProcessPhase = "stopping"
	agentRunnerProcessStopped  agentRunnerProcessPhase = "stopped"
)

// managedProcessIdentityFileName is the operator stop path's identity mirror. It
// is written inside the durable run root, which is exactly where the console
// stop command looks for it, so operator stop needs no new mechanism.
const managedProcessIdentityFileName = "managed-process.identity.json"

// managedProcessLaunchSlotFileName is the durable launch slot of one run root. It
// exists only while a launch owns the run root, so the ownership is decided by an
// atomic exclusive create rather than by a check followed by a write.
const managedProcessLaunchSlotFileName = "managed-process.launch.slot"

// agentRunnerProcessHandle is the live handle of one spawned process. The log
// files stay open until the run ends, so the launcher owns them and the handle
// closes them exactly once.
type agentRunnerProcessHandle struct {
	launchID  string
	requestID string
	process   *guard.ManagedProcess
	phase     agentRunnerProcessPhase
	startedAt time.Time
	grace     time.Duration
	// slotToken is the launch slot claim this handle owns, so its release can never
	// remove a slot another launch has taken in the meantime.
	slotToken string
	logs      []*os.File
}

// closeLogs closes the log artifacts once. A second call is a no-op.
func (handle *agentRunnerProcessHandle) closeLogs() {
	if handle == nil {
		return
	}
	for _, file := range handle.logs {
		if file != nil {
			_ = file.Close()
		}
	}
	handle.logs = nil
}

// agentRunnerProcessRegistry owns the live handles of one launcher. It never
// publishes decisions and never reads or writes store records.
type agentRunnerProcessRegistry struct {
	mu      sync.Mutex
	handles map[string]*agentRunnerProcessHandle
}

func newAgentRunnerProcessRegistry() *agentRunnerProcessRegistry {
	return &agentRunnerProcessRegistry{handles: map[string]*agentRunnerProcessHandle{}}
}

func (registry *agentRunnerProcessRegistry) register(handle *agentRunnerProcessHandle) error {
	if registry == nil {
		return fmt.Errorf("%w: nil registry", ErrAgentRunnerProcessMissing)
	}
	if handle == nil || !evidence.ValidID(handle.launchID) || !evidence.ValidID(handle.requestID) {
		return fmt.Errorf("%w: invalid handle binding", ErrAgentRunnerProcessIdentity)
	}
	if handle.process == nil {
		return fmt.Errorf("%w: handle without a managed process", ErrAgentRunnerProcessIdentity)
	}
	if handle.phase != agentRunnerProcessStarting && handle.phase != agentRunnerProcessRunning {
		return fmt.Errorf("%w: handle registered in phase %s", ErrAgentRunnerProcessIdentity, handle.phase)
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, exists := registry.handles[handle.launchID]; exists {
		return fmt.Errorf("%w: launch %s", ErrAgentRunnerProcessDuplicate, handle.launchID)
	}
	registry.handles[handle.launchID] = handle
	return nil
}

func (registry *agentRunnerProcessRegistry) lookup(launchID string) (*agentRunnerProcessHandle, bool) {
	if registry == nil {
		return nil, false
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	handle, found := registry.handles[launchID]
	return handle, found
}

// live reports whether any launch is still registered. The launcher refuses a
// second one, because the durable identity mirror is unique per run root and a
// second launch would overwrite the identity of the process that is still running.
func (registry *agentRunnerProcessRegistry) live() (string, bool) {
	if registry == nil {
		return "", false
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	for launchID := range registry.handles {
		return launchID, true
	}
	return "", false
}

// advance moves one live handle forward. `stopped` is terminal, so a handle that
// already reached it may never move again: a stopped process is never restarted.
func (registry *agentRunnerProcessRegistry) advance(launchID string, phase agentRunnerProcessPhase) error {
	if registry == nil {
		return fmt.Errorf("%w: nil registry", ErrAgentRunnerProcessMissing)
	}
	if !knownAgentRunnerProcessPhase(phase) {
		return fmt.Errorf("%w: unknown phase %s", ErrAgentRunnerProcessIdentity, phase)
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	handle, found := registry.handles[launchID]
	if !found {
		return fmt.Errorf("%w: launch %s", ErrAgentRunnerProcessMissing, launchID)
	}
	if handle.phase == agentRunnerProcessStopped && phase != agentRunnerProcessStopped {
		return fmt.Errorf("%w: launch %s already stopped", ErrAgentRunnerProcessIdentity, launchID)
	}
	handle.phase = phase
	return nil
}

func (registry *agentRunnerProcessRegistry) release(launchID string) {
	if registry == nil {
		return
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	delete(registry.handles, launchID)
}

func (registry *agentRunnerProcessRegistry) len() int {
	if registry == nil {
		return 0
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	return len(registry.handles)
}

func knownAgentRunnerProcessPhase(phase agentRunnerProcessPhase) bool {
	switch phase {
	case agentRunnerProcessStarting, agentRunnerProcessRunning, agentRunnerProcessStopping, agentRunnerProcessStopped:
		return true
	default:
		return false
	}
}

// agentRunnerProcessID encodes a guard process identity as the opaque ProcessID
// a launch receipt carries. The same canonical bytes are mirrored to disk for
// the operator stop path, so both readers see one representation.
func agentRunnerProcessID(identity guard.ProcessIdentity) (string, error) {
	if err := validateAgentRunnerProcessIdentity(identity); err != nil {
		return "", err
	}
	encoded, err := json.Marshal(identity)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAgentRunnerProcessIdentity, err)
	}
	canonical, err := evidence.Canonicalize(encoded)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAgentRunnerProcessIdentity, err)
	}
	return string(canonical), nil
}

// decodeAgentRunnerProcessID reverses agentRunnerProcessID and fails closed on
// anything a stop would have to guess about.
func decodeAgentRunnerProcessID(value string) (guard.ProcessIdentity, error) {
	var identity guard.ProcessIdentity
	if strings.TrimSpace(value) == "" {
		return identity, fmt.Errorf("%w: empty process id", ErrAgentRunnerProcessIdentity)
	}
	if err := evidence.DecodeStrictJSON([]byte(value), &identity); err != nil {
		return guard.ProcessIdentity{}, fmt.Errorf("%w: %v", ErrAgentRunnerProcessIdentity, err)
	}
	if err := validateAgentRunnerProcessIdentity(identity); err != nil {
		return guard.ProcessIdentity{}, err
	}
	return identity, nil
}

func validateAgentRunnerProcessIdentity(identity guard.ProcessIdentity) error {
	if identity.PID <= 0 || strings.TrimSpace(identity.StartToken) == "" {
		return fmt.Errorf("%w: incomplete process identity", ErrAgentRunnerProcessIdentity)
	}
	return nil
}

// writeAgentRunnerProcessIdentityMirror publishes the identity where the console
// stop path reads it. The write is atomic (temp file plus rename) so a reader
// never observes a half-written identity.
func writeAgentRunnerProcessIdentityMirror(runRoot string, identity guard.ProcessIdentity) error {
	if strings.TrimSpace(runRoot) == "" {
		return fmt.Errorf("%w: empty run root", ErrAgentRunnerProcessIdentity)
	}
	encoded, err := agentRunnerProcessID(identity)
	if err != nil {
		return err
	}
	target := filepath.Join(runRoot, managedProcessIdentityFileName)
	if err := os.MkdirAll(runRoot, 0o755); err != nil {
		return fmt.Errorf("%w: create %s: %v", ErrAgentRunnerProcessIdentity, runRoot, err)
	}
	temporary, err := os.CreateTemp(runRoot, "identity-*.tmp")
	if err != nil {
		return fmt.Errorf("%w: temp file in %s: %v", ErrAgentRunnerProcessIdentity, runRoot, err)
	}
	name := temporary.Name()
	if _, err := temporary.WriteString(encoded); err != nil {
		_ = temporary.Close()
		_ = os.Remove(name)
		return fmt.Errorf("%w: write %s: %v", ErrAgentRunnerProcessIdentity, name, err)
	}
	if err := temporary.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("%w: close %s: %v", ErrAgentRunnerProcessIdentity, name, err)
	}
	if err := os.Rename(name, target); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("%w: publish %s: %v", ErrAgentRunnerProcessIdentity, target, err)
	}
	return nil
}
