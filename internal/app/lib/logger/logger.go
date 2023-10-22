// package logger creates a logger instance
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

// contains a logger
var Log *zap.SugaredLogger

// NewLogger creating a logger
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

// RequestBatchJSONLogger logs all batch json requests
func RequestBatchJSONLogger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// reading the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// decode json to output a string to logs without extra spaces
		var b []models.BatchRequest
		err = json.Unmarshal([]byte(body), &b)
		if err != nil {
			log.Fatal(err)
		}
		// returning the data back
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		h(w, r)
		Log.Info("POST request",
			"URL", r.Host,
			"method", r.Method,
			"body", b,
		)
	})
}
