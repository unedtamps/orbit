package tmdb

import "fmt"

const (
	ImageBaseURL = "https://image.tmdb.org/t/p"
)

type MultiSearchResponse struct {
	Page            int                `json:"page"`
	Results         []SearchResult     `json:"results"`
	TotalPages     int                `json:"total_pages"`
	TotalResults   int                `json:"total_results"`
}

type SearchResult struct {
	ID            int      `json:"id"`
	MediaType     string   `json:"media_type"`
	Title         string   `json:"title"`
	Name          string   `json:"name"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	ReleaseDate   string   `json:"release_date"`
	FirstAirDate  string   `json:"first_air_date"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int      `json:"vote_count"`
	GenreIDs      []int    `json:"genre_ids"`
	Popularity    float64  `json:"popularity"`
	OriginalLanguage string `json:"original_language"`
	Adult         bool     `json:"adult"`
}

func (s SearchResult) DisplayTitle() string {
	if s.Title != "" {
		return s.Title
	}
	return s.Name
}

func (s SearchResult) DisplayDate() string {
	if s.ReleaseDate != "" {
		return s.ReleaseDate
	}
	return s.FirstAirDate
}

func (s SearchResult) PosterURL(size string) string {
	if s.PosterPath == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	return ImageBaseURL + "/" + size + s.PosterPath
}

func (s SearchResult) BackdropURL(size string) string {
	if s.BackdropPath == "" {
		return ""
	}
	if size == "" {
		size = "w1280"
	}
	return ImageBaseURL + "/" + size + s.BackdropPath
}

