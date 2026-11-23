//go:build windows

package guard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// The Windows ownership anchor is the Job object. It is created before the child
// exists and handed to CreateProcess through PROC_THREAD_ATTRIBUTE_JOB_LIST, so
// the child is inside the job at the instant it is created - before any code of
// the host runs again. With CREATE_SUSPENDED the child cannot execute a single
// instruction before guard resumes it, which is what makes "ownership before
// execution" a structural property instead of a timing assumption
// (docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// section 2 Q2). EXTENDED_STARTUPINFO_PRESENT is what makes the attribute list apply;
// without it CreateProcess silently ignores lpAttributeList.
const (
	jobObjectExtendedLimitInformation = 9
	jobObjectLimitKillOnJobClose      = 0x00002000
	processQueryLimitedInformation    = 0x1000
	createSuspended                   = 0x00000004
	createUnicodeEnvironment          = 0x00000400
	createExtendedStartupInfoPresent  = 0x00080000
	startfUseStdHandles               = 0x00000100
	procThreadAttributeJobList        = 0x0002000D
	procThreadAttributeHandleList     = 0x00020002
	handleFlagInherit                 = 0x00000001
	errorInvalidParameter             = syscall.Errno(87)
)

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW          = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject   = kernel32.NewProc("SetInformationJobObject")
	procTerminateJobObject        = kernel32.NewProc("TerminateJobObject")
	procOpenProcess               = kernel32.NewProc("OpenProcess")
	procGetProcessTimes           = kernel32.NewProc("GetProcessTimes")
	procGetExitCodeProcess        = kernel32.NewProc("GetExitCodeProcess")
	procIsProcessInJob            = kernel32.NewProc("IsProcessInJob")
	procInitializeProcThreadAttrs = kernel32.NewProc("InitializeProcThreadAttributeList")
	procUpdateProcThreadAttribute = kernel32.NewProc("UpdateProcThreadAttribute")
	procDeleteProcThreadAttrs     = kernel32.NewProc("DeleteProcThreadAttributeList")
	procCreateProcessW            = kernel32.NewProc("CreateProcessW")
	ntdll                         = syscall.NewLazyDLL("ntdll.dll")
	procNtResumeProcess           = ntdll.NewProc("NtResumeProcess")
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

// startupInfoEX is STARTUPINFOEXW: the plain STARTUPINFO followed by the
// attribute list pointer. syscall.StartupInfo has exactly the layout of the C
// struct (104 bytes on amd64, 112 for the extended form).
type startupInfoEX struct {
	StartupInfo   syscall.StartupInfo
	AttributeList uintptr
}

type processInformation struct {
	Process   syscall.Handle
	Thread    syscall.Handle
	ProcessID uint32
	ThreadID  uint32
}

// platformProcess owns the Job object handle. Several goroutines may stop and close
// concurrently (a Run racing two Terminate calls), so every handle use is
// serialised and close is idempotent.
type platformProcess struct {
	mu  sync.Mutex
	job syscall.Handle
	// process is the handle of the spawned process, captured only so resume can
	// reach it; the process handle itself is closed with the platform handle.
	process syscall.Handle
	// releases counts real Job handle releases. close only ever closes a live
	// handle, so the count stops at 1; a same-package test uses closeCount to pin
	// that the handle is released exactly once.
	releases int
	// envClearHook is an unexported test seam (F7): it observes the widened
	// environment block after the deferred clear ran, on both the success and the
	// failure path. Production code never sets it.
	envClearHook func(block []uint16)
}

func (process *platformProcess) jobHandle() (syscall.Handle, bool) {
	process.mu.Lock()
	defer process.mu.Unlock()
	return process.job, process.job != 0
}

func (process *platformProcess) processHandle() (syscall.Handle, bool) {
	process.mu.Lock()
	defer process.mu.Unlock()
	return process.process, process.process != 0
}

