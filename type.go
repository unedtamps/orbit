package main

import (
	"net/http"

	"github.com/webtor-io/go-jackett"
)

type Handler struct {
	*jackett.Client
}

type HanderI interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetTV(w http.ResponseWriter, r *http.Request)
}
