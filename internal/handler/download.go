package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	log.Printf("Download proxy: tracker=%s file=%s path=%.80s...", tracker, file, path)

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

	log.Printf("Download proxy: requesting %s/dl/%s/?jackett_apikey=***&path=...&file=%s", h.fetcher.APIURL(), tracker, file)

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, jackettURL, nil)
	if err != nil {
		log.Printf("Download proxy: failed to create request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", r.Header.Get("User-Agent"))
	req.Header.Set("Accept", r.Header.Get("Accept"))
	req.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Download proxy: request failed: %v", err)
		http.Error(w, "Failed to fetch from Jackett", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Download proxy: Jackett responded with status %d", resp.StatusCode)

	allowedHeaders := map[string]bool{
		"Content-Type":        true,
		"Content-Disposition": true,
		"Content-Length":      true,
		"Location":            true,
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

func (h *Handler) ResolveLink(w http.ResponseWriter, r *http.Request) {
	proxyURL := r.URL.Query().Get("url")
	if proxyURL == "" {
		writeJSONError(w, "url parameter is required", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(proxyURL, "/dl/") {
		writeJSONError(w, "url must be a /dl/ proxy URL", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(proxyURL)
	if err != nil {
		writeJSONError(w, "invalid url", http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 2 {
		writeJSONError(w, "invalid dl url path", http.StatusBadRequest)
		return
	}
	tracker := pathParts[1]
	query := parsed.Query()
	path := query.Get("path")
	file := query.Get("file")

	if path == "" {
		writeJSONError(w, "path is required", http.StatusBadRequest)
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

	log.Printf("Resolve link: following redirects for %s/dl/%s/", h.fetcher.APIURL(), tracker)

	var finalURL string
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				finalURL = req.URL.String()
			}
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, jackettURL, nil)
	if err != nil {
		writeJSONError(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Resolve link: request failed: %v", err)
		writeJSONError(w, "failed to resolve link", http.StatusBadGateway)
		return
	}
	resp.Body.Close()

	if finalURL != "" {
		log.Printf("Resolve link: resolved to %s", finalURL)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"url": finalURL})
		return
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if loc := resp.Header.Get("Location"); loc != "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"url": loc})
			return
		}
	}

	log.Printf("Resolve link: no redirect found, status %d", resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": jackettURL})
}
