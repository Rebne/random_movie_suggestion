package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rebne/movie_generator/data"
	"github.com/Rebne/movie_generator/handlers"
	"github.com/Rebne/movie_generator/models"

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
	idData       models.IDdata
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT = os.Getenv("PORT")
	// if FILEPATH == "" || PORT == "" {
	if PORT == "" {
		log.Fatal("Missing required environment variables from main.go")
	}
	db := data.InitDB()
	err = db.Ping()
	if err != nil {
		log.Fatal("Something went wrong with the database in main init()")
	}
	// idData, err = services.ReadIDData(FILEPATH)
	// if err != nil {
	// 	log.Fatal("Error reading JSON from file", err.Error())
	// }
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.HomeHandler)

	r.Get("/api/data/length", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTotalMovieCountHandler(w, r, &idData)
	})

	r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMovieDataHandler(w, r, &idData)
	})

	r.Post("/api/data/new", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateMovieListHandler(w, r, &idData)
	})

	r.Post("/generate", handlers.GenerateMovieCardHandler)

	r.Get("/secret/{token}/showlist", handlers.ShowMovieListHandler)

	r.Post("/secret/{token}/{action}/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.ManageMovieListHandler(w, r, &idData)
	})

	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/public/*", http.StripPrefix("/public/", fileServer))

	fmt.Printf("Server starting on :%s\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}
