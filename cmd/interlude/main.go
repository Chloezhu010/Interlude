package main

import (
	"fmt" // std lib, print to console
	"strings" // std lib, string manipulation
	"time" // std lib, time operations

	"github.com/shirou/gopsutil/v3/process" // process monitoring
)

func main() {
	fmt.Println("Watching for Claude Code...")

	for {
		// get all running processes (_ = ignore error)
		processes, _ := process.Processes()
		// loop through each process, _ = ignore index, p = process
		for _, p := range processes {
			// get process name
			cmdline, _ := p.Cmdline()
			// if process name contains "claude"
			if strings.Contains(strings.ToLower(cmdline), "claude"){
				// print the process name and its ID
				fmt.Printf("Detected: %s (PID: %d)\n", cmdline, p.Pid)
			}
		}
		// wait every 500ms to check again
		time.Sleep(500 * time.Millisecond)
	}
}