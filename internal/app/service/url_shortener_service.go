package service

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLShortenerService interface {
	AddNewURL(shortenedURLInfo entity.ShortenedURLInfo) (string, error)

	AddNewURLsInBatch(owner string, batch []dto.ShortenInBatchRequestItem) ([]*dto.ShortenInBatchResponseItem, error)

	GetURLByID(id string) (string, error)

	GetURLsByOwnerID(ownerID string) ([]entity.ShortenedURLInfo, error)

	PingDB() bool
}
