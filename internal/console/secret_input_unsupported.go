//go:build !windows

package console

import "errors"

var errSecureSecretInput = errors.New("secure terminal secret input unavailable")

func readTerminalSecret() ([]byte, error) {
	return nil, errSecureSecretInput
}
