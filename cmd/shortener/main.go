package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gogapopp/shortener/internal/app/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.MainHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetURLHandle)
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
