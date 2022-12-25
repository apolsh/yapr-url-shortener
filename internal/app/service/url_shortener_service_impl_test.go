package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/apolsh/yapr-url-shortener/internal/app/mocks"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	hostAddress = "http://localhost:8080"
	longURL1    = "http://longurl1.com"
	longURL2    = "http://longurl2.com"
)

type URLShortenerServiceSuite struct {
	suite.Suite
	repo    *mocks.MockURLRepository
	ctrl    *gomock.Controller
	service *URLShortenerServiceImpl
}

func TestURLShortenerServiceSuite(t *testing.T) {
	suite.Run(t, new(URLShortenerServiceSuite))
}

func (s *URLShortenerServiceSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.ctrl = ctrl
	s.repo = mocks.NewMockURLRepository(ctrl)

	s.service = NewURLShortenerService(s.repo, hostAddress)
}

func (s *URLShortenerServiceSuite) TestGetShortenURLFromID() {
	resp := s.service.GetShortenURLFromID(context.Background(), "123")

	assert.Equal(s.T(), hostAddress+"/"+"123", resp)
}

func (s *URLShortenerServiceSuite) TestGetURLByIDItemExist() {
	s.repo.EXPECT().GetByID(context.Background(), "123").Return(entity.ShortenedURLInfo{OriginalURL: longURL1}, nil)
	res, err := s.service.GetURLByID(context.Background(), "123")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), longURL1, res)
}

func (s *URLShortenerServiceSuite) TestGetURLByIDItemDoesNotExist() {
	s.repo.EXPECT().GetByID(context.Background(), "123").Return(entity.ShortenedURLInfo{}, ErrorItemIsDeleted)
	_, err := s.service.GetURLByID(context.Background(), "123")

	assert.Equal(s.T(), err, ErrorItemIsDeleted)
}

func (s *URLShortenerServiceSuite) TestGetURLsByOwnerIDWithSuccess() {
	owner := "owner"
	s.repo.EXPECT().GetAllByOwner(context.Background(), owner).Return([]entity.ShortenedURLInfo{
		{OriginalURL: longURL1, ID: "123"},
		{OriginalURL: longURL2, ID: "456"},
	}, nil)
	pairs, err := s.service.GetURLsByOwnerID(context.Background(), owner)

	expected := []dto.URLPair{
		{OriginalURL: longURL1, ShortURL: hostAddress + "/" + "123"},
		{OriginalURL: longURL2, ShortURL: hostAddress + "/" + "456"},
	}

	assert.NoError(s.T(), err)
	assert.True(s.T(), reflect.DeepEqual(expected, pairs))

}

func (s *URLShortenerServiceSuite) TestGetURLsByOwnerIDWithError() {
	owner := "owner"
	err := errors.New("new err")
	s.repo.EXPECT().GetAllByOwner(context.Background(), owner).Return([]entity.ShortenedURLInfo{}, err)
	_, errRes := s.service.GetURLsByOwnerID(context.Background(), owner)
	assert.Equal(s.T(), err, errRes)
}
