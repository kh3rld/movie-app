package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	TMDB *TMDBClient
	OMDB *OMDBClient
}

type SearchResult struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Year   string `json:"year"`
	Poster string `json:"poster"`
}

type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
}

type DetailResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Plot        string            `json:"plot"`
	Cast        []string          `json:"cast"`
	ReleaseDate string            `json:"release_date"`
	Poster      string            `json:"poster"`
	Ratings     map[string]string `json:"ratings"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	mediaType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	if q == "" || (mediaType != "movie" && mediaType != "tv") {
		writeError(w, http.StatusBadRequest, "Missing or invalid parameters")
		return
	}
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	respBytes, err := h.TMDB.Search(q, mediaType, page)
	if err != nil {
		writeError(w, http.StatusBadGateway, "Failed to fetch from TMDB: "+err.Error())
		return
	}

	results, total, totalPages, err := parseTMDBResults(respBytes, mediaType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to parse TMDB response")
		return
	}

	resp := SearchResponse{
		Results:    results,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// parseTMDBResults parses the TMDB API response and maps it to SearchResult slice.
func parseTMDBResults(data []byte, mediaType string) ([]SearchResult, int, int, error) {
	var raw struct {
		Results []struct {
			ID         int    `json:"id"`
			Title      string `json:"title"`
			Name       string `json:"name"`
			Release    string `json:"release_date"`
			FirstAir   string `json:"first_air_date"`
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalResults int `json:"total_results"`
		Page         int `json:"page"`
		TotalPages   int `json:"total_pages"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, 0, 0, err
	}
	results := make([]SearchResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		title := r.Title
		if mediaType == "tv" {
			title = r.Name
		}
		year := ""
		if r.Release != "" {
			year = strings.Split(r.Release, "-")[0]
		} else if r.FirstAir != "" {
			year = strings.Split(r.FirstAir, "-")[0]
		}
		poster := ""
		if r.PosterPath != "" {
			poster = "https://image.tmdb.org/t/p/w200" + r.PosterPath
		}
		results = append(results, SearchResult{
			ID:     strconv.Itoa(r.ID),
			Title:  title,
			Year:   year,
			Poster: poster,
		})
	}
	return results, raw.TotalResults, raw.TotalPages, nil
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mediaType := r.URL.Query().Get("type")
	if id == "" || (mediaType != "movie" && mediaType != "tv") {
		writeError(w, http.StatusBadRequest, "Missing or invalid parameters")
		return
	}

	// Fetch TMDB details
	tmdbData, err := h.TMDB.GetDetails(id, mediaType)
	if err != nil {
		writeError(w, http.StatusBadGateway, "Failed to fetch from TMDB: "+err.Error())
		return
	}
	detail, imdbID := parseTMDBDetail(tmdbData, mediaType)

	// Fetch OMDB details if imdbID is available
	if imdbID != "" && h.OMDB != nil {
		omdbData, err := h.OMDB.GetDetails(imdbID)
		if err == nil {
			mergeOMDBDetail(detail, omdbData)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

// parseTMDBDetail parses TMDB detail response and returns a DetailResponse and imdbID
func parseTMDBDetail(data []byte, mediaType string) (*DetailResponse, string) {
	var raw map[string]interface{}
	_ = json.Unmarshal(data, &raw)
	id := ""
	if v, ok := raw["id"].(float64); ok {
		id = strconv.Itoa(int(v))
	}
	title := getString(raw, "title")
	if mediaType == "tv" {
		title = getString(raw, "name")
	}
	plot := getString(raw, "overview")
	release := getString(raw, "release_date")
	if mediaType == "tv" {
		release = getString(raw, "first_air_date")
	}
	poster := ""
	if p, ok := raw["poster_path"].(string); ok && p != "" {
		poster = "https://image.tmdb.org/t/p/w500" + p
	}
	ratings := map[string]string{}
	if v, ok := raw["vote_average"].(float64); ok {
		ratings["tmdb"] = strconv.FormatFloat(v, 'f', 1, 64)
	}
	cast := []string{}
	if credits, ok := raw["credits"].(map[string]interface{}); ok {
		if castArr, ok := credits["cast"].([]interface{}); ok {
			for i, c := range castArr {
				if i >= 8 {
					break
				}
				if m, ok := c.(map[string]interface{}); ok {
					if n, ok := m["name"].(string); ok {
						cast = append(cast, n)
					}
				}
			}
		}
	}
	imdbID := getString(raw, "imdb_id")
	return &DetailResponse{
		ID:          id,
		Title:       title,
		Plot:        plot,
		Cast:        cast,
		ReleaseDate: release,
		Poster:      poster,
		Ratings:     ratings,
	}, imdbID
}

// mergeOMDBDetail merges OMDB data into DetailResponse
func mergeOMDBDetail(detail *DetailResponse, data []byte) {
	var raw map[string]interface{}
	_ = json.Unmarshal(data, &raw)
	if plot, ok := raw["Plot"].(string); ok && plot != "N/A" && plot != "" {
		detail.Plot = plot
	}
	if poster, ok := raw["Poster"].(string); ok && poster != "N/A" && poster != "" {
		detail.Poster = poster
	}
	if ratings, ok := raw["Ratings"].([]interface{}); ok {
		for _, r := range ratings {
			if m, ok := r.(map[string]interface{}); ok {
				src := getString(m, "Source")
				val := getString(m, "Value")
				if src != "" && val != "" {
					switch src {
					case "Internet Movie Database":
						detail.Ratings["imdb"] = val
					case "Rotten Tomatoes":
						detail.Ratings["rotten_tomatoes"] = val
					}
				}
			}
		}
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
