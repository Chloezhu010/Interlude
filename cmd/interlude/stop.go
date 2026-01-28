package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

func stopDaemonCmd() *cobra.Command {
	return &cobra.Command {
		Use: "stop",
		Short: "Stop the daemon",
		Run: func(cmd *cobra.Command, args []string) {
			// read PID from file
			pid, err := readPID()
			if err != nil {
				fmt.Println("Daemon is not running")
				return
			}
			// find the process
			process, err := os.FindProcess(pid)
			if err != nil {
				fmt.Println("Daemon is not running")
				return
			}
			// send SIGTERM (gracefull kill)
			err = process.Signal(syscall.SIGTERM)
			if err != nil {
				fmt.Printf("Failed to stop daemon: %v\n", err)
				return
			}
			// clean up PID file
			os.Remove(pidFile)
			fmt.Println("Daemon stopped")
		},
	}
}