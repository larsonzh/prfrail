//go:build linux || windows

package adapters

// managedAgentPolicyWriteProofDomain separates the recorded proof digest from
// every other evidence digest. The digest is what the assembly anchor carries, so
// two different observations can never be compared as equal by accident.
//
// P85 design rev 4, section 2: the version moved from :1 to :2 because the
// observed vector shape changed - it is now always the FOUR fixed slots
// [policy directory | parent directory | artifact | anchor], with an absent slot
// rendered as the canonical empty marker instead of being omitted. A digest of the
// old variable-length shape must never compare equal to a new one.
//
// P85 R=7 (user ruling of 2026-11-23): the constant lives in its own file behind
// the build tag the two observing halves carry (linux and windows). The remaining
// platform file never reads it, so keeping it in the shared file made it look
// unused in the darwin IDE view. This relocation fixes that diagnostic only; it is
// not a precedent for keeping code that is otherwise dead.
const managedAgentPolicyWriteProofDomain = "proofrail:agent-runner-policy-write-guard:2\n"
