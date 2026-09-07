//go:build windows

package applier

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplacePathOverwritesExistingFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "journal.json")
	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(dir, "replacement.tmp")
	if err := os.WriteFile(source, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := replacePath(source, target); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("unexpected content: %q", data)
	}
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatalf("source still present after replacement: %v", err)
	}
}

func TestReplacePathMissingSourceFails(t *testing.T) {
	dir := t.TempDir()
	if err := replacePath(filepath.Join(dir, "missing.tmp"), filepath.Join(dir, "target.json")); err == nil {
		t.Fatal("expected missing source to fail")
	}
}
