//go:build linux

package guard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type platformProcess struct{ pgid int }

func preparePlatformProcess(cmd *exec.Cmd) (platformProcess, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return platformProcess{}, nil
}

func (process *platformProcess) attach(pid int) error {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return err
	}
	if pgid != pid {
		return fmt.Errorf("unexpected process group %d for pid %d", pgid, pid)
	}
	process.pgid = pgid
	return nil
}

func (process *platformProcess) stop(_ ProcessIdentity, grace time.Duration) ([]string, error) {
	actions := []string{"signal-term-process-group"}
	if err := syscall.Kill(-process.pgid, syscall.SIGTERM); err != nil && !errors.Is(err, syscall.ESRCH) {
		return actions, err
	}
	deadline := time.Now().Add(grace)
	for time.Now().Before(deadline) {
		if !processGroupAlive(process.pgid) {
			return actions, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	if processGroupAlive(process.pgid) {
		actions = append(actions, "signal-kill-process-group")
		if err := syscall.Kill(-process.pgid, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
			return actions, err
		}
	}
	return actions, nil
}

func (process *platformProcess) terminate() error {
	if process.pgid == 0 {
		return nil
	}
	return syscall.Kill(-process.pgid, syscall.SIGKILL)
}

func (process *platformProcess) close() {}

func (process *platformProcess) verifyGone(ctx context.Context) error {
	for processGroupAlive(process.pgid) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
	return nil
}

func platformIdentity(pid int) (ProcessIdentity, error) {
	content, err := osReadProcStat(pid)
	if err != nil {
		return ProcessIdentity{}, err
	}
	return ProcessIdentity{PID: pid, StartToken: content}, nil
}

func platformIdentityAlive(identity ProcessIdentity) (bool, error) {
	token, err := osReadProcStat(identity.PID)
	if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ESRCH) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return token == identity.StartToken, nil
}

func osReadProcStat(pid int) (string, error) {
	data, err := syscall.ByteSliceFromString(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return "", err
	}
	content, err := readFileString(string(data[:len(data)-1]))
	if err != nil {
		return "", err
	}
	closeParen := strings.LastIndex(content, ") ")
	if closeParen < 0 {
		return "", errors.New("invalid /proc stat")
	}
	fields := strings.Fields(content[closeParen+2:])
	if len(fields) <= 19 {
		return "", errors.New("incomplete /proc stat")
	}
	if _, err := strconv.ParseUint(fields[19], 10, 64); err != nil {
		return "", err
	}
	return fields[19], nil
}

func processGroupAlive(pgid int) bool {
	err := syscall.Kill(-pgid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

func platformStopIdentity(identity ProcessIdentity, grace time.Duration) ([]string, error) {
	pgid, err := syscall.Getpgid(identity.PID)
	if errors.Is(err, syscall.ESRCH) {
		return []string{"already-stopped"}, nil
	}
	if err != nil {
		return []string{"signal-term-process-group"}, err
	}
	if pgid <= 0 {
		return []string{"signal-term-process-group"}, errors.New("invalid process group")
	}
	actions := []string{"signal-term-process-group"}
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil && !errors.Is(err, syscall.ESRCH) {
		return actions, err
	}
	deadline := time.Now().Add(grace)
	for time.Now().Before(deadline) {
		alive, err := platformIdentityAlive(identity)
		if errors.Is(err, syscall.ESRCH) {
			return actions, nil
		}
		if err != nil {
			return actions, err
		}
		if !alive {
			return actions, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	alive, err := platformIdentityAlive(identity)
	if errors.Is(err, syscall.ESRCH) || !alive {
		return actions, nil
	}
	if err != nil {
		return actions, err
	}
	actions = append(actions, "signal-kill-process-group")
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
		return actions, err
	}
	return actions, nil
}
