package service

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLShortenerServiceImpl struct {
	repository repository.URLRepository
}

func NewURLShortenerService(repo repository.URLRepository) URLShortenerService {
	return &URLShortenerServiceImpl{repository: repo}
}

func (r URLShortenerServiceImpl) AddNewURL(shortenedURLInfo *entity.ShortenedURLInfo) (string, error) {
	return r.repository.Save(shortenedURLInfo)
}

func (r URLShortenerServiceImpl) GetURLByID(id string) (string, error) {
	item, err := r.repository.GetByID(id)
	if item.IsDeleted() {
		return "", ErrorItemIsDeleted
	}
	return item.GetOriginalURL(), err
}

func (r URLShortenerServiceImpl) GetByOriginalURL(originalURL string) (*entity.ShortenedURLInfo, error) {
	return r.repository.GetByOriginalURL(originalURL)
}

func (r URLShortenerServiceImpl) GetURLsByOwnerID(ownerID string) ([]*entity.ShortenedURLInfo, error) {
	return r.repository.GetAllByOwner(ownerID)
}

func (r URLShortenerServiceImpl) PingDB() bool {
	return r.repository.Ping()
}

func (r URLShortenerServiceImpl) AddNewURLsInBatch(owner string, batch []*dto.ShortenInBatchRequestItem) ([]*dto.ShortenInBatchResponseItem, error) {
	return r.repository.SaveBatch(owner, batch)
}

func (r URLShortenerServiceImpl) DeleteURLsInBatch(owner string, ids []*string) error {
	return r.repository.DeleteURLsInBatch(owner, ids)
}
