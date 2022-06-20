package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/apolsh/yapr-url-shortener/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	urlShortenerStorage := repository.NewURLRepositoryInMemory(cfg.FileStoragePath)
	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage)
	chiHandler := handler.NewURLShortenerHandler(cfg.BaseURL, urlShortenerService)
	chiHandler.Register(router)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
