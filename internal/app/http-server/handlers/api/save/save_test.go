package save

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	mock_save "github.com/gogapopp/shortener/internal/app/http-server/handlers/api/save/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPostSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      map[string]string
		expectedBody string
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusCreated,
			reqBody:      map[string]string{"url": "https://practicum.yandex.ru"},
			expectedBody: `{"result":"http://localhost:8080/EwHXdJfB"}`,
		},
		{
			name:         "Test #2 fail",
			expectedCode: http.StatusBadRequest,
			reqBody:      map[string]string{"url": "invalid url"},
			expectedBody: "invalid request body\n",
		},
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockURLSaver := mock_save.NewMockURLSaver(mockCtrl)
	mockURLSaver.EXPECT().SaveURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockURLSaver.EXPECT().GetShortURL(gomock.Any()).AnyTimes()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := PostSaveJSONHandler(sugar, mockURLSaver, cfg)

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
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

			if tc.expectedBody != "invalid request body\n" && !strings.HasPrefix(w.Body.String(), `{"result":"http://localhost:8080/`) {
				t.Errorf("expected body %s, but got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}
