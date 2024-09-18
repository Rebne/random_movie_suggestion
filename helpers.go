package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ID struct {
	ID    int `json:"id"`
	Index int `json:"index"`
}

type IDdata struct {
	Length int  `json:"length"`
	IDs    []ID `json:"ids"`
}

type MovieData struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Plot    string `json:"Plot"`
	Runtime string `json:"Runtime"`
	Poster  string `json:"Poster"`
	Genre   string `json:"Genre"`
}

func fetchMovieData(id string) (MovieData, error) {
	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", id, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		return MovieData{}, err
	}
	defer resp.Body.Close()
	var data MovieData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return MovieData{}, err
	}
	if data.Title == "" {
		return MovieData{}, fmt.Errorf("error, no movie data returned from OMDb API")
	}
	return data, nil
}

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

func readIDData(filename string) (IDdata, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return IDdata{}, err
	}

	var data IDdata
	err = json.Unmarshal(file, &data)
	if err != nil {
		return IDdata{}, err
	}
	return data, nil
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

func writeIDdata(filename string, data IDdata) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
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
