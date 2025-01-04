package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OriginalImage struct {
	Source string `json:"source"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type RandomWikiResponse struct {
	Title         string        `json:"title"`
	OriginalImage OriginalImage `json:"originalimage,omitempty"`
	Description   string        `json:"description,omitempty"`
	Extract       string        `json:"extract"`
}

func NewRandomWikiResponse() (*RandomWikiResponse, error) {
	resp, err := http.Get("https://en.wikipedia.org/api/rest_v1/page/random/summary")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wikiResponse RandomWikiResponse
	err = json.Unmarshal(body, &wikiResponse)
	if err != nil {
		return nil, err
	}

	return &wikiResponse, nil
}
