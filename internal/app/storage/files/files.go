package files

import (
	"encoding/json"
	"os"

	"fmt"

	"github.com/gogapopp/shortener/internal/app/lib/models"
	"go.uber.org/zap"
)

var UUIDCounter int
var urlFileStorage []models.FileStorage

type storage struct {
	log             *zap.SugaredLogger
	urls            map[string]string
	fileStoragePath string
}

func NewStorage(fileStoragePath string, log *zap.SugaredLogger) (*storage, error) {
	const op = "storage.files.NewStorage"
	file, err := os.OpenFile(fileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Info(fmt.Sprintf("%s: %s", op, err))
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(fileStoragePath)
	if err != nil {
		log.Info(fmt.Sprintf("%s: %s", op, err))
		return nil, err
	}
	fileInfo, err := os.Stat(fileStoragePath)
	if err != nil {
		log.Info(fmt.Sprintf("%s: %s", op, err))
		return nil, err
	}
	if fileInfo.Size() != 0 {
		if err := json.Unmarshal(data, &urlFileStorage); err != nil {
			log.Info(fmt.Sprintf("%s: %s", op, err))
			return nil, err
		}
	}
	if len(urlFileStorage) != 0 {
		UUIDCounter = urlFileStorage[len(urlFileStorage)-1].UUID
	}

	return &storage{
		log:             log,
		urls:            make(map[string]string),
		fileStoragePath: fileStoragePath,
	}, nil
}

func (s *storage) SaveURL(baseURL, longURL, shortURL string) error {
	const op = "storage.files.SaveURL"
	s.urls[shortURL] = longURL
	file, err := os.OpenFile(s.fileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.log.Info(fmt.Sprintf("%s: %s", op, err))
		return err
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
		s.log.Info(fmt.Sprintf("%s: %s", op, err))
		return err
	}
	return os.WriteFile(s.fileStoragePath, data, 0666)
}

func (s *storage) GetURL(longURL string) (string, error) {
	const op = "storage.files.GetURL"
	for _, fileURL := range urlFileStorage {
		s.urls[fileURL.ShortURL] = fileURL.OriginalURL
	}
	shortURL, ok := s.urls[longURL]
	if !ok {
		return "", fmt.Errorf("url not found")
	}
	return shortURL, nil
}

func (s *storage) Ping() error {
	return fmt.Errorf("error ping DB")
}
