package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	voiceOverride := flag.String("v", "", "voice override")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "init" {
		if err := runInit(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	text := strings.Join(args, " ")
	if text == "" {
		fmt.Fprintln(os.Stderr, "usage: psay 'text to speak'")
		os.Exit(2)
	}

	settings, err := loadSettings()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	voice := settings.Voice
	if voice == "" {
		voice = defaultSettings.Voice
	}
	if *voiceOverride != "" {
		voice = *voiceOverride
	}
	if settings.LengthScale <= 0 {
		settings.LengthScale = 1.0
	}

	spoken := replaceLexicon(text, settings.Lexicon)
	if err := speak(voice, settings.LengthScale, spoken); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func speak(voice string, lengthScale float64, text string) error {
	psayDir, err := defaultDir()
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("", "psay-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	wav := filepath.Join(dir, "out.wav")

	piper := exec.Command("piper",
		"-m", voice,
		"--data-dir", filepath.Join(psayDir, "voices"),
		"--length-scale", fmt.Sprintf("%g", lengthScale),
		"-f", wav,
		"--", text,
	)
	if out, err := run(piper); err != nil {
		return fmt.Errorf("piper: %w\n%s\nhint: python -m piper.download_voices %s --download-dir %s",
			err, out, voice, filepath.Join(psayDir, "voices"))
	}

	_, err = run(exec.Command("afplay", wav))
	return err
}

func run(cmd *exec.Cmd) (string, error) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	}
	return "", nil
}

func defaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".psay"), nil
}
