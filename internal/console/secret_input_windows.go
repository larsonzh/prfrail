//go:build windows

package console

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const windowsEnableEchoInput = 0x0004
const maximumTerminalSecretBytes = 2560

var (
	errSecureSecretInput = errors.New("secure terminal secret input unavailable")
	kernel32SecretInput  = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode   = kernel32SecretInput.NewProc("GetConsoleMode")
	procSetConsoleMode   = kernel32SecretInput.NewProc("SetConsoleMode")
)

func readTerminalSecret() ([]byte, error) {
	handle := syscall.Handle(os.Stdin.Fd())
	var originalMode uint32
	result, _, _ := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&originalMode)))
	if result == 0 {
		return nil, errSecureSecretInput
	}
	result, _, _ = procSetConsoleMode.Call(uintptr(handle), uintptr(originalMode&^windowsEnableEchoInput))
	if result == 0 {
		return nil, errSecureSecretInput
	}
	defer procSetConsoleMode.Call(uintptr(handle), uintptr(originalMode))

	fmt.Fprint(os.Stderr, "Secret: ")
	secret := make([]byte, 0, 64)
	unit := []byte{0}
	defer clear(unit)
	for {
		count, err := os.Stdin.Read(unit)
		if count > 0 {
			if unit[0] == '\n' {
				break
			}
			if len(secret) >= maximumTerminalSecretBytes {
				clear(secret)
				return nil, errSecureSecretInput
			}
			secret = append(secret, unit[0])
		}
		if err != nil {
			if errors.Is(err, io.EOF) && len(secret) > 0 {
				break
			}
			clear(secret)
			return nil, errSecureSecretInput
		}
	}
	fmt.Fprintln(os.Stderr)
	if len(secret) > 0 && secret[len(secret)-1] == '\r' {
		secret = secret[:len(secret)-1]
	}
	return secret, nil
}
