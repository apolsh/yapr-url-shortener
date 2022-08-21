package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/config"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
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

	authCryptoProvider := crypto.NewCCMAES256CryptoProvider(cfg.AuthSecretKey)
	var urlShortenerStorage repository.URLRepository
	if cfg.DatabaseDSN != "" {
		urlShortenerStorage = repository.NewURLRepositoryPG(cfg.DatabaseDSN)
	} else {
		urlShortenerStorage = repository.NewURLRepositoryInMemory(cfg.FileStoragePath)
	}
	defer urlShortenerStorage.Close()

	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage)
	chiHandler := handler.NewURLShortenerHandler(cfg.BaseURL, urlShortenerService, authCryptoProvider)
	chiHandler.Register(router)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
