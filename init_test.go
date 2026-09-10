package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitWritesProtocol(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(".claude", "PSAY.md"))
	if err != nil {
		t.Fatalf(".claude/PSAY.md not written: %v", err)
	}
	if !strings.Contains(string(data), "Audio Announcement Protocol") {
		t.Error("embedded PSAY.md missing protocol heading")
	}
}

func TestRunInitKeepsExisting(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.MkdirAll(".claude", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(".claude", "PSAY.md"), []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(".claude", "PSAY.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom" {
		t.Errorf("existing .claude/PSAY.md overwritten with %q", data)
	}
}
