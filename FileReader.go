package main

import (
	"os"		// OS syscalls
	"io"		// EOF
	"log"		// Logs
	"fmt"		// I/O formatting
	"bufio"		// Buffered I/O
)

// FileReader reads file and saves its text to a string
func FileReader(path string) string {
	// Opening
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)

	if err != nil {
		log.Println("FileReader:", err)
		return ""
	}

	defer file.Close()

	// Creating a new reader
	reader := bufio.NewReader(file)

	msg := ""

	// Reading line by line
	for {
		input, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Reading:", err)
			return ""
		}

		msg += fmt.Sprint(input)
	}
	return msg
}