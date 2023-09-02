// package logger соддержит обработчик для логирования всех запросов на сервер
package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// NewLogger логирует каждый запрос
func NewLogger(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Infow("request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				log.Infow("request completed",
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration", time.Since(t1).String(),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
