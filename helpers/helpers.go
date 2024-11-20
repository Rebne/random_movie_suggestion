package helpers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Rebne/movie_generator/data"
	"github.com/Rebne/movie_generator/models"
)

func GetMovieIDs(data *models.IDdata) []string {
	movieIDs := make([]string, data.Length)
	for i, id := range data.IDs {
		movieIDs[i] = id.MovieID
	}
	return movieIDs
}

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func FormatRuntimeString(duration string) string {
	result := ""
	index := 0
	for index < len(duration) {
		if !IsDigit(duration[index]) {
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

func IsValidIMDbID(id string) bool {
	pattern := `^tt\d{7,8}$`
	match, _ := regexp.MatchString(pattern, id)
	return match
}

func GetNewIDs(prevLength int) (models.IDdata, error) {
	var temp []models.ID
	movieData, err := data.GetAllMoviesDB()
	if err != nil {
		return models.IDdata{}, err
	}
	for _, item := range movieData.IDs {
		if item.Index >= prevLength {
			temp = append(temp, item)
		}
	}
	if len(temp) == 0 {
		return models.IDdata{}, fmt.Errorf("error no new IDS found")
	}
	return models.IDdata{IDs: temp, Length: movieData.Length}, nil
}

func IdExists(data *models.IDdata, id string) bool {
	for _, existingID := range data.IDs {
		if existingID.MovieID == id {
			return true
		}
	}
	return false
}
