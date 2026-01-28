package main

import (
	"fmt"
	"os"
	"os/exec" // for spawning a child process
	"strconv" // convert int to string
	"syscall" // low level os call

	"github.com/spf13/cobra" // CLI framework
)

func startCmd() *cobra.Command {
	return &cobra.Command {
		Use: "start",
		Short: "Start the daemon",
		Run: func(cmd *cobra.Command, args []string) {
			// check if already running
			if isRunning() {
				fmt.Println("Daemon is already running")
				return
			}
			// if not, spawn background daemon
			spwanDaemon()
		},
	}
}

func spwanDaemon() {
	// open log file (create if missing, append to existing)
	logF, err := os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	// create the command
	cmd := exec.Command(os.Args[0], "run-daemon")
	// redirect child's output to log file
	cmd.Stdout = logF
	cmd.Stderr = logF
	// detach from terminal (survive terminal close)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	// start (parent continues immediately, child runs in background)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start daemon: %v\n", err)
		return
	}
	// save child's pid (stop it later)
	writePID(cmd.Process.Pid)
	fmt.Printf("Daemon started (PID: %d)\n", cmd.Process.Pid)
}

func writePID(pid int) {
	os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func runDaemonCmd() *cobra.Command {
	return &cobra.Command {
		Use: "run-daemon",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			runDaemon()
		},
	}
}