package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/encryptor"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/logger"
)

var BaseAddr string
var RunAddr string

func main() {
	flags := config.ParseFlags()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flags.FlagRunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		flags.FlagBaseAddr = envBaseAddr
	}
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := flags.FlagBaseAddr
	encryptor.GetBaseAddr(BaseAddr)
	RunAddr = flags.FlagRunAddr

	fmt.Println("Running the server at", RunAddr)
	RunServer()
}

// RunServer запускает сервер
func RunServer() {
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(handlers.PostShortURL))
		r.Get("/{id}", logger.ResponseLogger(handlers.GetHandleURL))
		r.Post("/api/shorten", logger.RequestLogger(handlers.PostJSONHandler))
	})

	log.Fatal(http.ListenAndServe(RunAddr, r))
}
