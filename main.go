package main

import (
	"fmt"
	"net/http"
	"os"

	_ "jackettest/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Swagger Example API
//	@version		2.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	0A2j9@example.com

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

	r.Get("/apidocs/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/apidocs/doc.json", hostURL)), //The url pointing to API definition
	))
	r.Get("/movies/{query}", handler.GetMovies)
	r.Get("/books/{query}", handler.GetBooks)
	r.Get("/tv/{query}", handler.GetTV)
	http.ListenAndServe(":9999", r)

}
