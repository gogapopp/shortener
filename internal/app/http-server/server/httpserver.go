package server

import (
	"crypto/tls"
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
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

func RunServer(cfg *config.Config, storage storage.Storage, log *zap.SugaredLogger) *http.Server {
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

	// настраиваем http сервер
	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
	}

	if cfg.HTTPSEnable {
		manager := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// создаём директорию для кэширования сертификатов
			Cache: autocert.DirCache("certs"),
			// адреса которые удовлетворяют сертификату
			HostPolicy: autocert.HostWhitelist(cfg.RunAddr),
		}
		server = &http.Server{
			TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
		}
		// запуск сервер с TLS сертификатом
		go func() {
			log.Info("Running the server at: ", cfg.RunAddr, " with TLS certificate")
			if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatal("error to start the server:", err)
			}
		}()
	} else {
		// запуск сервера без TLS сертификата
		go func() {
			log.Info("Running the server at: ", cfg.RunAddr)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal("error to start the server:", err)
			}
		}()
	}

	return server
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
