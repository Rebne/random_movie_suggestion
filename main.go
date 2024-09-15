package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const API_KEY = "3cbd1470"
const SECRET_TOKEN = "secret-token"
const FILENAME = "id_data.txt"
const PORT = 8080

var ids []string

func init() {
	ids = readFile(FILENAME)
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app := movieGenerator()
		app.Render(r.Context(), w)
	})

	r.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"ids": ids,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
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

	fmt.Printf("Server starting on :%d\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
