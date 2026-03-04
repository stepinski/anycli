package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anycli",
	Short: "Terminal client for AnythingLLM",
	Long: `anycli — chat with your knowledge base and codebase from the terminal.

  Zero-setup codebase chat, Unix-native piping, and a beautiful TUI.
  The Go alternative to anything-llm-cli.

Examples:
  anycli "how does auth work in this codebase?"
  git diff | anycli "write a commit message"
  cat error.log | anycli "what's wrong?"
  anycli tui`,

	SilenceUsage: true,
}

func Execute(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
