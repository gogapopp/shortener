package models

// используем в handlers
type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"result"`
}

// используем в config - filemanager.go
type ShortURL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
