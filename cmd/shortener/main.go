// package main реализует вызов всех компонентов нужных для работы сервера и запускает сервер
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/config"
	"github.com/gogapopp/shortener/internal/app/lib/logger"
	mygrpc "github.com/gogapopp/shortener/internal/app/network-server/handlers/grpc"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/batchsave"
	apisave "github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/save"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/stats"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/urlsdelete"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/api/userurls"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/ping"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/redirect"
	"github.com/gogapopp/shortener/internal/app/network-server/handlers/http/save"
	mwAuth "github.com/gogapopp/shortener/internal/app/network-server/middlewares/auth"
	mwGzip "github.com/gogapopp/shortener/internal/app/network-server/middlewares/gzip"
	mwLogger "github.com/gogapopp/shortener/internal/app/network-server/middlewares/logger"
	"github.com/gogapopp/shortener/internal/app/network-server/middlewares/subnet"
	"github.com/gogapopp/shortener/internal/app/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
)

// go build flags
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// build stdout
	buildStdout()
	cfg := config.ParseConfig()
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	// connecting to the storage
	storage, err := storage.NewRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// if postgres is connected, then do not forget to close the connection
	db, err := storage.Ping()
	if err == nil {
		defer db.Close()
	}
	// running http or https server depending on the config
	httpserver := RunHTTPServer(cfg, storage, log)
	// getting an rpc server
	grpcserver := RunGRPCServer(cfg, storage)

	// graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	grpcserver.GracefulStop()

	if err := httpserver.Shutdown(context.Background()); err != nil {
		log.Info("error shutdown the httpserver:", err)
	}
}

// buildStdout outputs the build version, build date, and build commit to the settings line
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

// RunGrpc starts grpc
func RunGRPCServer(cfg *config.Config, storage storage.Storage) *grpc.Server {
	// getting an grpc server
	grpcserver := mygrpc.NewGrpcServer(cfg, storage)
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal(err)
	}
	// launching the grpc server
	if err := grpcserver.Serve(listen); err != nil {
		log.Fatal(err)
	}
	return grpcserver
}

// RunServer initializes routers and middleweers, starts the server
func RunHTTPServer(cfg *config.Config, storage storage.Storage, log *zap.SugaredLogger) *http.Server {
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
	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
	}
	if cfg.HTTPSEnable {
		manager := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// creating a directory for caching certificates
			Cache: autocert.DirCache("certs"),
			// addresses that satisfy the certificate
			HostPolicy: autocert.HostWhitelist(cfg.RunAddr),
		}
		server = &http.Server{
			TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
		}
		// running a server with a TLS certificate
		go func() {
			log.Info("Running the server at: ", cfg.RunAddr, " with TLS certificate")
			if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatal("error to start the server:", err)
			}
		}()
	} else {
		// starting the server without a TLS certificate
		go func() {
			log.Info("Running the server at: ", cfg.RunAddr)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal("error to start the server:", err)
			}
		}()
	}

	return server
}

// pprofRoutes returns the handlers needed for profiling
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
