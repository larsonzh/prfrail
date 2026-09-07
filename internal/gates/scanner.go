package gates

import (
	"bytes"
	"context"
	"errors"
)

var ErrSecretDetected = errors.New("secret detected")

type SecretScanner struct {
	Markers [][]byte
}

func DefaultSecretScanner() SecretScanner {
	return SecretScanner{Markers: [][]byte{
		[]byte("-----BEGIN PRIVATE KEY-----"),
		[]byte("-----BEGIN OPENSSH PRIVATE KEY-----"),
		[]byte("AKIA"),
		[]byte("ghp_"),
		[]byte("github_pat_"),
	}}
}

func (scanner SecretScanner) Scan(ctx context.Context, _ string, data []byte) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	for _, marker := range scanner.Markers {
		if len(marker) > 0 && bytes.Contains(data, marker) {
			return "", ErrSecretDetected
		}
	}
	return "passed", nil
}
