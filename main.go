package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed speak.py
var speakScript string

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
	if settings.Speed <= 0 {
		settings.Speed = defaultSettings.Speed
	}

	spoken := replaceLexicon(text, settings.Lexicon)
	if err := speak(voice, settings.Speed, spoken); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func speak(voice string, speed float64, text string) error {
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

	python := filepath.Join(psayDir, "venv", "bin", "python")
	kokoro := exec.Command(python, "-c", speakScript,
		filepath.Join(psayDir, "kokoro", "kokoro-v1.0.onnx"),
		filepath.Join(psayDir, "kokoro", "voices-v1.0.bin"),
		voice,
		fmt.Sprintf("%g", speed),
		wav,
		text,
	)
	if out, err := run(kokoro); err != nil {
		return fmt.Errorf("kokoro: %w\n%s\nhint: run the psay install.sh to set up ~/.psay/venv and ~/.psay/kokoro", err, out)
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
