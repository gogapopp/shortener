package save

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	mock_save "github.com/gogapopp/shortener/internal/app/network-server/handlers/http/save/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      string
		expectedBody string
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusCreated,
			reqBody:      "https://practicum.yandex.ru/",
			expectedBody: "http://localhost:8080/ABCDEF\n",
		},
		{
			name:         "Test #2 fail",
			expectedCode: http.StatusBadRequest,
			reqBody:      "invalid url",
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

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := PostSaveHandler(sugar, mockURLSaver, cfg)

			req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tc.reqBody))
			if err != nil {
				assert.NoError(t, err)
			}
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

			if tc.expectedCode != http.StatusBadRequest && !strings.HasPrefix(w.Body.String(), "http://localhost:8080/") {
				t.Errorf("expected string has prefix http://localhost:8080/, but got %s", tc.expectedBody)
			}
		})
	}
}
