package repository

import (
	"errors"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

type URLRepository interface {
	Save(shortenedInfo entity.ShortenedURLInfo) (string, error)

	GetByID(id string) (entity.ShortenedURLInfo, error)

	GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error)
}

var ItemNotFoundError = errors.New("item not found")
