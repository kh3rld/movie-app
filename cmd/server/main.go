package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/kh3rld/movie-app/internal/api"
	"github.com/kh3rld/movie-app/internal/cache"
	"github.com/kh3rld/movie-app/internal/config"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := config.LoadConfig()
	tmdb := api.NewTMDBClient(cfg)
	omdb := api.NewOMDBClient(cfg)
	cache := cache.New()
	handler := &api.Handler{
		TMDB:  tmdb,
		OMDB:  omdb,
		Cache: cache,
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/search", handler.Search)
		r.Get("/detail", handler.Detail)
		r.Get("/trending", handler.Trending)
		r.Get("/recommendations", handler.Recommendations)
	})

	r.Get("/api/watchlist", handler.GetWatchlist)
	r.Post("/api/watchlist", handler.AddToWatchlist)
	r.Delete("/api/watchlist", handler.RemoveFromWatchlist)
	r.Patch("/api/watchlist", handler.MarkWatched)

	// Serve static files from ./web
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	log.Printf("Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
