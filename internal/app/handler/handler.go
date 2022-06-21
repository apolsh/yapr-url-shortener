package handler

import (
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
	if ct := r.Header.Get("Content-Type"); strings.Contains(ct, "text/html") || strings.Contains(ct, "text/plain") {
		body, err := io.ReadAll(r.Body)
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
	var body SaveURLBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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
