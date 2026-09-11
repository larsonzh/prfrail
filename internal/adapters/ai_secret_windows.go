//go:build windows

package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

const (
	windowsCredentialPrefix       = "windows-credential:"
	windowsCredentialTypeGeneric  = 1
	windowsCredentialPersistLocal = 2
	windowsErrorNotFound          = syscall.Errno(1168)
)

var (
	advapi32         = syscall.NewLazyDLL("advapi32.dll")
	procCredWriteW   = advapi32.NewProc("CredWriteW")
	procCredReadW    = advapi32.NewProc("CredReadW")
	procCredDeleteW  = advapi32.NewProc("CredDeleteW")
	procCredFree     = advapi32.NewProc("CredFree")
	ErrAISecretStore = errors.New("AI secret store operation failed")
)

type windowsCredential struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint16
	Comment            *uint16
	LastWrittenLow     uint32
	LastWrittenHigh    uint32
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        *uint16
	UserName           *uint16
}

type WindowsCredentialManager struct{}

func NewPlatformAISecretManager() AISecretManager {
	return WindowsCredentialManager{}
}

func (WindowsCredentialManager) Set(_ context.Context, reference string, secret []byte) error {
	target, err := windowsCredentialTarget(reference)
	if err != nil || len(secret) == 0 || len(secret) > 2560 {
		return fmt.Errorf("%w: invalid reference or secret size", ErrAISecretStore)
	}
	targetName, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return fmt.Errorf("%w: invalid target", ErrAISecretStore)
	}
	userName, _ := syscall.UTF16PtrFromString("ProofRail")
	credential := windowsCredential{
		Type:               windowsCredentialTypeGeneric,
		TargetName:         targetName,
		CredentialBlobSize: uint32(len(secret)),
		CredentialBlob:     &secret[0],
		Persist:            windowsCredentialPersistLocal,
		UserName:           userName,
	}
	result, _, callErr := procCredWriteW.Call(uintptr(unsafe.Pointer(&credential)), 0)
	if result == 0 {
		return windowsCredentialError("set", callErr)
	}
	return nil
}

func (WindowsCredentialManager) Resolve(_ context.Context, reference string) ([]byte, error) {
	credential, err := readWindowsCredential(reference)
	if err != nil {
		return nil, err
	}
	defer procCredFree.Call(uintptr(unsafe.Pointer(credential)))
	if credential.CredentialBlobSize == 0 || credential.CredentialBlob == nil {
		return nil, fmt.Errorf("%w: empty credential", ErrAISecretStore)
	}
	return append([]byte(nil), unsafe.Slice(credential.CredentialBlob, int(credential.CredentialBlobSize))...), nil
}

func (WindowsCredentialManager) Exists(_ context.Context, reference string) (bool, error) {
	credential, err := readWindowsCredential(reference)
	if errors.Is(err, ErrAISecretNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	procCredFree.Call(uintptr(unsafe.Pointer(credential)))
	return true, nil
}

func (WindowsCredentialManager) Delete(_ context.Context, reference string) error {
	target, err := windowsCredentialTarget(reference)
	if err != nil {
		return err
	}
	targetName, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return fmt.Errorf("%w: invalid target", ErrAISecretStore)
	}
	result, _, callErr := procCredDeleteW.Call(uintptr(unsafe.Pointer(targetName)), windowsCredentialTypeGeneric, 0)
	if result == 0 {
		return windowsCredentialError("delete", callErr)
	}
	return nil
}

func readWindowsCredential(reference string) (*windowsCredential, error) {
	target, err := windowsCredentialTarget(reference)
	if err != nil {
		return nil, err
	}
	targetName, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid target", ErrAISecretStore)
	}
	var credential *windowsCredential
	result, _, callErr := procCredReadW.Call(
		uintptr(unsafe.Pointer(targetName)),
		windowsCredentialTypeGeneric,
		0,
		uintptr(unsafe.Pointer(&credential)),
	)
	if result == 0 {
		return nil, windowsCredentialError("read", callErr)
	}
	return credential, nil
}

func windowsCredentialTarget(reference string) (string, error) {
	if !strings.HasPrefix(reference, windowsCredentialPrefix) {
		return "", fmt.Errorf("%w: unsupported reference", ErrAISecretStore)
	}
	target := strings.TrimPrefix(reference, windowsCredentialPrefix)
	if target == "" || len(target) > 32767 || strings.TrimSpace(target) != target || containsASCIIControl(target) {
		return "", fmt.Errorf("%w: invalid target", ErrAISecretStore)
	}
	return target, nil
}

func containsASCIIControl(value string) bool {
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return true
		}
	}
	return false
}

func windowsCredentialError(operation string, callErr error) error {
	if errors.Is(callErr, windowsErrorNotFound) {
		return ErrAISecretNotFound
	}
	return fmt.Errorf("%w: Windows credential %s", ErrAISecretStore, operation)
}
