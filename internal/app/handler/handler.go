package handler

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"io"
	"net/http"
	"strconv"
)

func NewHandler() *http.ServeMux {
	shortenerService := service.NewUrlShortenerService()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ct := r.Header.Get("Content-Type")
			if ct == "text/html; charset=UTF-8" || ct == "text/plain" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error while body reading", http.StatusInternalServerError)
				} else {
					urlId := shortenerService.AddNewUrl(string(body))
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(201)
					_, err := w.Write([]byte(strconv.Itoa(urlId)))
					if err != nil {
						http.Error(w, "Error while generating response", http.StatusInternalServerError)
					}
					return
				}
			}
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		case "GET":
			stringId := r.URL.Path[1:]
			id, err := strconv.Atoi(stringId)
			if err != nil {
				http.Error(w, "Invalid parameter", http.StatusBadRequest)
			} else {
				url := shortenerService.GetUrlById(id)
				if len(url) > 0 {
					http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				} else {
					http.NotFound(w, r)
				}
			}

		default:
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}

	})

	return mux
}
