package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed PSAY.md
var protocolDoc string

func runInit() error {
	const dir = ".claude"
	const path = dir + "/PSAY.md"
	if _, err := os.Stat(path); err == nil {
		fmt.Println(path, "already exists.")
	} else {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(protocolDoc), 0o644); err != nil {
			return err
		}
		fmt.Println("Wrote", path, "(audio announcement protocol).")
	}
	fmt.Println("Activate: add '@.claude/PSAY.md' to CLAUDE.md or AGENTS.md.")
	return nil
}
