//go:build windows && b3bnative && !b3bvariant

package adapters

// b3bVariantDegrade is the compile-time switch for the B3b variant arm. It is false in the
// default build, so the candidate primitives run exactly as written and no runtime branch is
// added to the writer. The variant build (build tag b3bvariant) flips it; see
// agent_runner_b3b_variant_arm_test.go for what that arm is for.
//
// The windows term is not decoration: the consumer (b3bPrimitiveFor) lives in a
// *_windows_test.go file, so without it this constant would be declared and unused in every
// non-Windows build of the package (reported as an unused-const warning).
const b3bVariantDegrade = false