// preparePlatformProcess builds the ownership anchor: the Job object the child is
// created in. The creation attributes themselves are built per spawn, because they
// also carry the exact list of handles the child inherits.
func preparePlatformProcess() (*platformProcess, error) {
	handle, _, callErr := procCreateJobObjectW.Call(0, 0)
	if handle == 0 {
		return nil, callErr
	}
	job := syscall.Handle(handle)
	info := extendedLimitInformation{}
	info.BasicLimitInformation.LimitFlags = jobObjectLimitKillOnJobClose
	ok, _, callErr := procSetInformationJobObject.Call(handle, jobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if ok == 0 {
		syscall.CloseHandle(job)
		return nil, callErr
	}
	return &platformProcess{job: job}, nil
}

// buildCreationAttributes builds the STARTUPINFOEX attribute list of one spawn:
//
//   - PROC_THREAD_ATTRIBUTE_JOB_LIST makes the child a member of the managed job at
//     the instant it is created;
//   - PROC_THREAD_ATTRIBUTE_HANDLE_LIST restricts inheritance to exactly the three
//     standard handles of this spawn. Without it bInheritHandles would hand the
//     child every inheritable handle in the host - including the pipe ends of a
//     spawn running concurrently in another goroutine - and a child that keeps a
//     stranger's pipe open would keep that spawn's reader from ever seeing
//     end-of-file. os/exec makes the same guarantee the same way.
func buildCreationAttributes(job syscall.Handle, inherited []syscall.Handle) (uintptr, []byte, error) {
	const attributeCount = 2
	var size uintptr
	_, _, _ = procInitializeProcThreadAttrs.Call(0, attributeCount, 0, uintptr(unsafe.Pointer(&size)))
	if size == 0 {
		return 0, nil, errors.New("initialize the process thread attribute list: no size")
	}
	buffer := make([]byte, size)
	list := uintptr(unsafe.Pointer(&buffer[0]))
	if ok, _, callErr := procInitializeProcThreadAttrs.Call(list, attributeCount, 0, uintptr(unsafe.Pointer(&size))); ok == 0 {
		return 0, nil, fmt.Errorf("initialize the process thread attribute list: %w", callErr)
	}
	jobValue := uintptr(job)
	if ok, _, callErr := procUpdateProcThreadAttribute.Call(
		list,
		0,
		procThreadAttributeJobList,
		uintptr(unsafe.Pointer(&jobValue)),
		unsafe.Sizeof(jobValue),
		0,
		0,
	); ok == 0 {
		procDeleteProcThreadAttrs.Call(list)
		return 0, nil, fmt.Errorf("attach the job list attribute: %w", callErr)
	}
	handleList := make([]syscall.Handle, len(inherited))
	copy(handleList, inherited)
	if len(handleList) > 0 {
		if ok, _, callErr := procUpdateProcThreadAttribute.Call(
			list,
			0,
			procThreadAttributeHandleList,
			uintptr(unsafe.Pointer(&handleList[0])),
			uintptr(len(handleList))*unsafe.Sizeof(handleList[0]),
			0,
			0,
		); ok == 0 {
			procDeleteProcThreadAttrs.Call(list)
			return 0, nil, fmt.Errorf("attach the handle list attribute: %w", callErr)
		}
	}
	return list, buffer, nil
}

// start creates the process suspended and already inside the managed job, with
// its standard streams wired to the spec writers.
func (process *platformProcess) start(spec *ProcessSpec) (*managedCommand, error) {
	job, ok := process.jobHandle()
	if !ok {
		return nil, errors.New("managed process job handle is closed")
	}
	// T3 (P83 fix round 7): an approved digest names ONE object, and an approved
	// object is only addressable by a path. A bare command name would be resolved
	// through PATH at spawn time - a lookup that happens later and can select a file
	// the approval never saw - so it is refused here instead of being resolved
	// behind the approval's back.
	var approved *approvedExecutable
	executable := ""
	if spec.ExecutableDigest != "" {
		if !strings.ContainsAny(spec.Command, `\/`) {
			return nil, fmt.Errorf("%w: the approved executable must be named by path", ErrExecutableBinding)
		}
		executable = spec.Command
		// T3: the handle is opened with a share mode that denies write and delete,
		// hashed through that handle and held across the whole verify -> spawn
		// sequence below. While it is held the image can neither be rewritten in
		// place nor replaced, so the path CreateProcessW resolves is the verified
		// object. A handle that cannot be opened exclusively is a refusal, never a
		// spawn by path.
		var openErr error
		approved, openErr = verifyApprovedExecutable(executable, spec)
		if openErr != nil {
			return nil, openErr
		}
		// The defer runs when this function returns, which is after CreateProcessW.
		defer func() { _ = approved.Close() }()
	} else {
		var resolveErr error
		executable, resolveErr = resolveManagedExecutable(spec.Command)
		if resolveErr != nil {
			return nil, resolveErr
		}
	}
	commandLine, err := buildCommandLine(spec.Command, spec.Args)
	if err != nil {
		return nil, err
	}
	environment, err := buildEnvironmentBlock(spec.Env)
	if err != nil {
		return nil, err
	}
	// F7: the block owns secrets in UTF-16 form, so it is cleared on every exit path
	// of this function - a refusal below, a failed CreateProcessW, and the success
	// path alike. clear() zeroes each element, so no copy of the secrets survives.
	//
	// The two defers are ordered so the observation seam runs AFTER the clear (LIFO)
	// and therefore sees the zeros the contract promises.
	defer func() {
		if process.envClearHook != nil {
			process.envClearHook(environment)
		}
	}()
	defer clear(environment)
	var environmentPtr *uint16
	if len(environment) > 0 {
		environmentPtr = &environment[0]
	}
	currentDirectory, err := optionalUTF16(spec.Dir)
	if err != nil {
		return nil, err
	}
	handles, err := newChildStdHandles(spec)
	if err != nil {
		return nil, err
	}
	// The child ends are inheritable, but only they are inherited (handle list); the
	// parent copies are closed as soon as the process exists, so the readers see
	// end-of-file when the child tree is gone.
	defer handles.closeParent()
	attributeList, attributeStorage, err := buildCreationAttributes(job, []syscall.Handle{handles.stdin, handles.stdout, handles.stderr})
	if err != nil {
		return nil, err
	}
	defer func() {
		procDeleteProcThreadAttrs.Call(attributeList)
		runtime.KeepAlive(attributeStorage)
	}()

	executablePtr, err := syscall.UTF16PtrFromString(executable)
	if err != nil {
		return nil, err
	}
	startup := startupInfoEX{}
	startup.StartupInfo.Cb = uint32(unsafe.Sizeof(startup))
	startup.StartupInfo.Flags = startfUseStdHandles
	startup.StartupInfo.StdInput = handles.stdin
	startup.StartupInfo.StdOutput = handles.stdout
	startup.StartupInfo.StdErr = handles.stderr
	startup.AttributeList = attributeList
	var procInfo processInformation
	created, _, callErr := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(executablePtr)),
		uintptr(unsafe.Pointer(commandLine)),
		0,
		0,
		1,
		uintptr(uint32(createSuspended|createUnicodeEnvironment|createExtendedStartupInfoPresent)),
		uintptr(unsafe.Pointer(environmentPtr)),
		uintptr(unsafe.Pointer(currentDirectory)),
		uintptr(unsafe.Pointer(&startup)),
		uintptr(unsafe.Pointer(&procInfo)),
	)
	if created == 0 {
		return nil, callErr
	}
	// The thread handle is not needed: resume goes through the process handle.
	syscall.CloseHandle(procInfo.Thread)
	process.mu.Lock()
	process.process = procInfo.Process
	process.mu.Unlock()

	command := &managedCommand{exitStatus: -1}
	// The closures read the process handle through the platform handle, so a handle
	// the platform already released can never be used again: Windows recycles handle
	// values, and a stale value could name an unrelated process.
	command.Process = &managedProcessHandle{
		pid: int(procInfo.ProcessID),
		killFn: func() error {
			handle, ok := process.processHandle()
			if !ok {
				return ErrInvalidProcess
			}
			return syscall.TerminateProcess(handle, 1)
		},
		waitFn: func() (*os.ProcessState, error) {
			handle, ok := process.processHandle()
			if !ok {
				return nil, ErrInvalidProcess
			}
			_, waitErr := syscall.WaitForSingleObject(handle, syscall.INFINITE)
			return nil, waitErr
		},
	}
	command.waitFn = func() error {
		handle, ok := process.processHandle()
		if !ok {
			handles.drain()
			return ErrInvalidProcess
		}
		_, waitErr := syscall.WaitForSingleObject(handle, syscall.INFINITE)
		var exitCode uint32
		if err := syscall.GetExitCodeProcess(handle, &exitCode); err == nil {
			command.exitStatus = int(exitCode)
		}
		// The output pipes end when the child tree is gone, so a caller that reads
		// the log artifacts after the run sees the complete stream.
		handles.drain()
		return waitErr
	}
	return command, nil
}

