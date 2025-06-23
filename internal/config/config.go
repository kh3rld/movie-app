package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	TMDBApiKey string
	OMDBApiKey string
}

func LoadConfig() *Config {
	// Load .env file from project root
	if err := godotenv.Load(filepath.Join(".env")); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	tmdb := os.Getenv("TMDB_API_KEY")
	omdb := os.Getenv("OMDB_API_KEY")
	if tmdb == "" || omdb == "" {
		log.Fatal("TMDB_API_KEY and OMDB_API_KEY must be set in environment or .env file")
	}
	return &Config{
		TMDBApiKey: tmdb,
		OMDBApiKey: omdb,
	}
}
