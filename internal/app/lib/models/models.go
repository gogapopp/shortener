// package models contains data structure models
package models

// models of structures
type (
	// Request is needed for json reading of body
	Request struct {
		// accepts LongURL
		URL string `json:"url"`
	}
	// Response is needed for a json response
	Response struct {
		// accepts shortURL
		ShortURL string `json:"result"`
	}
	// FileStorage is needed to write to a json file
	FileStorage struct {
		// unique request id
		UUID     int    `json:"uuid"`
		ShortURL string `json:"short_url"`
		// accepts LongURL
		OriginalURL string `json:"original_url"`
	}
	// BatchDatabaseResponse is needed for the server's batch response
	BatchDatabaseResponse struct {
		// accepts shortURL
		ShortURL string
		// accepts LongURL
		OriginalURL string
		// stores a unique ID
		CorrelationID string
		// user id
		UserID string
	}
	// Batch Request is needed for json reading of body
	BatchRequest struct {
		// stores a unique ID
		CorrelationID string `json:"correlation_id"`
		// accepts LongURL
		OriginalURL string `json:"original_url"`
	}
	// BatchResponse is needed for json response
	BatchResponse struct {
		// stores a unique ID
		CorrelationID string `json:"correlation_id"`
		// accepts shortURL
		ShortURL string `json:"short_url"`
	}
	// UserURLs
	UserURLs struct {
		// accepts LongURL
		OriginalURL string `json:"original_url"`
		// accepts shortURL
		ShortURL string `json:"short_url"`
	}
	// Stats
	Stasts struct {
		// number of shortener URLs in the service
		URLs int `json:"urls"`
		// number of users in the service
		UserIDs int `json:"users"`
	}
)
