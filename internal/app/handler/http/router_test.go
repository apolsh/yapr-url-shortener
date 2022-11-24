package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/apolsh/yapr-url-shortener/internal/app/mocks"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RouterSuite struct {
	suite.Suite
	shorts  *mocks.MockURLShortenerService
	handler http.Handler
	ctrl    *gomock.Controller
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))

}

var (
	cryptoProvider = crypto.NewAESCryptoProvider("secret")
	shortURL1      = "http://shorturl1.com/123"
	shortURL2      = "http://shorturl2.com/456"
	longURL1       = "http://longurl1.com"
	longURL2       = "http://longurl2.com"
)

func (s *RouterSuite) SetupTest() {
	r := chi.NewRouter()
	ctrl := gomock.NewController(s.T())
	s.ctrl = ctrl
	s.shorts = mocks.NewMockURLShortenerService(ctrl)

	NewRouter(r, s.shorts, cryptoProvider)
	s.handler = r
}

func (s *RouterSuite) TestGetShortenURLByIDWithSuccess() {
	s.shorts.EXPECT().GetURLByID("some_id").Return("http://rediercted.com/url", nil)
	req := httptest.NewRequest(http.MethodGet, "/some_id", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 307, resp.Code)
	assert.Equal(s.T(), "http://rediercted.com/url", resp.Header().Get("Location"))
}

func (s *RouterSuite) TestGetShortenURLByIDNotFound() {
	s.shorts.EXPECT().GetURLByID("some_id").Return("", repository.ErrorItemNotFound)
	req := httptest.NewRequest(http.MethodGet, "/some_id", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 404, resp.Code)
}

func (s *RouterSuite) TestGetShortenURLByIDItemDeleted() {
	s.shorts.EXPECT().GetURLByID("some_id").Return("", service.ErrorItemIsDeleted)
	req := httptest.NewRequest(http.MethodGet, "/some_id", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 410, resp.Code)
}

func (s *RouterSuite) TestGetShortenURLByIDNoID() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 405, resp.Code)
}

func (s *RouterSuite) TestPingWithOk() {
	s.shorts.EXPECT().PingDB().Return(true)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 200, resp.Code)
}

func (s *RouterSuite) TestPingNotOk() {
	s.shorts.EXPECT().PingDB().Return(false)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 500, resp.Code)
}

func (s *RouterSuite) TestGetShortenURLsByUserSomeFound() {
	s.shorts.EXPECT().GetURLsByOwnerID(gomock.Any()).Return([]dto.URLPair{
		{ShortURL: shortURL1, OriginalURL: longURL1},
		{ShortURL: shortURL2, OriginalURL: longURL2},
	}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 200, resp.Code)
	s.JSONEq(`[{"short_url":"http://shorturl1.com/123","original_url":"http://longurl1.com"},{"short_url":"http://shorturl2.com/456","original_url":"http://longurl2.com"}]`, string(body))

}

func (s *RouterSuite) TestGetShortenURLsByUserZeroResult() {
	s.shorts.EXPECT().GetURLsByOwnerID(gomock.Any()).Return([]dto.URLPair{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 204, resp.Code)
}

func (s *RouterSuite) TestGetShortenURLsByUserSomeError() {
	s.shorts.EXPECT().GetURLsByOwnerID(gomock.Any()).Return(nil, service.ErrorItemIsDeleted)
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 500, resp.Code)
}

func (s *RouterSuite) TestSaveShortenURLWrongContentType() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL1))
	req.Header.Set("Content-Type", "SOMERANDOMTYPE")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 400, resp.Code)
}

func (s *RouterSuite) TestSaveShortenURLWrongURLType() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set("Content-Type", "text/plain")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 400, resp.Code)
}

func (s *RouterSuite) TestSaveShortenURLNewURLSaved() {
	s.shorts.EXPECT().AddNewURL(gomock.Any()).Return("123", nil)
	s.shorts.EXPECT().GetShortenURLFromID("123").Return(shortURL1)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL1))
	req.Header.Set("Content-Type", "text/plain")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 201, resp.Code)
	assert.Equal(s.T(), shortURL1, string(body))
}

func (s *RouterSuite) TestSaveShortenURLAlreadySaved() {
	s.shorts.EXPECT().AddNewURL(gomock.Any()).Return("", repository.ErrorURLAlreadyStored)
	s.shorts.EXPECT().GetByOriginalURL(longURL1).Return(entity.ShortenedURLInfo{ID: "123"}, nil)
	s.shorts.EXPECT().GetShortenURLFromID("123").Return(shortURL1)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL1))
	req.Header.Set("Content-Type", "text/plain")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 409, resp.Code)
	assert.Equal(s.T(), shortURL1, string(body))
}

func (s *RouterSuite) TestSaveShortenURLsInBatchWrongType() {
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(longURL1))
	req.Header.Set("Content-Type", "SOMERANDOMTYPE")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 400, resp.Code)
}

func (s *RouterSuite) TestSaveShortenURLsInBatchWithSuccess() {

	s.shorts.EXPECT().AddNewURLsInBatch(gomock.Any(), gomock.Any()).Return([]dto.ShortenInBatchResponseItem{
		{CorrelationID: "1", ShortURL: shortURL1},
		{CorrelationID: "2", ShortURL: shortURL2},
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(
		`[
				{
					"correlation_id": "1",
					"original_url": "https://www.youtube.com/"
				},
				{
					"correlation_id": "2",
					"original_url": "https://www.twitch.tv/kochevnik"
				}
			] 
	`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 201, resp.Code)
	s.JSONEq(`[{"correlation_id":"1","short_url":"http://shorturl1.com/123"},{"correlation_id":"2","short_url":"http://shorturl2.com/456"}]`, string(body))
}

func (s *RouterSuite) TestSaveShortenURLJSONWithSuccess() {
	s.shorts.EXPECT().AddNewURL(gomock.Any()).Return("123", nil)
	s.shorts.EXPECT().GetShortenURLFromID("123").Return(shortURL1)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://golangdocs.com/golang-read-json-file"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 201, resp.Code)
	s.JSONEq(`{"result":"http://shorturl1.com/123"}`, string(body))

}

func (s *RouterSuite) TestSaveShortenURLJSONInvalidURL() {
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"blablabla"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 400, resp.Code)
}

func (s *RouterSuite) TestSaveShortenURLJSONErrorURLAlreadyStored() {
	s.shorts.EXPECT().AddNewURL(gomock.Any()).Return("", repository.ErrorURLAlreadyStored)
	s.shorts.EXPECT().GetByOriginalURL(longURL1).Return(entity.ShortenedURLInfo{ID: "123"}, nil)
	s.shorts.EXPECT().GetShortenURLFromID("123").Return(shortURL1)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"`+longURL1+`"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(s.T(), 409, resp.Code)
	s.JSONEq(`{"result":"http://shorturl1.com/123"}`, string(body))
}

func (s *RouterSuite) TestDeleteShortenURLsInBatchWithSuccess() {
	s.shorts.EXPECT().DeleteURLsInBatch(gomock.Any(), gomock.Any()).Return(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["123","456"]`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 202, resp.Code)
}

func (s *RouterSuite) TestDeleteShortenURLsInBatchWithError() {
	s.shorts.EXPECT().DeleteURLsInBatch(gomock.Any(), gomock.Any()).Return(service.ErrorItemIsDeleted)
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["123","456"]`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.handler.ServeHTTP(resp, req)

	assert.Equal(s.T(), 500, resp.Code)
}
