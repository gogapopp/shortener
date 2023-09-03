package batchsave

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	mock_batchsave "github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/batchsave/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostBatchJSONhHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      []models.BatchRequest
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusCreated,
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
			name:         "Test #2 invalid url",
			expectedCode: http.StatusBadRequest,
			reqBody: []models.BatchRequest{
				{
					CorrelationID: "3",
					OriginalURL:   "invalid url",
				},
			},
		},
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBatchSaver := mock_batchsave.NewMockBatchSaver(mockCtrl)
	mockBatchSaver.EXPECT().BatchInsertURL(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := PostBatchJSONhHandler(sugar, mockBatchSaver, cfg)

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				assert.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(reqBody))
			if err != nil {
				assert.NoError(t, err)
			}
			req.Header.Set("Content-Type", "application/json")
			cookie := &http.Cookie{
				Name:  "user_id",
				Value: "user_1|dXYCnu4AZYELoxU2SrRL6OEXUqvQ8+4SOD9Q/Rw0dxI=",
			}
			req.AddCookie(cookie)

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
