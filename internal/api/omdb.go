package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kh3rld/movie-app/internal/config"
)

type OMDBClient struct {
	APIKey string
}

func NewOMDBClient(cfg *config.Config) *OMDBClient {
	return &OMDBClient{APIKey: cfg.OMDBApiKey}
}

func (c *OMDBClient) GetDetails(imdbID string) ([]byte, error) {
	u, _ := url.Parse("https://www.omdbapi.com/")
	q := u.Query()
	q.Set("apikey", c.APIKey)
	q.Set("i", imdbID)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by OMDB")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OMDB error: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}
