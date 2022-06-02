package main

import (
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	const (
		appProtocol = "http"
		appDomain   = "localhost:8080"
	)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	urlShortenerStorage := repository.NewURLRepositoryInMemoryImpl()
	urlShortenerService := service.NewURLShortenerService(urlShortenerStorage)
	chiHandler := handler.NewURLShortenerHandler(fmt.Sprintf("%s://%s", appProtocol, appDomain), urlShortenerService)
	chiHandler.Register(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
