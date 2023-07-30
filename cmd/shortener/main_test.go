package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/api/batchsave"
	apisave "github.com/gogapopp/shortener/internal/app/http-server/handlers/api/save"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/api/urlsdelete"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/ping"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/redirect"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/save"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	"github.com/gogapopp/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func BenchmarkPostSaveHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := save.PostSaveHandler(sugar, storage, cfg)

	// Запуск бенчмарка
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// создание тестового запроса
		testURL := urlshortener.ShortenerURL(cfg.BaseAddr)
		body := strings.NewReader(testURL)
		req, err := http.NewRequest("POST", "/", body)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkPostSaveJSONHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := apisave.PostSaveJSONHandler(sugar, storage, cfg)

	// запуск бенчмарка
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// создание тестового запроса
		testURL := urlshortener.ShortenerURL(cfg.BaseAddr)
		body := testURL
		data := map[string]string{"url": body}
		reqBody, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(reqBody))
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkGetURLGetterHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := redirect.GetURLGetterHandler(sugar, storage, cfg)

	// сохраняем ссылку в бд
	handlerSave := save.PostSaveHandler(sugar, storage, cfg)
	testURL := urlshortener.ShortenerURL(cfg.BaseAddr)
	body := strings.NewReader(testURL)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		b.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handlerSave.ServeHTTP(rr, req)
	// Чтение тела ответа
	respBody, err := io.ReadAll(rr.Body)
	if err != nil {
		b.Fatal(err)
	}
	// запуск бенчмарка
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// создаём реквест
		req, err := http.NewRequest("GET", string(respBody), nil)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		// создание тестового запроса
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkGetPingDBHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := ping.GetPingDBHandler(sugar, storage, cfg)

	// создание тестового запроса
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		b.Fatal(err)
	}
	// запуск бенчмарка
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkPostBatchJSONhHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := batchsave.PostBatchJSONhHandler(sugar, storage, cfg)

	// запуск бенчмарка
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// создание тестового запроса
		type Data struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		data := []Data{
			{
				CorrelationID: urlshortener.ShortenerURL(cfg.BaseAddr),
				OriginalURL:   urlshortener.ShortenerURL(cfg.BaseAddr),
			},
			{
				CorrelationID: urlshortener.ShortenerURL(cfg.BaseAddr),
				OriginalURL:   urlshortener.ShortenerURL(cfg.BaseAddr),
			},
		}
		reqBody, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(reqBody))
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		// создание тестового запроса
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkDeleteHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
		// DatabasePath: "",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := urlsdelete.DeleteHandler(sugar, storage, cfg)

	data := []string{"MWmHmO",
		"/yRxA7V",
		"/NOMtJ6",
		"/88c078e7-452d-477b-8b55-70633284c97e"}
	reqBody, err := json.Marshal(data)
	if err != nil {
		b.Fatal(err)
	}
	req, err := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(reqBody))
	if err != nil {
		b.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:  "user_id",
		Value: "user_1|dXYCnu4AZYELoxU2SrRL6OEXUqvQ8+4SOD9Q/Rw0dxI=",
	}
	req.AddCookie(cookie)
	// запуск бенчмарка
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
