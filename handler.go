package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	jackett "github.com/webtor-io/go-jackett"
)

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
				return "Unknown"
			}
			return t.Format("Jan 2, 2006")
		},
		"categoryName": func(categories []uint) string {
			if len(categories) == 0 {
				return "Unknown"
			}
			catMap := map[uint]string{
				2000: "Movies",
				2010: "Movies/Foreign",
				2020: "Movies/Other",
				2030: "Movies/SD",
				2040: "Movies/HD",
				2045: "Movies/UHD",
				2050: "Movies/BluRay",
				2060: "Movies/3D",
				2070: "Movies/DVD",
				2080: "Movies/WEB-DL",
				3000: "Books",
				3030: "Books/Technical",
				5000: "TV",
				5020: "TV/SD",
				5030: "TV/HD",
				5040: "TV/UHD",
				5050: "TV/Other",
				5060: "TV/Sport",
				5070: "TV/Anime",
				5080: "TV/Documentary",
			}
			if name, ok := catMap[categories[0]]; ok {
				return name
			}
			return "Unknown"
		},
		"categoryClass": func(categories []uint) string {
			if len(categories) == 0 {
				return "unknown"
			}
			switch categories[0] {
			case 2000, 2010, 2020, 2030, 2040, 2045, 2050, 2060, 2070, 2080:
				return "movies"
			case 5000, 5020, 5030, 5040, 5050, 5060, 5070, 5080:
				return "tv"
			case 3000, 3030:
				return "books"
			default:
				return "unknown"
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
				return "Unknown"
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
				return "unknown"
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
		Client:    j,
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

	categoryTitles := map[string]string{
		"movies": "Movies",
		"tv":     "TV Series",
		"books":  "Books",
	}

	data := SearchResult{
		Query:         query,
		Category:      category,
		CategoryTitle: categoryTitles[category],
	}

	var results []jackett.Result
	var err error

	switch category {
	case "movies":
		results, err = h.Fetch(
			r.Context(),
			jackett.NewMovieSearch().WithCategories(2000).WithQuery(query).Build(),
		)
	case "tv":
		results, err = h.Fetch(
			r.Context(),
			jackett.NewTVSearch().WithCategories(5000, 5050, 5070, 5080).WithQuery(query).Build(),
		)
	case "books":
		results, err = h.Fetch(
			r.Context(),
			jackett.NewBookSearch().WithCategories(3030, 3000).WithQuery(query).Build(),
		)
	default:
		data.Error = "Invalid category"
	}

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
	results, err := h.Fetch(
		r.Context(),
		jackett.NewMovieSearch().WithCategories(2000).WithQuery(query).Build(),
	)
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
	results, err := h.Fetch(
		r.Context(),
		jackett.NewBookSearch().WithCategories(3030, 3000).WithQuery(query).Build(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

// GetTV godoc
//
//	  @Summary		Get TV Series
//		  @Description	Get TV Series (JSON API)
//		  @Tags			tv
//		  @Produce		json
//		  @Param			query	path	string	true	"Search query"
//		  @Success		200	{array}	jackett.Result
//		  @Router			/tv/{query} [get]
func (h *Handler) GetTV(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	results, err := h.Fetch(
		r.Context(),
		jackett.NewTVSearch().WithCategories(5000, 5050, 5070, 5080).WithQuery(query).Build(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, results)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
