package handlers

import (
	"io"
	"net/http"

	"github.com/knightsdd/shorturl2/internal/storage"
)

func GenShortUrl(storage storage.UrlStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid body", http.StatusBadGateway)
			return
		}
		postfix := storage.SaveValue(string(body))
		schema := "http"
		if r.TLS != nil {
			schema = "https"
		}
		shortUrl := schema + "://" + r.Host + "/" + postfix

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortUrl))
	}
}

func GetOriginalUrl(storage storage.UrlStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		prefix := r.PathValue("prefix")
		if originalUrl, ok := storage.GetValue(prefix); ok {
			w.Header().Set("Location", originalUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
