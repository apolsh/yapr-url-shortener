package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const preparedSaveBatchStatement = "SAVE_BATCH_STATEMENT"
const preparedMarkAsDeletedStatement = "MARK_AS_DELETED_STATEMENT"
const constraintOriginalURL = "shortened_urls_original_url_uindex"

type workerTask struct {
	ctx   context.Context
	query string
	args  []interface{}
}

func NewWorkerTask(ctx context.Context, query string, args ...interface{}) *workerTask {
	return &workerTask{ctx: ctx, query: query, args: args}
}

type AsyncDBTransactionWorker struct {
	workerTaskCh chan workerTask
}

func (w *AsyncDBTransactionWorker) executeTask(task *workerTask) {
	w.workerTaskCh <- *task
}

func newAsyncDBTransactionWorker(conn *pgxpool.Pool) *AsyncDBTransactionWorker {
	workerTaskCh := make(chan workerTask, runtime.NumCPU())
	worker := &AsyncDBTransactionWorker{workerTaskCh: workerTaskCh}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for task := range workerTaskCh {
				ctx := context.Background()
				tx, err := conn.BeginTx(task.ctx, pgx.TxOptions{})
				if err != nil {
					log.Println(err)
				}

				_, err = tx.Exec(ctx, task.query, task.args...)
				if err != nil {
					log.Println(err)
					_ = tx.Rollback(ctx)
				}

				if err := tx.Commit(ctx); err != nil {
					log.Println(err)
				}
			}
		}()
	}
	return worker
}

type URLRepositoryPG struct {
	DB          *pgxpool.Pool
	AsyncWorker *AsyncDBTransactionWorker
}

func NewURLRepositoryPG(databaseDSN string) (URLRepository, error) {
	conn, err := pgxpool.Connect(context.Background(), databaseDSN)
	if err != nil {
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}
	err = setupTable(conn)
	if err != nil {
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}

	asyncWorker := newAsyncDBTransactionWorker(conn)

	return &URLRepositoryPG{DB: conn, AsyncWorker: asyncWorker}, nil
}

func (repo URLRepositoryPG) Save(info *entity.ShortenedURLInfo) (string, error) {
	info.ID = nextID()

	q := "INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)"

	_, err := repo.DB.Exec(context.Background(), q, info.GetID(), info.GetOriginalURL(), info.GetOwner(), info.GetStatus())

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

func (repo *URLRepositoryPG) SaveBatch(owner string, batch []*dto.ShortenInBatchRequestItem) ([]*dto.ShortenInBatchResponseItem, error) {
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
		_, err := tx.Exec(ctx, preparedSaveBatchStatement, id, requestItem.OriginalURL, owner, 0)
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
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE id=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, id).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo URLRepositoryPG) GetByOriginalURL(url string) (*entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE original_url=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, url).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo URLRepositoryPG) GetAllByOwner(owner string) ([]*entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE owner=$1"

	rows, err := repo.DB.Query(context.Background(), q, owner)
	if err != nil {
		return nil, err
	}

	infos := make([]*entity.ShortenedURLInfo, 0)

	for rows.Next() {
		var info entity.ShortenedURLInfo
		err = rows.Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)
		if err != nil {
			return nil, err
		}
		infos = append(infos, &info)
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

func (repo *URLRepositoryPG) DeleteURLsInBatch(owner string, ids []*string) error {
	q := "UPDATE shortened_urls SET status = 1 WHERE owner = $1 AND id = ANY ($2)"
	task := NewWorkerTask(context.Background(), q, owner, ids)
	repo.AsyncWorker.executeTask(task)

	return nil
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
		id           varchar(20)        not null,
		original_url varchar            not null,
		owner        uuid               not null,
		status       smallint default 0 not null
	)`
	createIDIndexQ := "create unique index if not exists shortened_urls_id_uindex on shortened_urls (id)"
	createOriginalURLIndexQ := "create unique index if not exists shortened_urls_original_url_uindex on shortened_urls (original_url)"
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
	_, err = tx.Prepare(ctx, preparedSaveBatchStatement, "INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = tx.Prepare(ctx, preparedMarkAsDeletedStatement, "UPDATE shortened_urls SET status = 1 WHERE id = $1")
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
