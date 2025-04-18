package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	LOG_DIR = "/var/log/sessions/"
)

func main() {
	// Detect session type (SSH or SSM)
	sessionType := "ssh"
	if strings.Contains(os.Getenv("SSH_CONNECTION"), "ssm") {
		sessionType = "ssm"
	}

	// Unique session ID
	sessionID := time.Now().Format("20060102-150405") + "-" + os.Getenv("USER") + "-" + sessionType
	logFile := LOG_DIR + sessionID + ".log"

	// Ensure log directory exists
	if err := os.MkdirAll(LOG_DIR, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Start recording session
	cmd := exec.Command("script", "-f", logFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Recording session to %s", logFile)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to record session: %v", err)
	}

	log.Printf("Session recording saved: %s", logFile)
}
