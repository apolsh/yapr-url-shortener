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

// Server базовый интерфейс для серверов различного типа
type Server interface {
	Start() error
	StartTLS(config *tls.Config) error
	Stop(ctx context.Context) error
}

var _ Server = (*grpc.GRPCServer)(nil)
var _ Server = (*httpRouter.HTTPServer)(nil)

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

	httpServerConfigs := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	grpcServer := grpc.NewGRPCServer(cfg.GRPCServerAddress, urlShortenerService, authCryptoProvider, cfg.GetTrustedSubnet())
	httpServer := httpRouter.NewHTTPServer(httpServerConfigs)

	go func() {
		<-quit
		log.Info("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := grpcServer.Stop(ctx)
		if err != nil {
			log.Fatal(fmt.Errorf("could not gracefully shutdown the grpc server: %v", err))
		}

		if err := httpServer.Stop(ctx); err != nil {
			log.Fatal(fmt.Errorf("could not gracefully shutdown the http server: %v", err))
		}
		urlShortenerStorage.Close()
		close(done)
	}()

	if cfg.HTTPSEnabled {
		tlsConfig, err := getTLSConfig()
		if err != nil {
			log.Fatal(fmt.Errorf("could not get TLS configs %v", err))
		}

		go func() {
			err := grpcServer.StartTLS(tlsConfig)
			if err != nil {
				log.Fatal(fmt.Errorf("could not start grpc server: %v", err))
			}
		}()

		err = httpServer.StartTLS(tlsConfig)
		if err != nil {
			log.Fatal(fmt.Errorf("could not start grpc server: %v", err))
		}

	} else {
		go func() {
			err := grpcServer.Start()
			if err != nil {
				log.Fatal(fmt.Errorf("could not start grpc server: %v", err))
			}
		}()

		err = httpServer.Start()
		if err != nil {
			log.Fatal(fmt.Errorf("could not start grpc server: %v", err))
		}

	}

	<-done
	log.Info("Server stopped")
}

func getTLSConfig() (*tls.Config, error) {
	cer, err := tls.X509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
}
