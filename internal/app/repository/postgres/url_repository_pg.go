package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	preparedSaveBatchStatement     = "SAVE_BATCH_STATEMENT"
	preparedMarkAsDeletedStatement = "MARK_AS_DELETED_STATEMENT"
	constraintOriginalURL          = "shortened_urls_original_url_uindex"
)

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
	ctx := context.Background()
	go func() {
		for task := range workerTaskCh {
			thisTask := task
			go func() {
				tx, err := conn.BeginTx(thisTask.ctx, pgx.TxOptions{})
				if err != nil {
					log.Println(err)
				}

				_, err = tx.Exec(ctx, thisTask.query, thisTask.args...)
				if err != nil {
					log.Println(err)
					_ = tx.Rollback(ctx)
				}

				if err := tx.Commit(ctx); err != nil {
					log.Println(err)
				}
			}()
		}
	}()

	return worker
}

type URLRepositoryPG struct {
	DB          *pgxpool.Pool
	AsyncWorker *AsyncDBTransactionWorker
}

// NewURLRepositoryPG хранилище URL в СУБД Postgres, при создании происходит выполнения скрипта создания
// необходимых таблиц, а так же создается асинхронный воркер, который может выполнять запросы к БД в асинхронном режиме
func NewURLRepositoryPG(databaseDSN string) (repository.URLRepository, error) {
	conn, err := pgxpool.Connect(context.Background(), databaseDSN)
	if err != nil {
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}
	RunMigration(databaseDSN)
	err = setupPreparedStatements(conn)
	if err != nil {
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}

	asyncWorker := newAsyncDBTransactionWorker(conn)

	return &URLRepositoryPG{DB: conn, AsyncWorker: asyncWorker}, nil
}

func (repo URLRepositoryPG) Save(shortenedInfo entity.ShortenedURLInfo) (string, error) {
	shortenedInfo.ID = repository.NextID()

	q := "INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)"

	_, err := repo.DB.Exec(context.Background(), q, shortenedInfo.GetID(), shortenedInfo.GetOriginalURL(), shortenedInfo.GetOwner(), shortenedInfo.GetStatus())

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == constraintOriginalURL {
				return "", repository.ErrorURLAlreadyStored
			}
		}
		return "", err
	}
	return shortenedInfo.GetID(), nil
}

func (repo *URLRepositoryPG) SaveBatch(owner string, batch []dto.ShortenInBatchRequestItem) (map[string]string, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	response := make(map[string]string, len(batch))

	var id string
	for _, requestItem := range batch {
		id = repository.NextID()
		_, err := tx.Exec(ctx, preparedSaveBatchStatement, id, requestItem.OriginalURL, owner, 0)
		if err != nil {
			return nil, err
		}
		response[requestItem.CorrelationID] = id
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return response, nil
}

func (repo URLRepositoryPG) GetByID(id string) (entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE id=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, id).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		return entity.ShortenedURLInfo{}, err
	}

	return info, nil
}

func (repo URLRepositoryPG) GetByOriginalURL(url string) (entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE original_url=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(context.Background(), q, url).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		return entity.ShortenedURLInfo{}, err
	}

	return info, nil
}

func (repo URLRepositoryPG) GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE owner=$1"

	rows, err := repo.DB.Query(context.Background(), q, owner)
	if err != nil {
		return nil, err
	}

	infos := make([]entity.ShortenedURLInfo, 0)

	for rows.Next() {
		var info entity.ShortenedURLInfo
		err = rows.Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)
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

func (repo *URLRepositoryPG) DeleteURLsInBatch(owner string, ids []string) error {
	q := "UPDATE shortened_urls SET status = 1 WHERE owner = $1 AND id = ANY ($2)"
	task := NewWorkerTask(context.Background(), q, owner, ids)
	repo.AsyncWorker.executeTask(task)

	return nil
}

func (repo *URLRepositoryPG) Close() {
	repo.DB.Close()
}

func setupPreparedStatements(conn *pgxpool.Pool) error {
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

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
