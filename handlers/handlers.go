package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Rebne/movie_generator/data"
	"github.com/Rebne/movie_generator/helpers"
	"github.com/Rebne/movie_generator/models"
	"github.com/Rebne/movie_generator/services"
	"github.com/Rebne/movie_generator/web/views/home"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

var (
	API_KEY      string
	SECRET_TOKEN string
	FILEPATH     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	API_KEY = os.Getenv("API_KEY")
	SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
	FILEPATH = os.Getenv("FILEPATH")

	if API_KEY == "" || SECRET_TOKEN == "" || FILEPATH == "" {
		log.Fatal("Missing required environment variables")
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	app := home.App("ElsaRene Random Movie Generator")
	app.Render(r.Context(), w)
}

func GetTotalMovieCountHandler(w http.ResponseWriter, r *http.Request, idData *models.IDdata) {
	length, err := data.GetTableLengthDB()
	if err != nil {
		http.Error(w, "Error getting table length from database", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"length": length,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetMovieDataHandler(w http.ResponseWriter, r *http.Request, idData *models.IDdata) {
	total, err := data.GetTableLengthDB()
	if err != nil {
		http.Error(w, "Error getting table length from database", http.StatusInternalServerError)
	}
	ids, err := data.GetAllMovieIdsDB()
	if err != nil {
		http.Error(w, "Error getting table length from database", http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"total": total,
		"ids":   ids,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func UpdateMovieListHandler(w http.ResponseWriter, r *http.Request, idData *models.IDdata) {
	var requestData struct {
		CurrentLength string `json:"currentLength"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received currentLength: %s\n", requestData.CurrentLength)
	currLength, err := strconv.Atoi(requestData.CurrentLength)
	if err != nil {
		http.Error(w, "Error converting currentLength to integer", http.StatusBadRequest)
		return
	}
	newIDs, err := helpers.GetNewIDs(currLength, idData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"newLength": idData.Length,
		"newIDs":    newIDs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GenerateMovieCardHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	movieID := r.FormValue("movieID")
	data, err := services.FetchMovieData(movieID, API_KEY)
	if err != nil {
		fmt.Println("Error fetching movie data:", err)
		http.Error(w, "Failed to fetch movie data", http.StatusInternalServerError)
		return
	}

	component := home.MovieCard(data)
	component.Render(r.Context(), w)
}

func ShowMovieListHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token != SECRET_TOKEN {
		http.NotFound(w, r)
		return
	}

	fileContents, err := os.ReadFile(FILEPATH)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(fileContents)
}

func ManageMovieListHandler(w http.ResponseWriter, r *http.Request, idData *models.IDdata) {
	token := chi.URLParam(r, "token")
	if token != SECRET_TOKEN {
		http.NotFound(w, r)
		return
	}
	id := chi.URLParam(r, "id")
	action := chi.URLParam(r, "action")
	if !helpers.IsValidIMDbID(id) {
		http.Error(w, "Not a valid IMDb ID", http.StatusBadRequest)
		return
	}
	switch action {
	case "delete":
		err := services.RemoveID(idData, id, FILEPATH)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting ID: %s", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ID: %s DELETED FROM IDS\n", id)
	case "add":
		if !helpers.IdExists(idData, id) {
			http.Error(w, "ID already exists in the JSON data", http.StatusBadRequest)
			return
		}
		err := services.AddID(idData, id, FILEPATH)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding ID: %s", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ID: %s ADDED TO IDS\n", id)
	default:
		http.Error(w, fmt.Sprintf("Error: invalid action '%s' for id %s", action, id), http.StatusBadRequest)
	}
}
