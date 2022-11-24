package service

import (
	"errors"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLShortenerService interface {
	AddNewURL(shortenedURLInfo entity.ShortenedURLInfo) (string, error)

	AddNewURLsInBatch(owner string, batch []dto.ShortenInBatchRequestItem) ([]dto.ShortenInBatchResponseItem, error)

	GetURLByID(id string) (string, error)

	GetByOriginalURL(url string) (entity.ShortenedURLInfo, error)

	GetURLsByOwnerID(ownerID string) ([]dto.URLPair, error)

	PingDB() bool

	DeleteURLsInBatch(owner string, ids []string) error

	GetShortenURLFromID(id string) string
}

var ErrorItemIsDeleted = errors.New("item is marked as deleted")
