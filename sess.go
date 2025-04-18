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

	// Log session start
	logger.Println("[SESSION STARTED]", time.Now().Format(time.RFC3339))

	// Capture signals for session termination (Ctrl+C, SIGTERM, etc.)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Gracefully handle session close
	go func() {
		sigReceived := <-signalChannel
		logger.Printf("[SESSION CLOSED] %s - Received signal: %s\n", time.Now().Format(time.RFC3339), sigReceived)
		// Exit after logging session close
		os.Exit(0)
	}()

	// Simulate interactive session (you can replace this with real work)
	// You can also use this as a place for interactive sessions to run their tasks
	for {
		// Simulating a long-running process
		// In a real use case, replace this with your actual session tasks
		time.Sleep(1 * time.Second)
	}
}
