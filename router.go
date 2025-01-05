package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func runServer(esURL string, esApiKey string, port int) error {
	// Initialize Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
		APIKey:    esApiKey,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	// Check Elasticsearch connection
	res, err := es.Ping()
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to connect to Elasticsearch: %v", err)
	}
	log.Info().Msg("Connected to Elasticsearch")

	// Initialize Gin router
	router := gin.Default()

	router.Use(logger.SetLogger())

	// Add POST /createIndex route
	router.POST("/createIndex", func(c *gin.Context) {
		err := createArticleIndex(es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Index 'articles' created successfully"})
	})

	// Add POST /generate route
	router.POST("/generate", func(c *gin.Context) {
		// Parse query parameter
		countStr := c.Query("count")
		count, err := strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		// Generate and store documents
		articles := generateArticles(count)
		err = storeArticles(es, articles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store articles"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully added %d articles", count)})
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	log.Info().Msgf("Starting server on %s...", addr)
	return router.Run(addr)
}
