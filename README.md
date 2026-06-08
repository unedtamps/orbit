<p align="center">
  <img src="static/favicon.svg" width="120" alt="OrbitSearch Logo">
</p>

<h1 align="center">OrbitSearch</h1>

<p align="center">Content discovery API — Search for movies and TV series</p>

<p align="center">
  <a href="https://github.com/unedtamps/orbit/actions/workflows/build.yaml">
    <img src="https://github.com/unedtamps/orbit/actions/workflows/build.yaml/badge.svg" alt="Build Status">
  </a>
  <a href="https://hub.docker.com/r/unedotamps/orbit-search">
    <img src="https://img.shields.io/docker/pulls/unedotamps/orbit-search?label=Docker%20Pulls" alt="Docker Pulls">
  </a>
  <a href="https://github.com/unedtamps/orbit/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/unedtamps/orbit" alt="License">
  </a>
</p>

---

## Features

- **TMDB Integration** — Search movies & TV shows with rich metadata (posters, cast, reviews, seasons, episodes)
- **Trending Content** — Browse trending movies and TV shows from TMDB (weekly)
- **Torrent Search** — Find magnet links and torrents via Jackett across multiple trackers
- **Magnet Copy** — One-click copy magnet links to clipboard
- **Download Proxy** — Securely proxy download links through the server (API key hidden)
- **Slug Search** — Smart query generation with Unicode normalization for better torrent matches
- **Pagination** — Server-side and client-side pagination for search and magnet results
- **API Documentation** — Swagger/OpenAPI docs at `/apidocs/`

## Tech Stack

### Backend

| Technology | Purpose |
|------------|---------|
| [Go](https://go.dev/) | Server runtime |
| [Chi](https://github.com/go-chi/chi) | HTTP router |
| [go-jackett](https://github.com/webtor-io/go-jackett) | Jackett API client |
| [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) | Unicode normalization for slug generation |
| [Swaggo](https://github.com/swaggo/swag) | Swagger API documentation |

### Frontend

| Technology | Purpose |
|------------|---------|
| [HTMX](https://htmx.org/) | Partial page updates |
| [Alpine.js](https://alpinejs.dev/) | Interactive UI components |
| [Space Grotesk](https://fonts.google.com/specimen/Space+Grotesk) | Typography |
| [Font Awesome](https://fontawesome.com/) | Icons |

### Infrastructure

| Technology | Purpose |
|------------|---------|
| [Docker](https://www.docker.com/) | Containerized deployment |
| [Jackett](https://github.com/Jackett/Jackett) | Torrent indexer proxy |
| [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr) | Cloudflare bypass |
| [GitHub Actions](https://github.com/features/actions) | CI/CD pipeline |
| [Docker Hub](https://hub.docker.com/) | Image registry |

### APIs

| API | Purpose |
|-----|---------|
| [TMDB API](https://developer.themoviedb.org/) | Movie & TV metadata |
| [Jackett API](https://github.com/Jackett/Jackett) | Torrent search across trackers |

## Quick Start

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/unedtamps/orbit.git
cd orbit

# Configure environment
cp .envrc.example .env
# Edit .env with your API keys

# Start all services
docker compose up -d
```

### Running Locally

```bash
# Prerequisites: Go 1.25+, Jackett running on port 9117

# Install dependencies
go mod download

# Configure environment
export API_URL=http://localhost:9117
export API_KEY=your_jackett_api_key
export TMDB_API_KEY=your_tmdb_api_key

# Run the server
go run .
```

The server starts at `http://localhost:9999`

## Configuration

| Environment Variable | Required | Default | Description |
|---------------------|----------|---------|-------------|
| `API_URL` | Yes | — | Jackett server URL (e.g., `http://localhost:9117`) |
| `API_KEY` | Yes | — | Jackett API key |
| `TMDB_API_KEY` | Yes | — | TMDB API Bearer token |
| `HOST_URL` | No | `http://localhost:9999` | Public host URL (used for Swagger docs) |
| `PORT` | No | `9999` | Server port |
| `CORS_MAX_AGE` | No | `300` | CORS max age in seconds |
| `PROXY_TIMEOUT` | No | `30s` | Download proxy timeout |
| `STATIC_DIR` | No | `./static` | Static files directory |
| `TEMPLATE_GLOB` | No | `templates/*.html` | Template file glob pattern |

## API Endpoints

### Pages

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Home page with search and trending |
| `GET` | `/search?q=query` | Search results page |
| `GET` | `/movie/{id}` | Movie detail page |
| `GET` | `/tv/{id}` | TV show detail page |
| `GET` | `/tv/{id}/season/{season}` | Season detail page |

### API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=query` | TMDB multi-search (movies + TV) |
| `GET` | `/api/movie/{id}` | Movie details with credits & reviews |
| `GET` | `/api/tv/{id}` | TV show details with credits |
| `GET` | `/api/tv/{id}/season/{season}` | Season details with episodes |
| `GET` | `/api/movie/{id}/reviews` | Movie reviews (paginated) |
| `GET` | `/api/tv/{id}/reviews` | TV show reviews (paginated) |
| `GET` | `/api/trending/movies` | Trending movies |
| `GET` | `/api/trending/tv` | Trending TV shows |
| `GET` | `/api/resolve-link?url=...` | Resolve proxy download URL to magnet link |

### Torrent Search

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/magnet/movie/{id}?title=...&year=...` | Find magnets for a movie |
| `GET` | `/magnet/episode/{id}/s{season}/e{episode}` | Find magnets for an episode |
| `GET` | `/dl/{tracker}` | Download proxy (hides Jackett API key) |

## License

Apache 2.0
