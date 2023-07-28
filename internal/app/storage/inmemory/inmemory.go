package inmemory

import (
	"database/sql"
	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
)

type storage struct {
	urls      map[string]string
	urlsBatch map[string][]struct {
		OriginalURL string
		ShortURL    string
	}
}

func NewStorage() *storage {
	return &storage{
		urls: make(map[string]string),
		urlsBatch: make(map[string][]struct {
			OriginalURL string
			ShortURL    string
		}),
	}
}

func (s *storage) SaveURL(longURL, shortURL, correlationID string, userID string) error {
	s.urls[shortURL] = longURL
	return nil
}

func (s *storage) GetURL(shortURL, userID string) (string, error) {
	longURL, ok := s.urls[shortURL]
	if !ok {
		return "", fmt.Errorf("url not found")
	}
	return longURL, nil
}

func (s *storage) Ping() (*sql.DB, error) {
	return nil, fmt.Errorf("error ping DB")
}

func (s *storage) BatchInsertURL(urls []models.BatchDatabaseResponse, userID string) error {
	for _, url := range urls {
		s.urlsBatch[userID] = append(s.urlsBatch[userID], struct {
			OriginalURL string
			ShortURL    string
		}{
			OriginalURL: url.OriginalURL,
			ShortURL:    url.ShortURL,
		})
		s.urls[url.ShortURL] = url.OriginalURL
	}
	return nil
}

func (s *storage) GetShortURL(longURL string) string {
	return ""
}

func (s *storage) GetUserURLs(userID string) ([]models.UserURLs, error) {
	var result []models.UserURLs
	for _, url := range s.urlsBatch[userID] {
		result = append(result, models.UserURLs{
			OriginalURL: url.OriginalURL,
			ShortURL:    url.ShortURL,
		})
	}
	return result, nil
}
