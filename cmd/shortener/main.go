package main

import (
	"net/http"

	"github.com/knightsdd/shorturl2/internal/handlers"
	"github.com/knightsdd/shorturl2/internal/storage"
)

func main() {
	urlStorage := storage.GetStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", handlers.GenShortUrl(urlStorage))
	mux.HandleFunc("GET /{prefix}/{$}", handlers.GetOriginalUrl(urlStorage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
