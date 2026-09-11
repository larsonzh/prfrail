package console

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
)

type fakeSecretManager struct {
	values       map[string][]byte
	setValue     []byte
	setErr       error
	existsErr    error
	deleteErr    error
	resolveError error
}

func (manager *fakeSecretManager) Set(_ context.Context, reference string, secret []byte) error {
	manager.setValue = secret
	if manager.setErr != nil {
		return manager.setErr
	}
	if manager.values == nil {
		manager.values = make(map[string][]byte)
	}
	manager.values[reference] = append([]byte(nil), secret...)
	return nil
}

func (manager *fakeSecretManager) Resolve(_ context.Context, reference string) ([]byte, error) {
	if manager.resolveError != nil {
		return nil, manager.resolveError
	}
	value, found := manager.values[reference]
	if !found {
		return nil, adapters.ErrAISecretNotFound
	}
	return append([]byte(nil), value...), nil
}

func (manager *fakeSecretManager) Exists(_ context.Context, reference string) (bool, error) {
	if manager.existsErr != nil {
		return false, manager.existsErr
	}
	_, found := manager.values[reference]
	return found, nil
}

func (manager *fakeSecretManager) Delete(_ context.Context, reference string) error {
	if manager.deleteErr != nil {
		return manager.deleteErr
	}
	delete(manager.values, reference)
	return nil
}

func TestSecretCLISetStatusDeleteWithoutDisclosure(t *testing.T) {
	const canary = "proofrail-secret-cli-canary"
	const reference = "windows-credential:ProofRail/deepseek"
	manager := &fakeSecretManager{}
	cli := newTestCLI(t.TempDir())
	cli.secretManager = manager
	cli.readSecret = func() ([]byte, error) { return []byte(canary), nil }

	code, stdout, stderr := runCLI(t, cli, "secret", "set", "--ref", reference, "--json")
	if code != 0 || stderr != "" || strings.Contains(stdout, canary) {
		t.Fatalf("secret set failed or disclosed input: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, value := range manager.setValue {
		if value != 0 {
			t.Fatal("CLI did not clear the input secret after storage")
		}
	}
	code, stdout, stderr = runCLI(t, cli, "secret", "status", "--ref", reference, "--json")
	if code != 0 || stderr != "" || !strings.Contains(stdout, `"exists": true`) || strings.Contains(stdout, canary) {
		t.Fatalf("unexpected secret status: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	code, stdout, stderr = runCLI(t, cli, "secret", "delete", "--ref", reference, "--json")
	if code != 0 || stderr != "" || !strings.Contains(stdout, `"exists": false`) || strings.Contains(stdout, canary) {
		t.Fatalf("unexpected secret delete: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestSecretCLIRejectsPositionalValueWithoutReading(t *testing.T) {
	read := false
	cli := newTestCLI(t.TempDir())
	cli.secretManager = &fakeSecretManager{}
	cli.readSecret = func() ([]byte, error) {
		read = true
		return nil, errors.New("should not be called")
	}
	code, stdout, _ := runCLI(t, cli, "secret", "set", "--ref", "windows-credential:ProofRail/test", "secret-on-argv", "--json")
	if code != exitUsage || read || strings.Contains(stdout, "secret-on-argv") {
		t.Fatalf("positional secret was not rejected safely: code=%d read=%v stdout=%q", code, read, stdout)
	}
}

func TestSecretCLIRequiresReference(t *testing.T) {
	cli := newTestCLI(t.TempDir())
	for _, action := range []string{"set", "status", "delete"} {
		code, _, _ := runCLI(t, cli, "secret", action, "--json")
		if code != exitUsage {
			t.Fatalf("secret %s without reference returned %d", action, code)
		}
	}
}
