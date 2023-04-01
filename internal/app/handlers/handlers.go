package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

var urlMap = make(map[string]string)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// читаем тело реквеста
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, "Error reading request body")
			panic(err)
		}
		mainURL := string(body)
		// делаем из обычной ссылки сжатую
		shortURL := shortenerURL(mainURL)
		// получаем url path от новой сжатой ссылки /{id} и заполняем мапу
		parsedURL, err := url.Parse(shortURL)
		if err != nil {
			panic(err)
		}
		urlMap[parsedURL.Path] = mainURL
		// отправляем ответ
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, shortURL)
	} else {
		// если r.URL.Path не "/" то выполняем
		GetURLHandle(w, r)
	}
}

// функция "сжимает" строку
func shortenerURL(url string) string {
	id := uuid.New()
	return "http://localhost:8080/" + id.String()
}

func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
