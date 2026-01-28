package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func logsDaemonCmd() *cobra.Command {
	return &cobra.Command {
		Use: "logs",
		Short: "Show the daemon logs",
		Run: func(cmd *cobra.Command, args []string) {
			// read and print log file contents
			data, err := os.ReadFile(logFile)
			if err != nil {
				fmt.Printf("Failed to read logs: %v\n", err)
				return
			}
			fmt.Println(string(data))
		},
	}
}