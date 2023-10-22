// package global store contains global storage
// this code was needed to pass the tests of my training project
// in principle, you can do without it
package globalstore

import (
	"github.com/gogapopp/shortener/internal/app/lib/models"
)

// Store contains a repository for the links saved by the user
type Store struct {
	// [userId]: {"LongURL":long URL, "short URL":"short URL"}
	database map[string][]models.UserURLs
}

// declaring a global variable to write data from the handlers and concurrency packages to it
var GlobalStore *Store

func init() {
	GlobalStore = &Store{
		database: make(map[string][]models.UserURLs),
	}
}

// SaveURL To Database gets userID and saves models.URL accordingly
func (s *Store) SaveURLToDatabase(userID string, shortURL string, longURL string) {
	urls := s.GetURLsFromDatabase(userID)
	urls = append(urls, models.UserURLs{
		ShortURL:    shortURL,
		OriginalURL: longURL,
	})
	s.SaveURLsToDatabase(userID, urls)
}

// GetURLs From Database gets the models.URL structures according to the passed userUD
func (s *Store) GetURLsFromDatabase(userID string) []models.UserURLs {
	return s.database[userID]
}

// SaveURLs To Database saves the structure that was processed by the SetDeleteFlag function
func (s *Store) SaveURLsToDatabase(userID string, urls []models.UserURLs) {
	s.database[userID] = urls
}
