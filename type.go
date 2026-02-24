package main

import (
	"net/http"

	"github.com/webtor-io/go-jackett"
)

type HanderI interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetTV(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	SearchWeb(w http.ResponseWriter, r *http.Request)
	DownloadProxy(w http.ResponseWriter, r *http.Request)
}

type SearchResult struct {
	Query         string
	Category      string
	CategoryTitle string
	Results       []jackett.Result
	Count         int
	Error         string
}
