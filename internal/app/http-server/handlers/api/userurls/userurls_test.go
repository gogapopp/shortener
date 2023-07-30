package userurls

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_userurls "github.com/gogapopp/shortener/internal/app/http-server/handlers/api/userurls/mocks"
	"github.com/gogapopp/shortener/internal/app/lib/models"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGetURLsHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		userURLs     []models.UserURLs
		err          error
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusOK,
			userURLs: []models.UserURLs{
				{
					ShortURL:    "http://localhost:8080/1",
					OriginalURL: "https://practicum.yandex.ru",
				},
				{
					ShortURL:    "http://localhost:8080/2",
					OriginalURL: "https://google.com",
				},
			},
			err: nil,
		},
		{
			name:         "Test #2 no content",
			expectedCode: http.StatusNoContent,
			userURLs:     []models.UserURLs{},
			err:          errors.New("url not found"),
		},
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserURLsGetter := mock_userurls.NewMockUserURLsGetter(mockCtrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserURLsGetter.EXPECT().GetUserURLs(gomock.Any()).Return(tc.userURLs, tc.err).AnyTimes()

			handler := GetURLsHandler(sugar, mockUserURLsGetter, nil)

			req, err := http.NewRequest("GET", "/api/user/urls", nil)
			if err != nil {
				t.Fatal(err)
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

			if len(tc.userURLs) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected %d, got %d", tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
