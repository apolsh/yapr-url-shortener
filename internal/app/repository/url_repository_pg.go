package repository

import (
	"context"
	"errors"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const preparedSaveBatchStatement = "SAVE_BATCH_STATEMENT"
const constraintOriginalURL = "shortened_urls_original_url_uindex"

type URLRepositoryPG struct {
	DB *pgxpool.Pool
}

func NewURLRepositoryPG(databaseDSN string) URLRepository {
	conn, err := pgxpool.Connect(context.Background(), databaseDSN)
	if err != nil {
		panic(err)
	}
	err = setupTable(conn)
	if err != nil {
		panic(err)
	}

	return &URLRepositoryPG{DB: conn}
}

func (repo URLRepositoryPG) Save(info entity.ShortenedURLInfo) (string, error) {
	info.ID = nextID()

	q := "INSERT INTO shortened_urls (id, original_url, owner) VALUES ($1, $2, $3)"

	_, err := repo.DB.Exec(context.Background(), q, info.GetID(), info.GetOriginalURL(), info.GetOwner())

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == constraintOriginalURL {
				return "", ErrorURLAlreadyStored
			}
		}
		return "", err
	}
	return info.GetID(), nil
}

func (repo *URLRepositoryPG) SaveBatch(owner string, batch []dto.ShortenInBatchRequestItem) ([]*dto.ShortenInBatchResponseItem, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	response := make([]*dto.ShortenInBatchResponseItem, 0, len(batch))
	var id string
	for _, requestItem := range batch {
		id = nextID()
		responseItem := &dto.ShortenInBatchResponseItem{CorrelationID: requestItem.CorrelationID, ShortURL: id}
		_, err := tx.Exec(ctx, preparedSaveBatchStatement, id, requestItem.OriginalURL, owner)
		if err != nil {
			return nil, err
		}
		response = append(response, responseItem)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return response, nil
}

func (repo URLRepositoryPG) GetByID(id string) (*entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner FROM shortened_urls WHERE id=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, id).Scan(&info.ID, &info.OriginalURL, &info.Owner)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo URLRepositoryPG) GetByOriginalURL(url string) (*entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner FROM shortened_urls WHERE original_url=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, url).Scan(&info.ID, &info.OriginalURL, &info.Owner)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo URLRepositoryPG) GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner FROM shortened_urls WHERE owner=$1"

	rows, err := repo.DB.Query(context.Background(), q, owner)
	if err != nil {
		return nil, err
	}

	infos := make([]entity.ShortenedURLInfo, 0)

	for rows.Next() {
		var info entity.ShortenedURLInfo
		err = rows.Scan(&info.ID, &info.OriginalURL, &info.Owner)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
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

func setupTable(conn *pgxpool.Pool) error {
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	createTableQ := `create table if not exists shortened_urls
					(
						id           varchar(20) not null,
						original_url varchar     not null,
						owner        uuid        not null
					)`
	createIDIndexQ := "create unique index if not exists shortened_urls_id_uindex on shortened_urls (id)"
	createOriginalURLIndexQ := "create unique index index if not exists shortened_urls_original_url_uindex on shortened_urls (original_url)"
	_, err = tx.Exec(ctx, createTableQ)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, createIDIndexQ)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, createOriginalURLIndexQ)
	if err != nil {
		return err
	}
	_, err = tx.Prepare(ctx, preparedSaveBatchStatement, "INSERT INTO shortened_urls (id, original_url, owner) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
