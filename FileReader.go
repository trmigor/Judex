package main

import (
	"io/ioutil"
	"os"		// OS syscalls
	"log"		// Logs
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

	msg, err := ioutil.ReadAll(file)
	return string(msg)
}