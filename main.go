package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"myapp/utils"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	log.SetOutput(io.Discard)

	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file: %w", err)
		return
	}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on http://localhost:" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Error starting server:", err)
	}
}
