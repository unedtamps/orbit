package tmdb

import "fmt"

func (c *Client) MultiSearch(query string, page int) (*MultiSearchResponse, error) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf("%s/search/multi?query=%s&page=%d&language=en-US", baseURL, query, page)

	var result MultiSearchResponse
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
