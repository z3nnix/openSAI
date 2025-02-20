package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func appendMessageToFile(filename, message string) {
	if strings.Contains(message, "http://") || strings.Contains(message, "https://") || strings.Contains(message, "@") {
		log.Println("Message contains link or username! Skipping.")
		return
	}

	file, err := os.Open(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error opening file for reading:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == message {
			log.Println("Message already exists in the file! Skipping.")
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading file:", err)
		return
	}

	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		log.Println("Error writing to file:", err)
	} else {
		log.Println("Message successfully appended to the file.")
	}
}