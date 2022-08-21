package repository

import (
	"context"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/jackc/pgx/v4"
	"time"
)

type URLRepositoryPG struct {
	DB *pgx.Conn
}

func NewURLRepositoryPG(databaseDSN string) URLRepository {
	conn, err := pgx.Connect(context.Background(), databaseDSN)

	//db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		panic(err)
	}

	return &URLRepositoryPG{DB: conn}
}

func (repo URLRepositoryPG) Save(shortenedInfo entity.ShortenedURLInfo) (string, error) {
	return "", nil
}

func (repo URLRepositoryPG) GetByID(id string) (entity.ShortenedURLInfo, error) {
	return *entity.NewUnstoredShortenedURLInfo("", ""), nil
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
	err := repo.DB.Close(context.Background())
	if err != nil {
		fmt.Println("failed to close postgres connection: ", err.Error())
	}
}
