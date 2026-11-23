//go:build !windows && !linux

package adapters

// F01, remaining platforms: there is no reviewed mutation-rights primitive here,
// so the layout is unprovable. An unprovable layout is a refusal - the assembly
// must never dispatch on "probably fine".

import "fmt"

// managedAgentPolicyDomainWriteProofSlots is the fail-closed half for every
// platform without a reviewed mutation-rights primitive (P85 F01, design rev 3
// section 1; shape per design rev 4 section 2).
func managedAgentPolicyDomainWriteProofSlots(slots []managedAgentPolicyDomainSlot) (string, error) {
	return "", fmt.Errorf("%w: no managed-identity write proof exists on this platform, so the policy domain of %d observed slot(s) cannot be proven closed", ErrManagedAgentPolicyWriteProof, len(slots))
}
