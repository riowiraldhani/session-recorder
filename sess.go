package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	username := os.Getenv("USER") // or get from os/user
	timestamp := time.Now().Format("20060102-150405")
	logPath := "/var/log/sessions/" + timestamp + "-" + username + "-ssh.log"

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log.Printf("Recording session to %s", logPath)

	cmd := exec.Command("script", "-q", "-f", logPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	// 🔥 This will always run, even if session closes due to tab/browser exit
	logFile.WriteString("[SESSION CLOSED] " + time.Now().UTC().Format(time.RFC1123) + "\n")

	if err != nil {
		log.Printf("Script exited with error: %v", err)
	} else {
		log.Printf("Script exited cleanly.")
	}
}
