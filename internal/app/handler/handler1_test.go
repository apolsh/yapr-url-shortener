package handler

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/mock"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
)

//GetURLHandler
// with id -> not found
// with id -> found
// without id -> bad request

//GetUserURLsHandler
//found some
//500 error

//SaveURLHandler
// wrong contentType
// not URL
// empty body
// right url - right response

//SaveURLJSONHandler
// wrong contentType
// wrong structure
// not JSON
// not URL in JSON
// right structure - right response

type HandlerTestSuite struct {
	suite.Suite
	ctrl           gomock.Controller
	server         *httptest.Server
	urlServiceMock service.URLShortenerService
}

func (suite *HandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	s := mock.NewMockURLShortenerService(ctrl)
	p := mock.NewMockProvider(ctrl)
	handler := NewURLShortenerHandler("http://localhost:8080", s, p)
	r := chi.NewRouter()
	handler.Register(r)

	suite.server = httptest.NewServer(r)
}
