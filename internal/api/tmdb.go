package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kh3rld/movie-app/internal/config"
)

type TMDBClient struct {
	APIKey string
}

func NewTMDBClient(cfg *config.Config) *TMDBClient {
	return &TMDBClient{APIKey: cfg.TMDBApiKey}
}

func (c *TMDBClient) Search(query, mediaType string, page int) ([]byte, error) {
	base := "https://api.themoviedb.org/3/search/" + mediaType
	u, _ := url.Parse(base)
	q := u.Query()
	q.Set("api_key", c.APIKey)
	q.Set("query", query)
	q.Set("page", fmt.Sprintf("%d", page))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by TMDB")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDB error: %s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *TMDBClient) GetDetails(id, mediaType string) ([]byte, error) {
	base := "https://api.themoviedb.org/3/" + mediaType + "/" + id
	u, _ := url.Parse(base)
	q := u.Query()
	q.Set("api_key", c.APIKey)
	q.Set("append_to_response", "credits")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by TMDB")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDB error: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}
