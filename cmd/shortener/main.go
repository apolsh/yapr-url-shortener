package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"log"
	"net/http"
)

func main() {
	const baseURL = "localhost:8080"

	mux := handler.NewHandler(baseURL)
	s := &http.Server{Addr: baseURL, Handler: mux}
	log.Fatal(s.ListenAndServe())
}
