package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Start the session log with a session start message
	logger.Println("[SESSION STARTED]", time.Now().Format(time.RFC3339))

	// Create a channel to capture OS signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wait for the termination signal
	sigReceived := <-signalChannel
	logger.Printf("[SESSION CLOSED] %s - Received signal: %s\n", time.Now().Format(time.RFC3339), sigReceived)

	// Gracefully exit
	fmt.Println("Session closed. Logs have been saved.")
}
