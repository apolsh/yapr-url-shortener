package main

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"log"
	"net/http"
)

func main() {
	mux := handler.NewHandler()
	s := &http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(s.ListenAndServe())
}
