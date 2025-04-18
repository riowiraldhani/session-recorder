package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const LOG_DIR = "/var/log/sessions/"

func main() {
	// Avoid infinite recursion
	if os.Getenv("SCRIPT_SESSION") == "true" {
		return
	}

	// Get the current user, fallback to "unknown" if not set
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}

	// Determine session type based on environment
	sessionType := "ssh"
	if strings.Contains(os.Getenv("AWS_EXECUTION_ENV"), "AWS") {
		sessionType = "ssm"
	}

	// Generate a session ID based on the timestamp and user info
	sessionID := time.Now().Format("20060102-150405") + "-" + user + "-" + sessionType
	logFile := LOG_DIR + sessionID + ".log"

	// Log where the session will be recorded
	log.Printf("Recording session to %s\n", logFile)

	// Define a bash wrapper to capture the session's closure
	wrapper := `
trap 'echo "[SESSION CLOSED] $(date)" >> ` + logFile + `' EXIT
`

	// Start a shell with script to capture all output and actions
	cmd := exec.Command("script", "-f", "-q", "-c", wrapper+os.Getenv("SHELL"), logFile)
	cmd.Env = append(os.Environ(), "SCRIPT_SESSION=true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the shell session with the command wrapper
	if err := cmd.Run(); err != nil {
		log.Printf("Session ended with error: %v", err)
	}
}
