package fetcher

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/unedtamps/orbit/internal/model"

	jackett "github.com/webtor-io/go-jackett"
)

type Fetcher struct {
	client *jackett.Client
	apiURL string
	apiKey string
}

func New(client *jackett.Client, apiURL, apiKey string) *Fetcher {
	return &Fetcher{client: client, apiURL: apiURL, apiKey: apiKey}
}

func (f *Fetcher) APIURL() string { return f.apiURL }
func (f *Fetcher) APIKey() string { return f.apiKey }

func isMagnetLink(link string) bool {
	return strings.HasPrefix(strings.ToLower(link), "magnet:?")
}

func (f *Fetcher) isJackettLink(link string) bool {
	if f.apiURL == "" {
		return false
	}
	return strings.HasPrefix(link, f.apiURL)
}

// ToProxyURL converts a Jackett download URL to a proxy URL, hiding the API key.
//
// From: http://localhost:9117/dl/bitsearch/?jackett_apikey=...&path=...&file=...
// To:   /dl/bitsearch?path=...&file=...
func ToProxyURL(jackettURL string) string {
	if jackettURL == "" {
		return ""
	}
	parsed, err := url.Parse(jackettURL)
	if err != nil {
		return jackettURL
	}
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

// processResults converts Jackett API links to proxy URLs (hides API key).
func (f *Fetcher) processResults(
	ctx context.Context,
	results []jackett.Result,
) ([]jackett.Result, error) {
	for i := range results {
		if results[i].Link == "" {
			continue
		}
		if isMagnetLink(results[i].Link) {
			continue
		}
		if !f.isJackettLink(results[i].Link) {
			continue
		}
		results[i].Link = ToProxyURL(results[i].Link)
	}
	return results, nil
}

func (f *Fetcher) FetchMovies(ctx context.Context, query string) ([]jackett.Result, error) {
	results, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(model.MovieCategories...).
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	altResults, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(model.AltSearchCategories...).
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
			WithCategories(model.TVCategories...).
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	altResults, err := f.client.Fetch(
		ctx,
		jackett.NewMovieSearch().
			WithCategories(model.AltSearchCategories...).
			WithTrackers("rutor").
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return f.processResults(ctx, results)
	}

	return f.processResults(ctx, append(results, altResults...))
}

// Search does a generic Jackett search without category filters.
// Used for magnet link lookups (movies, episodes, etc).
func (f *Fetcher) Search(ctx context.Context, query string) ([]jackett.Result, error) {
	results, err := f.client.Fetch(
		ctx,
		jackett.NewRawSearch().
			WithQuery(query).
			Build(),
	)
	if err != nil {
		return nil, err
	}
	return f.processResults(ctx, results)
}

func (f *Fetcher) FetchByType(
	ctx context.Context,
	contentType model.ContentType,
	query string,
) ([]jackett.Result, error) {
	switch contentType {
	case model.ContentTypeMovies:
		return f.FetchMovies(ctx, query)
	case model.ContentTypeTV:
		return f.FetchTV(ctx, query)
	default:
		return nil, fmt.Errorf("unknown content type: %s", contentType)
	}
}
