package models

type Item struct {
	ID       string
	FullURL  string `json:"url"`
	ShortURL string `json:"short_url"`
}
