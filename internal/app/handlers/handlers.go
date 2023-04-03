package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

var urlMap = make(map[string]string)

func MainHandler(w http.ResponseWriter, r *http.Request) {
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
}

// сохраняем flags.FlagBaseAddr из main.go
var baseAddr string

func GetBaseAddr(str string) {
	baseAddr = str
}

// функция "сжимает" строку и возрващает айди
func shortenerURL(url string) string {
	id := uuid.New()
	addres := baseAddr
	if !strings.HasPrefix(baseAddr, "http://") {
		addres = "http://" + addres
	}
	if !strings.HasSuffix(baseAddr, "/") {
		addres = addres + "/"
	}
	return addres + id.String()
}

func GetURLHandle(w http.ResponseWriter, r *http.Request) {
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
