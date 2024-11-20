package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rebne/movie_generator/models"
	"github.com/joho/godotenv"
)

var API_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	API_KEY = os.Getenv("API_KEY")
	if API_KEY == "" {
		log.Fatal("Missing required environment variables")
	}
}

func FetchMovieData(id string) (models.MovieData, error) {
	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", id, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		return models.MovieData{}, err
	}
	defer resp.Body.Close()
	var data models.MovieData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return models.MovieData{}, err
	}
	if data.Title == "" {
		return models.MovieData{}, fmt.Errorf("error, no movie data returned from OMDb API")
	}
	return data, nil
}
