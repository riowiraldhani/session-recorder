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
	// Prevent recursion: don't run if already in a script session
	if os.Getenv("SCRIPT_SESSION") == "true" {
		return
	}

	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}

	sessionType := "ssh"
	if strings.Contains(os.Getenv("SSH_CONNECTION"), "ssm") {
		sessionType = "ssm"
	}

	sessionID := time.Now().Format("20060102-150405") + "-" + user + "-" + sessionType
	logFile := LOG_DIR + sessionID + ".log"

	log.Printf("Recording session to %s\n", logFile)

	cmd := exec.Command("script", "-f", "-q", "-c", os.Getenv("SHELL"), logFile)
	// Prevent recursion in subprocess
	cmd.Env = append(os.Environ(), "SCRIPT_SESSION=true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Error starting session recording: %v", err)
	}
}
