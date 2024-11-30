package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type YoutubeSearchResults struct {
	Kind       string `json:"kind"`
	Etag       string `json:"etag"`
	RegionCode string `json:"regionCode"`
	PageInfo   struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

// QueryYoutube returns youtube search results for query, with the query being wrapped in quotes
func QueryYoutube(query, apiKey string) (*YoutubeSearchResults, error) {
	if query == "" || apiKey == "" {
		return nil, errors.New("empty query or api key provided")
	}

	// Create the request
	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Add the query parameters
	q := req.URL.Query()
	q.Add("part", "id")
	q.Add("q", `"`+query+`"`)
	q.Add("key", apiKey)
	req.URL.RawQuery = q.Encode()

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var results YoutubeSearchResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &results, nil
}
