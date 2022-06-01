package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"log"
	"net/http"
)

type ApplicationConfig struct {
	baseURL          string
	shortenerService *service.URLShortenerService
	urlRepository    repository.URLRepository
}

func main() {
	const baseURL = "localhost:8080"
	config := ApplicationConfig{
		baseURL:          "localhost:8080",
		shortenerService: service.NewURLShortenerService(repository.NewURLRepositoryInMemoryImpl()),
	}

	mux := handler.NewHandler(config.baseURL, config.shortenerService)
	s := &http.Server{Addr: baseURL, Handler: mux}
	log.Fatal(s.ListenAndServe())
}
