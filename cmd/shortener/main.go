package main

import (
	"log"
	"net/http"

	"github.com/gogapopp/shortener/internal/app/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
