package http

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/apolsh/yapr-url-shortener/internal/logger"
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

var log = logger.LoggerOfComponent("http-router")

const (
	applicationJSON = "application/json; charset=utf-8"
)

const (
	realIPHeader = "X-Real-IP"
)

// Controller представляет собой http контроллер, умеющий обрабатывать запросы
type Controller struct {
	shortenService service.URLShortenerService
	trustedSubnet  *net.IPNet
}

// NewRouter конструктор для Controller
func NewRouter(r *chi.Mux, serviceImpl service.URLShortenerService, provider crypto.CryptographicProvider, trustedSubnet *net.IPNet) {
	c := &Controller{shortenService: serviceImpl, trustedSubnet: trustedSubnet}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(customMiddleware.AuthMiddleware(provider))

	r.Route("/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/ping", c.PingDB)
			r.Get("/{urlID}", c.GetShortenURLByID)
			r.Get("/api/user/urls", c.GetShortenURLsByUser)
			r.Post("/", c.SaveShortenURL)
		})
		r.With(customMiddleware.JSONFilterMiddleware).Group(func(r chi.Router) {
			r.Route("/api", func(r chi.Router) {
				r.Get("/internal/stats", c.GetAppStats)
				r.Post("/shorten/batch", c.SaveShortenURLsInBatch)
				r.Post("/shorten", c.SaveShortenURLJSON)
				r.Delete("/user/urls", c.DeleteShortenURLsInBatch)
			})
		})
	})

}

// PingDB проверяет работу хранилища URL
func (c *Controller) PingDB(w http.ResponseWriter, r *http.Request) {
	ok := c.shortenService.PingDB(r.Context())
	if ok {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetShortenURLByID производит редирект на сохраненный ранее в хранилище URL
func (c *Controller) GetShortenURLByID(w http.ResponseWriter, r *http.Request) {
	if urlID := chi.URLParam(r, "urlID"); urlID != "" {
		foundURL, err := c.shortenService.GetURLByID(r.Context(), urlID)
		if err != nil {
			log.Error(err)
			if errors.Is(repository.ErrorItemNotFound, err) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(service.ErrorItemIsDeleted, err) {
				http.Error(w, "", http.StatusGone)
				return
			}
		}
		http.Redirect(w, r, foundURL, http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "Invalid parameter", http.StatusMethodNotAllowed)
}

// GetShortenURLsByUser возвращает список пар (короткий + длинный) URL пользователя
func (c *Controller) GetShortenURLsByUser(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	urlPairs, err := c.shortenService.GetURLsByOwnerID(r.Context(), ownerID)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(urlPairs) == 0 {
		w.WriteHeader(204)
		return
	}

	setContentType(w, applicationJSON)
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(urlPairs); err != nil {
		log.Error(err)
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

// SaveShortenURL принимает запрос в виде простого текста, сохраняет URL в хранилище
func (c *Controller) SaveShortenURL(w http.ResponseWriter, r *http.Request) {
	if !isValidContentType(r, "text/plain", "text", "application/x-gzip") {
		http.Error(w, invalidContentTypeError, http.StatusBadRequest)
		return
	}

	urlString, err := extractTextBody(r)
	if err != nil {
		log.Error(err)
		http.Error(w, bodyReadingError, http.StatusInternalServerError)
		return
	}
	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		log.Error(err)
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}

	var urlID string
	statusCode := 201
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	urlID, err = c.shortenService.AddNewURL(r.Context(), *entity.NewUnstoredShortenedURLInfo(ownerID, urlString))
	if err != nil {
		log.Error(err)
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := c.shortenService.GetByOriginalURL(r.Context(), urlString)
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
	_, err = w.Write([]byte(c.shortenService.GetShortenURLFromID(r.Context(), urlID)))
	if err != nil {
		log.Error(err)
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

// SaveShortenURLsInBatch сохраняет сразу несколько URL в хранилище за один запрос
func (c *Controller) SaveShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	var body []dto.ShortenInBatchRequestItem

	err := extractJSONBody(r, &body)
	if err != nil {
		log.Error(err)
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	batch, err := c.shortenService.AddNewURLsInBatch(r.Context(), ownerID, body)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setContentType(w, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(batch); err != nil {
		log.Error(err)
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
	}
}

// SaveShortenURLJSON принимает запрос в виде JSON, сохраняет URL в хранилище
func (c *Controller) SaveShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	var body saveURLBody
	err := extractJSONBody(r, &body)
	if err != nil {
		log.Error(err)
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(body.URL)
	if err != nil {
		log.Error(err)
		http.Error(w, parseURLError, http.StatusBadRequest)
		return
	}
	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	statusCode := 201
	info := *entity.NewUnstoredShortenedURLInfo(ownerID, body.URL)

	urlID, err := c.shortenService.AddNewURL(r.Context(), info)
	if err != nil {
		if errors.Is(err, repository.ErrorURLAlreadyStored) {
			info, err := c.shortenService.GetByOriginalURL(r.Context(), body.URL)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			urlID = info.GetID()
			statusCode = 409
		}
	}
	setContentType(w, applicationJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&saveURLResponse{Result: c.shortenService.GetShortenURLFromID(r.Context(), urlID)}); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
}

// DeleteShortenURLsInBatch помечает URL в хранилище как удаленный
func (c *Controller) DeleteShortenURLsInBatch(w http.ResponseWriter, r *http.Request) {
	var ids []string
	err := extractJSONBody(r, &ids)
	if err != nil {
		log.Error(err)
		http.Error(w, decodeRequestBodyError, http.StatusBadRequest)
		return
	}

	ownerID := r.Context().Value(customMiddleware.OwnerID).(string)
	err = c.shortenService.DeleteURLsInBatch(r.Context(), ownerID, ids)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// GetAppStats получить статистику приложения
func (c *Controller) GetAppStats(w http.ResponseWriter, r *http.Request) {
	stringIP := r.Header.Get(realIPHeader)

	ip := net.ParseIP(stringIP)
	if ip == nil {
		log.Error(fmt.Errorf("unable to parse %s to IP", stringIP))
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if !c.trustedSubnet.Contains(ip) {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	statistic, err := c.shortenService.GetAppStatistic(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	setContentType(w, applicationJSON)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statistic); err != nil {
		http.Error(w, encodeResponseBodyError, http.StatusInternalServerError)
	}
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
			log.Error(err)
		}
	}(reader)
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func extractTextBody(r *http.Request) (string, error) {
	var reader io.ReadCloser
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Error(err)
			return "", err
		}
		reader = gz
	} else {
		reader = r.Body
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Error(err)
		}
	}(reader)
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(body), nil
}

type saveURLBody struct {
	URL string `json:"url"`
}

type saveURLResponse struct {
	Result string `json:"result"`
}
