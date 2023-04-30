package gzipMiddleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipMiddleware проверяет сжат ли запрос и возвращает сжат
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// может ли клиент принять сжатый файл
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		// проверяем сжат ли запрос, расшифровываем
		if r.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = reader
		}

		// сжимаем ответ и ставим заголовок
		w.Header().Set("Content-Encoding", "gzip")
		zr := gzip.NewWriter(w)
		defer zr.Close()

		zrw := gzipResponseWriter{Writer: zr, ResponseWriter: w}
		h.ServeHTTP(zrw, r)
	})
}
