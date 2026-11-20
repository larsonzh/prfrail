package console

// T1/T2/T3 shared base: ONE resolver of a trusted ledger source.
//
// Findings F2 (T1), F3 (T2) and F6 (T3) all rest on the same question - which
// authorization-ledger file is the approval source? - so the answer is decided
// here once instead of three times:
//
//   - the source is DERIVED from a durable root the product already owns (the
//     durable run root for an AgentRunner run, the chain-configuration directory
//     for the probe), never taken from a caller-selected flag;
//   - a caller-supplied path is accepted only when it names exactly that
//     derived file; a disagreement is a refusal, not an override;
//   - when the root is missing the derivation is impossible and the refusal is
//     explicit, never a silent fallback to an in-memory copy.
//
// The helper says nothing about HOW the file is read: callers re-read it at the
// moment of use, which is the second half of the same rule.

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// trustedLedgerPath resolves the ledger the durable root implies and refuses a
// caller-selected path that disagrees with it. callerPath is expected to be
// already resolved (absolute), because the CLI layers resolve their flags against
// the invocation directory before they hand anything to the assembly.
//
// An empty callerPath keeps the operator flow unchanged: the derived location is
// the answer, and no flag is required.
func trustedLedgerPath(durableRoot, fileName, callerPath string) (string, error) {
	root := strings.TrimSpace(durableRoot)
	if root == "" {
		return "", fmt.Errorf("the trusted durable root is not configured, so the %s approval source cannot be derived", fileName)
	}
	derived, err := filepath.Abs(filepath.Join(root, fileName))
	if err != nil {
		return "", fmt.Errorf("resolve the trusted %s path: %v", fileName, err)
	}
	derived = filepath.Clean(derived)
	if strings.TrimSpace(callerPath) == "" {
		return derived, nil
	}
	selected, err := filepath.Abs(filepath.Clean(callerPath))
	if err != nil {
		return "", fmt.Errorf("resolve the caller-selected %s path: %v", fileName, err)
	}
	if !trustedPathEqual(selected, derived) {
		return "", fmt.Errorf("the caller-selected %s path %q disagrees with the trusted %q implied by %q", fileName, filepath.Clean(selected), derived, root)
	}
	return derived, nil
}

// trustedPathEqual compares two resolved paths the way the host filesystem does:
// case-insensitively on Windows, byte-exactly elsewhere.
func trustedPathEqual(left, right string) bool {
	left, right = filepath.Clean(left), filepath.Clean(right)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(left, right)
	}
	return left == right
}
