//go:build windows

package guard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
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

type platformProcess struct{ job syscall.Handle }

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
	handle, err := syscall.OpenProcess(processSetQuota|processTerminate|processSuspendResume, false, uint32(pid))
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(handle)
	ok, _, callErr := procAssignProcessToJob.Call(uintptr(process.job), uintptr(handle))
	if ok == 0 {
		return callErr
	}
	status, _, _ := procNtResumeProcess.Call(uintptr(handle))
	if status != 0 {
		return fmt.Errorf("resume managed process: NTSTATUS 0x%x", status)
	}
	return nil
}

func (process *platformProcess) stop(_ ProcessIdentity, _ time.Duration) ([]string, error) {
	ok, _, callErr := procTerminateJobObject.Call(uintptr(process.job), 1)
	if ok == 0 {
		return []string{"terminate-job-object"}, callErr
	}
	return []string{"terminate-job-object"}, nil
}

func (process *platformProcess) terminate() error {
	ok, _, callErr := procTerminateJobObject.Call(uintptr(process.job), 1)
	if ok == 0 {
		return callErr
	}
	return nil
}

func (process *platformProcess) close() {
	if process.job != 0 {
		_ = syscall.CloseHandle(process.job)
		process.job = 0
	}
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
