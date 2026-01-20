package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	jackett "github.com/webtor-io/go-jackett"
)

func NewHandler() (HanderI, error) {
	apiUrl := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	j, err := jackett.New(jackett.Settings{
		ApiURL: apiUrl,
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, err
	}
	return &Handler{Client: j}, nil
}

// GetMovies godoc
//
//		@Summary		Get Movies
//		@Description	Get Movies
//		@Tags			movies
//		@Produce		json
//		@Param			query	path	string	true	"query"
//	 @Success		200	{object} jackett.Result
//		@Router			/movies/{query} [get]
func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	result, err := h.Fetch(r.Context(), jackett.NewMovieSearch().WithCategories(2000).WithQuery(query).Build())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, result)
}

// GetBooks godoc
//
//		@Summary		Get Books
//		@Description	Get Books
//		@Tags			books
//		@Produce		json
//		@Param			query	path	string	true	"query"
//	 @Success		200	{object} jackett.Result
//		@Router			/books/{query} [get]
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	result, err := h.Fetch(r.Context(), jackett.NewBookSearch().WithCategories(3030, 3000).WithQuery(query).Build())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, result)
}

// GetTV godoc
//
//		   @Summary		Get TV
//		   @Description	Get TV by query
//		   @Tags			tv
//		   @Produce		json
//		   @Param			query	path	string	true	"query"
//	  @Success		200	{object} jackett.Result
//		   @Router			/tv/{query} [get]
func (h *Handler) GetTV(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	result, err := h.Fetch(r.Context(), jackett.NewTVSearch().WithCategories(5000, 5050, 5070, 5080).WithQuery(query).Build())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
