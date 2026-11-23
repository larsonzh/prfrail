//go:build !windows && !linux

package adapters

// P85 F21, platform stub: the install source rule is implemented for Windows and
// Linux only. Every other Unix - darwin included - refuses instead of inventing an
// equivalent rule: a probe on a platform whose immutability rule is unknown must
// fail closed, not fall back to "assume immutable".
//
// The symbol below is the one the shared files and the platform-tagged tests
// reference, so an unsupported host still compiles and every entry point refuses
// at run time. installSourceInterpreter is declared by the LINUX half only:
// nothing outside it calls it, so it is absent here on purpose (it would be
// dead code on darwin).

import (
	"fmt"
	"runtime"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// verifyInstallSourceImmutability always refuses on an unsupported platform.
//
// The signature is the same one the Windows and Linux halves carry, and leaf has
// the same meaning in all three: the handle the executable was opened through. This
// half refuses before any permission fact is observed, so the handle is accepted and
// intentionally unused - a caller holding it gets the same fail-closed verdict it
// would get by name, never a different one.
func verifyInstallSourceImmutability(executable string, leaf *evidence.BoundFile) (AIInstallSourceImmutability, error) {
	_ = leaf
	return AIInstallSourceImmutability{}, fmt.Errorf("%w: %s has no install source rule, so %q cannot be proven immutable",
		ErrAIInstallSourceProof, runtime.GOOS, executable)
}
