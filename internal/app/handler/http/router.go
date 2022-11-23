package http

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	customMiddleware "github.com/apolsh/yapr-url-shortener/internal/app/middleware"
)

const (
	decodeRequestBodyError  = "failed to decode body"
	encodeResponseBodyError = "failed to encode body response"
	invalidContentTypeError = "invalid content type"
	bodyReadingError        = "Error while body reading"
	parseURLError           = "Passed value is not valid URL"
)

const (
	applicationJSON = "application/json; charset=utf-8"
)

type controller struct {
	shortenService service.URLShortenerService
}

func NewRouter(r *chi.Mux, serviceImpl service.URLShortenerService, provider crypto.CryptographicProvider) {
	c := &controller{shortenService: serviceImpl}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.CompressResponse)
	r.Use(customMiddleware.AuthMiddleware(provider))

	r.Route("/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/ping", c.PingDB)
			r.Get("/{urlID}", c.GetShortenURLByID)
			r.Get("/api/user/urls", c.GetShortenURLsByUser)
			r.Post("/", c.SaveShortenURL)
		})
		r.With(customMiddleware.JsonFilterMiddleware).Group(func(r chi.Router) {
			r.Route("/api", func(r chi.Router) {
				r.Post("/shorten/batch", c.SaveShortenURLsInBatch)
				r.Post("/shorten", c.SaveShortenURLJSON)
				r.Delete("/user/urls", c.DeleteShortenURLsInBatch)
			})
		})
	})

}

func (c *controller) PingDB(w http.ResponseWriter, r *http.Request) {
	ok := c.shortenService.PingDB()
	if ok {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *controller) GetShortenURLByID(w http.ResponseWriter, r *http.Request) {
	if urlID := chi.URLParam(r, "urlID"); urlID != "" {
		foundURL, err := c.shortenService.GetURLByID(urlID)
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

func (c *controller) GetShortenURLsByUser(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	urlPairs, err := c.shortenService.GetURLsByOwnerID(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(urlPairs) == 0 {
		w.WriteHeader(204)
		return
	}

	setContentType(w, applicationJSON)
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(urlPairs); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

func (c *controller) SaveShortenURL(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "text/plain", "text", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	urlString, err := extractTextBody(r)
	if err != nil {
		http.Error(w, bodyReadingError, http.StatusInternalServerError)
		return
	}
	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}

	var urlID string
	statusCode := 201
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	urlID, err = c.shortenService.AddNewURL(*entity.NewUnstoredShortenedURLInfo(ownerID, urlString))
	if err != nil {
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := c.shortenService.GetByOriginalURL(urlString)
			urlID = info.GetID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			statusCode = 409
		}
	}
	setContentType(w, "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(c.shortenService.GetShortenURLFromId(urlID)))
	if err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

func (c *controller) SaveShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	var body []dto.ShortenInBatchRequestItem

	err := extractJSONBody(r, &body)
	if err != nil {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	batch, err := c.shortenService.AddNewURLsInBatch(ownerID, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setContentType(w, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(batch); err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
	}
}

func (c *controller) SaveShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	var body SaveURLBody
	err := extractJSONBody(r, &body)
	if err != nil {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(body.URL)
	if err != nil {
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	statusCode := 201
	urlID, err := c.shortenService.AddNewURL(*entity.NewUnstoredShortenedURLInfo(ownerID, body.URL))
	if err != nil {
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := c.shortenService.GetByOriginalURL(body.URL)
			urlID = info.GetID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			statusCode = 409
		}
	}
	setContentType(w, applicationJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&SaveURLResponse{Result: c.shortenService.GetShortenURLFromId(urlID)}); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

func (c *controller) DeleteShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	var ids []string
	err := extractJSONBody(r, &ids)
	if err != nil {
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	err = c.shortenService.DeleteURLsInBatch(ownerID, ids)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func setContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
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

func extractJSONBody(r *http.Request, v interface{}) error {
	var reader io.ReadCloser
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		reader = gz
	}
	reader = r.Body
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			//log.Err(err).Msg(err.Error())
		}
	}(reader)
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return err
	}
	return nil
}

func extractTextBody(r *http.Request) (string, error) {
	var reader io.ReadCloser
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return "", err
		}
		reader = gz
	}
	reader = r.Body
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			//log.Err(err).Msg(err.Error())
		}
	}(reader)
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type SaveURLBody struct {
	URL string `json:"url"`
}

type SaveURLResponse struct {
	Result string `json:"result"`
}
