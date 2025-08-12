package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/knightsdd/shorturl2/internal/config"
	"github.com/knightsdd/shorturl2/internal/models"
	"github.com/knightsdd/shorturl2/internal/storage"
)

func GenShortUrl(storage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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
		shortUrl := config.ServerEndpoint + postfix

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortUrl))
	}
}

func GetOriginalUrl(storage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
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

func ShortUrlAPI(storage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.ShortenAPIRequest{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, `"Invalid request body"`, http.StatusBadRequest)
			return
		}
		if req.Url == nil || *req.Url == "" {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, `"field 'URL' is required"`, http.StatusBadRequest)
			return
		}
		postfix := storage.SaveValue(*req.Url)
		shortUrl := config.ServerEndpoint + postfix

		res := models.ShortenAPIResponse{
			Result: shortUrl,
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(res); err != nil {
			http.Error(w, "Internal serrver error", http.StatusInternalServerError)
			return
		}
	}
}
