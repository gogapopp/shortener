package routes

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/middlewares"
)

// Routes инициализирует логгер и реализует роуты сервера
func Routes() *chi.Mux {
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(middlewares.GzipMiddleware(handlers.PostShortURL)))
		r.Get("/{id}", logger.ResponseLogger(middlewares.GzipMiddleware(handlers.GetHandleURL)))
		r.Get("/ping", logger.ResponseLogger(handlers.GetPingDatabase))
		r.Post("/api/shorten", logger.RequestJSONLogger(middlewares.GzipMiddleware(handlers.PostJSONHandler)))
	})

	return r
}
