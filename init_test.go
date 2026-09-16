package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func harnessPaths(home string) []string {
	paths := make([]string, len(harnessFiles))
	for i, rel := range harnessFiles {
		paths[i] = filepath.Join(home, rel)
	}
	return paths
}

func TestRunInitInstallsToExistingOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	keep := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(keep), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keep, []byte("my rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runInit(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(keep)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), psayStart) || !strings.HasPrefix(string(data), "my rules\n") {
		t.Errorf("existing file not extended in place: %q", data)
	}
	for _, p := range harnessPaths(home) {
		if p == keep {
			continue
		}
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s should not have been created", p)
		}
	}
}

func TestRunInitUpdatesStaleBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	stale := "my rules\n" + psayStart + "\n\nstale protocol\n" + psayEnd + "\ntrailing notes\n"
	if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runInit(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.HasPrefix(got, "my rules\n") || !strings.HasSuffix(got, "trailing notes\n") {
		t.Errorf("surrounding content not preserved: %q", got)
	}
	if strings.Contains(got, "stale protocol") {
		t.Errorf("stale block survived: %q", got)
	}
	if !strings.Contains(got, "Audio Announcement Protocol") || strings.Count(got, psayStart) != 1 {
		t.Errorf("current block not installed exactly once: %q", got)
	}
}

func TestRunInitIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	for _, p := range harnessPaths(home) {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
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
