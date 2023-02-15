package service

import (
	"context"
	"fmt"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

// URLShortenerServiceImpl реализация сервиса по работе с URL
type URLShortenerServiceImpl struct {
	repository  repository.URLRepository
	hostAddress string
}

var _ URLShortenerService = new(URLShortenerServiceImpl)

// NewURLShortenerService конструктор
func NewURLShortenerService(repo repository.URLRepository, hostAddress string) *URLShortenerServiceImpl {
	return &URLShortenerServiceImpl{repository: repo, hostAddress: hostAddress}
}

// AddNewURL сохраняет URL в хранилище
func (r *URLShortenerServiceImpl) AddNewURL(ctx context.Context, shortenedURLInfo entity.ShortenedURLInfo) (string, error) {
	return r.repository.Save(ctx, shortenedURLInfo)
}

// GetURLByID возвращает оригинальный URL найденный по идентификатору
func (r *URLShortenerServiceImpl) GetURLByID(ctx context.Context, id string) (string, error) {
	item, err := r.repository.GetByID(ctx, id)
	if item.IsDeleted() {
		return "", ErrorItemIsDeleted
	}
	return item.GetOriginalURL(), err
}

// GetByOriginalURL возвращает entity.ShortenedURLInfo найденный по оригинальному URL
func (r *URLShortenerServiceImpl) GetByOriginalURL(ctx context.Context, originalURL string) (entity.ShortenedURLInfo, error) {
	return r.repository.GetByOriginalURL(ctx, originalURL)
}

// GetURLsByOwnerID  возвращает массив пар (укороченная + оригинальная ссылка) найденные по владельцу URL
func (r *URLShortenerServiceImpl) GetURLsByOwnerID(ctx context.Context, ownerID string) ([]dto.URLPair, error) {
	urlInfos, err := r.repository.GetAllByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	urlPairs := make([]dto.URLPair, 0, len(urlInfos))
	for _, v := range urlInfos {
		urlPairs = append(urlPairs, v.ToURLPair(r.hostAddress))
	}
	return urlPairs, nil

}

// PingDB проверяет работоспособность хранилища на основе которого работает URLShortenerService
func (r *URLShortenerServiceImpl) PingDB(ctx context.Context) bool {
	return r.repository.Ping(ctx)
}

// AddNewURLsInBatch сохраняет массив URL в хранилище
func (r *URLShortenerServiceImpl) AddNewURLsInBatch(ctx context.Context, owner string, batch []dto.ShortenInBatchRequestItem) ([]dto.ShortenInBatchResponseItem, error) {
	correlationToID, err := r.repository.SaveBatch(ctx, owner, batch)
	if err != nil {
		return nil, err
	}
	response := make([]dto.ShortenInBatchResponseItem, 0, len(correlationToID))
	for k, v := range correlationToID {
		response = append(response, dto.ShortenInBatchResponseItem{CorrelationID: k, ShortURL: r.GetShortenURLFromID(ctx, v)})
	}
	return response, nil
}

// DeleteURLsInBatch помечает как удаленные URL, переданные в списке и принадлежащие указанному пользователю
func (r *URLShortenerServiceImpl) DeleteURLsInBatch(ctx context.Context, owner string, ids []string) error {
	return r.repository.DeleteURLsInBatch(ctx, owner, ids)
}

// GetShortenURLFromID создает укороченный URL основываясь на идентификаторе сохраненного URL
func (r *URLShortenerServiceImpl) GetShortenURLFromID(_ context.Context, id string) string {
	return fmt.Sprintf("%s/%s", r.hostAddress, id)
}

// GetAppStatistic получить статистику приложения
func (r *URLShortenerServiceImpl) GetAppStatistic(ctx context.Context) (dto.AppStatisticItem, error) {
	return r.repository.GetAppStatistic(ctx)
}
