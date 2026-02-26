package db

type Config struct {
	Id        string `json:"id"`
	Path      string `json:"path"`
	Alias     string `json:"alias"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
