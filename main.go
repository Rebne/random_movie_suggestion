package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const API_KEY = "3cbd1470"
const SECRET_TOKEN = "secret-token"
const FILENAME = "id_data.txt"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app := movieGenerator()
		app.Render(r.Context(), w)
	})

	r.HandleFunc("/secret/{token}/{action}/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		if token != "your-secret-token" {
			http.NotFound(w, r)
			return
		}
		id := vars["id"]
		action := vars["action"]
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

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
