package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/unedtamps/orbit/internal/fetcher"
)

type MagnetHandler struct {
	fetcher  *fetcher.Fetcher
	template *template.Template
}

func NewMagnetHandler(f *fetcher.Fetcher, tmpl *template.Template) *MagnetHandler {
	return &MagnetHandler{fetcher: f, template: tmpl}
}

func (h *MagnetHandler) GetMovieMagnets(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	year := r.URL.Query().Get("year")
	if title == "" {
		http.Error(w, "title parameter is required", http.StatusBadRequest)
		return
	}

	query := title
	if year != "" {
		// Extract just the year (first 4 chars) from full date like "2022-03-01"
		if len(year) >= 4 {
			query = fmt.Sprintf("%s %s", title, year[:4])
		}
	}

	log.Printf("Magnet search: %q", query)
	results, err := h.fetcher.Search(r.Context(), query)
	if err != nil {
		log.Printf("Magnet search error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "magnet_results.html", results); err != nil {
		log.Printf("Magnet template error: %v", err)
		http.Error(w, "Failed to render results", http.StatusInternalServerError)
	}
}

func (h *MagnetHandler) GetEpisodeMagnets(w http.ResponseWriter, r *http.Request) {
	showName := r.URL.Query().Get("name")
	season := r.URL.Query().Get("season")
	episode := r.URL.Query().Get("episode")
	if showName == "" || season == "" || episode == "" {
		http.Error(w, "name, season, and episode parameters are required", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("%s S%sE%s", showName, season, episode)
	log.Printf("Magnet search: %q", query)
	results, err := h.fetcher.Search(r.Context(), query)
	if err != nil {
		log.Printf("Magnet search error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "magnet_results.html", results); err != nil {
		log.Printf("Magnet template error: %v", err)
		http.Error(w, "Failed to render results", http.StatusInternalServerError)
	}
}
