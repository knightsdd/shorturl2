package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/knightsdd/shorturl2/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenShortUrl(t *testing.T) {
	tstorage := storage.UrlStorage{
		"Basgf21hA": "https://testsite.one",
		"hAGd3am6a": "https://website.q.two",
	}
	type want struct {
		status      int
		contentType string
	}
	tests := []struct {
		name    string
		storage storage.UrlStorage
		method  string
		body    string
		want    want
	}{
		{
			name:    "Test #1 success",
			storage: tstorage,
			method:  http.MethodPost,
			body:    "https://test-case-url.com",
			want: want{
				status:      http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name:    "Test #2 method not allow",
			storage: tstorage,
			method:  http.MethodPatch,
			body:    "https://test-case-url.com",
			want: want{
				status:      http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
		{
			name:    "Test #3 method not allow",
			storage: tstorage,
			method:  http.MethodGet,
			body:    "",
			want: want{
				status:      http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.method, `/`, body)
			rw := httptest.NewRecorder()

			GenShortUrl(tt.storage)(rw, request)
			response := rw.Result()

			assert.Equal(t, tt.want.status, response.StatusCode, "некорректный статус")
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "некорректный заголовок")

			if response.StatusCode == http.StatusCreated {
				rawBody, err := io.ReadAll(response.Body)
				defer response.Body.Close()

				require.NoError(t, err, "ошибка при чтении тела ответа")
				body := string(rawBody)
				postfix := body[len(body)-8:]
				url, ok := tt.storage[postfix]
				require.True(t, ok, "В хранилище нет требуемой ссылки")
				assert.Equal(t, url, tt.body)
			}
		})
	}
}

func TestGetOriginalUrl(t *testing.T) {
	tstorage := storage.UrlStorage{
		"Basgf21h": "https://testsite.one",
		"hAGd3am6": "https://website.q.two",
	}
	type want struct {
		status   int
		location string
	}
	tests := []struct {
		name      string
		storage   storage.UrlStorage
		targetUrl string
		method    string
		want      want
	}{
		{
			name:      "Test #1 success",
			storage:   tstorage,
			targetUrl: `Basgf21h`,
			method:    http.MethodGet,
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://testsite.one",
			},
		},
		{
			name:      "Test #2 success",
			storage:   tstorage,
			targetUrl: `hAGd3am6`,
			method:    http.MethodGet,
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://website.q.two",
			},
		},
		{
			name:      "Test #3 method not allow",
			storage:   tstorage,
			targetUrl: `hAGd3am6`,
			method:    http.MethodPost,
			want: want{
				status:   http.StatusMethodNotAllowed,
				location: "",
			},
		},
		{
			name:      "Test #4 bad request",
			storage:   tstorage,
			targetUrl: `XXXxxx`,
			method:    http.MethodGet,
			want: want{
				status:   http.StatusBadRequest,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, `/`+tt.targetUrl, nil)
			request.SetPathValue("prefix", tt.targetUrl)
			rw := httptest.NewRecorder()

			GetOriginalUrl(tt.storage)(rw, request)
			response := rw.Result()

			require.Equal(t, tt.want.status, response.StatusCode, "Некорректный статус")
			if tt.want.status == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, response.Header.Get("Location"), "Некорректный заголовок")
			}
		})
	}
}
