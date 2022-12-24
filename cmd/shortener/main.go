package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/config"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	httpRouter "github.com/apolsh/yapr-url-shortener/internal/app/handler/http"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// buildVersion - версия сборки.
// buildDate - дата сборки
// buildCommit - комментарий сборки
func main() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	cfg := config.Load()

	router := chi.NewRouter()

	authCryptoProvider := crypto.NewAESCryptoProvider(cfg.AuthSecretKey)
	var urlShortenerStorage repository.URLRepository
	var err error
	if cfg.DatabaseDSN != "" {
		urlShortenerStorage, err = repository.NewURLRepositoryPG(cfg.DatabaseDSN)
	} else {
		urlShortenerStorage, err = repository.NewURLRepositoryInMemory(make(map[string]entity.ShortenedURLInfo), cfg.FileStoragePath)
	}
	if err != nil {
		panic(err)
	}
	defer urlShortenerStorage.Close()

	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage, cfg.BaseURL)
	httpRouter.NewRouter(router, urlShortenerService, authCryptoProvider)

	router.Mount("/debug", middleware.Profiler())

	done := make(chan bool)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-quit
		log.Println("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", cfg.ServerAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", cfg.ServerAddress, err)
	}

	<-done
	logger.Println("Server stopped")
}
