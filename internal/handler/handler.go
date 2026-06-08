package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"strings"

	"github.com/unedtamps/orbit/internal/fetcher"
	"github.com/unedtamps/orbit/internal/tmdb"
)

type Handler struct {
	fetcher  *fetcher.Fetcher
	tmdb     *tmdb.Client
	template *template.Template
}

func New(f *fetcher.Fetcher, tm *tmdb.Client, tmpl *template.Template) *Handler {
	return &Handler{fetcher: f, tmdb: tm, template: tmpl}
}

func LoadTemplates(glob string) *template.Template {
	funcMap := template.FuncMap{
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
		"sub": func(a, b uint) uint {
			if a > b {
				return a - b
			}
			return 0
		},
		"first": func(n int, slice interface{}) interface{} {
			v := reflect.ValueOf(slice)
			if v.Kind() != reflect.Slice {
				return slice
			}
			if v.Len() <= n {
				return slice
			}
			return v.Slice(0, n).Interface()
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(glob)
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	return tmpl
}
