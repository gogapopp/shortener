// package inmemory реализация интерфейса Storage для записи в файл
package inmemory

import (
	"database/sql"
	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
)

// storage хранилище ссылок
type storage struct {
	urls      map[string]string
	urlsBatch map[string][]struct {
		OriginalURL string
		ShortURL    string
	}
}

// NewStorage создаёт хранилище storage
func NewStorage() *storage {
	return &storage{
		urls: make(map[string]string),
		urlsBatch: make(map[string][]struct {
			OriginalURL string
			ShortURL    string
		}),
	}
}

// SaveURL сохраняет ссылки в хранилище
func (s *storage) SaveURL(longURL, shortURL, correlationID string, userID string) error {
	s.urls[shortURL] = longURL
	return nil
}

// GetURL получает ссылку из хранилища
func (s *storage) GetURL(shortURL, userID string) (bool, string, error) {
	longURL, ok := s.urls[shortURL]
	if !ok {
		return false, "", fmt.Errorf("url not found")
	}
	return false, longURL, nil
}

// Ping() проверяет подключение к базе данных
func (s *storage) Ping() (*sql.DB, error) {
	return nil, fmt.Errorf("error ping DB")
}

// BatchInsertURL реализует batch запись скоращённых ссылок в хранилище
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

// GetShortURL получает короткую ссылку из хранилища
func (s *storage) GetShortURL(longURL string) string {
	return ""
}

// GetUserURLs возвращает ссылки которые сохранял определённый пользователь
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

// SetDeleteFlag реализует логику удаления ссылок из хранилища
func (s *storage) SetDeleteFlag(IDs []string, userID string) error {
	return nil
}
