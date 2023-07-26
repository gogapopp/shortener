package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func NewLogger(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Infow(
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"request_id", middleware.GetReqID(r.Context()),
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

// // Создаём собственную реализацию метода http.ResponseWriter
// type (
// 	responseData struct {
// 		status int
// 		size   int
// 	}

// 	loggingResponseWriter struct {
// 		http.ResponseWriter
// 		responseData *responseData
// 	}
// )

// func (r *loggingResponseWriter) Write(b []byte) (int, error) {
// 	size, err := r.ResponseWriter.Write(b)
// 	r.responseData.size += size
// 	return size, err
// }

// func (r *loggingResponseWriter) WriteHeader(statusCode int) {
// 	r.ResponseWriter.WriteHeader(statusCode)
// 	r.responseData.status = statusCode
// }

// // ResponseLogger логирует GET запрос
// func ResponseLogger(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		responseData := &responseData{
// 			status: 0,
// 			size:   0,
// 		}

// 		lw := loggingResponseWriter{
// 			ResponseWriter: w,
// 			responseData:   responseData,
// 		}

// 		h(&lw, r)

// 		Log.Info("GET request",
// 			zap.Int("status", responseData.status),
// 			zap.Int("size", responseData.size),
// 		)
// 	})
// }

// // RequestLogger логирует POST запрос
// func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		// читаем боди запоса
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// возвращаем данные обратно
// 		r.Body = io.NopCloser(bytes.NewBuffer(body))

// 		h(w, r)
// 		duration := time.Since(start)
// 		Log.Infow("POST request",
// 			"URL", r.Host,
// 			"method", r.Method,
// 			"duration", duration.Nanoseconds(),
// 			"body", string(body),
// 		)
// 	})
// }

// // RequestJSONLogger логирует POST json запрос
// func RequestJSONLogger(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		// читаем боди запоса
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// декодируем json чтоб выводить строку в логи без лишних пробелов
// 		var b models.Request
// 		err = json.Unmarshal([]byte(body), &b)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// возвращаем данные обратно
// 		r.Body = io.NopCloser(bytes.NewBuffer(body))

// 		h(w, r)
// 		duration := time.Since(start)
// 		Log.Infow("POST request",
// 			"URL", r.Host,
// 			"method", r.Method,
// 			"duration", duration.Nanoseconds(),
// 			"body", b.URL,
// 		)
// 	})
// }

// // RequestBatchJSONLogger логирует POST json запрос
// func RequestBatchJSONLogger(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		// читаем боди запоса
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// декодируем json чтоб выводить строку в логи без лишних пробелов
// 		var b []models.BatchRequest
// 		err = json.Unmarshal([]byte(body), &b)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// возвращаем данные обратно
// 		r.Body = io.NopCloser(bytes.NewBuffer(body))

// 		h(w, r)
// 		duration := time.Since(start)
// 		Log.Info("POST request",
// 			"URL", r.Host,
// 			"method", r.Method,
// 			"duration", duration.Nanoseconds(),
// 		)
// 	})
// }

// // ResponseDELETELogger логирует DELETE запрос
// func ResponseDELETELogger(h http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// читаем боди запоса
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// декодируем json чтоб выводить строку в логи без лишних пробелов
// 		var IDs []string
// 		err = json.Unmarshal(body, &IDs)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// возвращаем данные обратно
// 		r.Body = io.NopCloser(bytes.NewBuffer(body))

// 		h(w, r)
// 		Log.Infow("POST request",
// 			"URL", r.Host,
// 			"method", r.Method,
// 			"body", IDs,
// 		)
// 	})
// }
