package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetMovies godoc
//
//	@Summary		Get Movies
//	@Description	Get Movies (JSON API)
//	@Tags			movies
//	@Produce		json
//	@Param			query	path	string	true	"Search query"
//	@Success		200	{array}	jackett.Result
//	@Router			/api/movies/search/{query} [get]
func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.fetcher.FetchMovies(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

// GetTV godoc
//
//	@Summary		Get TV Series
//	@Description	Get TV Series (JSON API)
//	@Tags			tv
//	@Produce		json
//	@Param			query	path	string	true	"Search query"
//	@Success		200	{array}	jackett.Result
//	@Router			/api/tv/search/{query} [get]
func (h *Handler) GetTV(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.fetcher.FetchTV(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
