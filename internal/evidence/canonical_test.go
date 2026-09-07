package evidence

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalVectors(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "contracts", "vectors", "canonical.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		VectorID        string          `json:"vectorId"`
		Input           json.RawMessage `json:"input"`
		DomainSeparator string          `json:"domainSeparator"`
		Canonical       string          `json:"canonical"`
		Digest          string          `json:"digest"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	for _, vector := range vectors {
		t.Run(vector.VectorID, func(t *testing.T) {
			canonical, err := Canonicalize(vector.Input)
			if err != nil {
				t.Fatal(err)
			}
			if string(canonical) != vector.Canonical {
				t.Fatalf("canonical = %q, want %q", canonical, vector.Canonical)
			}
			if digest := Digest(vector.DomainSeparator, canonical); digest != vector.Digest {
				t.Fatalf("digest = %s, want %s", digest, vector.Digest)
			}
		})
	}
}

func TestCanonicalizeRejectsAmbiguousInput(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{"duplicate key", []byte(`{"a":1,"a":2}`), ErrDuplicateKey},
		{"trailing value", []byte(`{} {}`), ErrInvalidJSON},
		{"invalid UTF-8", []byte{'"', 0xff, '"'}, ErrInvalidJSON},
		{"unpaired high surrogate", []byte(`"\ud800"`), ErrInvalidJSON},
		{"unpaired low surrogate", []byte(`"\udc00"`), ErrInvalidJSON},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Canonicalize(test.input)
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestCanonicalNumbersFollowECMAScriptThresholds(t *testing.T) {
	canonical, err := Canonicalize([]byte(`[1e-7,1e-6,1e20,1e21,-0]`))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(canonical), `[1e-7,0.000001,100000000000000000000,1e+21,0]`; got != want {
		t.Fatalf("canonical = %s, want %s", got, want)
	}
}
