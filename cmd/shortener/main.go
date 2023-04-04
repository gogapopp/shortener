package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/config"
	"github.com/gogapopp/shortener/internal/app/handlers"
)

var BaseAddr string

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
	handlers.GetBaseAddr(BaseAddr)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.MainHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle)
		})
	})
	log.Fatal(http.ListenAndServe(flags.FlagRunAddr, r))
}
