package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DownloadProxy(w http.ResponseWriter, r *http.Request) {
	tracker := chi.URLParam(r, "tracker")
	if tracker == "" {
		http.Error(w, "Tracker is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	path := query.Get("path")
	file := query.Get("file")

	if path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	jackettURL := fmt.Sprintf(
		"%s/dl/%s/?jackett_apikey=%s&path=%s",
		h.fetcher.APIURL(),
		tracker,
		h.fetcher.APIKey(),
		path,
	)
	if file != "" {
		jackettURL = fmt.Sprintf("%s&file=%s", jackettURL, url.QueryEscape(file))
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, jackettURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	req.Header.Set("Accept", r.Header.Get("Accept"))
	req.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch from Jackett", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	allowedHeaders := map[string]bool{
		"Content-Type":        true,
		"Content-Disposition": true,
		"Content-Length":      true,
	}
	for key, values := range resp.Header {
		if allowedHeaders[key] {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
