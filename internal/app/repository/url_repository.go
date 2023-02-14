//go:generate mockgen -destination=../mocks/url_repository_mock.go -package=mocks github.com/apolsh/yapr-url-shortener/internal/app/repository URLRepository

package repository

import (
	"context"
	"errors"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/rs/xid"
)

type URLRepository interface {
	// Save сохраняет URL в хранилище
	Save(ctx context.Context, shortenedInfo entity.ShortenedURLInfo) (string, error)

	// SaveBatch сохраняет массив URL в хранилище
	SaveBatch(ctx context.Context, owner string, batch []dto.ShortenInBatchRequestItem) (map[string]string, error)

	// GetByID возвращает ShortenedURLInfo найденный по идентификатору
	GetByID(ctx context.Context, id string) (entity.ShortenedURLInfo, error)

	// GetByOriginalURL возвращает ShortenedURLInfo найденный по оригинальному URL
	GetByOriginalURL(ctx context.Context, url string) (entity.ShortenedURLInfo, error)

	// GetAllByOwner  возвращает массив ShortenedURLInfo найденный по владельцу URL
	GetAllByOwner(ctx context.Context, owner string) ([]entity.ShortenedURLInfo, error)

	// Close инициирует завершение процессов в хранилища
	Close()

	// Ping проверяет работоспособность хранилища
	Ping(ctx context.Context) bool

	// DeleteURLsInBatch помечает как удаленные URL, переданные в списке и принадлежащие указанному пользователю
	DeleteURLsInBatch(ctx context.Context, owner string, ids []string) error
}

var ErrorItemNotFound = errors.New("item not found")

var ErrorURLAlreadyStored = errors.New("provided URL is already stored")

func NextID() string {
	return xid.New().String()
}
