package main

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"testing"
)

func psayDir(t *testing.T) string {
	t.Helper()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(home, ".psay")
}

func TestLoadSettingsAutoInit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	s, err := loadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if s.Voice != defaultSettings.Voice || s.Speed != defaultSettings.Speed {
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

	want := Settings{Voice: "af_bella", Speed: 0.9, Lexicon: map[string]string{"db": "dee bee"}}
	if err := saveSettings(want); err != nil {
		t.Fatal(err)
	}

	got, err := loadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if got.Voice != want.Voice || got.Speed != want.Speed || !maps.Equal(got.Lexicon, want.Lexicon) {
		t.Errorf("loadSettings() = %+v, want %+v", got, want)
	}
}

func TestLoadSettingsMigratesPiperFormat(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	path := filepath.Join(psayDir(t), "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"voice": "en_US-lessac-medium", "length_scale": 1.15, "lexicon": {"db": "database"}}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := loadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if s.Voice != defaultSettings.Voice || s.Speed != defaultSettings.Speed {
		t.Errorf("migration kept piper values: %+v", s)
	}
	if !maps.Equal(s.Lexicon, map[string]string{"db": "database"}) {
		t.Errorf("migration lost lexicon: %v", s.Lexicon)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if _, present := raw["length_scale"]; present {
		t.Errorf("migrated file still has length_scale: %s", data)
	}
}
