package handler

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/middleware"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
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

const decodeRequestBodyError = "failed to decode body"
const encodeResponseBodyError = "failed to encode body response"
const invalidContentTypeError = "invalid content type"
const bodyReadingError = "Error while body reading"
const parseURLError = "Passed value is not valid URL"

type handler struct {
	address            string
	service            service.URLShortenerService
	authCryptoProvider crypto.CryptographicProvider
}

func isValidContentType(r *http.Request, allowedTypes ...string) bool {
	actualContentType := r.Header.Get("Content-Type")

	for _, allowed := range allowedTypes {
		if strings.Contains(actualContentType, allowed) {
			return true
		}
	}
	return false
}

func NewURLShortenerHandler(appAddress string, serviceImpl service.URLShortenerService, provider crypto.CryptographicProvider) Handler {
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
		r.Delete("/api/user/urls", h.DeleteShortenURLsInBatch)
	})
}

func (h *handler) DeleteShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "application/json", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	reader, err := getBodyReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	var ids []*string

	if err := json.NewDecoder(reader).Decode(&ids); err != nil || ids == nil || len(ids) == 0 {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}
	//ownerID := r.Context().Value(middleware.OwnerID).(string)
	//err = h.service.DeleteURLsInBatch(ownerID, ids)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) PostShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "application/json", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	reader, err := getBodyReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	var body []*dto.ShortenInBatchRequestItem

	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
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
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(batch); err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
	}
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
		foundURL, err := h.service.GetURLByID(urlID)
		if errors.Is(repository.ErrorItemNotFound, err) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(service.ErrorItemIsDeleted, err) {
			http.Error(w, "", http.StatusGone)
			return
		}
		http.Redirect(w, r, foundURL, http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "Invalid parameter", http.StatusMethodNotAllowed)
}

func (h *handler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(middleware.OwnerID).(string)
	shortenedURLSInfos, err := h.service.GetURLsByOwnerID(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(shortenedURLSInfos) == 0 {
		w.WriteHeader(204)
		return
	}

	responseBody := make([]*GetUserURLsResponse, 0, len(shortenedURLSInfos))

	//for _, info := range shortenedURLSInfos {
	//	responseBody = append(responseBody, NewGetUserURLsResponse(h.createShortURLFromID(info.GetID()), info.GetOriginalURL()))
	//}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

func (h *handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "text/plain", "text", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	reader, err := getBodyReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, bodyReadingError, http.StatusInternalServerError)
		return
	}
	urlString := string(body)
	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}

	var urlID string
	statusCode := 201
	ownerID := r.Context().Value(middleware.OwnerID).(string)
	urlID, err = h.service.AddNewURL(entity.NewUnstoredShortenedURLInfo(ownerID, urlString))
	if err != nil {
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := h.service.GetByOriginalURL(urlString)
			urlID = info.GetID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			statusCode = 409
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(h.createShortURLFromID(urlID)))
	if err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

func (h *handler) SaveURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "application/json", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	reader, err := getBodyReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	var body SaveURLBody

	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(body.URL)
	if err != nil {
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}
	ownerID := r.Context().Value(middleware.OwnerID).(string)
	statusCode := 201
	urlID, err := h.service.AddNewURL(entity.NewUnstoredShortenedURLInfo(ownerID, body.URL))
	if err != nil {
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := h.service.GetByOriginalURL(body.URL)
			urlID = info.GetID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			statusCode = 409
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&SaveURLResponse{Result: h.createShortURLFromID(urlID)}); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

type GetUserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewGetUserURLsResponse(shortURL, originalURL string) *GetUserURLsResponse {
	return &GetUserURLsResponse{ShortURL: shortURL, OriginalURL: originalURL}
}

func (h *handler) createShortURLFromID(shortURLID string) string {
	return fmt.Sprintf("%s/%s", h.address, shortURLID)
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
