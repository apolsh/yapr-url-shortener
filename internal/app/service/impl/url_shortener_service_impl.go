package impl

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/app/service"
)

type URLShortenerServiceImpl struct {
	repository repository.URLRepository
}

func NewURLShortenerService(repo repository.URLRepository) service.URLShortenerService {
	return &URLShortenerServiceImpl{repository: repo}
}

func (r URLShortenerServiceImpl) AddNewURL(shortenedURLInfo entity.ShortenedURLInfo) (string, error) {
	return r.repository.Save(shortenedURLInfo)
}

func (r URLShortenerServiceImpl) GetURLByID(id string) (string, error) {
	byID, err := r.repository.GetByID(id)
	return byID.GetOriginalURL(), err
}

func (r URLShortenerServiceImpl) GetURLsByOwnerID(ownerID string) ([]entity.ShortenedURLInfo, error) {
	return r.repository.GetAllByOwner(ownerID)
}
