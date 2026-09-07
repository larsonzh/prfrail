package tickets

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const failureFingerprintDomain = "proofrail:failure-fingerprint:1\n"

var (
	ErrInvalidFingerprintInput = errors.New("invalid fingerprint input")
	codePattern                = regexp.MustCompile(`^[a-z][a-z0-9-]*(?:\.[a-z][a-z0-9-]*)+$`)
)

type Subject struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

type FingerprintInput struct {
	Code         string  `json:"code"`
	Subject      Subject `json:"subject"`
	FailurePoint string  `json:"failurePoint"`
}

func Fingerprint(input FingerprintInput) (string, error) {
	if err := validateFingerprintInput(input); err != nil {
		return "", err
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.Canonicalize(payload)
	if err != nil {
		return "", err
	}
	return evidence.Digest(failureFingerprintDomain, canonical), nil
}

func validateFingerprintInput(input FingerprintInput) error {
	if !codePattern.MatchString(input.Code) {
		return fmt.Errorf("%w: invalid error code", ErrInvalidFingerprintInput)
	}
	switch input.Subject.Kind {
	case "system", "chain", "task", "step", "hook", "adapter", "store":
	default:
		return fmt.Errorf("%w: invalid subject kind", ErrInvalidFingerprintInput)
	}
	if !evidence.ValidID(input.Subject.ID) || !evidence.ValidID(input.FailurePoint) {
		return fmt.Errorf("%w: invalid subject ID or failure point", ErrInvalidFingerprintInput)
	}
	return nil
}
