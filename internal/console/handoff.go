package console

import (
	"errors"
	"fmt"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const handoffSecretDomain = "proofrail:handoff-secret-input:1\n"

var (
	ErrInvalidHandoffInput = errors.New("invalid handoff input")
	ErrSecureInputRequired = errors.New("secure handoff input required")
)

type InputPolicy string

const (
	StructuredInput   InputPolicy = "structured"
	SecretDirectInput InputPolicy = "secret-direct"
)

type InputChannel string

const (
	TerminalInput InputChannel = "terminal"
	ModelInput    InputChannel = "model"
)

type PromptField struct {
	ID     string
	Secret bool
}

type HandoffInputSpec struct {
	Policy InputPolicy
	Fields []PromptField
}

type InputCapability struct {
	TerminalSecretInput bool
}

type PromptValue struct {
	ID    string
	Value string
}

type CollectedInput struct {
	Structured     map[string]string
	SecretEvidence []string
}

func ValidateHandoffInputSpec(spec HandoffInputSpec) error {
	if len(spec.Fields) == 0 {
		return fmt.Errorf("%w: empty prompt fields", ErrInvalidHandoffInput)
	}
	fields := make(map[string]PromptField, len(spec.Fields))
	for _, field := range spec.Fields {
		if !evidence.ValidID(field.ID) {
			return fmt.Errorf("%w: invalid prompt id", ErrInvalidHandoffInput)
		}
		if _, exists := fields[field.ID]; exists {
			return fmt.Errorf("%w: duplicate prompt id", ErrInvalidHandoffInput)
		}
		fields[field.ID] = field
	}
	switch spec.Policy {
	case StructuredInput:
		for _, field := range spec.Fields {
			if field.Secret {
				return fmt.Errorf("%w: structured input cannot include secret fields", ErrInvalidHandoffInput)
			}
		}
	case SecretDirectInput:
		for _, field := range spec.Fields {
			if !field.Secret {
				return fmt.Errorf("%w: secret-direct input requires all fields to be secret", ErrInvalidHandoffInput)
			}
		}
	default:
		return fmt.Errorf("%w: unknown input policy", ErrInvalidHandoffInput)
	}
	return nil
}

func PrepareHandoffInput(spec HandoffInputSpec, values []PromptValue, channel InputChannel, capability InputCapability) (CollectedInput, error) {
	if err := ValidateHandoffInputSpec(spec); err != nil {
		return CollectedInput{}, err
	}
	if channel != TerminalInput && channel != ModelInput {
		return CollectedInput{}, fmt.Errorf("%w: unknown input channel", ErrInvalidHandoffInput)
	}
	if spec.Policy == SecretDirectInput {
		if channel != TerminalInput || !capability.TerminalSecretInput {
			return CollectedInput{}, ErrSecureInputRequired
		}
	}
	provided := make(map[string]string, len(values))
	for _, value := range values {
		if _, exists := provided[value.ID]; exists {
			return CollectedInput{}, fmt.Errorf("%w: duplicate prompt value", ErrInvalidHandoffInput)
		}
		if err := validatePromptValue(value.Value); err != nil {
			return CollectedInput{}, err
		}
		provided[value.ID] = value.Value
	}
	collected := CollectedInput{
		Structured:     map[string]string{},
		SecretEvidence: []string{},
	}
	for _, field := range spec.Fields {
		value, ok := provided[field.ID]
		if !ok {
			return CollectedInput{}, fmt.Errorf("%w: missing prompt value %q", ErrInvalidHandoffInput, field.ID)
		}
		delete(provided, field.ID)
		if field.Secret {
			collected.SecretEvidence = append(collected.SecretEvidence, evidence.Digest(handoffSecretDomain, []byte(field.ID+"\n"+value)))
			continue
		}
		collected.Structured[field.ID] = value
	}
	if len(provided) != 0 {
		return CollectedInput{}, fmt.Errorf("%w: undeclared prompt values", ErrInvalidHandoffInput)
	}
	return collected, nil
}

func validatePromptValue(value string) error {
	if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\x00\r") || len(value) > 4096 {
		return fmt.Errorf("%w: invalid prompt value", ErrInvalidHandoffInput)
	}
	return nil
}
