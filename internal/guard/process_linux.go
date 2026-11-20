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
	"sync"
	"syscall"
	"time"
)

// platformProcess owns the process group of a managed tree. Linux has no OS
// handle to release, so close is a structural no-op; it still records the
// release so the "released exactly once" rule is observable and enforced the
// same way it is on Windows.
type platformProcess struct {
	pgid int
	mu   sync.Mutex
	// released marks that the handle ownership was given up; releases counts the
	// real transitions, so it stops at 1.
	released bool
	releases int
}

// procStat is the part of /proc/<pid>/stat this package needs: the scheduling
// state, the process group and the start-time token that keeps an identity
// unambiguous across pid reuse.
type procStat struct {
	state      byte
	group      int
	startToken string
}

// preparePlatformProcess builds the Linux ownership plan. There is no creation-time
// handle to hold on Linux: ownership is the process group recorded by start and
// verified by attach.
//
// Pdeathsig is deliberately not used here (design section 5 risk 2): the kernel signals
// it when the *creating thread* exits, not when the host process exits, so it
// cannot be the crash-recovery anchor and would only add a second, partial rule.
// The recovery anchor is the durable starting marker plus the adoption-kill scan
// (adapters.RecoverStarting); a crash therefore leaves the orphan running only
// until the next start scans the run root.
func preparePlatformProcess() (*platformProcess, error) {
	return &platformProcess{}, nil
}

// start creates the process in its own process group. Linux has no suspend step,
// so the group is the ownership anchor and the child runs from the first instant.
func (process *platformProcess) start(spec *ProcessSpec) (*managedCommand, error) {
	// T3 (P83 fix round 7): when an approved digest exists the child is created
	// from the very descriptor that was hashed, not from the name that was hashed.
	// The name is resolved exactly once - by the open below - and the spawn resolves
	// the descriptor instead: execveat(fd, "", AT_EMPTY_PATH) semantics, expressed
	// through the kernel's own descriptor form /proc/self/fd/<n>, because the Go
	// standard library exposes no execveat wrapper and this repository takes no
	// third-party dependency. A host that cannot serve that form refuses; it never
	// falls back to re-resolving the path, which would launch whatever the name
	// points at now.
	commandPath := spec.Command
	var approved *approvedExecutable
	if spec.ExecutableDigest != "" {
		if !strings.Contains(spec.Command, "/") {
			return nil, fmt.Errorf("%w: the approved executable must be named by path", ErrExecutableBinding)
		}
		opened, err := verifyApprovedExecutable(spec.Command, spec)
		if err != nil {
			return nil, err
		}
		approved = opened
		// Held until this function returns, i.e. until the child exists.
		defer func() { _ = approved.Close() }()
		if approved.file.Fd() == ^uintptr(0) {
			return nil, fmt.Errorf("%w: the verified descriptor is unavailable", ErrExecutableBinding)
		}
		// E (P83 fix round 8): the verified descriptor is handed to the child through
		// ExtraFiles instead of a raw syscall.Dup. A dup clears close-on-exec, so the
		// duplicated descriptor was inheritable by EVERY process started around the
		// same time - the verified executable leaked into unrelated concurrent
		// children. ExtraFiles marks the descriptor inheritable for this child alone
		// and places it at descriptor 3 (the first free slot above stdin/stdout/stderr),
		// so the child is still created from the very descriptor that was hashed -
		// execveat(fd, "", AT_EMPTY_PATH) semantics through the kernel's own
		// /proc/self/fd/<n> form - and no other child can see it.
		commandPath = "/proc/self/fd/3"
	}
	// argv[0] stays the name the caller wrote: the child is created from the
	// descriptor but still sees its own path as argv[0], exactly as before.
	cmd := exec.Command(spec.Command, spec.Args...)
	if approved != nil {
		// T3: only the approved branch replaces the resolved path with the verified
		// descriptor. Rewriting cmd.Path unconditionally would discard the PATH
		// lookup exec.Command already performed for an unapproved bare command.
		cmd.Path = commandPath
		cmd.ExtraFiles = []*os.File{approved.file}
	}
	cmd.Dir = spec.Dir
	cmd.Env = spec.Env
	cmd.Stdout = spec.Stdout
	cmd.Stderr = spec.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		// T3: no path fallback. Re-resolving spec.Command here would create the
		// child from whatever the name addresses at this instant, which is precisely
		// the substitution the binding exists to refuse.
		if approved != nil {
			return nil, fmt.Errorf("%w: start the verified descriptor: %v", ErrExecutableBinding, err)
		}
		return nil, err
	}
	command := &managedCommand{exitStatus: -1}
	command.Process = &managedProcessHandle{
		pid:    cmd.Process.Pid,
		killFn: func() error { return cmd.Process.Kill() },
		waitFn: func() (*os.ProcessState, error) { return cmd.Process.Wait() },
	}
	command.waitFn = func() error {
		waitErr := cmd.Wait()
		if cmd.ProcessState != nil {
			command.exitStatus = cmd.ProcessState.ExitCode()
		}
		return waitErr
	}
	return command, nil
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

// openApprovedExecutable opens the approved executable read-only without
// following a link and without blocking on a non-regular object (P83 T3). The
// descriptor is the object the spawn is created from, so the open is also the only
// name resolution that happens for an approved executable. O_NONBLOCK only
// matters for a FIFO or a device node, which the bound open refuses right after.
func openApprovedExecutable(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
}

// resume is a no-op on Linux: there is no suspended state to leave, the process
// group is already the ownership anchor, and running remains the caller's last
// startup step.
func (process *platformProcess) resume() error { return nil }

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

func (process *platformProcess) close() {
	process.mu.Lock()
	defer process.mu.Unlock()
	if process.released {
		return
	}
	process.released = true
	process.releases++
}

// closeCount reports how many times the handle ownership was released.
func (process *platformProcess) closeCount() int {
	process.mu.Lock()
	defer process.mu.Unlock()
	return process.releases
}

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

// platformStopIdentity stops a persisted identity through a process-group signal.
// The bound parameter only belongs to platforms whose stop command is an external
// program; here the grace period remains the only window, so it is ignored.
func platformStopIdentity(_ context.Context, identity ProcessIdentity, grace, _ time.Duration) ([]string, error) {
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
