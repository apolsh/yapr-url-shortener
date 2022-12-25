//go:generate mockgen -destination=../mocks/url_shortener_service_mock.go -package=mocks github.com/apolsh/yapr-url-shortener/internal/app/service URLShortenerService
package service

import (
	"context"
	"errors"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLShortenerService interface {
	//AddNewURL сохраняет URL в хранилище
	AddNewURL(ctx context.Context, shortenedURLInfo entity.ShortenedURLInfo) (string, error)

	// AddNewURLsInBatch сохраняет массив URL в хранилище
	AddNewURLsInBatch(ctx context.Context, owner string, batch []dto.ShortenInBatchRequestItem) ([]dto.ShortenInBatchResponseItem, error)

	// GetURLByID возвращает оригинальный URL найденный по идентификатору
	GetURLByID(ctx context.Context, id string) (string, error)

	// GetByOriginalURL возвращает entity.ShortenedURLInfo найденный по оригинальному URL
	GetByOriginalURL(ctx context.Context, url string) (entity.ShortenedURLInfo, error)

	// GetURLsByOwnerID  возвращает массив пар (укороченная + оригинальная ссылка) найденные по владельцу URL
	GetURLsByOwnerID(ctx context.Context, ownerID string) ([]dto.URLPair, error)

	// PingDB проверяет работоспособность хранилища на основе которого работает URLShortenerService
	PingDB(ctx context.Context) bool

	// DeleteURLsInBatch помечает как удаленные URL, переданные в списке и принадлежащие указанному пользователю
	DeleteURLsInBatch(ctx context.Context, owner string, ids []string) error

	//GetShortenURLFromID создает укороченный URL основываясь на идентификаторе сохраненного URL
	GetShortenURLFromID(ctx context.Context, id string) string
}

var ErrorItemIsDeleted = errors.New("item is marked as deleted")
