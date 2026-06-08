package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/unedtamps/orbit/internal/tmdb"
)

type TMDBHandler struct {
	client *tmdb.Client
}

func NewTMDBHandler(client *tmdb.Client) *TMDBHandler {
	return &TMDBHandler{client: client}
}

func writeJSONError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func cleanTMDBError(err error) string {
	msg := err.Error()
	if strings.Contains(msg, "context deadline exceeded") || strings.Contains(msg, "Client.Timeout") {
		return "TMDB API timed out. Please try again."
	}
	if strings.Contains(msg, "connection refused") {
		return "Cannot connect to TMDB API. Please try again later."
	}
	if strings.Contains(msg, "no such host") {
		return "Cannot reach TMDB API. Check your internet connection."
	}
	if strings.Contains(msg, "TMDB API error") {
		return "TMDB API returned an error. Please try again."
	}
	return "Failed to fetch data from TMDB. Please try again."
}

func (h *TMDBHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSONError(w, "query parameter q is required", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	result, err := h.client.MultiSearch(query, page)
	if err != nil {
		log.Printf("TMDB search error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	result, err := h.client.GetMovieDetails(id)
	if err != nil {
		log.Printf("TMDB movie error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetTV(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, "invalid tv id", http.StatusBadRequest)
		return
	}

	result, err := h.client.GetTVDetails(id)
	if err != nil {
		log.Printf("TMDB TV error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetSeason(w http.ResponseWriter, r *http.Request) {
	tvID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, "invalid tv_id", http.StatusBadRequest)
		return
	}
	season, err := strconv.Atoi(chi.URLParam(r, "season"))
	if err != nil {
		writeJSONError(w, "invalid season number", http.StatusBadRequest)
		return
	}

	result, err := h.client.GetSeasonDetails(tvID, season)
	if err != nil {
		log.Printf("TMDB season error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetMovieReviews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	result, err := h.client.GetMovieReviews(id, page)
	if err != nil {
		log.Printf("TMDB movie reviews error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetTVReviews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, "invalid tv id", http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	result, err := h.client.GetTVReviews(id, page)
	if err != nil {
		log.Printf("TMDB TV reviews error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "week"
	}

	result, err := h.client.GetTrendingMovies(window)
	if err != nil {
		log.Printf("TMDB trending movies error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *TMDBHandler) GetTrendingTV(w http.ResponseWriter, r *http.Request) {
	window := r.URL.Query().Get("window")
	if window == "" {
		window = "week"
	}

	result, err := h.client.GetTrendingTV(window)
	if err != nil {
		log.Printf("TMDB trending TV error: %v", err)
		writeJSONError(w, cleanTMDBError(err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
