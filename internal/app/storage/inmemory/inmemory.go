package inmemory

import (
	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
)

type storage struct {
	urls map[string]string
}

func NewStorage() *storage {
	return &storage{
		urls: make(map[string]string),
	}
}

func (s *storage) SaveURL(longURL, shortURL, correlationID string) error {
	s.urls[shortURL] = longURL
	return nil
}

func (s *storage) GetURL(shortURL string) (string, error) {
	longURL, ok := s.urls[shortURL]
	if !ok {
		return "", fmt.Errorf("url not found")
	}
	return longURL, nil
}

func (s *storage) Ping() error {
	return fmt.Errorf("error ping DB")
}

func (s *storage) BatchInsertURL(urls []models.BatchDatabaseResponse) error {
	return nil
}
