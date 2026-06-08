package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/unedtamps/orbit/docs"
	"github.com/unedtamps/orbit/internal/config"
	"github.com/unedtamps/orbit/internal/fetcher"
	"github.com/unedtamps/orbit/internal/handler"
	"github.com/unedtamps/orbit/internal/tmdb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/webtor-io/go-jackett"
)

//	@title			OrbitSearch API
//	@version		1.0
//	@description	Content discovery API - Search for movies and TV series
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	j, err := jackett.New(jackett.Settings{
		ApiURL: cfg.APIURL,
		ApiKey: cfg.APIKey,
	})
	if err != nil {
		log.Fatalf("Failed to create Jackett client: %v", err)
	}

	f := fetcher.New(j, cfg.APIURL, cfg.APIKey)
	tmdbClient := tmdb.NewClient(cfg.TMDBAPIKey)
	tmpl := handler.LoadTemplates(cfg.TemplateGlob)
	h := handler.New(f, tmdbClient, tmpl)
	tmdbH := handler.NewTMDBHandler(tmdbClient)
	magnetH := handler.NewMagnetHandler(f, tmpl)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           cfg.CORSMaxAge,
	}))

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", h.Index)
	r.Get("/search", h.SearchPage)
	r.Get("/movie/{id}", h.MovieDetailPage)
	r.Get("/tv/{id}", h.TVDetailPage)
	r.Get("/tv/{id}/season/{season}", h.SeasonPage)

	r.Get("/api/search", tmdbH.Search)
	r.Get("/api/movie/{id}", tmdbH.GetMovie)
	r.Get("/api/tv/{id}", tmdbH.GetTV)
	r.Get("/api/tv/{id}/season/{season}", tmdbH.GetSeason)
	r.Get("/api/movie/{id}/reviews", tmdbH.GetMovieReviews)
	r.Get("/api/tv/{id}/reviews", tmdbH.GetTVReviews)

	r.Get("/api/trending/movies", tmdbH.GetTrendingMovies)
	r.Get("/api/trending/tv", tmdbH.GetTrendingTV)

	r.Get("/magnet/movie/{id}", magnetH.GetMovieMagnets)
	r.Get("/magnet/episode/{id}/s{season}/e{episode}", magnetH.GetEpisodeMagnets)

	r.Get("/api/movies/search/{query}", h.GetMovies)
	r.Get("/api/tv/search/{query}", h.GetTV)
	r.Get("/dl/{tracker}", h.DownloadProxy)
	r.Get("/api/resolve-link", h.ResolveLink)

	r.Get("/apidocs/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/apidocs/doc.json", cfg.HostURL)),
	))

	addr := ":" + cfg.Port
	fmt.Printf("OrbitSearch starting on http://localhost%s\n", addr)
	fmt.Printf("Web Interface: http://localhost%s/\n", addr)
	fmt.Printf("API Docs: http://localhost%s/apidocs/\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
