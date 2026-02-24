package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	jackett "github.com/webtor-io/go-jackett"
)

type Handler struct {
	*Fetcher
	templates *template.Template
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"safeURL": func(u string) template.URL {
			return template.URL(u)
		},
		"lower": strings.ToLower,
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
		"proxyURL": func(jackettURL string) string {
			// Convert Jackett download URL to our proxy URL
			// From: http://localhost:9117/dl/bitsearch/?jackett_apikey=...&path=...&file=...
			// To: /dl/bitsearch?path=...&file=...
			if jackettURL == "" {
				return ""
			}
			parsed, err := url.Parse(jackettURL)
			if err != nil {
				return jackettURL
			}
			// Extract tracker from path (/dl/bitsearch/...)
			pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(pathParts) < 2 || pathParts[0] != "dl" {
				return jackettURL
			}
			tracker := pathParts[1]
			query := parsed.Query()
			path := query.Get("path")
			file := query.Get("file")
			if path == "" {
				return jackettURL
			}
			newQuery := url.Values{}
			newQuery.Set("path", path)
			if file != "" {
				newQuery.Set("file", file)
			}
			return fmt.Sprintf("/dl/%s?%s", tracker, newQuery.Encode())
		},
		"formatSize": func(bytes uint64) string {
			if bytes == 0 {
				return "0 B"
			}
			sizes := []string{"B", "KB", "MB", "GB", "TB"}
			fbytes := float64(bytes)
			var i int
			for i = 0; i < len(sizes)-1 && fbytes >= 1024; i++ {
				fbytes /= 1024
			}
			return fmt.Sprintf("%.2f %s", fbytes, sizes[i])
		},
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return "Other"
			}
			return t.Format("Jan 2, 2006")
		},
		"categoryName": func(categories []uint) string {
			if len(categories) == 0 {
				return "Other"
			}
			catMap := map[uint]string{
				// Movies
				2000:   "Movies",
				2010:   "Movies/Foreign",
				2020:   "Movies/Other",
				2030:   "Movies/SD",
				2040:   "Movies/HD",
				2045:   "Movies/UHD",
				2050:   "Movies/BluRay",
				2060:   "Movies/3D",
				2070:   "Movies/DVD",
				2080:   "Movies/WEB-DL",
				8000:   "Movies/Other",
				100001: "Movies/Other",
				100211: "Movies/HD",
				100467: "Movies/UHD",
				100507: "Movies/HD",
				// TV
				5000:   "TV",
				5020:   "TV/SD",
				5030:   "TV/HD",
				5040:   "TV/UHD",
				5050:   "TV/Other",
				5060:   "TV/Sport",
				5070:   "TV/Anime",
				5080:   "TV/Documentary",
				100002: "TV/Other",
				100205: "TV/HD",
				100212: "TV/HD",
				105852: "TV/Other",
				112972: "TV/Anime",
				143862: "TV/Other",
				// Books
				3000:   "Books",
				3030:   "Books/Technical",
				7000:   "Books/Other",
				7010:   "Books/Mags",
				7020:   "Books/EBook",
				7030:   "Books/Comics",
				100102: "Books/Other",
				100601: "Books/Other",
				112812: "Books/EBook",
				115634: "Books/EBook",
				121527: "Books/EBook",
				122878: "Books/EBook",
				124764: "Books/EBook",
				126347: "Books/EBook",
				132586: "Books/EBook",
				145689: "Books/EBook",
				147656: "Books/EBook",
				148364: "Books/EBook",
				158575: "Books/EBook",
				160428: "Books/EBook",
				165359: "Books/EBook",
			}
			if name, ok := catMap[categories[0]]; ok {
				return name
			}
			return "Other"
		},
		"categoryClass": func(categories []uint) string {
			if len(categories) == 0 {
				return "other"
			}
			switch categories[0] {
			// Movies
			case 2000, 2010, 2020, 2030, 2040, 2045, 2050, 2060, 2070, 2080, 8000,
				100001, 100211, 100467, 100507:
				return "movies"
			// TV
			case 5000, 5020, 5030, 5040, 5050, 5060, 5070, 5080,
				100002, 100205, 100212, 105852, 112972, 143862:
				return "tv"
			// Books
			case 3000, 3030, 7000, 7010, 7020, 7030,
				100102, 100601, 112812, 115634, 121527, 122878, 124764, 126347,
				132586, 145689, 147656, 148364, 158575, 160428, 165359:
				return "books"
			default:
				return "other"
			}
		},
		"detectQuality": func(title string) string {
			lower := strings.ToLower(title)
			switch {
			case strings.Contains(lower, "2160p") || strings.Contains(lower, "4k") || strings.Contains(lower, "uhd"):
				return "UHD"
			case strings.Contains(lower, "1080p") || strings.Contains(lower, "bluray") || strings.Contains(lower, "720p"):
				return "HD"
			case strings.Contains(lower, "web-dl") || strings.Contains(lower, "webdl") || strings.Contains(lower, "webrip"):
				return "HD"
			case strings.Contains(lower, "hdtv"):
				return "HD"
			default:
				return "Other"
			}
		},
		"qualityClass": func(title string) string {
			lower := strings.ToLower(title)
			switch {
			case strings.Contains(lower, "2160p") || strings.Contains(lower, "4k") || strings.Contains(lower, "uhd"):
				return "uhd"
			case strings.Contains(lower, "1080p") || strings.Contains(lower, "bluray") || strings.Contains(lower, "720p"):
				return "hd"
			case strings.Contains(lower, "web-dl") || strings.Contains(lower, "webdl") || strings.Contains(lower, "webrip"):
				return "hd"
			case strings.Contains(lower, "hdtv"):
				return "hd"
			default:
				return "other"
			}
		},
		"leechers": func(peers, seeders uint) uint {
			if peers > seeders {
				return peers - seeders
			}
			return 0
		},
	}
}

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

	tmpl, err := template.New("").Funcs(templateFuncs()).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		Fetcher:   NewFetcher(j, apiUrl, apiKey),
		templates: tmpl,
	}, nil
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := h.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SearchWeb(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	contentType, valid := ParseContentType(category)
	if !valid {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	data := SearchResult{
		Query:         query,
		Category:      category,
		CategoryTitle: ContentTypeTitle(contentType),
	}

	results, err := h.Fetcher.FetchByType(r.Context(), contentType, query)
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Results = results
		data.Count = len(results)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Check if this is an HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		// Return just the results partial for HTMX
		err = h.templates.ExecuteTemplate(w, "search_results.html", data)
	} else {
		// Return full page with results embedded for direct visits
		err = h.templates.ExecuteTemplate(w, "search_page.html", data)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetMovies godoc
//
//	@Summary		Get Movies
//	@Description	Get Movies (JSON API)
//	@Tags			movies
//	@Produce		json
//	@Param			query	path	string	true	"Search query"
//	@Success		200	{array}	jackett.Result
//	@Router			/movies/{query} [get]
func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.Fetcher.FetchMovies(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

// GetBooks godoc
//
//	@Summary		Get Books
//	@Description	Get Books (JSON API)
//	@Tags			books
//	@Produce		json
//	@Param			query	path	string	true	"Search query"
//	@Success		200	{array}	jackett.Result
//	@Router			/books/{query} [get]
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.Fetcher.FetchBooks(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

// GetTV godoc
//
//	  @Summary			Get TV Series
//		  @Description	Get TV Series (JSON API)
//		  @Tags			tv
//		  @Produce		json
//		  @Param			query	path	string	true	"Search query"
//		  @Success		200	{array}	jackett.Result
//		  @Router			/tv/{query} [get]
func (h *Handler) GetTV(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.Fetcher.FetchTV(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

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
		h.apiURL,
		tracker,
		h.apiKey,
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
			// Don't follow redirects - we'll pass them through to the client
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch from Jackett", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