type ReviewResponse struct {
	ID            int      `json:"id"`
	Page          int      `json:"page"`
	Results       []Review `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

type Review struct {
	ID            string        `json:"id"`
	Author        string        `json:"author"`
	AuthorDetails AuthorDetails `json:"author_details"`
	Content       string        `json:"content"`
	CreatedAt     string        `json:"created_at"`
	URL           string        `json:"url"`
}

type AuthorDetails struct {
	Name        string   `json:"name"`
	Username    string   `json:"username"`
	AvatarPath  string   `json:"avatar_path"`
	Rating      *float64 `json:"rating"`
}

func (a AuthorDetails) AvatarURL() string {
	if a.AvatarPath == "" {
		return ""
	}
	if len(a.AvatarPath) > 1 && a.AvatarPath[0] == '/' && a.AvatarPath[1] == 'h' {
		return a.AvatarPath[1:]
	}
	return ImageBaseURL + "/w45" + a.AvatarPath
}

func (r Review) RatingDisplay() string {
	if r.AuthorDetails.Rating == nil {
		return ""
	}
	rating := *r.AuthorDetails.Rating
	if rating == 0 {
		return ""
	}
	return fmt.Sprintf("%.1f", rating)
}

type MovieDetails struct {
	Adult               bool            `json:"adult"`
	BackdropPath        string          `json:"backdrop_path"`
	Budget              int64           `json:"budget"`
	Genres              []Genre         `json:"genres"`
	Homepage            string          `json:"homepage"`
	ID                  int             `json:"id"`
	IMDbID              string          `json:"imdb_id"`
	OriginalLanguage    string          `json:"original_language"`
	OriginalTitle       string          `json:"original_title"`
	Overview            string          `json:"overview"`
	Popularity          float64         `json:"popularity"`
	PosterPath          string          `json:"poster_path"`
	ProductionCompanies []Company       `json:"production_companies"`
	ReleaseDate         string          `json:"release_date"`
	Revenue             int64           `json:"revenue"`
	Runtime             int             `json:"runtime"`
	Status              string          `json:"status"`
	Tagline             string          `json:"tagline"`
	Title               string          `json:"title"`
	VoteAverage         float64         `json:"vote_average"`
	VoteCount           int             `json:"vote_count"`
	Credits             *Credits        `json:"credits,omitempty"`
	Reviews             *ReviewResponse `json:"reviews,omitempty"`
}

func (m MovieDetails) PosterURL(size string) string {
	if m.PosterPath == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	return ImageBaseURL + "/" + size + m.PosterPath
}

func (m MovieDetails) BackdropURL(size string) string {
	if m.BackdropPath == "" {
		return ""
	}
	if size == "" {
		size = "w1280"
	}
	return ImageBaseURL + "/" + size + m.BackdropPath
}

type TVDetails struct {
	BackdropPath      string          `json:"backdrop_path"`
	EpisodeRunTime    []int           `json:"episode_run_time"`
	FirstAirDate      string          `json:"first_air_date"`
	Genres            []Genre         `json:"genres"`
	Homepage          string          `json:"homepage"`
	ID                int             `json:"id"`
	InProduction      bool            `json:"in_production"`
	Languages         []string        `json:"languages"`
	LastAirDate       string          `json:"last_air_date"`
	Name              string          `json:"name"`
	Networks          []Network       `json:"networks"`
	NumberOfEpisodes  int             `json:"number_of_episodes"`
	NumberOfSeasons   int             `json:"number_of_seasons"`
	OriginCountry     []string        `json:"origin_country"`
	OriginalLanguage  string          `json:"original_language"`
	OriginalName      string          `json:"original_name"`
	Overview          string          `json:"overview"`
	Popularity        float64         `json:"popularity"`
	PosterPath        string          `json:"poster_path"`
	Seasons           []Season        `json:"seasons"`
	Status            string          `json:"status"`
	Tagline           string          `json:"tagline"`
	Type              string          `json:"type"`
	VoteAverage       float64         `json:"vote_average"`
	VoteCount         int             `json:"vote_count"`
	Credits           *Credits        `json:"credits,omitempty"`
}

func (t TVDetails) PosterURL(size string) string {
	if t.PosterPath == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	return ImageBaseURL + "/" + size + t.PosterPath
}

func (t TVDetails) BackdropURL(size string) string {
	if t.BackdropPath == "" {
		return ""
	}
	if size == "" {
		size = "w1280"
	}
	return ImageBaseURL + "/" + size + t.BackdropPath
}

type Season struct {
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

func (s Season) PosterURL(size string) string {
	if s.PosterPath == "" {
		return ""
	}
	if size == "" {
		size = "w200"
	}
	return ImageBaseURL + "/" + size + s.PosterPath
}

type SeasonDetails struct {
	ID           int       `json:"id"`
	AirDate      string    `json:"air_date"`
	Episodes     []Episode `json:"episodes"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	PosterPath   string    `json:"poster_path"`
	SeasonNumber int       `json:"season_number"`
}

func (s SeasonDetails) PosterURL(size string) string {
	if s.PosterPath == "" {
		return ""
	}
	if size == "" {
		size = "w200"
	}
	return ImageBaseURL + "/" + size + s.PosterPath
}

type Episode struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Overview         string  `json:"overview"`
	AirDate          string  `json:"air_date"`
	EpisodeNumber    int     `json:"episode_number"`
	SeasonNumber     int     `json:"season_number"`
	StillPath        string  `json:"still_path"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Runtime          int     `json:"runtime"`
}

func (e Episode) StillURL(size string) string {
	if e.StillPath == "" {
		return ""
	}
	if size == "" {
		size = "w300"
	}
	return ImageBaseURL + "/" + size + e.StillPath
}

func (e Episode) SeasonEpisodeCode() string {
	return fmt.Sprintf("S%02dE%02d", e.SeasonNumber, e.EpisodeNumber)
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type Network struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
	Order       int    `json:"order"`
}

func (c CastMember) ProfileURL(size string) string {
	if c.ProfilePath == "" {
		return ""
	}
	if size == "" {
		size = "w185"
	}
	return ImageBaseURL + "/" + size + c.ProfilePath
}

type CrewMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	Job         string `json:"job"`
	ProfilePath string `json:"profile_path"`
}
