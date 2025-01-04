package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type Article struct {
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	Image             string    `json:"image"`
	AltTextImage      string    `json:"alt_text_image"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	Language          string    `json:"language"`
	Tags              []string  `json:"tags"`
}

func NewArticle() Article {
	wikiResponse, err := NewRandomWikiResponse()
	AvailableTags := []string{
		"History", "Science", "Technology", "Art",
		"Culture", "Geography", "Politics", "Sports",
		"Health", "Education", "Entertainment", "Nature",
		"Economy", "Religion", "Philosophy", "Society",
	}
	AvailableLanguages := []string{"fr", "en", "de", ""}

	if err != nil {
		log.Panic().Err(err).Msg("Failed to get random wiki response")
	}

	// Randomly select up to 5 tags
	numTags := rand.Intn(5)
	selectedTags := make(map[string]struct{})
	for len(selectedTags) < numTags {
		tag := AvailableTags[rand.Intn(len(AvailableTags))]
		selectedTags[tag] = struct{}{}
	}

	tags := make([]string, 0, len(selectedTags))
	for tag := range selectedTags {
		tags = append(tags, tag)
	}

	// Randomly select a language
	language := AvailableLanguages[rand.Intn(len(AvailableLanguages))]

	return Article{
		Title:             wikiResponse.Title,
		Content:           wikiResponse.Extract,
		Image:             wikiResponse.OriginalImage.Source,
		AltTextImage:      wikiResponse.Description,
		CreationTimestamp: time.Now(),
		Language:          language,
		Tags:              tags,
	}
}

func generateArticles(count int) []Article {
	articles := make([]Article, count)

	// Leave content empty for user to implement
	for i := 0; i < count; i++ {
		articles[i] = NewArticle()
	}
	return articles
}

func storeArticles(es *elasticsearch.Client, articles []Article) error {
	ctx := context.Background()

	for _, article := range articles {
		// Convert the article to JSON
		data, err := json.Marshal(article)
		if err != nil {
			return fmt.Errorf("failed to marshal article: %v", err)
		}

		// Index the article into Elasticsearch
		res, err := es.Index(
			"articles",                      // Index name
			strings.NewReader(string(data)), // Document body
			es.Index.WithContext(ctx),       // Context for request
			es.Index.WithRefresh("true"),    // Immediate refresh for visibility
		)
		if err != nil {
			return fmt.Errorf("failed to index article: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("error indexing article: %s", res.String())
		}
	}

	return nil
}

func encodeArticle(article Article) *strings.Reader {
	// Convert Article to JSON
	data, _ := json.Marshal(article)
	return strings.NewReader(string(data))
}

func createArticleIndex(es *elasticsearch.Client) error {
	ctx := context.Background()

	// Define the index mapping
	mapping := `{
		"mappings": {
			"properties": {
				"title": { "type": "text" },
				"content": { "type": "text" },
				"image": { "type": "binary" },
				"alt_text_image": { "type": "text" },
				"creation_timestamp": { "type": "date" },
				"language": { "type": "keyword" },
				"tags": { "type": "keyword" }
			}
		}
	}`

	// Create the index
	res, err := es.Indices.Create(
		"articles",
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
		es.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}
	defer res.Body.Close()

	// Check response
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}
