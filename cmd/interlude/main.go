package main

import (
	"fmt" // std lib, print to console
	"os" // std lib, file operations: readfile, userhome
	"strconv" // std lib, string-to-int conversion
	"strings" // std lib, string manipulation
	"time" // std lib, time operations

	// "github.com/shirou/gopsutil/v3/process" // process monitoring
)

var (
	home, _ = os.UserHomeDir()
	statusFile = home + "/.interlude/status"
	startTimeFile = home + "/.interlude/start_time"
)

const threshold = 5 * time.Second

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

/* Main function */
func main() {
	fmt.Println("Interlude daemon started. Listening for Claude Code...")
	// create a ticker that ticks every 250ms
	ticker := time.NewTicker(250 * time.Millisecond)
	// init trackers
	wasActive := false
	notified := false
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
			startTime := getStartTime()
			if time.Since(startTime) > threshold {
				fmt.Println("Still working... Stay focused!")
				notified = true
			}
		}
		// transition: active to idle
		if !active && wasActive {
			duration := time.Since(getStartTime())
			fmt.Printf("Done! Took %v\n", duration)
		}
		wasActive = active
	}
}
