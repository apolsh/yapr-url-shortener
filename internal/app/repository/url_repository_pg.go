package repository

import (
	"context"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type URLRepositoryPG struct {
	DB *pgxpool.Pool
}

func NewURLRepositoryPG(databaseDSN string) URLRepository {
	conn, err := pgxpool.Connect(context.Background(), databaseDSN)
	if err != nil {
		panic(err)
	}

	return &URLRepositoryPG{DB: conn}
}

func (repo URLRepositoryPG) Save(shortenedInfo entity.ShortenedURLInfo) (string, error) {
	return "", nil
}

func (repo URLRepositoryPG) GetByID(id string) (entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner FROM public.shortened_urls WHERE id=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, id).Scan(&info.ID, &info.OriginalURL, &info.Owner)

	if err != nil {
		return info, nil
	}

	return info, nil
}

func (repo URLRepositoryPG) GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error) {
	return make([]entity.ShortenedURLInfo, 0), nil
}

func (repo *URLRepositoryPG) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := repo.DB.Ping(ctx); err != nil {
		return false
	} else {
		return true
	}
}

func (repo *URLRepositoryPG) Close() {
	repo.DB.Close()
}
