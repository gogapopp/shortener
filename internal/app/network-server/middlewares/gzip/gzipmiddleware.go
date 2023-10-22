// package gzip contains an implementation of the request compression method
package gzip

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// GZipMiddleware checks if the request is compressed and returns compressed
func GzipMiddleware(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("gzip middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w
			// checking the content type in the request
			contentType := r.Header.Get("Content-Type")
			// if it is application/json or text/html, then we allow to compress
			if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
				// check if the client supports compression
				acceptEncoding := r.Header.Get("Accept-Encoding")
				supportGzip := strings.Contains(acceptEncoding, "gzip")
				// implement compression if the client supports gzip compression
				if supportGzip {
					cw := newCompressWriter(w)
					w.Header().Set("Content-Encoding", "gzip")
					ow = cw
					defer cw.Close()
				}
			}
			// check if the received request is encrypted
			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			// if encrypted, then we read and write in the body
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}
			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
