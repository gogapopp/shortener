package config

import (
	"encoding/json"
	"log"
	"os"
)

var pathStorage string
var UUIDCounter int
var URLSMap = make(map[string]string)

// принимает значение path storage из main
func GetStoragePath(str string) {
	pathStorage = str
}

type ShortURL struct {
	Uuid        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var ShortURLStorage []ShortURL

func CreateFile() {
	file, err := os.OpenFile(pathStorage, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func Save() error {
	data, err := json.MarshalIndent(ShortURLStorage, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathStorage, data, 0666)
}

func Load() error {
	data, err := os.ReadFile(pathStorage)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &ShortURLStorage); err != nil {
		return err
	}
	UUIDCounter = ShortURLStorage[len(ShortURLStorage)-1].Uuid
	return nil
}

func RestoreURL() {
	for _, urls := range ShortURLStorage {
		URLSMap[urls.ShortURL] = urls.OriginalURL
	}
}
