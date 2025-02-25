package main

import (
    "bufio"
    "log"
    "os"
    "strings"
)

func appendMessageToFile(filename string, message string, chatID string) {
    log.Printf("check ChatID: %s", chatID)

    if strings.Contains(message, "http://") || 
       strings.Contains(message, "https://") || 
       strings.Contains(message, "@") {
        log.Println("Links found. Skip")
        return
    }

    if isChatNeglected("neglected.bot", chatID) {
        log.Printf("Chat: %s ignore-list. Skip", chatID)
        return
    }

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Cannot open file %s: %v", filename, err)
        return
    }
    defer file.Close()

    if _, err := file.WriteString(message + "\n"); err != nil {
        log.Printf("Write error %s: %v", filename, err)
    } else {
        log.Printf("Message write in %s", filename)
    }
}

func isChatNeglected(filename, targetID string) bool {
    file, err := os.Open(filename)
    if err != nil {
        if os.IsNotExist(err) {
            log.Println("neglected.bot not found")
        }
        return false
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == targetID {
            log.Printf("%s", targetID)
            return true
        }
    }
    return false
}