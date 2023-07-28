package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/lib/models"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// NewLogger создаём логгер
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	Sugar := logger.Sugar()
	Log = Sugar

	return Sugar, nil
}

func RequestBatchJSONLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		Log.Info("POST request",
			"URL", r.Host,
			"method", r.Method,
			"body", b,
		)
	})
}
