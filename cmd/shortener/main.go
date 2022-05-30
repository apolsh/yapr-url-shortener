package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"log"
	"net/http"
)

func main() {
	const baseUrl = "localhost:8080"

	mux := handler.NewHandler(baseUrl)
	s := &http.Server{Addr: baseUrl, Handler: mux}
	log.Fatal(s.ListenAndServe())
}
