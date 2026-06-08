package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SearchPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Query":  query,
		"IsHome": false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// If this is an HTMX request, return just the results fragment
	if r.Header.Get("HX-Request") == "true" {
		if err := h.template.ExecuteTemplate(w, "search_results_partial.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, return the full page
	if err := h.template.ExecuteTemplate(w, "search.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) MovieDetailPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := h.tmdb.GetMovieDetails(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch movie: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Movie":  movie,
		"IsHome": false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "movie_detail.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) TVDetailPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid TV ID", http.StatusBadRequest)
		return
	}

	tv, err := h.tmdb.GetTVDetails(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch TV show: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"TV":     tv,
		"TVID":   id,
		"IsHome": false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "tv_detail.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SeasonPage(w http.ResponseWriter, r *http.Request) {
	tvID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid TV ID", http.StatusBadRequest)
		return
	}
	seasonNum, err := strconv.Atoi(chi.URLParam(r, "season"))
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	season, err := h.tmdb.GetSeasonDetails(tvID, seasonNum)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch season: %v", err), http.StatusInternalServerError)
		return
	}

	tv, err := h.tmdb.GetTVDetails(tvID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch TV show: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"TVName":  tv.Name,
		"TVID":    tvID,
		"Season":  season,
		"IsHome":  false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "season_detail.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
