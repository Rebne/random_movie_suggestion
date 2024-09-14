package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const API_KEY = "3cbd1470"

func main() {
	fmt.Println(API_KEY)

	component := hello("World")
	component.Render(context.Background(), os.Stdout)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := hello("World")
		component.Render(r.Context(), w)
	})

	r.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
