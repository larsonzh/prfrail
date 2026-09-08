//go:build !linux && !windows

package guard

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

type platformProcess struct{}

func preparePlatformProcess(_ *exec.Cmd) (platformProcess, error) {
	return platformProcess{}, errors.New("managed process trees are unsupported on this platform")
}
func (platformProcess) attach(int) error { return errors.New("unsupported platform") }
func (platformProcess) stop(ProcessIdentity, time.Duration) ([]string, error) {
	return nil, errors.New("unsupported platform")
}
func (platformProcess) terminate() error { return nil }
func (platformProcess) close()           {}
func (platformProcess) verifyGone(context.Context) error {
	return errors.New("unsupported platform")
}
func platformIdentity(int) (ProcessIdentity, error) {
	return ProcessIdentity{}, errors.New("unsupported platform")
}
func platformIdentityAlive(ProcessIdentity) (bool, error) {
	return false, errors.New("unsupported platform")
}

func platformStopIdentity(ProcessIdentity, time.Duration) ([]string, error) {
	return nil, errors.New("unsupported platform")
}
