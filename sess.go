package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"os/exec"
)

func main() {
	// Define the log file path and create/open the log file
	logFilePath := fmt.Sprintf("/var/log/sessions/%s-%s-%s.log", time.Now().Format("20060102"), time.Now().Format("150405"), "session")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new logger
	logger := log.New(logFile, "", log.LstdFlags)

	// Log session start
	logger.Println("[SESSION STARTED]", time.Now().Format(time.RFC3339))

	// Capture signals for session termination (Ctrl+C, SIGTERM, etc.)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Run a bash shell as a subprocess to capture all commands, including 'sudo su'
	cmd := exec.Command("/bin/bash") // You can use '/bin/bash' or any other shell
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the shell session
	if err := cmd.Start(); err != nil {
		logger.Printf("Error starting shell session: %v\n", err)
		os.Exit(1)
	}

	// Go routine to handle session closure
	go func() {
		sigReceived := <-signalChannel
		logger.Printf("[SESSION CLOSED] %s - Received signal: %s\n", time.Now().Format(time.RFC3339), sigReceived)
		cmd.Process.Kill() // Ensure that the subprocess (bash/shell) is killed
		os.Exit(0)
	}()

	// Wait for the shell to exit
	if err := cmd.Wait(); err != nil {
		logger.Printf("Error waiting for shell process: %v\n", err)
	}

}
