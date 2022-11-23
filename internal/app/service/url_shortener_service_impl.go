package service

import (
	"fmt"
	"log"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLShortenerServiceImpl struct {
	repository  repository.URLRepository
	hostAddress string
}

var _ URLShortenerService = new(URLShortenerServiceImpl)

func NewURLShortenerService(repo repository.URLRepository, hostAddress string) *URLShortenerServiceImpl {
	return &URLShortenerServiceImpl{repository: repo, hostAddress: hostAddress}
}

func (r *URLShortenerServiceImpl) AddNewURL(shortenedURLInfo entity.ShortenedURLInfo) (string, error) {
	return r.repository.Save(shortenedURLInfo)
}

func (r *URLShortenerServiceImpl) GetURLByID(id string) (string, error) {
	item, err := r.repository.GetByID(id)
	log.Println(item)
	if item.IsDeleted() {
		return "", ErrorItemIsDeleted
	}
	return item.GetOriginalURL(), err
}

func (r *URLShortenerServiceImpl) GetByOriginalURL(originalURL string) (entity.ShortenedURLInfo, error) {
	return r.repository.GetByOriginalURL(originalURL)
}

func (r *URLShortenerServiceImpl) GetURLsByOwnerID(ownerID string) ([]dto.URLPair, error) {
	urlInfos, err := r.repository.GetAllByOwner(ownerID)
	if err != nil {
		return nil, err
	}
	urlPairs := make([]dto.URLPair, 0, len(urlInfos))
	for _, v := range urlInfos {
		urlPairs = append(urlPairs, v.ToURLPair(r.hostAddress))
	}
	return urlPairs, nil

}

func (r *URLShortenerServiceImpl) PingDB() bool {
	return r.repository.Ping()
}

func (r *URLShortenerServiceImpl) AddNewURLsInBatch(owner string, batch []dto.ShortenInBatchRequestItem) ([]dto.ShortenInBatchResponseItem, error) {
	correlationToID, err := r.repository.SaveBatch(owner, batch)
	if err != nil {
		return nil, err
	}
	response := make([]dto.ShortenInBatchResponseItem, 0, len(correlationToID))
	for k, v := range correlationToID {
		response = append(response, dto.ShortenInBatchResponseItem{CorrelationID: k, ShortURL: r.GetShortenURLFromID(v)})
	}
	return response, nil
}

func (r *URLShortenerServiceImpl) DeleteURLsInBatch(owner string, ids []string) error {
	return r.repository.DeleteURLsInBatch(owner, ids)
}

func (r *URLShortenerServiceImpl) GetShortenURLFromID(id string) string {
	return fmt.Sprintf("%s/%s", r.hostAddress, id)
}
