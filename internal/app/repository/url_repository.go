package repository

import (
	"errors"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/rs/xid"
)

type URLRepository interface {
	Save(shortenedInfo entity.ShortenedURLInfo) (string, error)

	GetByID(id string) (entity.ShortenedURLInfo, error)

	GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error)

	Close()

	Ping() bool
}

var ErrorItemNotFound = errors.New("item not found")

func nextID() string {
	return xid.New().String()
}
