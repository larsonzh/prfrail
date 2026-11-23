//go:build !linux

package console

// P85 F02/F06/F20: the frozen-policy publication on platforms without a
// handle-relative create/link primitive.
//
// P85 F02/F06/F20 is a Linux finding: there the directory can be bound and the
// creation, the hard link, the removal and the fsync can all be expressed relative
// to that descriptor, so the fixed shape lives in
// agent_runner_policy_publish_linux.go. This half keeps the path-based steps:
//
//   - Windows cannot create or hard-link relative to a directory handle with the
//     standard library, so the final-name resolution stays path-based. P85 R7
//     registers that residue as an INDEPENDENT S1 boundary of the user-mode path
//     publications (user ruling of 2026-11-23): it is not an F04 exemption and must
//     not be cited as one. The identity is re-proved immediately before and after
//     the path operations, so a directory swapped in between refuses instead of
//     receiving the artifact; the creation and link steps themselves are the
//     registered remainder;
//   - darwin and the remaining unix targets refuse directory binding outright, so
//     the verified binding never produces a handle here and this function is
//     unreachable on them.
//
// The identity is still asserted around the path operations: the name the caller
// derived must still resolve to the object the binding produced, before and after
// the publication, so a directory swapped in between refuses instead of receiving
// the artifact. The no-replace discipline is identical on both halves (F09: no
// os.Rename fallback anywhere).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// agentRunnerPolicyTemporaryPattern is the same-prefix pattern the handle-based
	// half uses, so a leftover from either shape is recognizable.
	agentRunnerPolicyTemporaryPattern = ".agent-runner-policy-*"
)

// publishAgentRunnerPolicyFileIn publishes one payload without replacing an
// existing name, asserting that directory still addresses the bound object.
func publishAgentRunnerPolicyFileIn(bound *evidence.BoundDir, fileID evidence.FileID, directory, base string, payload []byte, label string) error {
	path := filepath.Join(directory, base)
	if err := agentRunnerAssertDirectoryIdentity(bound, fileID, directory, label); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, agentRunnerPolicyTemporaryPattern)
	if err != nil {
		return fmt.Errorf("%w: create the %s temporary file: %v", ErrInvalidConfig, label, err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	writeErr := temporary.Chmod(0o600)
	if writeErr == nil {
		_, writeErr = temporary.Write(payload)
	}
	if writeErr == nil {
		writeErr = temporary.Sync()
	}
	if closeErr := temporary.Close(); writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		return fmt.Errorf("%w: write the %s temporary file: %v", ErrInvalidConfig, label, writeErr)
	}
	if err := os.Link(temporaryPath, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: the %s %s already exists; a frozen policy is never replaced", ErrInvalidConfig, label, path)
		}
		return fmt.Errorf("%w: publish the %s %s without replacing it: %v", ErrInvalidConfig, label, path, err)
	}
	// The publication itself is the durability point on Windows (no directory
	// handle can be flushed through Go), and the identity is asserted again so a
	// swap during the path operations is visible instead of silently accepted.
	if err := agentRunnerAssertDirectoryIdentity(bound, fileID, directory, label); err != nil {
		return err
	}
	if err := bound.SyncDirectory(); err != nil {
		return fmt.Errorf("%w: sync the %s directory %s: %v", ErrInvalidConfig, label, directory, err)
	}
	return nil
}

// agentRunnerAssertDirectoryIdentity proves the name the caller derived still
// resolves to the object the verified binding produced.
//
// P85 design rev 4, section 3: the re-check is a RE-BIND through the same
// root-chain walk the binding uses (evidence.OpenBoundDirVerified), not a
// name-based evidence.FileIDAt. A junction or symbolic link introduced in the
// directory path after the bind is therefore refused by the walk instead of being
// followed to whatever object it points at, which is the shape that made the old
// comparison agree with itself (the registered path-based publication boundary
// cannot back this position).
func agentRunnerAssertDirectoryIdentity(bound *evidence.BoundDir, fileID evidence.FileID, directory, label string) error {
	if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
		return fmt.Errorf("%w: the %s directory %s is no longer the approved object: %v", ErrInvalidConfig, label, directory, err)
	}
	verified, verifiedID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		return fmt.Errorf("%w: the %s directory %s can no longer be re-bound on its root chain: %v", ErrInvalidConfig, label, directory, err)
	}
	defer func() { _ = verified.Close() }()
	if verifiedID != fileID {
		return fmt.Errorf("%w: the %s directory %s no longer names the object that was bound (%s != %s)", ErrInvalidConfig, label, directory, verifiedID, fileID)
	}
	return nil
}
