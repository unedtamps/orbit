package tmdb

import "fmt"

func (c *Client) GetTVDetails(id int) (*TVDetails, error) {
	url := fmt.Sprintf("%s/tv/%d?append_to_response=credits&language=en-US", baseURL, id)

	var result TVDetails
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetSeasonDetails(tvID int, seasonNumber int) (*SeasonDetails, error) {
	url := fmt.Sprintf("%s/tv/%d/season/%d?language=en-US", baseURL, tvID, seasonNumber)

	var result SeasonDetails
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetTVReviews(id int, page int) (*ReviewResponse, error) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf("%s/tv/%d/reviews?page=%d&language=en-US", baseURL, id, page)

	var result ReviewResponse
	if err := c.doRequest(url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
