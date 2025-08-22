package models

// TMDBEpisode 表示TMDB的一集信息
type TMDBEpisode struct {
	EpisodeNumber int    `json:"episode_number"`
	AirDate       string `json:"air_date"`
}

// TMDBSeason 表示TMDB的一季信息
type TMDBSeason struct {
	SeasonNumber int           `json:"season_number"`
	Episodes     []TMDBEpisode `json:"episodes"`
}

// TMDBShow 表示TMDB的剧集信息
type TMDBShow struct {
	Name         string       `json:"name"`
	FirstAirDate string       `json:"first_air_date"`
	Type         string       `json:"type"`
	Seasons      []TMDBSeason `json:"seasons"`
}

// TMDBMovie 表示TMDB的电影信息
type TMDBMovie struct {
	Title         string `json:"title"`
	ReleaseDate   string `json:"release_date"`
	OriginalTitle string `json:"original_title"`
}

// TMDBError TMDB API错误响应
type TMDBError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}
