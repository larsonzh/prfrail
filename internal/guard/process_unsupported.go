//go:build !linux && !windows

package guard

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type platformProcess struct{}

func preparePlatformProcess() (*platformProcess, error) {
	return nil, errors.New("managed process trees are unsupported on this platform")
}
func (process *platformProcess) start(spec *ProcessSpec) (*managedCommand, error) {
	// T3 (P83 fix round 7): darwin has no descriptor-exec reachable from the
	// standard library (no execveat and no /proc descriptor form), so an approved
	// executable cannot be launched as the object that was verified. It refuses
	// instead of silently re-resolving the path.
	if spec.ExecutableDigest != "" {
		return nil, fmt.Errorf("%w: descriptor exec is unavailable on this platform", ErrExecutableBinding)
	}
	return nil, errors.New("unsupported platform")
}

func (platformProcess) attach(int) error { return errors.New("unsupported platform") }
func (platformProcess) resume() error    { return errors.New("unsupported platform") }
func (platformProcess) stop(ProcessIdentity, time.Duration) ([]string, error) {
	return nil, errors.New("unsupported platform")
}
func (platformProcess) terminate() error { return nil }
func (platformProcess) close()           {}
func (platformProcess) closeCount() int  { return 0 }
func (platformProcess) verifyGone(context.Context) error {
	return errors.New("unsupported platform")
}
func platformIdentity(int) (ProcessIdentity, error) {
	return ProcessIdentity{}, errors.New("unsupported platform")
}
func platformIdentityAlive(ProcessIdentity) (bool, error) {
	return false, errors.New("unsupported platform")
}

func platformStopIdentity(context.Context, ProcessIdentity, time.Duration, time.Duration) ([]string, error) {
	return nil, errors.New("unsupported platform")
}
