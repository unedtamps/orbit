package main

import (
	"fmt"
	"net/http"
	"os"

	_ "jackettest/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			OrbitSearch API
//	@version		1.0
//	@description	Content discovery API - Search for movies, TV series, and books
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/
func main() {
	handler, err := NewHandler()
	hostURL := os.Getenv("HOST_URL")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Static files (CSS)
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Web Interface Routes (HTMX)
	r.Get("/", handler.Index)
	r.Get("/search", handler.SearchWeb) // HTMX endpoint returns HTML fragment

	// API Routes (JSON)
	r.Get("/movies/{query}", handler.GetMovies)
	r.Get("/books/{query}", handler.GetBooks)
	r.Get("/tv/{query}", handler.GetTV)

	// Swagger Documentation
	r.Get("/apidocs/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/apidocs/doc.json", hostURL)),
	))

	port := ":9999"
	fmt.Printf("OrbitSearch starting on http://localhost%s\n", port)
	fmt.Printf("Web Interface: http://localhost%s/\n", port)
	fmt.Printf("API Docs: http://localhost%s/apidocs/\n", port)

	http.ListenAndServe(port, r)
}
