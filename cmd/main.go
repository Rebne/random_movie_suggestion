package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rebne/movie_generator/handlers"

	// "github.com/Rebne/movie_generator/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var (
	API_KEY      string
	SECRET_TOKEN string
	FILEPATH     string
	PORT         string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Missing required environment variables from main.go")
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.HomeHandler)

	r.Get("/api/data/length", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTotalMovieCountHandler(w, r)
	})

	r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMovieDataHandler(w, r)
	})

	r.Post("/api/data/new", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMovieListHandler(w, r)
	})

	r.Post("/generate", handlers.GenerateMovieCardHandler)

	r.Get("/secret/{token}/showlist", handlers.ShowMovieListHandler)

	r.Get("/secret/{token}/{action}/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.ManageMovieListHandler(w, r)
	})

	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/public/*", http.StripPrefix("/public/", fileServer))

	fmt.Printf("Server starting on :%s\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}
