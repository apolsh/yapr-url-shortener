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
	"syscall"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/config"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/handler/grpc"
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
	httpRouter.NewRouter(router, urlShortenerService, authCryptoProvider, cfg.GetTrustedSubnet())

	router.Mount("/debug", middleware.Profiler())

	done := make(chan bool)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT)

	httpServer := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-quit
		log.Info("httpServer is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpServer.SetKeepAlivesEnabled(false)
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal(fmt.Errorf("could not gracefully shutdown the httpServer: %v", err))
		}
		urlShortenerStorage.Close()
		close(done)
	}()

	log.Info("Server is ready to handle requests at " + cfg.ServerAddress)

	go func() {
		startHTTPServer(cfg, httpServer)
	}()

	_, starter := grpc.GetServerStarter(":8081", urlShortenerService, authCryptoProvider, cfg.GetTrustedSubnet())

	go starter()

	<-done
	log.Info("Server stopped")
}

func startHTTPServer(cfg config.Config, server *http.Server) {
	if cfg.HTTPSEnabled {
		cer, err := tls.X509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Error(err)
			return
		}

		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}

		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("could not listen on %s: %v", cfg.ServerAddress, err))
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("could not listen on %s: %v", cfg.ServerAddress, err))
		}
	}
}
