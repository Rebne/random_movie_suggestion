package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const API_KEY = "3cbd1470"

func main() {
	fmt.Println(API_KEY)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	r.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