// childStdHandles holds the three standard handles of the child plus the parent
// side of the output pipes.
type childStdHandles struct {
	stdin  syscall.Handle
	stdout syscall.Handle
	stderr syscall.Handle
	// files are the parent-side pipe ends. They are closed through their *os.File:
	// closing the raw handle directly would leave the finalizer free to close the
	// same numeric handle after Windows had recycled it.
	files []*os.File
	// handles are the raw null-device handles this spawn opened itself.
	handles []syscall.Handle
	readers []pipeReader
	copies  sync.WaitGroup
}

// pipeReader is one parent-side output pipe that is copied into the spec writer.
type pipeReader struct {
	target io.Writer
	read   *os.File
}

// newChildStdHandles wires the standard streams. A nil writer means the child
// writes to NUL, exactly like an unset exec.Cmd stream: nothing is buffered, so a
// readable process never blocks on a full pipe nobody drains.
func newChildStdHandles(spec *ProcessSpec) (*childStdHandles, error) {
	handles := &childStdHandles{}
	nul, err := openNullHandle()
	if err != nil {
		return nil, err
	}
	handles.stdin = nul
	handles.handles = append(handles.handles, nul)
	for _, stream := range []struct {
		writer io.Writer
		child  *syscall.Handle
	}{
		{writer: spec.Stdout, child: &handles.stdout},
		{writer: spec.Stderr, child: &handles.stderr},
	} {
		if stream.writer == nil {
			quiet, openErr := openNullHandle()
			if openErr != nil {
				handles.closeParent()
				return nil, openErr
			}
			*stream.child = quiet
			handles.handles = append(handles.handles, quiet)
			continue
		}
		read, write, pipeErr := os.Pipe()
		if pipeErr != nil {
			handles.closeParent()
			return nil, pipeErr
		}
		child := syscall.Handle(write.Fd())
		if inheritErr := syscall.SetHandleInformation(child, handleFlagInherit, handleFlagInherit); inheritErr != nil {
			_ = read.Close()
			_ = write.Close()
			handles.closeParent()
			return nil, inheritErr
		}
		*stream.child = child
		handles.files = append(handles.files, write)
		handles.readers = append(handles.readers, pipeReader{target: stream.writer, read: read})
	}
	for _, reader := range handles.readers {
		handles.copies.Add(1)
		go func(reader pipeReader) {
			defer handles.copies.Done()
			_, _ = io.Copy(reader.target, reader.read)
			_ = reader.read.Close()
		}(reader)
	}
	return handles, nil
}

