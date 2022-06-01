package handler

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const contentTypeTextPlain = "text/plain; charset=utf-8"
const contentTypeTextHTML = "text/html; charset=utf-8"

var testURL1 = "https://riptutorial.com/go/example/2570/http-hello-world-with-custom-server-and-mux"
var contentTypeHeader = map[string]string{
	"Content-Type": "text/plain",
}

type MockURLRepository struct {
	Storage map[int]string
}

func (receiver *MockURLRepository) Save(url string) int {
	id := len(receiver.Storage)
	receiver.Storage[id] = url
	return id
}

func (receiver MockURLRepository) GetByID(id int) string {
	s := receiver.Storage[id]
	return s
}

type expected struct {
	contentType string
	statusCode  int
	body        string
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		request    string
		headers    map[string]string
		body       string
		storedUrls map[int]string
		expected   expected
	}{
		{
			name:       "Adding single URL",
			method:     http.MethodPost,
			request:    "/",
			headers:    contentTypeHeader,
			body:       testURL1,
			storedUrls: make(map[int]string),
			expected: expected{
				contentType: contentTypeTextPlain,
				statusCode:  201,
				body:        "",
			},
		},
		{
			name:       "Get single URL",
			method:     http.MethodGet,
			request:    "/0",
			headers:    contentTypeHeader,
			body:       "",
			storedUrls: map[int]string{0: "https://www.google.com/"},
			expected: expected{
				contentType: contentTypeTextHTML,
				statusCode:  307,
				body:        "<a href=\"https://www.google.com/\">Temporary Redirect</a>.\n\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var request *http.Request
			if test.method == http.MethodGet {
				request = httptest.NewRequest(test.method, test.request, nil)
			} else {
				request = httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
			}
			for k, v := range test.headers {
				request.Header.Set(k, v)
			}
			recorder := httptest.NewRecorder()
			h := NewHandler("localhost:8080", service.NewURLShortenerService(&MockURLRepository{test.storedUrls}))
			h.ServeHTTP(recorder, request)
			result := recorder.Result()

			assert.Equal(t, test.expected.statusCode, result.StatusCode)
			assert.Equal(t, test.expected.contentType, result.Header.Get("Content-Type"))

			if len(test.expected.body) > 0 {
				body, err := ioutil.ReadAll(result.Body)
				require.NoError(t, err)
				assert.Equal(t, test.expected.body, string(body))
				err = result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
