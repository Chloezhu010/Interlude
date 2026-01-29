package main

import (
	"github.com/spf13/cobra" // CLI framework
)

// startCmd runs the daemon in the foreground
// Shows logs in terminal, then TUI appears when Claude is active past threshold
func startCmd() *cobra.Command {
	return &cobra.Command {
		Use: "start",
		Short: "Start Interlude and listen for Claude Code",
		Run: func(cmd *cobra.Command, args []string) {
			runDaemon() // run directly in foreground
		},
	}
}