package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var (
	API_KEY      string
	SECRET_TOKEN string
	FILENAME     string
	PORT         string
	idData       IDdata
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	API_KEY = os.Getenv("API_KEY")
	SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
	FILENAME = os.Getenv("FILENAME")
	PORT = os.Getenv("PORT")

	if API_KEY == "" || SECRET_TOKEN == "" || FILENAME == "" || PORT == "" {
		log.Fatal("Missing required environment variables")
	}
	idData, err = readIDData(FILENAME)
	if err != nil {
		log.Fatal("Error reading JSON from file", err.Error())
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		app := app("ElsaRene Random Movie Generator")
		app.Render(r.Context(), w)
	})

	r.Get("/api/data/length", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"length": idData.Length,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"ids": idData.IDs,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	r.Post("/generate", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		movieID := r.FormValue("movieID")
		data, err := fetchMovieData(movieID)
		if err != nil {
			fmt.Println("Error fetching movie data:", err)
			http.Error(w, "Failed to fetch movie data", http.StatusInternalServerError)
			return
		}

		component := movieCard(data)
		component.Render(r.Context(), w)
	})

	r.Get("/secret/{token}/{action}/{id}", func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token != SECRET_TOKEN {
			http.NotFound(w, r)
			return
		}
		id := chi.URLParam(r, "id")
		action := chi.URLParam(r, "action")
		switch action {
		case "delete":
			err := removeID(idData, id, FILENAME)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error deleting ID: %s", err), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ID: %s DELETED FROM IDS", id)
		case "add":
			// Also appending to current working server instance
			err := addID(idData, id, FILENAME)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error adding ID: %s", err), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ID: %s ADDED TO IDS", id)
		default:
			http.Error(w, fmt.Sprintf("Error: invalid action '%s' for id %s", action, id), http.StatusBadRequest)
		}
	})

	fileServer := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fileServer))

	fmt.Printf("Server starting on :%s\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}
