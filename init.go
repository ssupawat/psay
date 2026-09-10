package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed PSAY.md
var protocolDoc string

func runInit() error {
	const path = "PSAY.md"
	if _, err := os.Stat(path); err == nil {
		fmt.Println("PSAY.md already exists in this directory.")
	} else {
		if err := os.WriteFile(path, []byte(protocolDoc), 0o644); err != nil {
			return err
		}
		fmt.Println("Wrote PSAY.md (audio announcement protocol).")
	}
	fmt.Println("Activate: add '@PSAY.md' to CLAUDE.md, or 'Read PSAY.md for the audio announcement protocol.' to AGENTS.md.")
	return nil
}
