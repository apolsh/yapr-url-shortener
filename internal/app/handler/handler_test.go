package handler

import (
	"encoding/json"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/config"
	"github.com/apolsh/yapr-url-shortener/internal/app/mock"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var saveURLHeaders = map[string]string{
	"Content-Type": "text/plain",
}
var saveURLJSONHeaders = map[string]string{
	"Content-Type": "application/json",
}
var testURL1 = "https://riptutorial.com/go/example/2570/http-hello-world-with-custom-server-and-mux"
var testURL2 = "https://bitfieldconsulting.com/golang/map-declaring-initializing"

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

func executeSaveURLRequest(t *testing.T, server *httptest.Server, requestBody string, headers map[string]string) (*http.Response, string) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodPost, server.URL, strings.NewReader(requestBody))

	require.NoError(t, err)
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response, string(body)
}

func executeSaveURLJSONRequest(t *testing.T, server *httptest.Server, requestBody string, headers map[string]string) (*http.Response, string) {
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodPost, server.URL+"/api/shorten", strings.NewReader(requestBody))

	require.NoError(t, err)
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response, string(body)
}

var cfg config.Config

func getConfig() config.Config {
	if cfg == (config.Config{}) {
		cfg = config.Load()
	}
	return cfg
}

func TestHandler_GetURLHandler(t *testing.T) {
	alreadyStoredURLs := make(map[int]string)
	alreadyStoredURLs[0] = testURL1
	alreadyStoredURLs[1] = testURL2
	cfg := getConfig()

	repositoryMock := mock.NewURLRepositoryMock(alreadyStoredURLs)
	h := NewURLShortenerHandler(cfg.BaseURL, service.NewURLShortenerService(repositoryMock))
	r := chi.NewRouter()
	h.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	t.Run("Get existing URL", func(t *testing.T) {
		response, _ := executeGetURLRequest(t, server, "/0")
		err := response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, testURL1, response.Request.URL.String())

		t.Run("Extract not exist id", func(t *testing.T) {
			response, body := executeGetURLRequest(t, server, "/100")
			err := response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, 404, response.StatusCode)
			assert.Equal(t, "404 page not found\n", body)
		})

		t.Run("Extract empty id", func(t *testing.T) {
			response, _ := executeGetURLRequest(t, server, "/")
			err := response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, 405, response.StatusCode)
		})

		t.Run("Extract wrong id type", func(t *testing.T) {
			response, body := executeGetURLRequest(t, server, "/abc")
			err := response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, 400, response.StatusCode)
			assert.Equal(t, "Invalid parameter\n", body)
		})
	})
}

func TestHandler_SaveURLHandler(t *testing.T) {
	cfg := getConfig()
	h := NewURLShortenerHandler(cfg.BaseURL, service.NewURLShortenerService(mock.NewURLRepositoryMock(make(map[int]string))))
	r := chi.NewRouter()
	h.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	t.Run("Adding and extracting multiple URL's", func(t *testing.T) {
		response, stringBody := executeSaveURLRequest(t, server, testURL1, saveURLHeaders)
		err := response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, fmt.Sprintf("%s/0", cfg.BaseURL), stringBody)

		response, stringBody = executeSaveURLRequest(t, server, testURL2, saveURLHeaders)
		err = response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, fmt.Sprintf("%s/1", cfg.BaseURL), stringBody)
	})

	t.Run("Add with invalid content-type", func(t *testing.T) {
		response, body := executeSaveURLRequest(t, server, testURL2, map[string]string{"Content-Type": "multipart/form-data"})
		err := response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Invalid Content-Type: multipart/form-data\n", body)
	})
}

func TestHandler_SaveURLJSONHandler(t *testing.T) {
	cfg := getConfig()
	h := NewURLShortenerHandler(cfg.BaseURL, service.NewURLShortenerService(mock.NewURLRepositoryMock(make(map[int]string))))
	r := chi.NewRouter()
	h.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	t.Run("Adding and extracting multiple URL's", func(t *testing.T) {

		marshal, err := json.Marshal(&SaveURLBody{URL: testURL1})
		require.NoError(t, err)

		response, stringBody := executeSaveURLJSONRequest(t, server, string(marshal), saveURLJSONHeaders)
		err = response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 201, response.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"result":"%s/0"}`, cfg.BaseURL), stringBody)

		marshal, err = json.Marshal(&SaveURLBody{URL: testURL2})
		require.NoError(t, err)

		response, stringBody = executeSaveURLJSONRequest(t, server, string(marshal), saveURLJSONHeaders)
		err = response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 201, response.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"result":"%s/1"}`, cfg.BaseURL), stringBody)
	})

	//t.Run("Add with invalid content-type", func(t *testing.T) {
	//	response, body := executeSaveURLRequest(t, server, testURL2, map[string]string{"Content-Type": "application/json"})
	//	err := response.Body.Close()
	//	require.NoError(t, err)
	//	assert.Equal(t, 400, response.StatusCode)
	//	assert.Equal(t, "Invalid Content-Type\n", body)
	//})
}

func TestHandler_CommonTests(t *testing.T) {
	cfg := getConfig()
	h := NewURLShortenerHandler(cfg.BaseURL, service.NewURLShortenerService(mock.NewURLRepositoryMock(make(map[int]string))))
	r := chi.NewRouter()
	h.Register(r)

	server := httptest.NewServer(r)
	defer server.Close()

	t.Run("Send wrong http Method", func(t *testing.T) {
		var request *http.Request
		var err error

		request, err = http.NewRequest(http.MethodPut, server.URL, strings.NewReader(testURL1))

		require.NoError(t, err)
		for k, v := range saveURLHeaders {
			request.Header.Set(k, v)
		}

		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)

		err = response.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 405, response.StatusCode)
	})
}
