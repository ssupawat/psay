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

// Appends the protocol inline (no @-imports) so every harness loads it verbatim.
func runInit() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	body := strings.Replace(protocolDoc, "# Audio Announcement Protocol", "## Audio Announcement Protocol", 1)
	block := psayStart + "\n\n" + body + psayEnd + "\n"

	installed := 0
	for _, rel := range harnessFiles {
		path := filepath.Join(home, rel)
		existing, err := os.ReadFile(path)
		if err == nil && strings.Contains(string(existing), psayStart) {
			fmt.Println("already installed:", rel)
			continue
		}
		if err := appendBlock(path, block, existing); err != nil {
			return err
		}
		fmt.Println("installed:", rel)
		installed++
	}
	fmt.Printf("%d/%d harness files updated\n", installed, len(harnessFiles))
	return nil
}

func appendBlock(path, block string, existing []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
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
