package model

import (
	jackett "github.com/webtor-io/go-jackett"
)

type ContentType string

const (
	ContentTypeMovies ContentType = "movies"
	ContentTypeTV     ContentType = "tv"
)

type SearchResult struct {
	Query         string
	Category      string
	CategoryTitle string
	Results       []jackett.Result
	Count         int
	Error         string
}

func ContentTypeTitle(ct ContentType) string {
	titles := map[ContentType]string{
		ContentTypeMovies: "Movies",
		ContentTypeTV:     "TV Series",
	}
	if title, ok := titles[ct]; ok {
		return title
	}
	return ""
}

func ParseContentType(s string) (ContentType, bool) {
	switch s {
	case "movies":
		return ContentTypeMovies, true
	case "tv":
		return ContentTypeTV, true
	default:
		return "", false
	}
}
