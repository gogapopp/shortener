package gzipMiddleware

import (
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// проверяем тип запроса, если json или html то используем сжатие
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := w.Header().Get("Content-Type")
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
		return w.Writer.Write(b)
	} else {
		return w.ResponseWriter.Write(b)
	}
}
