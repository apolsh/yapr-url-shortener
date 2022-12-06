package http

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
)

var (
	s service.URLShortenerService
	c crypto.CryptographicProvider
)

func ExampleController_PingDB() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_GetShortenURLByID() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodGet, "/some_id", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_GetShortenURLsByUser() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_SaveShortenURL() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL1))
	req.Header.Set("Content-Type", "text/plain")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_SaveShortenURLsInBatch() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(
		`[
				{
					"correlation_id": "1",
					"original_url": "https://www.youtube.com/"
				},
				{
					"correlation_id": "2",
					"original_url": "https://www.twitch.tv/kochevnik"
				}
			] 
	`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_SaveShortenURLJSON() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://golangdocs.com/golang-read-json-file"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}

func ExampleController_DeleteShortenURLsInBatch() {
	r := chi.NewRouter()
	NewRouter(r, s, c)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["123","456"]`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
}
