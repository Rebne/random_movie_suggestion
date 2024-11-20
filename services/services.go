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

func AddID(dataSet *models.IDdata, movieID string, FILEPATH string) error {
	movieData, err := FetchMovieData(movieID)
	if err != nil {
		return fmt.Errorf("failed to fetch moviedata from OMD")
	}
	dataSet.IDs = append(dataSet.IDs, models.ID{MovieID: movieID, Index: dataSet.Length, Title: movieData.Title})
	dataSet.Length++
	return WriteIdData(FILEPATH, dataSet)
}

func RemoveID(dataSet *models.IDdata, movieID string, FILEPATH string) error {
	for i := range dataSet.IDs {
		if dataSet.IDs[i].MovieID == movieID {
			dataSet.IDs = append(dataSet.IDs[:i], dataSet.IDs[i+1:]...)
			dataSet.Length--
			dataSet.ReIndexMovieIDs()
			return WriteIdData(FILEPATH, dataSet)
		}
	}
	return fmt.Errorf("ID not found")
}

func ReadIDData(FILEPATH string) (models.IDdata, error) {
	file, err := os.ReadFile(FILEPATH)
	if err != nil {
		return models.IDdata{}, err
	}

	var data models.IDdata
	err = json.Unmarshal(file, &data)
	if err != nil {
		return models.IDdata{}, err
	}
	return data, nil
}

func WriteIdData(FILEPATH string, data *models.IDdata) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(FILEPATH, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
