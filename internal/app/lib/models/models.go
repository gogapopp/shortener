package models

type (
	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		ShortURL string `json:"result"`
	}

	FileStorage struct {
		UUID        int    `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	BatchDatabaseResponse struct {
		ShortURL      string
		OriginalURL   string
		CorrelationID string
		UserID        string
	}

	BatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	BatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	UserURLs struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}
)
