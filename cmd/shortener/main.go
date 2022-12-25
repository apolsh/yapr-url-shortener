package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
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
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/inmemory"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/postgres"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/apolsh/yapr-url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"

	log = logger.LoggerOfComponent("main")
)

//go:embed tls/cert.pem
var tlsCert []byte

//go:embed tls/key.pem
var tlsKey []byte

// buildVersion - версия сборки
// buildDate - дата сборки
// buildCommit - комментарий сборки.
func main() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	cfg := config.Load()
	logger.SetGlobalLevel(cfg.LogLevel)

	router := chi.NewRouter()

	authCryptoProvider := crypto.NewAESCryptoProvider(cfg.AuthSecretKey)
	var urlShortenerStorage repository.URLRepository
	var err error
	if cfg.DatabaseDSN != "" {
		urlShortenerStorage, err = postgres.NewURLRepositoryPG(cfg.DatabaseDSN)
	} else {
		urlShortenerStorage, err = inmemory.NewURLRepositoryInMemory(make(map[string]entity.ShortenedURLInfo), cfg.FileStoragePath)
	}
	if err != nil {
		log.Fatal(err)
	}

	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage, cfg.BaseURL)
	httpRouter.NewRouter(router, urlShortenerService, authCryptoProvider)

	router.Mount("/debug", middleware.Profiler())

	done := make(chan bool)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-quit
		log.Info("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(fmt.Errorf("Could not gracefully shutdown the server: %v\n", err))
		}
		urlShortenerStorage.Close()
		close(done)
	}()

	log.Info("Server is ready to handle requests at " + cfg.ServerAddress)

	if cfg.HTTPSEnabled {
		cer, err := tls.X509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Error(err)
			return
		}

		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}

		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("Could not listen on %s: %v\n", cfg.ServerAddress, err))
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("Could not listen on %s: %v\n", cfg.ServerAddress, err))
		}
	}

	<-done
	log.Info("Server stopped")
}
