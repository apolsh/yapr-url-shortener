package handler

import (
	"encoding/json"
	"github.com/apolsh/yapr-url-shortener/internal/app/mocks"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const dummyURL1 = "https://google.com"
const dummyURL2 = "https://ya.ru"

var dummyShortenedURLInfo1 = entity.NewShortenedURLInfo("0", "owner", dummyURL1)
var dummyShortenedURLInfo2 = entity.NewShortenedURLInfo("1", "owner", dummyURL2)

type HandlerTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	server         *httptest.Server
	urlServiceMock *mocks.MockURLShortenerService
}

func (suite *HandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.ctrl = ctrl
}

func (suite *HandlerTestSuite) BeforeTest(suiteName, testName string) {
	s := mocks.NewMockURLShortenerService(suite.ctrl)
	p := mocks.NewMockCryptographicProvider(suite.ctrl)
	handler := NewURLShortenerHandler("http://localhost:8080", s, p)
	r := chi.NewRouter()
	handler.Register(r)
	p.EXPECT().Encrypt(gomock.Any()).Return("encrypted_value")

	suite.urlServiceMock = s
	suite.server = httptest.NewServer(r)
}

func (suite *HandlerTestSuite) TestSaveURLHandlerValidRequest() {
	id := "id1"
	suite.urlServiceMock.EXPECT().AddNewURL(gomock.Any()).Return(id, nil)

	response, shortenedURL := executeSaveURL(suite.T(), suite.server, dummyURL1)
	defer response.Body.Close()

	assert.Equal(suite.T(), 201, response.StatusCode)
	assert.Equal(suite.T(), "http://localhost:8080/"+id, shortenedURL)
}

func (suite *HandlerTestSuite) TestSaveURLHandlerEmptyBody() {
	response, _ := executeSaveURL(suite.T(), suite.server, "")
	defer response.Body.Close()

	assert.Equal(suite.T(), 400, response.StatusCode)
}

func (suite *HandlerTestSuite) TestSaveURLJSONHandlerValidRequest() {
	id := "id1"
	suite.urlServiceMock.EXPECT().AddNewURL(gomock.Any()).Return(id, nil)

	response, jsonResult, _ := executeSaveURLJSON(suite.T(), suite.server, dummyURL1)
	defer response.Body.Close()

	assert.Equal(suite.T(), 201, response.StatusCode)
	assert.Equal(suite.T(), "http://localhost:8080/"+id, jsonResult.Result)
}

func (suite *HandlerTestSuite) TestSaveURLJSONHandlerEmptyBody() {
	response, _, errorMessage := executeSaveURLJSON(suite.T(), suite.server, "")
	defer response.Body.Close()

	assert.Equal(suite.T(), 400, response.StatusCode)
	assert.Equal(suite.T(), "Passed value is not valid URL\n", errorMessage)
}

func (suite *HandlerTestSuite) TestGetUserURLsHandlerWithExistingURLs() {
	suite.urlServiceMock.EXPECT().GetURLsByOwnerID(gomock.Any()).Return([]*entity.ShortenedURLInfo{dummyShortenedURLInfo1, dummyShortenedURLInfo2}, nil)

	response, urLsResponses := executeGetUserURLsRequest(suite.T(), suite.server)

	redirectResponse := getRedirectResponse(response)
	defer redirectResponse.Body.Close()

	assert.Equal(suite.T(), 200, response.StatusCode)
	assert.Equal(suite.T(), 2, len(urLsResponses))
}

func (suite *HandlerTestSuite) TestGetUserURLsHandlerNotFoundURLs() {
	suite.urlServiceMock.EXPECT().GetURLsByOwnerID(gomock.Any()).Return([]*entity.ShortenedURLInfo{}, nil)

	response, urLsResponses := executeGetUserURLsRequest(suite.T(), suite.server)

	redirectResponse := getRedirectResponse(response)
	defer redirectResponse.Body.Close()

	assert.Equal(suite.T(), 204, response.StatusCode)
	assert.Equal(suite.T(), 0, len(urLsResponses))
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithExistingURL() {
	suite.urlServiceMock.EXPECT().GetURLByID(gomock.Any()).Return(dummyURL1, nil)

	response, _ := executeGetURLRequest(suite.T(), suite.server, "/0")
	redirectResponse := getRedirectResponse(response)
	defer redirectResponse.Body.Close()

	assert.Equal(suite.T(), http.StatusTemporaryRedirect, redirectResponse.StatusCode)
	assert.Equal(suite.T(), dummyURL1, redirectResponse.Header.Get("Location"))
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithURLNotExist() {
	suite.urlServiceMock.EXPECT().GetURLByID(gomock.Any()).Return("", repository.ErrorItemNotFound)

	response, _ := executeGetURLRequest(suite.T(), suite.server, "/0")
	defer response.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, response.StatusCode)
}

func (suite *HandlerTestSuite) TestGetURLHandlerWithoutID() {
	response, _ := executeGetURLRequest(suite.T(), suite.server, "/")
	defer response.Body.Close()

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

func executeSaveURLJSON(t *testing.T, server *httptest.Server, urlToSave string) (*http.Response, *SaveURLResponse, string) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodPost, server.URL+"/api/shorten", strings.NewReader(`{"url":"`+urlToSave+`"}`))
	require.NoError(t, err)
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	if response.StatusCode == 201 {
		var parsedBody SaveURLResponse
		err = json.Unmarshal(body, &parsedBody)
		require.NoError(t, err)

		return response, &parsedBody, ""
	}

	return response, nil, string(body)
}

func executeSaveURL(t *testing.T, server *httptest.Server, urlToSave string) (*http.Response, string) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodPost, server.URL+"/", strings.NewReader(urlToSave))
	require.NoError(t, err)
	request.Header.Add("Content-Type", "text")

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response, string(body)
}

func executeGetUserURLsRequest(t *testing.T, server *httptest.Server) (*http.Response, []GetUserURLsResponse) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodGet, server.URL+"/api/user/urls", nil)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	var body []GetUserURLsResponse

	if response.ContentLength != 0 {
		err = json.NewDecoder(response.Body).Decode(&body)
		require.NoError(t, err)
	}

	defer response.Body.Close()
	return response, body
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
