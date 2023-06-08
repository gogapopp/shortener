package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gogapopp/shortener/internal/app/models"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует логер с установленным уровнем логирования
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

// Создаём собственную реализацию метода http.ResponseWriter
type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// ResponseLogger логирует GET запрос
func ResponseLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h(&lw, r)

		Log.Info("GET request",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}

// RequestLogger логирует POST запрос
func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		h(w, r)
		duration := time.Since(start)
		Log.Info("POST request",
			zap.String("URL", r.Host),
			zap.String("method", r.Method),
			zap.Int64("duration", duration.Nanoseconds()),
			zap.String("body", string(body)),
		)
	})
}

// RequestJSONLogger логирует POST json запрос
func RequestJSONLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// декодируем json чтоб выводить строку в логи без лишних пробелов
		var b models.Request
		err = json.Unmarshal([]byte(body), &b)
		if err != nil {
			log.Fatal(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		h(w, r)
		duration := time.Since(start)
		Log.Info("POST request",
			zap.String("URL", r.Host),
			zap.String("method", r.Method),
			zap.Int64("duration", duration.Nanoseconds()),
			zap.String("body", b.URL),
		)
	})
}

// RequestBatchJSONLogger логирует POST json запрос
func RequestBatchJSONLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// декодируем json чтоб выводить строку в логи без лишних пробелов
		var b []models.BatchRequest
		err = json.Unmarshal([]byte(body), &b)
		if err != nil {
			log.Fatal(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		h(w, r)
		duration := time.Since(start)
		Log.Info("POST request",
			zap.String("URL", r.Host),
			zap.String("method", r.Method),
			zap.Int64("duration", duration.Nanoseconds()),
		)
	})
}

// ResponseDELETELogger логирует DELETE запрос
func ResponseDELETELogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// декодируем json чтоб выводить строку в логи без лишних пробелов
		var IDs []string
		err = json.Unmarshal(body, &IDs)
		if err != nil {
			log.Fatal(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		h(w, r)
		Log.Info("POST request",
			zap.String("URL", r.Host),
			zap.String("method", r.Method),
			zap.Strings("body", IDs),
		)
	})
}
