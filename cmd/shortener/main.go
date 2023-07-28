package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/config"
	apisave "github.com/gogapopp/shortener/internal/app/http-server/handlers/api/save"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/ping"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/redirect"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/save"
	mwAuth "github.com/gogapopp/shortener/internal/app/http-server/middlewares/auth"
	mwGzip "github.com/gogapopp/shortener/internal/app/http-server/middlewares/gzip"
	mwLogger "github.com/gogapopp/shortener/internal/app/http-server/middlewares/logger"
	"github.com/gogapopp/shortener/internal/app/lib/logger"
	"github.com/gogapopp/shortener/internal/app/storage"
)

func main() {
	ctx := context.Background()
	// парсим конфиг
	cfg := config.ParseConfig()
	// инициализируем логер
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// подключаем хранилище
	storage, err := storage.NewRepo(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// подключаем роуты и мидлвееры
	r := chi.NewRouter()
	r.Use(mwAuth.AuthMiddleware(log))
	r.Use(mwGzip.GzipMiddleware(log))
	r.Use(mwLogger.NewLogger(log))
	r.Route("/", func(r chi.Router) {
		r.Post("/", save.PostSaveHandler(log, storage, cfg))
		r.Get("/{id}", redirect.GetURLGetterHandler(log, storage, cfg))
		r.Post("/api/shorten", apisave.PostSaveJSONHandler(log, storage, cfg))
		r.Get("/ping", ping.GetPingDBHandler(log, storage, cfg))
		// r.Post("/api/shorten/batch", handlers.PostBatchJSONhHandler)
		// r.Get("/api/user/urls", handlers.GetURLs)
		// r.Delete("/api/user/urls", handlers.DeleteShortURLs)
	})

	// запускаем сервер
	log.Info("Running the server at: ", "addres: ", cfg.RunAddr)
	log.Fatal(http.ListenAndServe(cfg.RunAddr, r))
}
