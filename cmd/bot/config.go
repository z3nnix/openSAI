package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readToken(path string) (string, error) {
	botToken, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(string(botToken))
	if len(token) == 0 {
		return "", fmt.Errorf("token file is empty")
	}

	return token, nil
}

func readNames(path string) (map[string]bool, error) {
	namesFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer namesFile.Close()

	names := make(map[string]bool)
	scanner := bufio.NewScanner(namesFile)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names[name] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}