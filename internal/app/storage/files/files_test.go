package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	s, err := NewStorage("files.json")
	assert.NoError(t, err)

	// сохраняем ссылку в хранилище
	err = s.SaveURL("https://www.example.com", "short", "correlationID", "userID")
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

	// удаляем созданный для теста файл
	err = os.Remove("files.json")
	assert.NoError(t, err)
}
