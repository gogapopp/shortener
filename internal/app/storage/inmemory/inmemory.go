package inmemory

import (
	"fmt"
)

type storage struct {
	urls map[string]string
}

func NewStorage() *storage {
	return &storage{
		urls: make(map[string]string),
	}
}

func (s *storage) SaveURL(baseURL, longURL, shortURL, correlationID string) error {
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
