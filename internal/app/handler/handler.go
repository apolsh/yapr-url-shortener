package handler

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/middleware"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
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
	address            string
	service            service.URLShortenerService
	authCryptoProvider crypto.Provider
}

func isValidContentType(contentType string, allowedTypes ...string) bool {
	for _, allowed := range allowedTypes {
		if strings.Contains(contentType, allowed) {
			return true
		}
	}
	return false
}

func NewURLShortenerHandler(appAddress string, serviceImpl service.URLShortenerService, provider crypto.Provider) Handler {
	return &handler{
		address:            appAddress,
		service:            serviceImpl,
		authCryptoProvider: provider,
	}
}

func (h *handler) Register(router *chi.Mux) {
	router.Use(middleware.CompressResponse)
	router.Use(middleware.AuthMiddleware(h.authCryptoProvider))
	router.Route("/", func(r chi.Router) {
		r.Get("/ping", h.PingDB)
		r.Get("/{urlID}", h.GetURLHandler)
		r.Post("/api/shorten/batch", h.PostShortenURLsInBatch)
		r.Get("/api/user/urls", h.GetUserURLsHandler)
		r.Post("/", h.SaveURLHandler)
		r.Post("/api/shorten", h.SaveURLJSONHandler)

	})
}

func (h *handler) PostShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	if isValidContentType(r.Header.Get("Content-Type"), "application/json", "application/x-gzip") {
		reader, err := getBodyReader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()

		var body []dto.ShortenInBatchRequestItem

		if err := json.NewDecoder(reader).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		ownerID := r.Context().Value(middleware.OwnerID).(string)
		batch, err := h.service.AddNewURLsInBatch(ownerID, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, item := range batch {
			item.ShortURL = h.createShortURLFromID(item.ShortURL)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(batch); err != nil {
			http.Error(w, "Error while generating response", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Invalid Content-Type: "+r.Header.Get("Content-Type"), http.StatusBadRequest)
}

func (h *handler) PingDB(w http.ResponseWriter, r *http.Request) {
	ok := h.service.PingDB()
	if ok {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if urlID := chi.URLParam(r, "urlID"); urlID != "" {
		url, err := h.service.GetURLByID(urlID)
		if errors.Is(repository.ErrorItemNotFound, err) {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "Invalid parameter", http.StatusMethodNotAllowed)
}

type GetUserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewGetUserURLsResponse(shortURL, originalURL string) *GetUserURLsResponse {
	return &GetUserURLsResponse{ShortURL: shortURL, OriginalURL: originalURL}
}

func (h *handler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(middleware.OwnerID).(string)
	shortenedURLSInfos, err := h.service.GetURLsByOwnerID(ownerID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	if len(shortenedURLSInfos) == 0 {
		w.WriteHeader(204)
		return
	}

	response := make([]GetUserURLsResponse, 0, len(shortenedURLSInfos))

	for _, info := range shortenedURLSInfos {
		response = append(response, *NewGetUserURLsResponse(h.createShortURLFromID(info.GetID()), info.GetOriginalURL()))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
	}
}

func (h *handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if isValidContentType(r.Header.Get("Content-Type"), "text", "text/plain", "application/x-gzip") {
		reader, err := getBodyReader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, "Error while body reading", http.StatusInternalServerError)
			return
		}
		urlString := string(body)
		_, err = url.ParseRequestURI(urlString)
		if err != nil {
			http.Error(w, "Passed value is not valid URL", http.StatusBadRequest)
			return
		}

		ownerID := r.Context().Value(middleware.OwnerID).(string)
		urlID, _ := h.service.AddNewURL(*entity.NewUnstoredShortenedURLInfo(ownerID, urlString))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(201)
		_, err = w.Write([]byte(h.createShortURLFromID(urlID)))
		if err != nil {
			http.Error(w, "Error while generating response", http.StatusInternalServerError)
		}
		return

	}
	http.Error(w, "Invalid Content-Type: "+r.Header.Get("Content-Type"), http.StatusBadRequest)
}

func (h *handler) createShortURLFromID(shortURLID string) string {
	return fmt.Sprintf("%s/%s", h.address, shortURLID)
}

func (h *handler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	if isValidContentType(r.Header.Get("Content-Type"), "application/json", "application/x-gzip") {
		reader, err := getBodyReader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()

		var body SaveURLBody

		if err := json.NewDecoder(reader).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, err = url.ParseRequestURI(body.URL)
		if err != nil {
			http.Error(w, "Passed value is not valid URL", http.StatusBadRequest)
			return
		}
		ownerID := r.Context().Value(middleware.OwnerID).(string)
		urlID, _ := h.service.AddNewURL(*entity.NewUnstoredShortenedURLInfo(ownerID, body.URL))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(&SaveURLResponse{Result: h.createShortURLFromID(urlID)}); err != nil {
			http.Error(w, "Error while generating response", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Invalid Content-Type: "+r.Header.Get("Content-Type"), http.StatusBadRequest)
}

func getBodyReader(r *http.Request) (io.ReadCloser, error) {
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		return gz, nil
	}
	return r.Body, nil
}
