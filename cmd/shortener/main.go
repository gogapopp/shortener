package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/logger"
	"github.com/gogapopp/shortener/internal/app/middlewares/gzipmw"
	mwLogger "github.com/gogapopp/shortener/internal/app/middlewares/logger"
	"github.com/gogapopp/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	// инициализируем логер
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// парсим флаги
	config.InitializeServerConfig()
	// пытается подключится к базе данных, если не получается то пытаеимся записывать в файл, если не получается то в память
	if err := config.SetupDatabaseAndFilemanager(ctx); err != nil {
		log.Info("ошибка: ", err)
		os.Exit(1)
	}
	defer storage.DB().Close()

	// роуты
	r := chi.NewRouter()
	r.Use(gzipmw.GzipMiddleware())
	r.Use(mwLogger.NewLogger(log))
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostShortURL)
		r.Get("/{id}", handlers.GetHandleURL)
		r.Get("/ping", handlers.GetPingDatabase)
		r.Post("/api/shorten", handlers.PostJSONHandler)
		r.Post("/api/shorten/batch", handlers.PostBatchJSONhHandler)
		r.Get("/api/user/urls", handlers.GetURLs)
		r.Delete("/api/user/urls", handlers.DeleteShortURLs)
	})

	// запускаем сервер
	logger.Log.Info("Running the server at", zap.String("addres", config.RunAddr))
	log.Fatal(http.ListenAndServe(config.RunAddr, r))
}
