package main

import (
	"context"

	jackett "github.com/webtor-io/go-jackett"
)

type Fetcher struct {
	client *jackett.Client
}

func NewFetcher(client *jackett.Client) *Fetcher {
	return &Fetcher{client: client}
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
		return results, nil
	}

	return append(results, altResults...), nil
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

	return append(results, altResults...), nil
}

func (f *Fetcher) FetchBooks(ctx context.Context, query string) ([]jackett.Result, error) {
	return f.client.Fetch(
		ctx,
		jackett.NewBookSearch().
			WithCategories(3030, 7000, 7010, 7020, 7030, 158575, 115634, 160428, 145689, 112812,
				122878, 147656, 124764, 126347, 165359, 165359, 148364, 132586, 121527, 100102, 100601).
			WithQuery(query).
			Build(),
	)
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
