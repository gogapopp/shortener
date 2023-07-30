// package main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/api/batchsave"
	apisave "github.com/gogapopp/shortener/internal/app/http-server/handlers/api/save"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/api/urlsdelete"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/api/userurls"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/ping"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/redirect"
	"github.com/gogapopp/shortener/internal/app/http-server/handlers/save"
	mwAuth "github.com/gogapopp/shortener/internal/app/http-server/middlewares/auth"
	mwGzip "github.com/gogapopp/shortener/internal/app/http-server/middlewares/gzip"
	mwLogger "github.com/gogapopp/shortener/internal/app/http-server/middlewares/logger"
	"github.com/gogapopp/shortener/internal/app/lib/logger"
	"github.com/gogapopp/shortener/internal/app/storage"
)

// main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
func main() {
	// парсим конфиг
	cfg := config.ParseConfig()
	// инициализируем логер
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// подключаем хранилище
	storage, err := storage.NewRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// проверяем установленно ли подключение к бд
	db, err := storage.Ping()
	if err == nil {
		defer db.Close()
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
		r.Post("/api/shorten/batch", logger.RequestBatchJSONLogger(batchsave.PostBatchJSONhHandler(log, storage, cfg)))
		r.Get("/api/user/urls", userurls.GetURLsHandler(log, storage, cfg))
		r.Delete("/api/user/urls", urlsdelete.DeleteHandler(log, storage, cfg))
	})
	r.Mount("/debug/pprof", pprofRoutes())

	// запускаем сервер
	log.Info("Running the server at: ", "addres: ", cfg.RunAddr)
	log.Fatal(http.ListenAndServe(cfg.RunAddr, r))
}

// pprofRoutes возвращает хендлеры нужные для профилирования
func pprofRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/block", pprof.Handler("block"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)
	return r
}