// closeParent closes the parent copies of the handles the child inherits. It is
// idempotent, and it must run as soon as CreateProcess returned: as long as the
// parent holds a write end, a reader never sees end-of-file.
func (handles *childStdHandles) closeParent() {
	if handles == nil {
		return
	}
	for _, file := range handles.files {
		_ = file.Close()
	}
	handles.files = nil
	for _, handle := range handles.handles {
		syscall.CloseHandle(handle)
	}
	handles.handles = nil
}

// drain waits for the output copies to finish.
func (handles *childStdHandles) drain() {
	if handles == nil {
		return
	}
	handles.copies.Wait()
}

// openNullHandle returns an inheritable handle to the null device.
func openNullHandle() (syscall.Handle, error) {
	handle, err := syscall.Open("NUL", syscall.O_RDWR, 0)
	if err != nil {
		return 0, err
	}
	if err := syscall.SetHandleInformation(handle, handleFlagInherit, handleFlagInherit); err != nil {
		syscall.CloseHandle(handle)
		return 0, err
	}
	return handle, nil
}

// resolveManagedExecutable mirrors exec.Command: a bare name is looked up so the
// launcher keeps working with both an absolute path and a PATH-resolved command.
func resolveManagedExecutable(name string) (string, error) {
	if strings.ContainsAny(name, `\/`) {
		return name, nil
	}
	return exec.LookPath(name)
}

