package ping

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_ping "github.com/gogapopp/shortener/internal/app/http-server/handlers/ping/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGetPingDBHandler(t *testing.T) {
	cases := []struct {
		name         string
		expectedCode int
		expectedErr  error
	}{
		{
			name:         "Test #1 success",
			expectedCode: http.StatusOK,
			expectedErr:  nil,
		},
		// {
		// 	name:         "Test #2 internal server error",
		// 	expectedCode: http.StatusInternalServerError,
		// 	expectedErr:  errors.New("error ping DB"),
		// },
	}

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	sugar := log.Sugar()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDBPinger := mock_ping.NewMockDBPinger(mockCtrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockDBPinger.EXPECT().Ping().Return(nil, tc.expectedErr).AnyTimes()

			handler := GetPingDBHandler(sugar, mockDBPinger, nil)

			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				t.Fatal(err)
			}

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
