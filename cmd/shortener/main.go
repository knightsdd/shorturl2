package main

import (
	"net/http"

	"github.com/knightsdd/shorturl2/internal/handlers"
	"github.com/knightsdd/shorturl2/internal/storage"

	"github.com/go-chi/chi/v5"
)

// example with only net http mux server
// func main() {
// 	urlStorage := storage.GetStorage()
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/{$}", handlers.GenShortUrl(urlStorage))
// 	mux.HandleFunc("GET /{prefix}/{$}", handlers.GetOriginalUrl(urlStorage))

// 	err := http.ListenAndServe(`:8080`, mux)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func MainRouter(storage storage.UrlStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/{$}", handlers.GenShortUrl(storage))
	r.Get("/{prefix}/{$}", handlers.GetOriginalUrl(storage))
	return r
}

func main() {
	urlStorage := storage.GetStorage()
	r := MainRouter(urlStorage)
	http.ListenAndServe(":8080", r)
}
