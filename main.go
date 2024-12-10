package main

import (
	"log"
	"myapp/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Suppress Gin's default logging if needed
	gin.SetMode(gin.ReleaseMode)

	// Connect to the database
	err := utils.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Unable to connect to the Database: %v", err)
	}

	// Create a Gin router
	router := gin.Default()

	// Routes
	router.GET("/title/:imdb_id", func(c *gin.Context) {
		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'imdb_id'"})
			return
		}

		data, err := utils.Scrape(imdbID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	router.GET("/search", func(c *gin.Context) {
		query := c.Query("title")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter 'title'"})
			return
		}

		data, err := utils.Search(strings.ReplaceAll(query, " ", "+"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	router.GET("/title/:imdb_id/cast", func(c *gin.Context) {
		imdbID := c.Param("imdb_id")
		if imdbID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter 'imdb_id'"})
			return
		}

		data, err := utils.ScrapeCast(imdbID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
