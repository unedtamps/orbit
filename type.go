package main

import (
	"html/template"
	"net/http"

	"github.com/webtor-io/go-jackett"
)

type Handler struct {
	*jackett.Client
	templates *template.Template
}

type HanderI interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetTV(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	SearchWeb(w http.ResponseWriter, r *http.Request)
}

// SearchResult represents the data passed to the search template
type SearchResult struct {
	Query         string
	Category      string
	CategoryTitle string
	Results       []jackett.Result
	Count         int
	Error         string
}
