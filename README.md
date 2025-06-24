# Movie/TV Show Discovery Web App

A web application for discovering movies and TV shows, powered by TMDB and OMDB APIs.

### 1. Clone the repository
```sh
git clone https://github.com/kh3rld/movie-app.git
cd movie-app
```

### 2. Set up environment variables
Create a `.env` file in the project root:
```
TMDB_API_KEY=your_tmdb_api_key
OMDB_API_KEY=your_omdb_api_key
```

### 3. Install dependencies
```sh
go mod tidy
```

### 4. Run server
```sh
go run cmd/server/main.go
```

### 5. View app
Open `web/index.html` in your browser, or visit `http://localhost:8080` if served by the backend.


## API Docs

### `/api/search`
- **GET**: Search for movies or TV shows
- **Query params:**
  - `q` (string, required): Search query
  - `type` (string, required): `movie` or `tv`
  - `page` (int, optional): Page number
  - `genre` (string, optional): Genre ID
- **Response:**
```json
{
  "results": [ { "id": "...", "title": "...", "year": "...", "poster": "..." } ],
  "total": 100,
  "page": 1,
  "total_pages": 10
}
```

### `/api/detail`
- **GET**: Get detailed info for a movie or TV show
- **Query params:**
  - `id` (string, required): TMDB ID
  - `type` (string, required): `movie` or `tv`
- **Response:**
```json
{
  "id": "...",
  "title": "...",
  "plot": "...",
  "cast": ["..."],
  "release_date": "...",
  "poster": "...",
  "ratings": { "tmdb": "...", "imdb": "...", "rotten_tomatoes": "..." }
}
```

### `/api/genres`
- **GET**: List genres for movies or TV shows
- **Query params:**
  - `type` (string, required): `movie` or `tv`
- **Response:**
```json
{
  "genres": [ { "id": "...", "name": "..." } ]
}
```


## Contributing
1. Fork the repo and create your feature branch (`git checkout -b feature/your-feature`)
2. Commit your changes with clear messages
3. Push to the branch (`git push origin feature/your-feature`)
4. Open a Pull Request


## License
[MIT](LICENSE)
