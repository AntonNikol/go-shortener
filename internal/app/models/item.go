package models

type Item struct {
	ID       string
	FullURL  string `json:"url"`
	ShortURL string `json:"short_url"`
	UserID   string
}

type ItemResponse struct {
	ID       string `json:"-"`
	FullURL  string `json:"original_url"`
	ShortURL string `json:"short_url"`
	UserID   string `json:"-"`
}

type ItemList struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url,omitempty"`
	ShortURL    string `json:"short_url"`
}
