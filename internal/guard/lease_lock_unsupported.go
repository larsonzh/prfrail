//go:build !windows && !unix

package guard

// P84 fix round 8 (Rev 3b): a platform with neither LockFileEx nor flock(2) has no way
// to make the resource lock exclusive between processes, so taking it is refused
// outright. A refusal is the fail-closed answer: the caller reports it instead of
// writing a lease record it never actually serialised.

import (
	"fmt"
	"runtime"
)

// leaseLock is a placeholder so the shared acquisition code compiles here; it can never
// be returned because createLeaseLock always fails.
type leaseLock struct{}

func createLeaseLock(string) (*leaseLock, error) {
	return nil, fmt.Errorf("lease locking is not supported on %s", runtime.GOOS)
}

func (*leaseLock) release() {}
