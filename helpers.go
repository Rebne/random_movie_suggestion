package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func formatRuntimeString(duration string) string {
	result := ""
	index := 0
	for index < len(duration) {
		if !isDigit(duration[index]) {
			break
		}
		result += string(duration[index])
		index++
	}

	durationAsInt, err := strconv.Atoi(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting string %s\n", result)
		return strings.ReplaceAll(duration, " ", "")
	}
	hours := durationAsInt / 60
	minutes := durationAsInt % 60

	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func isValidIMDbID(id string) bool {
	pattern := `^tt\d{7,8}$`
	match, _ := regexp.MatchString(pattern, id)
	return match
}

func readFile(filename string) ([]string, error) {
	var ids []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func appendToFile(line string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if !isValidIMDbID(line) {
		return fmt.Errorf("invalid IMDb ID: %s", line)
	}
	_, err = file.WriteString(line + "\n")
	if err != nil {
		return err
	}
	return nil
}

func deleteFromFile(lineToDelete string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if line != lineToDelete {
			lines = append(lines, line)
		} else {
			found = true
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("line to delete not found: %s", lineToDelete)
	}

	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}
