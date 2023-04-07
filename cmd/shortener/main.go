package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/handlers"
	"github.com/gogapopp/shortener/internal/app/shortener"
)

var BaseAddr string
var RunAddr string

// init() инизиализирует переменные окружения переданные через flags
func init() {
	flags := config.ParseFlags()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flags.FlagRunAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		flags.FlagBaseAddr = envBaseAddr
	}
	// передаём FlagBaseAddr в handlers.go (функция записывает значение в переменную которая находится в пакете handlers)
	BaseAddr := flags.FlagBaseAddr
	shortener.GetBaseAddr(BaseAddr)
	RunAddr = flags.FlagRunAddr
}

func main() {
	RunServer()
}

// RunServer запускает сервер
func RunServer() {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostShortURL)
		r.Get("/{id}", handlers.GetHandleURL)
	})

	log.Fatal(http.ListenAndServe(RunAddr, r))
}
