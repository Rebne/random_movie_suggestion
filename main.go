package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const API_KEY = "3cbd1470"
const SECRET_TOKEN = "secret-token"
const FILENAME = "id_data.txt"
const PORT = 8080

var ids []string

type MovieData struct {
	Poster  string `json:"Poster"`
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Plot    string `json:"Plot"`
	Runtime string `json:"Runtime"`
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
	return data, nil
}

func init() {
	ids = readFile(FILENAME)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		app := movieGenerator()
		app.Render(r.Context(), w)
	})

	r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"ids": ids,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	r.Post("/clicked", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		testMovieData := MovieData{
			Poster: "https://m.media-amazon.com/images/M/MV5BNzA5ZDNlZWMtM2NhNS00NDJjLTk4NDItYTRmY2EwMWZlMTY3XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
			Title:  "Lorde of the Rings The Fellowship of the Ring",
		}
		// id := r.FormValue("id")
		// if id == "" {
		// 	http.Error(w, "No ID provided", http.StatusBadRequest)
		// 	return
		// }
		// fmt.Printf("Received ID: %s\n", id)
		// data, err := fetchMovieData(id)
		// if err != nil {
		// 	http.Error(w, "Failed to fetch movie data", http.StatusInternalServerError)
		// 	return
		// }
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<img id='image' src='%s' alt='Movie %s'>", testMovieData.Poster, testMovieData.Title)
	})

	r.Get("/secret/{token}/{action}/{id}", func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		if token != SECRET_TOKEN {
			http.NotFound(w, r)
			return
		}
		id := chi.URLParam(r, "id")
		action := chi.URLParam(r, "action")
		if action == "delete" {
			deleteFromFile(id, FILENAME)
			fmt.Fprintf(w, "Deleted %s\n", id)
		} else if action == "add" {
			appendToFile(id, FILENAME)
			fmt.Fprintf(w, "Added %s\n", id)
		} else {
			http.NotFound(w, r)
		}
	})

	fileServer := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fileServer))

	fmt.Printf("Server starting on :%d\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
