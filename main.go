package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myapp/utils"
	"net/http"
	"strings"
)

func main() {
	log.SetOutput(io.Discard)
	err := utils.ConnectToDatabase()
	if err != nil {
		log.Fatal("Unable to connect to the Database", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		imdbID := r.URL.Query().Get("imdb_id")
		if imdbID == "" {
			http.Error(w, "Missing query parameter 'imdb_id'", http.StatusBadRequest)
			return
		}

		data, err := utils.Scrape(imdbID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Missing query parameter 'query'", http.StatusBadRequest)
		}

		data, err := utils.Search(strings.ReplaceAll(query, " ", "+"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	http.HandleFunc("/cast", func(w http.ResponseWriter, r *http.Request) {

		imdbID := r.URL.Query().Get("imdb_id")
		if imdbID == "" {
			http.Error(w, "Missing query parameter 'imdb_id'", http.StatusBadRequest)
			return
		}

		data, err := utils.ScrapeCast(imdbID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	// Start the server
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
