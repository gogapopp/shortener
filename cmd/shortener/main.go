// package main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/logger"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/api/batchsave"
	apisave "github.com/gogapopp/shortener/internal/app/network-server/handlers/api/save"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/api/stats"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/api/urlsdelete"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/api/userurls"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/ping"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/redirect"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/save"
	mwAuth "github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	mwGzip "github.com/gogapopp/shortener/internal/app/network-server/middlewares/gzip"
	mwLogger "github.com/gogapopp/shortener/internal/app/network-server/middlewares/logger"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/subnet"
	"github.com/gogapopp/shortener/internal/app/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// go build flags
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
func main() {
	// build stdout
	buildStdout()
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

	// получаем http или https сервер в зависимости от конфига
	httpserver := RunHTTPServer(cfg, storage, log)

	// реализация graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	// проверяем установленно ли подключение к бд
	db, err := storage.Ping()
	if err == nil {
		if err := db.Close(); err != nil {
			log.Info("error close db conn:", err)
		}
	}

	if err := httpserver.Shutdown(context.Background()); err != nil {
		log.Info("error shutdown the server:", err)
	}
}

// buildStdout выводит в строку терминала build version, build date, build commit
func buildStdout() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// RunServer инициализирует роуты и мидлвееры, запускает сервер
func RunHTTPServer(cfg *config.Config, storage storage.Storage, log *zap.SugaredLogger) *http.Server {
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
		r.With(subnet.SubnetMiddleware(log, cfg.TrustedSubnet)).Get("/api/internal/stats", stats.GetStat(log, storage, cfg))
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
