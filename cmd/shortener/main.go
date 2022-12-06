package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/apolsh/yapr-url-shortener/internal/app/config"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	httpRouter "github.com/apolsh/yapr-url-shortener/internal/app/handler/http"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	buildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", BuildDate)
	fmt.Println("Build commit: ", BuildCommit)

	cfg := config.Load()

	router := chi.NewRouter()

	authCryptoProvider := crypto.NewAESCryptoProvider(cfg.AuthSecretKey)
	var urlShortenerStorage repository.URLRepository
	var err error
	if cfg.DatabaseDSN != "" {
		urlShortenerStorage, err = repository.NewURLRepositoryPG(cfg.DatabaseDSN)
	} else {
		urlShortenerStorage, err = repository.NewURLRepositoryInMemory(cfg.FileStoragePath)
	}
	if err != nil {
		panic(err)
	}
	defer urlShortenerStorage.Close()

	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage, cfg.BaseURL)
	httpRouter.NewRouter(router, urlShortenerService, authCryptoProvider)

	router.Mount("/debug", middleware.Profiler())

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
