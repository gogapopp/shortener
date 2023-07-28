package redirect

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/shortener/internal/app/config"
	mock_redirect "github.com/gogapopp/shortener/internal/app/http-server/handlers/redirect/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGetURLGetterHandler(t *testing.T) {
	cases := []struct {
		name         string
		reqURL       string
		expectedBody string
		mockError    error
		expectedCode int
	}{
		{
			name:         "Test #1: success",
			reqURL:       "http://localhost:8080/EwHXdJfB",
			expectedBody: "",
			mockError:    nil,
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "Test #2: fail",
			reqURL:       "http://localhost:8080/Grbg34",
			expectedBody: "url not found\n",
			mockError:    errors.New("url not found\n"),
			expectedCode: http.StatusBadRequest,
		},
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockURLGetter := mock_redirect.NewMockURLGetter(mockCtrl)

	cfg := &config.Config{
		BaseAddr: "http://localhost:8080/",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockURLGetter.EXPECT().GetURL(tc.reqURL, gomock.Any()).Return(tc.expectedBody, tc.mockError)

			handler := GetURLGetterHandler(sugar, mockURLGetter, cfg)

			req, err := http.NewRequest("GET", tc.reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			cookie := &http.Cookie{
				Name:  "user_id",
				Value: "user_1|dXYCnu4AZYELoxU2SrRL6OEXUqvQ8+4SOD9Q/Rw0dxI=",
			}
			req.AddCookie(cookie)

			w := httptest.NewRecorder()
			handler(w, req)

			if w.Body.String() != tc.expectedBody {
				t.Errorf("expected %s, got %s", tc.expectedBody, w.Body.String())
			}

			if w.Code != tc.expectedCode {
				t.Errorf("expected %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}
