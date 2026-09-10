package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Settings struct {
	Voice       string            `json:"voice"`
	LengthScale float64           `json:"length_scale"`
	Lexicon     map[string]string `json:"lexicon"`
}

var defaultSettings = Settings{
	Voice:       "en_US-lessac-medium",
	LengthScale: 1.0,
	Lexicon: map[string]string{
		"psay": "p say",
		"aws":  "A W S",
		"db":   "database",
		"err":  "error",
		"nil":  "null",
		"repo": "repository",
		"auth": "authentication",
	},
}

func settingsPath() (string, error) {
	dir, err := defaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

func loadSettings() (Settings, error) {
	path, err := settingsPath()
	if err != nil {
		return Settings{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if err := saveSettings(defaultSettings); err != nil {
			return Settings{}, err
		}
		return defaultSettings, nil
	}
	if err != nil {
		return Settings{}, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return s, nil
}

func saveSettings(s Settings) error {
	path, err := settingsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
