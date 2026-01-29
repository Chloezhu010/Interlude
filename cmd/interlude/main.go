package main

import (
	"fmt" // for printing errors
	"os" // for os.Exit

	"github.com/spf13/cobra" // CLI framework
)

func main() {
	// create root command
	rootCmd := &cobra.Command {
		Use: "interlude",
		Short: "Interlude is a tool to help you stay focused on your work when using Claude Code",
	}
	// add subcommands
	rootCmd.AddCommand(startCmd())
	// execute root command with error handling
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}