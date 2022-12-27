package models

type Item struct {
	ID       string
	FullURL  string `json:"full_url"`
	ShortURL string `json:"short_url"`
}
