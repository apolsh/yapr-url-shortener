package test

import (
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/handler"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var contentTypeHeader = map[string]string{
	"Content-Type": "text/plain",
}
var emptyHeaders = make(map[string]string)
var testURL1 = "https://riptutorial.com/go/example/2570/http-hello-world-with-custom-server-and-mux"
var testURL2 = "https://bitfieldconsulting.com/golang/map-declaring-initializing"

func executeTestRequest(t *testing.T, server *httptest.Server, method, path, requestBody string, headers map[string]string) (*http.Response, string) {
	var request *http.Request
	var err error
	if method == http.MethodGet {
		request, err = http.NewRequest(method, server.URL+path, nil)
	} else {
		request, err = http.NewRequest(method, server.URL+path, strings.NewReader(requestBody))
	}

	require.NoError(t, err)
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	err = response.Body.Close()
	require.NoError(t, err)

	return response, string(body)
}

func TestIntegration(t *testing.T) {
	addr := "localhost:8080"
	h := handler.NewHandler(addr, service.NewURLShortenerService(repository.NewURLRepositoryInMemoryImpl()))
	server := httptest.NewServer(h)
	defer server.Close()

	t.Run("Adding and extracting multiple URL's", func(t *testing.T) {
		response, stringBody := executeTestRequest(t, server, http.MethodPost, "", testURL1, contentTypeHeader)
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, fmt.Sprintf("http://%s/0", addr), stringBody)

		response, stringBody = executeTestRequest(t, server, http.MethodPost, "", testURL2, contentTypeHeader)
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, fmt.Sprintf("http://%s/1", addr), stringBody)

		response, stringBody = executeTestRequest(t, server, http.MethodGet, "/0", "", emptyHeaders)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, testURL1, response.Request.URL.String())

		response, stringBody = executeTestRequest(t, server, http.MethodGet, "/1", "", emptyHeaders)
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, testURL2, response.Request.URL.String())
	})

	t.Run("Extract wrong id", func(t *testing.T) {
		response, body := executeTestRequest(t, server, http.MethodGet, "/100", "", emptyHeaders)
		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, "404 page not found\n", body)
	})

	t.Run("Extract empty id", func(t *testing.T) {
		response, body := executeTestRequest(t, server, http.MethodGet, "/", "", emptyHeaders)
		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Invalid parameter\n", body)
	})

	t.Run("Extract wrong id type", func(t *testing.T) {
		response, body := executeTestRequest(t, server, http.MethodGet, "/abc", "", emptyHeaders)
		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Invalid parameter\n", body)
	})

	t.Run("Add with invalid content-type", func(t *testing.T) {
		response, body := executeTestRequest(t, server, http.MethodPost, "", testURL2, map[string]string{"Content-Type": "application/json"})
		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Invalid Content-Type\n", body)
	})

	t.Run("Send wrong http Method", func(t *testing.T) {
		response, body := executeTestRequest(t, server, http.MethodPut, "", testURL2, map[string]string{"Content-Type": "application/json"})
		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Invalid request\n", body)
	})

}
