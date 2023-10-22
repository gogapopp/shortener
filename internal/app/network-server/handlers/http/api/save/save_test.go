package save

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/urlshortener"
	mock_save "github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/save/mocks"
	"github.com/gogapopp/shortener/internal/app/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      map[string]string
		expectedBody string
		checkCookie  bool
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
				assert.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
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

			if tc.expectedBody != "invalid request body\n" && !strings.HasPrefix(w.Body.String(), `{"result":"http://localhost:8080/`) {
				t.Errorf("expected body %s, but got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}
func BenchmarkPostSaveJSONHandler(b *testing.B) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
		// DatabasePath: "host=localhost user=postgres database=postgres port=5432 password=123",
	}
	storage, _ := storage.NewRepo(cfg)
	// инициализация хендлера
	handler := PostSaveJSONHandler(sugar, storage, cfg)

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
			assert.NoError(b, err)
		}
		b.StartTimer()

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
