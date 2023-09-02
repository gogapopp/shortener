// package models содержит в себе модели струтур данных
package models

// модели структур
type (
	// Request нужен для json прочтения body
	Request struct {
		// принимает longURL
		URL string `json:"url"`
	}
	// Response нужен для json ответа
	Response struct {
		// принимает shortURL
		ShortURL string `json:"result"`
	}
	// FileStorage нужен для записи в файл формата json
	FileStorage struct {
		// уникальный айди запроса
		UUID     int    `json:"uuid"`
		ShortURL string `json:"short_url"`
		// принимает longURL
		OriginalURL string `json:"original_url"`
	}
	// BatchDatabaseResponse нужен для batch ответа сервера
	BatchDatabaseResponse struct {
		// принимает shortURL
		ShortURL string
		// принимает longURL
		OriginalURL string
		// хранит уникальный айди
		CorrelationID string
		// айди пользователя
		UserID string
	}
	// BatchRequest нужен для json прочтения body
	BatchRequest struct {
		// хранит уникальный айди
		CorrelationID string `json:"correlation_id"`
		// принимает longURL
		OriginalURL string `json:"original_url"`
	}
	// BatchResponse нужен для json ответа
	BatchResponse struct {
		// хранит уникальный айди
		CorrelationID string `json:"correlation_id"`
		// принимает shortURL
		ShortURL string `json:"short_url"`
	}
	// UserURLs
	UserURLs struct {
		// принимает longURL
		OriginalURL string `json:"original_url"`
		// принимает shortURL
		ShortURL string `json:"short_url"`
	}
	// Stats
	Stasts struct {
		// количество сокращённых URL в сервисе
		URLs int `json:"urls"`
		// количество пользователей в сервисе
		UserIDs int `json:"users"`
	}
)
