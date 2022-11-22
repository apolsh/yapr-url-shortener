package repository

import (
	"errors"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/rs/xid"
)

type URLRepository interface {
	Save(shortenedInfo *entity.ShortenedURLInfo) (string, error)

	SaveBatch(owner string, batch []*dto.ShortenInBatchRequestItem) (map[string]string, error)

	GetByID(id string) (*entity.ShortenedURLInfo, error)

	GetByOriginalURL(url string) (*entity.ShortenedURLInfo, error)

	GetAllByOwner(owner string) ([]*entity.ShortenedURLInfo, error)

	Close()

	Ping() bool

	DeleteURLsInBatch(owner string, ids []string) error
}

var ErrorItemNotFound = errors.New("item not found")

var ErrorURLAlreadyStored = errors.New("provided URL is already stored")

func nextID() string {
	return xid.New().String()
}
