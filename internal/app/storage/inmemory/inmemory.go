package inmemory

import (
	"fmt"

	"go.uber.org/zap"
)

type storage struct {
	log  *zap.SugaredLogger
	urls map[string]string
}

func NewStorage(log *zap.SugaredLogger) *storage {
	return &storage{
		log:  log,
		urls: make(map[string]string),
	}
}

func (s *storage) SaveURL(baseURL, longURL, shortURL string) {
	s.urls[shortURL] = longURL
}

func (s *storage) GetURL(longURL string) (string, error) {
	shortURL, ok := s.urls[longURL]
	if !ok {
		return "", fmt.Errorf("url not found")
	}
	return shortURL, nil
}
