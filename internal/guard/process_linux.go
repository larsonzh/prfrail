//go:build linux

package guard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type platformProcess struct{ pgid int }

// procStat is the part of /proc/<pid>/stat this package needs: the scheduling
// state, the process group and the start-time token that keeps an identity
// unambiguous across pid reuse.
type procStat struct {
	state      byte
	group      int
	startToken string
}

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
	stat, err := readProcStat(pid)
	if err != nil {
		return ProcessIdentity{}, err
	}
	return ProcessIdentity{PID: pid, StartToken: stat.startToken}, nil
}

func platformIdentityAlive(identity ProcessIdentity) (bool, error) {
	stat, err := readProcStat(identity.PID)
	if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ESRCH) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if procStateTerminal(stat.state) {
		// An exited process that nobody has reaped yet keeps its /proc entry and its
		// original start token, but it cannot execute or write anything. Reading it as
		// alive made every stop that does not own the child handle look uncertain: the
		// identity-only stop path cannot reap the child itself, so waiting for the
		// zombie to disappear would always run out the grace period.
		return false, nil
	}
	return stat.startToken == identity.StartToken, nil
}

// procStateTerminal reports whether a /proc state letter means the process has
// finished: Z is a zombie awaiting reaping, X/x is dead.
func procStateTerminal(state byte) bool {
	return state == 'Z' || state == 'X' || state == 'x'
}

func readProcStat(pid int) (procStat, error) {
	data, err := syscall.ByteSliceFromString(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return procStat{}, err
	}
	content, err := readFileString(string(data[:len(data)-1]))
	if err != nil {
		return procStat{}, err
	}
	closeParen := strings.LastIndex(content, ") ")
	if closeParen < 0 {
		return procStat{}, errors.New("invalid /proc stat")
	}
	fields := strings.Fields(content[closeParen+2:])
	if len(fields) <= 19 {
		return procStat{}, errors.New("incomplete /proc stat")
	}
	if len(fields[0]) != 1 {
		return procStat{}, errors.New("invalid /proc stat state")
	}
	group, err := strconv.Atoi(fields[2])
	if err != nil {
		return procStat{}, err
	}
	startToken := fields[19]
	if _, err := strconv.ParseUint(startToken, 10, 64); err != nil {
		return procStat{}, err
	}
	return procStat{state: fields[0][0], group: group, startToken: startToken}, nil
}

func processGroupAlive(pgid int) bool {
	err := syscall.Kill(-pgid, 0)
	switch {
	case err == nil:
		// The group still exists, but it may hold nothing except unreaped zombies.
		return processGroupHasLiveMember(pgid)
	case errors.Is(err, syscall.EPERM):
		// The group exists but belongs to another user: /proc cannot show its members,
		// so the conservative answer is that it may still be alive.
		return true
	default:
		return false
	}
}

// processGroupHasLiveMember reports whether any member of the group can still
// run. A group whose members have all exited still answers signal 0 until their
// parent reaps them, so the signal answer alone cannot prove a stop.
func processGroupHasLiveMember(pgid int) bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		// Without /proc the signal probe is the best evidence available.
		return true
	}
	for _, entry := range entries {
		pid, convErr := strconv.Atoi(entry.Name())
		if convErr != nil || pid <= 0 {
			continue
		}
		stat, statErr := readProcStat(pid)
		if statErr != nil || stat.group != pgid {
			continue
		}
		if !procStateTerminal(stat.state) {
			return true
		}
	}
	return false
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
