package console

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestPrepareHandoffInputStructured(t *testing.T) {
	spec := HandoffInputSpec{
		Policy: StructuredInput,
		Fields: []PromptField{{ID: "goal"}, {ID: "scope"}},
	}
	result, err := PrepareHandoffInput(spec, []PromptValue{{ID: "goal", Value: "fix failing test"}, {ID: "scope", Value: "internal/chain"}}, ModelInput, InputCapability{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SecretEvidence) != 0 || result.Structured["goal"] != "fix failing test" || result.Structured["scope"] != "internal/chain" {
		t.Fatalf("unexpected structured result: %+v", result)
	}
}

func TestPrepareHandoffInputRejectsUndeclaredPrompt(t *testing.T) {
	spec := HandoffInputSpec{Policy: StructuredInput, Fields: []PromptField{{ID: "goal"}}}
	_, err := PrepareHandoffInput(spec, []PromptValue{{ID: "goal", Value: "do work"}, {ID: "extra", Value: "unexpected"}}, TerminalInput, InputCapability{TerminalSecretInput: true})
	if !errors.Is(err, ErrInvalidHandoffInput) {
		t.Fatalf("expected undeclared prompt rejection, got: %v", err)
	}
}

func TestPrepareHandoffInputRejectsUnknownChannel(t *testing.T) {
	spec := HandoffInputSpec{Policy: StructuredInput, Fields: []PromptField{{ID: "goal"}}}
	_, err := PrepareHandoffInput(spec, []PromptValue{{ID: "goal", Value: "do work"}}, InputChannel("clipboard"), InputCapability{})
	if !errors.Is(err, ErrInvalidHandoffInput) {
		t.Fatalf("expected unknown channel rejection, got: %v", err)
	}
}

func TestPrepareHandoffInputSecretDirectRequiresSecureTerminal(t *testing.T) {
	spec := HandoffInputSpec{Policy: SecretDirectInput, Fields: []PromptField{{ID: "token", Secret: true}}}
	values := []PromptValue{{ID: "token", Value: "secret-canary-123"}}
	if _, err := PrepareHandoffInput(spec, values, ModelInput, InputCapability{TerminalSecretInput: true}); !errors.Is(err, ErrSecureInputRequired) {
		t.Fatalf("expected model-channel rejection, got: %v", err)
	}
	if _, err := PrepareHandoffInput(spec, values, TerminalInput, InputCapability{TerminalSecretInput: false}); !errors.Is(err, ErrSecureInputRequired) {
		t.Fatalf("expected capability rejection, got: %v", err)
	}
}

func TestPrepareHandoffInputSecretDirectDoesNotExposeSecret(t *testing.T) {
	canary := "secret-canary-123"
	spec := HandoffInputSpec{Policy: SecretDirectInput, Fields: []PromptField{{ID: "token", Secret: true}}}
	result, err := PrepareHandoffInput(spec, []PromptValue{{ID: "token", Value: canary}}, TerminalInput, InputCapability{TerminalSecretInput: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SecretEvidence) != 1 || len(result.Structured) != 0 {
		t.Fatalf("unexpected secret-direct result: %+v", result)
	}
	if strings.Contains(fmt.Sprintf("%+v", result), canary) {
		t.Fatal("secret canary leaked into collected output")
	}
}
