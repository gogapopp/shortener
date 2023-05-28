package middlewares

import (
	"net/http"
	"strings"
)

// GzipMiddleware проверяет сжат ли запрос и возвращает сжат
func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		// проверяем тип контента в запроса
		contentType := r.Header.Get("Content-Type")
		// если это application/json или text/html то разрешаем сжать
		if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
			// проверяем поддерживает ли клиент сжатие
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportGzip := strings.Contains(acceptEncoding, "gzip")
			// реализуем сжатие если клиент поддерживает сжатие gzip
			if supportGzip {
				cw := newCompressWriter(w)
				w.Header().Set("Content-Encoding", "gzip")
				ow = cw
				defer cw.Close()
			}
		}

		// проверяем зашифрован ли полученный запрос
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		// если зашифрован то читаем и записываем в боди
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
