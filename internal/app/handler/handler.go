package handler

import (
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func NewHandler(adress string) *http.ServeMux {
	shortenerService := service.NewURLShortenerService()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ct := r.Header.Get("Content-Type")
			if strings.Contains(ct, "text/html") || strings.Contains(ct, "text/plain") {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error while body reading", http.StatusInternalServerError)
				} else {
					urlID := shortenerService.AddNewURL(string(body))
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(201)
					_, err := w.Write([]byte(fmt.Sprintf("http://%s/%d", adress, urlID)))
					if err != nil {
						http.Error(w, "Error while generating response", http.StatusInternalServerError)
					}
					return
				}
			}
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		case "GET":
			stringID := r.URL.Path[1:]
			id, err := strconv.Atoi(stringID)
			if err != nil {
				http.Error(w, "Invalid parameter", http.StatusBadRequest)
			} else {
				url := shortenerService.GetURLByID(id)
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
