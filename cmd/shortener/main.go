package main

import (
	"fmt"
	"net/http"

	"github.com/knightsdd/shorturl2/internal/config"
	"github.com/knightsdd/shorturl2/internal/handlers"
	"github.com/knightsdd/shorturl2/internal/storage"

	"github.com/go-chi/chi/v5"
)

func MainRouter(storage storage.UrlStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handlers.GenShortUrl(storage))
	r.Get("/{prefix}", handlers.GetOriginalUrl(storage))
	return r
}

func main() {
	config.ParseServerFlags()
	runAddress := config.GetServerRunAddress()
	urlStorage := storage.GetStorage()
	r := MainRouter(urlStorage)

	fmt.Println("Сервер запущен на:", runAddress)
	http.ListenAndServe(runAddress, r)
}
