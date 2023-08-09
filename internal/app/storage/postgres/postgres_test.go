package postgres

import (
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	cfg := config.ParseConfig()
	s, err := NewStorage(cfg.DatabasePath)
	assert.NoError(t, err)

	// сохраняем ссылку в хранилище
	err = s.SaveURL("https://www.example.com", "short", "correlationID", "userID")
	assert.NoError(t, err)
	// сохраняем ссылку в хранилище но ожидаем ошибку т.к. ссылка уже сохранена
	err = s.SaveURL("https://www.example.com", "short", "correlationID", "userID")
	assert.Error(t, err)

	urls := []models.BatchDatabaseResponse{
		{ShortURL: "short1", OriginalURL: "https://www.example1.com", CorrelationID: "correlationID1"},
		{ShortURL: "short2", OriginalURL: "https://www.example2.com", CorrelationID: "correlationID2"},
	}
	err = s.BatchInsertURL(urls, "userID")
	assert.NoError(t, err)

	// получаем ссылку из хранилища
	found, longURL, err := s.GetURL("short", "userID")
	assert.NoError(t, err)
	// переменная found содержит значение удалена ли ссылка
	assert.False(t, found)
	assert.Equal(t, "https://www.example.com", longURL)

	// проверяем получение ссылки с некорректным значением
	found, longURL, err = s.GetURL("notfound", "userID")
	assert.Error(t, err)
	assert.False(t, found)
	assert.Equal(t, "", longURL)
}
