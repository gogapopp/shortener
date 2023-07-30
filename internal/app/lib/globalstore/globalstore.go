// package globalstore содержит глобальное хранилище
package globalstore

import (
	"github.com/gogapopp/shortener/internal/app/lib/models"
)

// был создан ради того чтоб пройти тест 14 инкремента

type Store struct {
	database map[string][]models.UserURLs
}

// объявляем глобальную переменную, чтоб записывать в неё данные из пакетов handlers и concurrency
var GlobalStore *Store

func init() {
	GlobalStore = &Store{
		database: make(map[string][]models.UserURLs),
	}
}

// SaveURLToDatabase получает userID и соответсвенно ему сохраняет models.URL
func (s *Store) SaveURLToDatabase(userID string, shortURL string, longURL string) {
	urls := s.GetURLsFromDatabase(userID)
	urls = append(urls, models.UserURLs{
		ShortURL:    shortURL,
		OriginalURL: longURL,
	})
	s.SaveURLsToDatabase(userID, urls)
}

// GetURLsFromDatabase получает структуры models.URL соответственно переданному userUD
func (s *Store) GetURLsFromDatabase(userID string) []models.UserURLs {
	return s.database[userID]
}

// SaveURLsToDatabase сохраняет структуру, которая была обработана функцией SetDeleteFlag
func (s *Store) SaveURLsToDatabase(userID string, urls []models.UserURLs) {
	s.database[userID] = urls
}
