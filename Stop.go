package main

import (
	"fmt"		// Output formatting
	"os"		// OS syscalls
	"log"		// Logs
	"context"	// Database disconnection
)

// Stop handles os.Kill and os.Interrupt signals and ends the program if sigs != nil, otherwise just ends the program
func Stop(sigs chan os.Signal, logFile *os.File) {
	// Receiving signals
	var sig os.Signal
	if sigs != nil {
		sig = <-sigs
	}

	switch sig {
	case os.Kill:
		log.Println("Received os.Kill")
	case os.Interrupt:
		log.Println("Received os.Interrupt")
	}

	// Database disconnection
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")

	// Closing the log file
	log.Println("STOPPED")
	fmt.Fprintln(logFile)
	logFile.Close()

	fmt.Println()
	os.Exit(0)
}