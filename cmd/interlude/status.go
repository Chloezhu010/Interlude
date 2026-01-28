package main

import (
	"fmt"
	"os"
	"strconv" // convert int to string
	"strings"
	"syscall" // low level os call

	"github.com/spf13/cobra" // CLI framework
)

func statusDaemonCmd() *cobra.Command {
	return &cobra.Command {
		Use: "status",
		Short: "Show the daemon status",
		Run: func(cmd *cobra.Command, args []string) {
			if isRunning(){
				pid, _ := readPID()
				fmt.Printf("Daemon is running (PID %d)\n", pid)
			} else {
				fmt.Println("Daemon is not running")
			}
		},
	}
}

/* Check if the daemon process is alive */
func isRunning() bool {
	// get the PID from the file
	pid, err := readPID()
	if err != nil {
		return false
	}
	// create the process struct
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// check if the process is alive by sending signal 0
	// return nil if process exists, error if not
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

/* read the PID from the file */
func readPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
} 