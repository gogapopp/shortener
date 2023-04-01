package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandler(t *testing.T) {
	testCases := []struct {
		urlPath      string
		method       string
		expectedCode int
		contentType  string
		expectedBody string
	}{
		{method: http.MethodPost, urlPath: "/", expectedCode: http.StatusCreated},
		{method: http.MethodGet, urlPath: "/random-key", expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.urlPath, nil)
			w := httptest.NewRecorder()

			mainHandler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}
