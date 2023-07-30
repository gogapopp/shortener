package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/shortener/internal/app/lib/models"
)

func TestRequestBatchJSONLogger(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      []models.BatchRequest
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusOK,
			reqBody: []models.BatchRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "https://practicum.yandex.ru",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "https://google.com",
				},
			},
		},
		{
			name:         "Test #2 success",
			expectedCode: http.StatusBadRequest,
			reqBody:      nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := NewLogger()
			if err != nil {
				t.Fatal(err)
			}
			Log = logger

			handler := RequestBatchJSONLogger(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected %d, got %d", tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
