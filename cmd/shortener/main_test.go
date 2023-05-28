package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	testCases := []struct {
		urlPath      string
		method       string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, urlPath: "/", expectedCode: http.StatusBadRequest, expectedBody: "Invalid URL\n"},
		{method: http.MethodGet, urlPath: "/unavailable-key", expectedCode: http.StatusBadRequest, expectedBody: "Link does not exist\n"},
		{method: http.MethodPost, urlPath: "/api/shorten", expectedCode: http.StatusBadRequest, expectedBody: "Invalid URL\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.urlPath, nil)
			w := httptest.NewRecorder()

			if tc.method == http.MethodPost {
				handlers.PostShortURL(w, r)
			} else if tc.method == http.MethodGet {
				handlers.GetHandleURL(w, r)
			}

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}
