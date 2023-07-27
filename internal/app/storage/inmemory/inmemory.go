package inmemory

import (
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
	s.urls[longURL] = shortURL
}
