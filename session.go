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
	if os.Getenv("SCRIPT_SESSION") == "true" {
		return
	}

	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}

	sessionType := "ssh"
	if strings.Contains(os.Getenv("AWS_EXECUTION_ENV"), "AWS") {
		sessionType = "ssm"
	}

	sessionID := time.Now().Format("20060102-150405") + "-" + user + "-" + sessionType
	logFile := LOG_DIR + sessionID + ".log"

	log.Printf("Recording session to %s\n", logFile)

	// Bash wrapper to catch logout and cleanly exit
	wrapper := `
trap 'echo "[SESSION CLOSED] $(date)"' EXIT
` + os.Getenv("SHELL")

	// Run the wrapped shell with `script`
	cmd := exec.Command("script", "-f", "-q", "-c", wrapper, logFile)
	cmd.Env = append(os.Environ(), "SCRIPT_SESSION=true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Session ended with error: %v", err)
	}
}
