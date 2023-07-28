package files

import (
	"encoding/json"
	"os"

	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
)

var UUIDCounter int
var urlFileStorage []models.FileStorage

type storage struct {
	urls            map[string]string
	fileStoragePath string
}

func NewStorage(fileStoragePath string) (*storage, error) {
	const op = "storage.files.NewStorage"
	file, err := os.OpenFile(fileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	defer file.Close()

	data, err := os.ReadFile(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	fileInfo, err := os.Stat(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	if fileInfo.Size() != 0 {
		if err := json.Unmarshal(data, &urlFileStorage); err != nil {
			return nil, fmt.Errorf("%s: %s", op, err)
		}
	}
	if len(urlFileStorage) != 0 {
		UUIDCounter = urlFileStorage[len(urlFileStorage)-1].UUID
	}

	return &storage{
		urls:            make(map[string]string),
		fileStoragePath: fileStoragePath,
	}, nil
}

func (s *storage) SaveURL(longURL, shortURL, correlationID string) error {
	const op = "storage.files.SaveURL"
	s.urls[shortURL] = longURL
	file, err := os.OpenFile(s.fileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	defer file.Close()

	UUIDCounter++
	urlFileStorage = append(urlFileStorage, models.FileStorage{
		UUID:        UUIDCounter,
		ShortURL:    shortURL,
		OriginalURL: longURL,
	})
	data, err := json.MarshalIndent(urlFileStorage, "", "   ")
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return os.WriteFile(s.fileStoragePath, data, 0666)
}

func (s *storage) GetURL(shortURL string) (string, error) {
	for _, fileURL := range urlFileStorage {
		s.urls[fileURL.ShortURL] = fileURL.OriginalURL
	}
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

func (s *storage) GetShortURL(longURL string) string {
	return ""
}
