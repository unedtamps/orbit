package tmdb

import "fmt"

func (c *Client) GetMovieDetails(id int) (*MovieDetails, error) {
	url := fmt.Sprintf("%s/movie/%d?append_to_response=credits,reviews&language=en-US", baseURL, id)

	var result MovieDetails
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetMovieReviews(id int, page int) (*ReviewResponse, error) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf("%s/movie/%d/reviews?page=%d&language=en-US", baseURL, id, page)

	var result ReviewResponse
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
