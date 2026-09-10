package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunInitWritesProtocol(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("PSAY.md")
	if err != nil {
		t.Fatalf("PSAY.md not written: %v", err)
	}
	if !strings.Contains(string(data), "Audio Announcement Protocol") {
		t.Error("embedded PSAY.md missing protocol heading")
	}
}

func TestRunInitKeepsExisting(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.WriteFile("PSAY.md", []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("PSAY.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "custom" {
		t.Errorf("existing PSAY.md overwritten with %q", data)
	}
}
