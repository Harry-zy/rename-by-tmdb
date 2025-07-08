package config

// Config 配置结构
type Config struct {
	TMDBAPI struct {
		Key      string `json:"key"`
		Language string `json:"language"`
	} `json:"tmdb_api"`
	RemoteServer struct {
		BaseURL    string `json:"base_url"`
		AuthToken  string `json:"auth_token"`
		UploadMode bool   `json:"upload_mode"`
	} `json:"remote_server"`
}
