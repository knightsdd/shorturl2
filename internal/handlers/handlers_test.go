package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/knightsdd/shorturl2/internal/config"
	"github.com/knightsdd/shorturl2/internal/models"
	"github.com/knightsdd/shorturl2/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-chi/chi/v5"
)

func TestGenShortUrl(t *testing.T) {
	tstorage := storage.InMemoryStorage{
		"Basgf21hA": "https://testsite.one",
		"hAGd3am6a": "https://website.q.two",
	}
	type want struct {
		status      int
		contentType string
	}
	tests := []struct {
		name    string
		storage storage.InMemoryStorage
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
	tstorage := storage.InMemoryStorage{
		"Basgf21h": "https://testsite.one",
		"hAGd3am6": "https://website.q.two",
		"sAMd3an8": "https://website.q.two/media/files/2025-01-01/s",
	}
	type want struct {
		status   int
		location string
	}
	tests := []struct {
		name      string
		storage   storage.InMemoryStorage
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
			name:      "Test #3 success",
			storage:   tstorage,
			targetUrl: `sAMd3an8`,
			method:    http.MethodGet,
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://website.q.two/media/files/2025-01-01/s",
			},
		},
		{
			name:      "Test #4 method not allow",
			storage:   tstorage,
			targetUrl: `hAGd3am6`,
			method:    http.MethodPost,
			want: want{
				status:   http.StatusMethodNotAllowed,
				location: "",
			},
		},
		{
			name:      "Test #5 bad request",
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

// MockUrlStorage - мок реализации UrlStorage для тестов
type MockUrlStorage struct {
	SaveValueFunc func(url string) string
}

func (m *MockUrlStorage) SaveValue(url string) string {
	if m.SaveValueFunc != nil {
		return m.SaveValueFunc(url)
	}
	return "abc123"
}

func (m *MockUrlStorage) GetValue(key string) (value string, ok bool) {
	return "", true
}

func setupRouter(storage storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/shorten", ShortUrlAPI(storage))
	return r
}

func TestShortUrlAPI_Success(t *testing.T) {
	// Настраиваем мок хранилища
	mockStorage := &MockUrlStorage{
		SaveValueFunc: func(url string) string {
			if url != "https://example.com" {
				t.Errorf("expected 'https://example.com', got '%s'", url)
			}
			return "test123"
		},
	}

	// Создаем тестовый запрос
	requestBody := models.ShortenAPIRequest{Url: stringPtr("https://example.com")}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder
	rr := httptest.NewRecorder()

	// Настраиваем роутер и вызываем хендлер
	r := setupRouter(mockStorage)
	r.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Проверяем заголовок Content-Type
	contentType := rr.Header().Get("content-type")
	if contentType != "application/json" {
		t.Errorf("content type header does not match: got %v want %v",
			contentType, "application/json")
	}

	// Проверяем тело ответа
	var response models.ShortenAPIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	expected := config.ServerEndpoint + "test123"
	if response.Result != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Result, expected)
	}
}

func TestShortUrlAPI_InvalidMethod(t *testing.T) {
	mockStorage := &MockUrlStorage{}

	testMethods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range testMethods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/api/shorten", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r := setupRouter(mockStorage)
			r.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("for method %s: got status %v want %v",
					method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestShortUrlAPI_InvalidJSON(t *testing.T) {
	mockStorage := &MockUrlStorage{}

	// Невалидный JSON
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer([]byte("{invalid}")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := setupRouter(mockStorage)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestShortUrlAPI_MissingURLField(t *testing.T) {
	mockStorage := &MockUrlStorage{}

	// Запрос без поля URL
	requestBody := map[string]interface{}{"not_url": "test"}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := setupRouter(mockStorage)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestShortUrlAPI_EmptyURLField(t *testing.T) {
	mockStorage := &MockUrlStorage{}

	// Запрос с пустым URL
	requestBody := models.ShortenAPIRequest{Url: stringPtr("")}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := setupRouter(mockStorage)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Вспомогательная функция для создания указателя на строку
func stringPtr(s string) *string {
	return &s
}
