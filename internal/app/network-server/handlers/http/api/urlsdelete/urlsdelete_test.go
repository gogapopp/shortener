package urlsdelete

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_urlsdelete "github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/urlsdelete/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		reqBody      []string
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusAccepted,
			reqBody:      []string{"6qxTVvsy", "RTfd56hn", "Jlfd67ds"},
		},
		{
			name:         "Test #1 success",
			expectedCode: http.StatusAccepted,
			reqBody:      []string{""},
		},
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockURLDeleter := mock_urlsdelete.NewMockURLDeleter(mockCtrl)
	mockURLDeleter.EXPECT().SetDeleteFlag(gomock.Any(), gomock.Any()).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := DeleteHandler(sugar, mockURLDeleter, nil)

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				assert.NoError(t, err)
			}

			req, err := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(reqBody))
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