// approvedExecutableOpenReparsePoint is FILE_FLAG_OPEN_REPARSE_POINT: the approved
// name is opened as itself, so the object class is read from the handle the spawn
// will use instead of from whatever a redirect points at (P85 F04).
const approvedExecutableOpenReparsePoint = 0x00200000

// openApprovedExecutable opens the approved executable with FILE_SHARE_READ and
// nothing else (P83 T3).
//
// FILE_SHARE_READ denies every other opener the write and delete access it would
// need to rewrite the image in place or to replace it by a rename, so as long as
// this handle is open the name keeps resolving to the object that was hashed. The
// handle is deliberately not FILE_FLAG_BACKUP_SEMANTICS: the approved object has
// to be a file, and a directory is refused by the open itself.
//
// P85 F04 (final start-name binding): the name is opened with
// FILE_FLAG_OPEN_REPARSE_POINT and the object class is read from THAT handle. A
// reparse point at the approved name (a symlink or junction standing in for the
// executable) is refused instead of being followed, and so is a directory. The
// kernel still re-resolves the whole name chain inside CreateProcessW; that window
// is a registered residual, narrowed on the proofrail side by the ancestor-chain
// proof (adapters: every component judged, no reparse point anywhere on the chain).
func openApprovedExecutable(path string) (*os.File, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, syscall.FILE_SHARE_READ, nil, syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL|approvedExecutableOpenReparsePoint, 0)
	if err != nil {
		return nil, err
	}
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		syscall.CloseHandle(handle)
		return nil, err
	}
	if info.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0 || info.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
		syscall.CloseHandle(handle)
		return nil, fmt.Errorf("%w: the approved executable name resolves to a reparse point or a directory", ErrExecutableBinding)
	}
	return os.NewFile(uintptr(handle), path), nil
}

// buildCommandLine builds the command line the same way os/exec does: every
// argument escaped individually. argv[0] stays the name as the caller wrote it.
func buildCommandLine(name string, args []string) (*uint16, error) {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, syscall.EscapeArg(name))
	for _, arg := range args {
		parts = append(parts, syscall.EscapeArg(arg))
	}
	return syscall.UTF16PtrFromString(strings.Join(parts, " "))
}

// environmentBlockInitialCapacity is the capacity rule of the environment block:
// eight UTF-16 units per entry plus the terminating NUL.
//
// T7 (P83 fix round 6): the rule is a function so a test can compute the exact
// capacity and drive the builder to the boundary deterministically - one entry
// whose encoded length equals this capacity reaches the terminator-headroom branch
// without relying on a lucky combination of entries. Before this, the branch was
// only reachable by accident and no counterexample exercised it.
func environmentBlockInitialCapacity(entries int) int {
	return entries*8 + 1
}

