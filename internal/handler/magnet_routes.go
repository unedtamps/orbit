package handler

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/unedtamps/orbit/internal/fetcher"

	jackett "github.com/webtor-io/go-jackett"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var multiHyphen = regexp.MustCompile(`-{2,}`)

var stripCombining = transform.Chain(
	norm.NFD,
	transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}),
	norm.NFC,
)

func normalizeUnicode(s string) string {
	result, _, _ := transform.String(stripCombining, strings.ToLower(s))
	return result
}

func slugify(s string) string {
	s = normalizeUnicode(s)
	s = strings.ReplaceAll(s, "ß", "ss")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r == ' ' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	s = multiHyphen.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func dedupe(a, b []jackett.Result) []jackett.Result {
	seen := make(map[string]int, len(a))
	merged := make([]jackett.Result, 0, len(a)+len(b))

	for _, r := range a {
		key := dedupeKey(r)
		seen[key] = len(merged)
		merged = append(merged, r)
	}

	for _, r := range b {
		key := dedupeKey(r)
		if idx, ok := seen[key]; ok {
			if r.Seeders > merged[idx].Seeders {
				merged[idx] = r
			}
		} else {
			seen[key] = len(merged)
			merged = append(merged, r)
		}
	}

	return merged
}

func dedupeKey(r jackett.Result) string {
	if r.InfoHash != "" {
		return r.InfoHash
	}
	return strings.ToLower(r.Title)
}

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

	query := slugify(title)
	if year != "" && len(year) >= 4 {
		query = query + "-" + year[:4]
	}

	log.Printf("Magnet search: %q", query)
	results, err := h.fetcher.Search(r.Context(), query)
	if err != nil {
		log.Printf("Magnet search error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeResults(w, results)
}

func (h *MagnetHandler) GetEpisodeMagnets(w http.ResponseWriter, r *http.Request) {
	showName := r.URL.Query().Get("name")
	season := r.URL.Query().Get("season")
	episode := r.URL.Query().Get("episode")
	episodeTitle := r.URL.Query().Get("title")
	if showName == "" || season == "" || episode == "" {
		http.Error(w, "name, season, and episode parameters are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	query1 := slugify(showName) + "-s" + season + "e" + episode
	log.Printf("Magnet search (slug): %q", query1)
	results1, err := h.fetcher.Search(ctx, query1)
	if err != nil {
		log.Printf("Magnet search error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []jackett.Result
	if episodeTitle != "" {
		query2 := slugify(showName) + "-" + slugify(episodeTitle)
		log.Printf("Magnet search (title): %q", query2)
		results2, err := h.fetcher.Search(ctx, query2)
		if err != nil {
			log.Printf("Magnet search error (title query): %v", err)
			results = results1
		} else {
			results = dedupe(results1, results2)
		}
	} else {
		results = results1
	}

	h.writeResults(w, results)
}

func (h *MagnetHandler) writeResults(w http.ResponseWriter, results []jackett.Result) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "magnet_results.html", results); err != nil {
		log.Printf("Magnet template error: %v", err)
		http.Error(w, "Failed to render results", http.StatusInternalServerError)
	}
}
