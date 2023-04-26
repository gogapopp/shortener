package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/models"
)

var urlMap = make(map[string]string)

// PostShortURL получает ссылку в body и присваивает ей уникальный ключ, значение хранит в мапе "key": "url"
func PostShortURL(w http.ResponseWriter, r *http.Request) {
	// читаем тело реквеста
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusMethodNotAllowed)
		return
	}
	mainURL := string(body)
	// делаем из обычной ссылки сжатую
	shortURL := encryptor.ShortenerURL(mainURL)
	// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	// сохраняем значение в мапу
	urlMap[parsedURL.Path] = mainURL
	// отправляем ответ
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortURL)
}

// GetHandleURL проверяет валидная ссылка или нет, если валидная то редиректит по адрессу
func GetHandleURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	// проверяем есть ли значение в мапе
	if _, ok := urlMap[id]; ok {
		w.Header().Add("Location", urlMap[id])
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// PostJSONHandler принимает url: url и возвращает result: shortUrl
func PostJSONHandler(w http.ResponseWriter, r *http.Request) {
	// декодируем данные из тела запроса
	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusMethodNotAllowed)
		return
	}

	shortURL := encryptor.ShortenerURL(req.URL)
	// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	// сохраняем значение в мапу
	urlMap[parsedURL.Path] = req.URL

	// передаём значение в ответ
	var resp models.Response
	resp.ShortURL = shortURL

	// устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
