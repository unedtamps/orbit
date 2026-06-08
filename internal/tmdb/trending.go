package tmdb

import "fmt"

func (c *Client) GetTrendingMovies(timeWindow string) (*MultiSearchResponse, error) {
	if timeWindow == "" {
		timeWindow = "week"
	}
	url := fmt.Sprintf("%s/trending/movie/%s?language=en-US", baseURL, timeWindow)

	var result MultiSearchResponse
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetTrendingTV(timeWindow string) (*MultiSearchResponse, error) {
	if timeWindow == "" {
		timeWindow = "week"
	}
	url := fmt.Sprintf("%s/trending/tv/%s?language=en-US", baseURL, timeWindow)

	var result MultiSearchResponse
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
