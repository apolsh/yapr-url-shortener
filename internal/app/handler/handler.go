package handler

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/middleware"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SaveURLBody struct {
	URL string `json:"url"`
}

type SaveURLResponse struct {
	Result string `json:"result"`
}

type Handler interface {
	Register(router *chi.Mux)
}

type handler struct {
	address string
	service service.URLShortenerService
}

func NewURLShortenerHandler(appAddress string, serviceImpl service.URLShortenerService) Handler {
	return &handler{
		address: appAddress,
		service: serviceImpl,
	}
}

func (h *handler) Register(router *chi.Mux) {
	router.Use(middleware.GZIPHandle)
	router.Route("/", func(r chi.Router) {
		r.Get("/{urlID}", h.GetURLHandler)
		r.Post("/", h.SaveURLHandler)
		r.Post("/api/shorten", h.SaveURLJSONHandler)
	})
}

func (h *handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if urlID := chi.URLParam(r, "urlID"); urlID != "" {
		id, err := strconv.Atoi(urlID)
		if err != nil {
			http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		}
		url := h.service.GetURLByID(id)
		if url != "" {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
		return
	}
	http.Error(w, "Invalid parameter", http.StatusBadRequest)
}

func (h *handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); strings.Contains(ct, "text") || strings.Contains(ct, "text/plain") || strings.Contains(ct, "application/x-gzip") {
		var reader io.Reader
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, "Error while body reading", http.StatusInternalServerError)
		} else {
			urlID := h.service.AddNewURL(string(body))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(201)
			_, err := w.Write([]byte(fmt.Sprintf("%s/%d", h.address, urlID)))
			if err != nil {
				http.Error(w, "Error while generating response", http.StatusInternalServerError)
			}
			return
		}

	}
	http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
}

func (h *handler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	var body SaveURLBody

	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	urlID := h.service.AddNewURL(body.URL)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	responseURL := fmt.Sprintf("%s/%d", h.address, urlID)
	if err := json.NewEncoder(w).Encode(&SaveURLResponse{Result: responseURL}); err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
	}
}
