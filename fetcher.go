package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	jackett "github.com/webtor-io/go-jackett"
)

type Fetcher struct {
	client *jackett.Client
	apiURL string
	apiKey string
}

func NewFetcher(client *jackett.Client, apiURL, apiKey string) *Fetcher {
	return &Fetcher{client: client, apiURL: apiURL, apiKey: apiKey}
}

func isMagnetLink(link string) bool {
	return strings.HasPrefix(strings.ToLower(link), "magnet:?")
}

func (f *Fetcher) isJackettLink(link string) bool {
	if f.apiURL == "" {
		return false
	}
	return strings.HasPrefix(link, f.apiURL)
}

func (f *Fetcher) toProxyURL(jackettURL string) string {
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
}

// processResults converts Jackett API links to proxy URLs (hides API key)
func (f *Fetcher) processResults(
	ctx context.Context,
	results []jackett.Result,
) ([]jackett.Result, error) {
	for i := range results {
		if results[i].Link == "" {
			continue
		}
		// Skip if already magnet link
		if isMagnetLink(results[i].Link) {
			continue
		}
		// Only process if it's a Jackett API link
		if !f.isJackettLink(results[i].Link) {
			continue
		}
		// Convert to proxy URL instead of fetching magnet (faster)
		results[i].Link = f.toProxyURL(results[i].Link)
	}
	return results, nil
}

func (f *Fetcher) FetchMovies(ctx context.Context, query string) ([]jackett.Result, error) {
	results, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(2000, 100467, 2020, 2030, 2040, 2045, 2060, 100211, 100507, 100001).
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	altResults, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(8000).
			WithTrackers("rutor").
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return f.processResults(ctx, results)
	}

	return f.processResults(ctx, append(results, altResults...))
}

func (f *Fetcher) FetchTV(ctx context.Context, query string) ([]jackett.Result, error) {
	results, err := f.client.Fetch(
		ctx,
		jackett.NewTVSearch().
			WithCategories(5000, 5050, 5070, 5080, 143862, 105852, 112972, 100205, 100212, 100002).
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	altResults, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(8000).
			WithTrackers("rutor").
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return f.processResults(ctx, results)
	}

	return f.processResults(ctx, append(results, altResults...))
}

func (f *Fetcher) FetchBooks(ctx context.Context, query string) ([]jackett.Result, error) {
	results, err := f.client.Fetch(
		ctx,
		jackett.NewBookSearch().
			WithCategories(3030, 7000, 7010, 7020, 7030, 158575, 115634, 160428, 145689, 112812,
				122878, 147656, 124764, 126347, 165359, 165359, 148364, 132586, 121527, 100102, 100601).
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}
	return f.processResults(ctx, results)
}

type ContentType string

const (
	ContentTypeMovies ContentType = "movies"
	ContentTypeTV     ContentType = "tv"
	ContentTypeBooks  ContentType = "books"
)

func (f *Fetcher) FetchByType(
	ctx context.Context,
	contentType ContentType,
	query string,
) ([]jackett.Result, error) {
	switch contentType {
	case ContentTypeMovies:
		return f.FetchMovies(ctx, query)
	case ContentTypeTV:
		return f.FetchTV(ctx, query)
	case ContentTypeBooks:
		return f.FetchBooks(ctx, query)
	default:
		return nil, nil
	}
}

func ContentTypeTitle(ct ContentType) string {
	titles := map[ContentType]string{
		ContentTypeMovies: "Movies",
		ContentTypeTV:     "TV Series",
		ContentTypeBooks:  "Books",
	}
	if title, ok := titles[ct]; ok {
		return title
	}
	return ""
}

func AllContentTypes() []ContentType {
	return []ContentType{ContentTypeMovies, ContentTypeTV, ContentTypeBooks}
}

func ParseContentType(s string) (ContentType, bool) {
	switch s {
	case "movies":
		return ContentTypeMovies, true
	case "tv":
		return ContentTypeTV, true
	case "books":
		return ContentTypeBooks, true
	default:
		return "", false
	}
}
