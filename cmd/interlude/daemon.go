package main

import (
	"fmt" // print to console
	"os" // file operations: readfile, userhome
	"os/signal" // catch os signals eg. ctrl C
	"strconv" // string-to-int conversion
	"strings" // string manipulation
	"syscall" // low-level os signals (sigterm, sigint)
	"time" // time operations

	// "github.com/shirou/gopsutil/v3/process" // process monitoring
)

var (
	home, _ = os.UserHomeDir()
	statusFile = home + "/.interlude/status" // status file path: "active" or "idle" written by cc hooks
	startTimeFile = home + "/.interlude/start_time" // start time file path: Unix timestamp written by cc hooks
	pidFile = home + "/.interlude/interlude.pid" // pid file path: PID of the daemon process
	logsFile = home + "/.interlude/interlude.log" // logs file path: logs of the daemon process
)

const threshold = 5 * time.Second

/* Main daemon loop
 * called by "interlude run-daemon"
*/
func runDaemon() {
	fmt.Println("Interlude daemon started. Listening for Claude Code...")
// graceful shutdown setup
	// create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1) // buffer size 1
	// listen for SIGTERM and SIGINT
	// SIGTERM: from kill or "interlude stop"
	// SIGINT: from ctrl C
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	// goroutine run in background to handle signals
	go func() {
		<- sigChan // block until signal is received
		fmt.Println("Interlude shutting down...")
		// cleanup
		os.Remove(pidFile)
		os.Exit(0)
	}() // call it immediately


// main detection loop
	// create a ticker that ticks every 250ms
	ticker := time.NewTicker(250 * time.Millisecond)
	// init trackers
	wasActive := false
	notified := false
	var sessionStartTime time.Time
	// loop with 250ms interval
	for range ticker.C {
		// check current status
		active := isClaudeActive()
		// transition: idle to active
		if active && !wasActive {
			fmt.Println("Claude started working...")
			notified = false
		}
		// active and past threshold - show intervention UI
		if active && !notified {
			sessionStartTime = getStartTime()
			if time.Since(sessionStartTime) > threshold {
				fmt.Println("Still working... Stay focused!")
				notified = true
			}
		}
		// transition: active to idle
		if !active && wasActive {
			duration := time.Since(sessionStartTime)
			fmt.Printf("Claude Done! It took %v\n", duration)
		}
		wasActive = active
	}

}

/* Help function to check if Claude Code is actively responding
 * @return: true if CC is actively running, false otherwise
*/
func isClaudeActive() bool {
	data, err := os.ReadFile(statusFile)
	if err != nil { // file not found
		return false
	}
	return strings.TrimSpace(string(data)) == "active"
}

/* Helper function to get the start time written by on-start script
 * @return: start time as time.Time, or now if not found
*/
func getStartTime() time.Time {
	data, err := os.ReadFile(startTimeFile)
	if err != nil { // file not found
		fmt.Println("No start time file found. Starting now.")
		return time.Now()
	}
	timestamp, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		fmt.Println("Invalid start time format. Starting now.")
		return time.Now()
	}
	return time.Unix(timestamp, 0)
}
