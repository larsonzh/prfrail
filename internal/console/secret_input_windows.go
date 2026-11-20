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
	secret, err := readTerminalSecretFrom(os.Stdin)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr)
	return secret, nil
}

// readTerminalSecretFrom is the F7 test seam of the terminal reader: it is the
// whole input loop without the console-mode dance, so a test can drive it with a
// plain reader and observe the buffers a growth abandons.
func readTerminalSecretFrom(reader io.Reader) ([]byte, error) {
	secret := make([]byte, 0, 64)
	unit := []byte{0}
	defer clear(unit)
	for {
		count, err := reader.Read(unit)
		if count > 0 {
			if unit[0] == '\n' {
				break
			}
			if len(secret) >= maximumTerminalSecretBytes {
				clear(secret)
				return nil, errSecureSecretInput
			}
			// F7 (P83 fix round): a full buffer is copied into a larger one, and the
			// array it leaves behind still holds the terminal input typed so far. It
			// is cleared before the slice is repointed, so a growth cannot leave a
			// readable pre-growth copy behind.
			if len(secret) == cap(secret) {
				grown := make([]byte, len(secret), cap(secret)*2)
				copy(grown, secret)
				clear(secret)
				reportTerminalSecretBuffer("discarded", secret)
				secret = grown
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
	if len(secret) > 0 && secret[len(secret)-1] == '\r' {
		secret = secret[:len(secret)-1]
	}
	return secret, nil
}

// terminalSecretBufferObserver is the F7 test seam of the terminal reader.
// Production leaves it nil; a test sets it to collect the buffer a growth
// abandons, so deleting the clear() above turns the test red.
var terminalSecretBufferObserver func(stage string, buffer []byte)

// reportTerminalSecretBuffer hands one buffer to the F7 test seam.
func reportTerminalSecretBuffer(stage string, buffer []byte) {
	if terminalSecretBufferObserver != nil {
		terminalSecretBufferObserver(stage, buffer)
	}
}
