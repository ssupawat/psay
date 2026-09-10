package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed PSAY.md
var protocolDoc string

var harnessFiles = []string{
	".codex/AGENTS.md",
	".claude/CLAUDE.md",
	".pi/agent/AGENTS.md",
	".zcode/AGENTS.md",
}

const psayStart = "<!-- psay:start -->"
const psayEnd = "<!-- psay:end -->"

// Installs into instruction files that already exist; never creates files
// for harnesses the user doesn't use. Appends inline (no @-imports) so every
// harness loads it verbatim.
func runInit() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	body := strings.Replace(protocolDoc, "# Audio Announcement Protocol", "## Audio Announcement Protocol", 1)
	block := psayStart + "\n\n" + body + psayEnd + "\n"

	installed, skipped := 0, 0
	for _, rel := range harnessFiles {
		path := filepath.Join(home, rel)
		existing, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			fmt.Println("skipped (not present):", rel)
			skipped++
			continue
		}
		if err != nil {
			return err
		}
		if strings.Contains(string(existing), psayStart) {
			fmt.Println("already installed:", rel)
			continue
		}
		if err := appendBlock(path, block, existing); err != nil {
			return err
		}
		fmt.Println("installed:", rel)
		installed++
	}
	fmt.Printf("%d installed, %d skipped, %d already installed\n", installed, skipped, len(harnessFiles)-installed-skipped)
	return nil
}

func appendBlock(path, block string, existing []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		if _, err := f.Write([]byte("\n")); err != nil {
			return err
		}
	}
	_, err = f.WriteString("\n" + block)
	return err
}
