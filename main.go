package main

import (
	"encoding/json"
	"io"
	"log"
	"myapp/utils"
	"net/http"
)

func main() {
	log.SetOutput(io.Discard)
	err := utils.ConnectToDatabase()
	if err != nil {
		log.Fatal("Unable to connect to the Database", err)
		return
	}

	http.HandleFunc("/title", func(w http.ResponseWriter, r *http.Request) {
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
		query := r.URL.Query().Get("title")
		if query == "" {
			http.Error(w, "Missing query parameter 'title'", http.StatusBadRequest)
			return
		}

		data, err := utils.Search(query)
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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
