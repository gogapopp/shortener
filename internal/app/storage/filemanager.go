package storage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gogapopp/shortener/internal/app/models"
)

var pathStorage string
var UUIDCounter int
var URLSMap = make(map[string]string)

var ShortURLStorage []models.ShortURL

// GetStoragePath принимает значение path storage из main
func GetStoragePath(str string) {
	pathStorage = str
}

// CreateFile создаёт файл с названием из pathStorage
func CreateFile() {
	file, err := os.OpenFile(pathStorage, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

// Save записывает поля json из стуктуры ShortURLStorage в файл
func Save() error {
	data, err := json.MarshalIndent(ShortURLStorage, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathStorage, data, 0666)
}

// Load читает файл с сохранёнными ссылками и записываем в структуру ShortURLStorage
func Load() error {
	data, err := os.ReadFile(pathStorage)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &ShortURLStorage); err != nil {
		return err
	}
	UUIDCounter = ShortURLStorage[len(ShortURLStorage)-1].UUID
	return nil
}

// RestoreURL записывает данные из структуры ShortURLStorage в мапу URLSMap
func RestoreURL() {
	for _, urls := range ShortURLStorage {
		URLSMap[urls.ShortURL] = urls.OriginalURL
	}
}
