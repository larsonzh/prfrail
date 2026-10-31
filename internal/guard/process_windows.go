//go:build windows

package guard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	jobObjectExtendedLimitInformation = 9
	jobObjectLimitKillOnJobClose      = 0x00002000
	processQueryLimitedInformation    = 0x1000
	processSetQuota                   = 0x0100
	processTerminate                  = 0x0001
	processSuspendResume              = 0x0800
	createSuspended                   = 0x00000004
	errorInvalidParameter             = syscall.Errno(87)
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW        = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject = kernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJob      = kernel32.NewProc("AssignProcessToJobObject")
	procTerminateJobObject      = kernel32.NewProc("TerminateJobObject")
	procOpenProcess             = kernel32.NewProc("OpenProcess")
	procGetProcessTimes         = kernel32.NewProc("GetProcessTimes")
	procGetExitCodeProcess      = kernel32.NewProc("GetExitCodeProcess")
	ntdll                       = syscall.NewLazyDLL("ntdll.dll")
	procNtResumeProcess         = ntdll.NewProc("NtResumeProcess")
)

type basicLimitInformation struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

type ioCounters struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

type extendedLimitInformation struct {
	BasicLimitInformation basicLimitInformation
	IoInfo                ioCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

// platformProcess owns the Job object handle. Several goroutines may stop and
// close concurrently (a Run racing two Terminate calls), so every handle use is
// serialised and close is idempotent.
type platformProcess struct {
	mu  sync.Mutex
	job syscall.Handle
	// releases counts real Job handle releases. close only ever closes a live
	// handle, so the count stops at 1; a same-package test uses closeCount to pin
	// that the handle is released exactly once.
	releases int
}

func (process *platformProcess) handle() (syscall.Handle, bool) {
	process.mu.Lock()
	defer process.mu.Unlock()
	return process.job, process.job != 0
}

func preparePlatformProcess(cmd *exec.Cmd) (platformProcess, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createSuspended}
	handle, _, callErr := procCreateJobObjectW.Call(0, 0)
	if handle == 0 {
		return platformProcess{}, callErr
	}
	job := syscall.Handle(handle)
	info := extendedLimitInformation{}
	info.BasicLimitInformation.LimitFlags = jobObjectLimitKillOnJobClose
	ok, _, callErr := procSetInformationJobObject.Call(handle, jobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if ok == 0 {
		syscall.CloseHandle(job)
		return platformProcess{}, callErr
	}
	return platformProcess{job: job}, nil
}

func (process *platformProcess) attach(pid int) error {
	job, ok := process.handle()
	if !ok {
		return errors.New("managed process job handle is closed")
	}
	handle, err := syscall.OpenProcess(processSetQuota|processTerminate|processSuspendResume, false, uint32(pid))
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(handle)
	assigned, _, callErr := procAssignProcessToJob.Call(uintptr(job), uintptr(handle))
	if assigned == 0 {
		return callErr
	}
	status, _, _ := procNtResumeProcess.Call(uintptr(handle))
	if status != 0 {
		return fmt.Errorf("resume managed process: NTSTATUS 0x%x", status)
	}
	return nil
}

func (process *platformProcess) stop(_ ProcessIdentity, _ time.Duration) ([]string, error) {
	process.mu.Lock()
	defer process.mu.Unlock()
	if process.job == 0 {
		// The handle is already released: closing the Job object already killed
		// the tree (kill-on-job-close), so the stop is proven by construction.
		return []string{"job-object-released"}, nil
	}
	ok, _, callErr := procTerminateJobObject.Call(uintptr(process.job), 1)
	if ok == 0 {
		return []string{"terminate-job-object"}, callErr
	}
	return []string{"terminate-job-object"}, nil
}

func (process *platformProcess) terminate() error {
	process.mu.Lock()
	defer process.mu.Unlock()
	if process.job == 0 {
		return nil
	}
	ok, _, callErr := procTerminateJobObject.Call(uintptr(process.job), 1)
	if ok == 0 {
		return callErr
	}
	return nil
}

func (process *platformProcess) close() {
	process.mu.Lock()
	defer process.mu.Unlock()
	if process.job != 0 {
		_ = syscall.CloseHandle(process.job)
		process.job = 0
		process.releases++
	}
}

// closeCount reports how many times the Job handle was really released.
func (process *platformProcess) closeCount() int {
	process.mu.Lock()
	defer process.mu.Unlock()
	return process.releases
}

func (process *platformProcess) verifyGone(_ context.Context) error { return nil }

func platformIdentity(pid int) (ProcessIdentity, error) {
	handle, _, callErr := procOpenProcess.Call(processQueryLimitedInformation, 0, uintptr(uint32(pid)))
	if handle == 0 {
		return ProcessIdentity{}, callErr
	}
	defer syscall.CloseHandle(syscall.Handle(handle))
	var creation, exit, kernel, user syscall.Filetime
	ok, _, callErr := procGetProcessTimes.Call(handle, uintptr(unsafe.Pointer(&creation)), uintptr(unsafe.Pointer(&exit)), uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if ok == 0 {
		return ProcessIdentity{}, callErr
	}
	token := uint64(creation.HighDateTime)<<32 | uint64(creation.LowDateTime)
	return ProcessIdentity{PID: pid, StartToken: fmt.Sprintf("%016x", token)}, nil
}

func platformIdentityAlive(identity ProcessIdentity) (bool, error) {
	handle, _, callErr := procOpenProcess.Call(processQueryLimitedInformation, 0, uintptr(uint32(identity.PID)))
	if handle == 0 {
		if errors.Is(callErr, errorInvalidParameter) {
			return false, nil
		}
		return false, callErr
	}
	defer syscall.CloseHandle(syscall.Handle(handle))
	var exitCode uint32
	ok, _, callErr := procGetExitCodeProcess.Call(handle, uintptr(unsafe.Pointer(&exitCode)))
	if ok == 0 {
		return false, callErr
	}
	if exitCode != 259 {
		return false, nil
	}
	current, err := platformIdentity(identity.PID)
	if err != nil {
		if errors.Is(err, errorInvalidParameter) {
			return false, nil
		}
		return false, err
	}
	return current.StartToken == identity.StartToken, nil
}

// stopCommandResult carries one finished taskkill invocation back to the caller
// that is still waiting for it.
type stopCommandResult struct {
	output []byte
	err    error
}

// platformStopIdentity stops a persisted process identity through taskkill.
//
// taskkill is a detached child: it is never killed when it outlives bound, because
// a half-run taskkill leaves the tree partially stopped and its result unknown. The
// goroutine hands its result to a buffered channel so a late completion is still
// reclaimed (exactly one CombinedOutput/Wait per command) without blocking forever.
// On expiry the caller gets errStopCommandExpired and must verify the tree itself.
func platformStopIdentity(ctx context.Context, identity ProcessIdentity, _, bound time.Duration) ([]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	actions := []string{"taskkill-process-tree"}
	if bound <= 0 {
		bound = identityStopCommandBound
	}
	command := exec.Command("taskkill", "/PID", strconv.Itoa(identity.PID), "/T", "/F")
	done := make(chan stopCommandResult, 1)
	go func() {
		output, err := command.CombinedOutput()
		done <- stopCommandResult{output: output, err: err}
	}()
	timer := time.NewTimer(bound)
	defer timer.Stop()
	var result stopCommandResult
	select {
	case result = <-done:
	case <-ctx.Done():
		// Same shape as an expiry: the detached command may still finish the job,
		// so this reports "unknown" rather than a failure.
		return append(actions, "taskkill-detached"), fmt.Errorf("%w: taskkill for pid %d was left running: %v", errStopCommandExpired, identity.PID, ctx.Err())
	case <-timer.C:
		return append(actions, "taskkill-detached"), fmt.Errorf("%w: taskkill did not finish for pid %d within %s", errStopCommandExpired, identity.PID, bound)
	}
	if result.err != nil {
		text := strings.ToLower(string(result.output))
		if strings.Contains(text, "not found") || strings.Contains(text, "no running instance") {
			return []string{"already-stopped"}, nil
		}
		trimmed := strings.TrimSpace(string(result.output))
		if trimmed == "" {
			return actions, result.err
		}
		return actions, fmt.Errorf("taskkill failed: %w (%s)", result.err, trimmed)
	}
	return actions, nil
}
