package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id1            = "123"
	id2            = "456"
	baseURL        = "http://localhost:8080/"
	databaseDSN    = "postgresql://yaprurlshortener:yaPR_ttuss@localhost:5432/yapr-url-shortener"
	longURL1       = "https://www.youtube.com/watch?v=c_l5Bz9hmFA&t=1212s"
	longURL2       = "https://github.com/go-resty/resty"
	authCookieName = "sessionId"
)

var (
	cryptoProvider = crypto.NewAESCryptoProvider("very_secret_key")
	ownerID        = uuid.New()
)

type IntegrationSuite struct {
	suite.Suite
	db     *pgxpool.Pool
	client *resty.Client
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
	suite.Run(t, new(IntegrationSuite))
	suite.Run(t, new(IntegrationSuite))
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	s.client = resty.New()
	ctx := context.Background()
	connect, err := pgxpool.Connect(ctx, databaseDSN)
	if err != nil {
		fmt.Println(err.Error())
	}
	s.db = connect
}

func (s *IntegrationSuite) SetupTest() {
	_, _ = s.db.Exec(context.Background(), "TRUNCATE \"shortened_urls\"")
}

func (s *IntegrationSuite) TestSaveShortenURL() {
	res, _ := s.client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(longURL1).
		Post(baseURL)

	assert.Equal(s.T(), 201, res.StatusCode())
	assert.True(s.T(), strings.Contains(string(res.Body()), baseURL))
}

func (s *IntegrationSuite) TestGetShortenURLByID() {
	//language=postgresql
	_, _ = s.db.Exec(context.Background(), `INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)`, id1, longURL1, ownerID, 0)
	res, _ := s.client.R().Get(baseURL + id1)
	originalResponse := res.RawResponse.Request.Response
	assert.Equal(s.T(), 307, originalResponse.StatusCode)
	assert.Equal(s.T(), longURL1, originalResponse.Header.Get("Location"))
}

func (s *IntegrationSuite) TestGetShortenURLsByUser() {
	sessionIDValue := cryptoProvider.Encrypt(ownerID[:])
	//language=postgresql
	_, _ = s.db.Exec(context.Background(), `INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)`, id1, longURL1, ownerID, 0)
	//language=postgresql
	_, _ = s.db.Exec(context.Background(), `INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)`, id2, longURL2, ownerID, 0)
	res, _ := s.client.R().SetCookie(&http.Cookie{Name: authCookieName, Value: sessionIDValue}).Get(baseURL + "api/user/urls")

	assert.Equal(s.T(), 200, res.StatusCode())
	assert.JSONEq(s.T(), `[{"short_url":"http://localhost:8080/123","original_url":"https://www.youtube.com/watch?v=c_l5Bz9hmFA\u0026t=1212s"},{"short_url":"http://localhost:8080/456","original_url":"https://github.com/go-resty/resty"}]`, string(res.Body()))
}

func (s *IntegrationSuite) TestSaveShortenURLsInBatch() {
	res, _ := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`[
				{
					"correlation_id": "1",
					"original_url": "https://www.youtube.com/watch?v=c_l5Bz9hmFA\u0026t=1212s"
				},
				{
					"correlation_id": "2",
					"original_url": "https://github.com/go-resty/resty"
				}
			] 
	`).Post(baseURL + "api/shorten/batch")

	assert.Equal(s.T(), 201, res.StatusCode())
}

func (s *IntegrationSuite) TestSaveShortenURLJSON() {
	res, _ := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"url":"https://golangdocs.com/golang-read-json-file"}`).
		Post(baseURL + "api/shorten")

	assert.Equal(s.T(), 201, res.StatusCode())
}

func (s *IntegrationSuite) TestDeleteShortenURLsInBatch() {
	sessionIDValue := cryptoProvider.Encrypt(ownerID[:])
	//language=postgresql
	_, _ = s.db.Exec(context.Background(), `INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)`, id1, longURL1, ownerID, 0)
	//language=postgresql
	_, _ = s.db.Exec(context.Background(), `INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)`, id2, longURL2, ownerID, 0)
	res, _ := s.client.R().
		SetCookie(&http.Cookie{Name: authCookieName, Value: sessionIDValue}).
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`["%s","%s"]`, id1, id2)).
		Delete(baseURL + "api/user/urls")

	assert.Equal(s.T(), 202, res.StatusCode())

	wg := sync.WaitGroup{}
	wg.Add(1)
	time.AfterFunc(5*time.Second, func() {
		res, _ := s.client.R().Get(baseURL + id1)
		assert.Equal(s.T(), 410, res.StatusCode())
		res, _ = s.client.R().Get(baseURL + id2)
		assert.Equal(s.T(), 410, res.StatusCode())
		wg.Done()
	})
	wg.Wait()
}
