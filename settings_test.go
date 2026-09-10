package main

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettingsAutoInit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	s, err := loadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if s.Voice != defaultSettings.Voice || s.LengthScale != defaultSettings.LengthScale {
		t.Errorf("auto-init returned %+v, want defaults", s)
	}
	if !maps.Equal(s.Lexicon, defaultSettings.Lexicon) {
		t.Errorf("auto-init lexicon = %v, want %v", s.Lexicon, defaultSettings.Lexicon)
	}

	data, err := os.ReadFile(filepath.Join(psayDir(t), "settings.json"))
	if err != nil {
		t.Fatalf("settings.json not created: %v", err)
	}
	var written Settings
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatal(err)
	}
	if written.Voice != defaultSettings.Voice {
		t.Errorf("written voice = %q, want %q", written.Voice, defaultSettings.Voice)
	}
}

func TestLoadSettingsReadsExisting(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	want := Settings{Voice: "en_GB-alba-medium", LengthScale: 0.8, Lexicon: map[string]string{"db": "dee bee"}}
	if err := saveSettings(want); err != nil {
		t.Fatal(err)
	}

	got, err := loadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if got.Voice != want.Voice || got.LengthScale != want.LengthScale || !maps.Equal(got.Lexicon, want.Lexicon) {
		t.Errorf("loadSettings() = %+v, want %+v", got, want)
	}
}

func psayDir(t *testing.T) string {
	t.Helper()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(home, ".psay")
}
