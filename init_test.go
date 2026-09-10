package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitInstallsToAllHarnesses(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	for _, rel := range harnessFiles {
		data, err := os.ReadFile(filepath.Join(home, rel))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		if !strings.Contains(string(data), psayStart) || !strings.Contains(string(data), "Audio Announcement Protocol") {
			t.Errorf("%s missing protocol block", rel)
		}
	}
}

func TestRunInitIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude", "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), psayStart); got != 1 {
		t.Errorf("protocol block appears %d times, want 1", got)
	}
}

func TestRunInitPreservesExisting(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".codex", "AGENTS.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("my rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "my rules\n") {
		t.Errorf("existing content not preserved: %q", data)
	}
}
