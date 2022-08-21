package handler

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/mock"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

const dummyURL1 = "https://google.com"

type HandlerTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	server         *httptest.Server
	urlServiceMock *mock.MockURLShortenerService
}

func (suite *HandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.ctrl = ctrl
}

func (suite *HandlerTestSuite) BeforeTest(suiteName, testName string) {
	s := mock.NewMockURLShortenerService(suite.ctrl)
	p := mock.NewMockProvider(suite.ctrl)
	handler := NewURLShortenerHandler("http://localhost:8080", s, p)
	r := chi.NewRouter()
	handler.Register(r)
	p.EXPECT().Encrypt(gomock.Any()).Return("encrypted_value")

	suite.urlServiceMock = s
	suite.server = httptest.NewServer(r)
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithExistingURL() {
	suite.urlServiceMock.EXPECT().GetURLByID(gomock.Any()).Return(dummyURL1, nil)

	response, _ := executeGetURLRequest(suite.T(), suite.server, "/0")
	redirectResponse := getRedirectResponse(response)

	assert.Equal(suite.T(), http.StatusTemporaryRedirect, redirectResponse.StatusCode)
	assert.Equal(suite.T(), dummyURL1, redirectResponse.Header.Get("Location"))
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithURLNotExist() {
	suite.urlServiceMock.EXPECT().GetURLByID(gomock.Any()).Return("", repository.ItemNotFoundError)

	response, _ := executeGetURLRequest(suite.T(), suite.server, "/0")

	assert.Equal(suite.T(), http.StatusNotFound, response.StatusCode)
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithoutID() {
	response, _ := executeGetURLRequest(suite.T(), suite.server, "/")

	assert.Equal(suite.T(), http.StatusMethodNotAllowed, response.StatusCode)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func executeGetURLRequest(t *testing.T, server *httptest.Server, path string) (*http.Response, string) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodGet, server.URL+path, nil)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response, string(body)
}

func getRedirectResponse(response *http.Response) *http.Response {
	tempResponse := response
	var preTempResponse *http.Response

	for tempResponse != nil {
		preTempResponse = tempResponse
		tempResponse = tempResponse.Request.Response
	}
	return preTempResponse
}