// buildEnvironmentBlock encodes the child environment. A nil slice means "inherit
// the host environment", which is what an unset exec.Cmd environment does.
//
// F7 (CONTRACTS.md:216 "the adapter only ever receives a resolved, short-lived
// byte copy and clears it after use"): the block is returned as a slice the caller
// owns rather than as a raw pointer, because the block carries the child's
// secrets and the caller has to clear it. The previous signature returned
// &block[0], which had no owner: the buffer stayed live, uncleared, on both the
// success and every failure path.
//
// P83 FIX ROUND 4 (Item D): the per-entry UTF-16 slices this function allocates are
// intermediates that also carry secrets, and growing the block abandons a backing
// array that still holds them. Both are now owned by this function and cleared by
// one deferred helper, so no exit path - success, conversion error, or a grown
// block - can leave a readable copy behind.
//
// T7 (P83 fix round 6): every abandoned array and the partial block of the failure
// path is handed to the observer seam AFTER it was zeroed, so a test can hold the
// backing array and assert the clearing instead of merely asserting that the call
// returned an error.
func buildEnvironmentBlock(env []string) ([]uint16, error) {
	if env == nil {
		return nil, nil
	}
	scratch := make([][]uint16, 0, len(env))
	// The closure captures scratch by reference: a deferred call would snapshot the
	// empty slice header and clear nothing.
	defer func() { clearEnvironmentScratch(scratch) }()
	block := make([]uint16, 0, environmentBlockInitialCapacity(len(env)))
	for _, entry := range env {
		converted, err := syscall.UTF16FromString(entry)
		if err != nil {
			// The partially built block already holds secrets; clear it before the
			// error returns. The deferred helper clears every entry converted so far.
			clear(block)
			// T7: the zeroed partial block is reported through the same seam the growth
			// paths use, so the failure path has a direct counterexample.
			reportEnvironmentBlockScratch("partial-block", block)
			return nil, err
		}
		scratch = append(scratch, converted)
		// The seam observes the entry before the deferred clear runs, so a test can
		// hold this very backing array and assert after the call that it was zeroed.
		reportEnvironmentBlockScratch("entry", converted)
		if needed := len(block) + len(converted); needed > cap(block) {
			grown := make([]uint16, len(block), needed*2)
			copy(grown, block)
			// The previous backing array still carries the entries written so far, so
			// it is cleared before the slice is repointed at the new array.
			clear(block)
			reportEnvironmentBlockScratch("discarded-block", block)
			block = grown
		}
		block = append(block, converted...)
	}
	// F7 (P83 fix round): the terminating NUL is reserved before it is appended. A
	// slice that happens to be exactly full would otherwise make append allocate a
	// new backing array and abandon this one while it still holds the child's
	// environment - a reallocation the clear() calls above cannot see.
	//
	// T7 (P83 fix round 6): this branch is reachable deterministically through
	// environmentBlockInitialCapacity: entries whose encoded length equals the
	// capacity of the initial block arrive here with len == cap.
	if len(block) == cap(block) {
		grown := make([]uint16, len(block), len(block)*2+1)
		copy(grown, block)
		clear(block)
		reportEnvironmentBlockScratch("discarded-block-terminator", block)
		block = grown
	}
	block = append(block, 0)
	return block, nil
}

// clearEnvironmentScratch zeroes every intermediate UTF-16 slice a build
// allocated. F7 (Item D): it is called from one deferred site so the success path
// and every early return share the same guarantee instead of relying on each
// branch to remember.
func clearEnvironmentScratch(slices [][]uint16) {
	for _, slice := range slices {
		clear(slice)
	}
	runtime.KeepAlive(slices)
}

// environmentBlockScratchObserver is the F7 test seam of this file. Production
// leaves it nil; a test sets it to collect every intermediate UTF-16 slice and
// every backing array a growth abandons. The collected slices share their backing
// arrays with the builder, so a test can assert after the call that each one was
// zeroed - deleting any clear() in buildEnvironmentBlock turns it red.
var environmentBlockScratchObserver func(stage string, slice []uint16)

// reportEnvironmentBlockScratch hands one slice to the F7 test seam.
func reportEnvironmentBlockScratch(stage string, slice []uint16) {
	if environmentBlockScratchObserver != nil {
		environmentBlockScratchObserver(stage, slice)
	}
}

// optionalUTF16 encodes an optional path argument.
func optionalUTF16(value string) (*uint16, error) {
	if value == "" {
		return nil, nil
	}
	return syscall.UTF16PtrFromString(value)
}

// attach verifies the ownership anchor instead of establishing it: the child was
// created inside the job, so a process outside it means the creation did not
// honour the attribute list and the spawn must not proceed.
func (process *platformProcess) attach(pid int) error {
	job, ok := process.jobHandle()
	if !ok {
		return errors.New("managed process job handle is closed")
	}
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(handle)
	var inJob int32
	verified, _, callErr := procIsProcessInJob.Call(uintptr(handle), uintptr(job), uintptr(unsafe.Pointer(&inJob)))
	if verified == 0 {
		return callErr
	}
	if inJob == 0 {
		return fmt.Errorf("managed process %d was created outside its job", pid)
	}
	return nil
}

// resume lets the process run for the first time. It is the last step of the
// startup order, so nothing the child does can happen before the ownership anchor
// and the OnSpawn hook are complete.
func (process *platformProcess) resume() error {
	handle, ok := process.processHandle()
	if !ok {
		return errors.New("managed process handle is closed")
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
	if process.process != 0 {
		syscall.CloseHandle(process.process)
		process.process = 0
	}
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
