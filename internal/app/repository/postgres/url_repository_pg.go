package postgres

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	preparedSaveBatchStatement     = "SAVE_BATCH_STATEMENT"
	preparedMarkAsDeletedStatement = "MARK_AS_DELETED_STATEMENT"
	constraintOriginalURL          = "shortened_urls_original_url_uindex"
)

var (
	workerLogger = logger.LoggerOfComponent("db-async-worker")
	dbLogger     = logger.LoggerOfComponent("url-db-pg")
)

type workerTask struct {
	query string
	args  []interface{}
}

func newWorkerTask(query string, args ...interface{}) *workerTask {
	return &workerTask{query: query, args: args}
}

type asyncDBTransactionWorker struct {
	workerTaskCh chan workerTask
	wg           sync.WaitGroup
}

func (w *asyncDBTransactionWorker) executeTask(task *workerTask) {
	w.wg.Add(1)
	w.workerTaskCh <- *task
}

func (w *asyncDBTransactionWorker) close() {
	waitWithTimeout(&w.wg, 5*time.Second)
}

func newAsyncDBTransactionWorker(conn *pgxpool.Pool) *asyncDBTransactionWorker {
	workerTaskCh := make(chan workerTask, runtime.NumCPU())
	worker := &asyncDBTransactionWorker{workerTaskCh: workerTaskCh}
	ctx := context.Background()
	go func() {
		for task := range workerTaskCh {
			thisTask := task
			go func() {
				tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
				if err != nil {
					workerLogger.Error(err)
				}

				_, err = tx.Exec(ctx, thisTask.query, thisTask.args...)
				if err != nil {
					workerLogger.Error(err)
					_ = tx.Rollback(ctx)
				}

				if err := tx.Commit(ctx); err != nil {
					workerLogger.Error(err)
				}
				worker.wg.Done()
			}()
		}
	}()

	return worker
}

// URLRepositoryPG хранилище URL в postgres
type URLRepositoryPG struct {
	DB          *pgxpool.Pool
	AsyncWorker *asyncDBTransactionWorker
}

// NewURLRepositoryPG хранилище URL в СУБД Postgres, при создании происходит выполнения скрипта создания
// необходимых таблиц, а так же создается асинхронный воркер, который может выполнять запросы к БД в асинхронном режиме
func NewURLRepositoryPG(databaseDSN string) (repository.URLRepository, error) {
	conn, err := pgxpool.Connect(context.Background(), databaseDSN)
	if err != nil {
		dbLogger.Error(err)
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}
	RunMigration(databaseDSN)
	err = setupPreparedStatements(conn)
	if err != nil {
		dbLogger.Error(err)
		return nil, fmt.Errorf(`repository initialization error: %w`, err)
	}

	asyncWorker := newAsyncDBTransactionWorker(conn)

	return &URLRepositoryPG{DB: conn, AsyncWorker: asyncWorker}, nil
}

// Save сохранить новый URL
func (repo URLRepositoryPG) Save(ctx context.Context, shortenedInfo entity.ShortenedURLInfo) (string, error) {
	shortenedInfo.ID = repository.NextID()

	q := "INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)"

	_, err := repo.DB.Exec(ctx, q, shortenedInfo.GetID(), shortenedInfo.GetOriginalURL(), shortenedInfo.GetOwner(), shortenedInfo.GetStatus())

	var pgErr *pgconn.PgError
	if err != nil {
		dbLogger.Error(err)
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == constraintOriginalURL {
				return "", repository.ErrorURLAlreadyStored
			}
		}
		return "", err
	}
	return shortenedInfo.GetID(), nil
}

// SaveBatch сохранить сразу несколько новых URL
func (repo *URLRepositoryPG) SaveBatch(ctx context.Context, owner string, batch []dto.ShortenInBatchRequestItem) (map[string]string, error) {
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
			dbLogger.Error(err)
			return nil, err
		}
		response[requestItem.CorrelationID] = id
	}

	if err := tx.Commit(ctx); err != nil {
		dbLogger.Error(err)
		return nil, err
	}
	return response, nil
}

// GetByID получить URL по ID
func (repo URLRepositoryPG) GetByID(ctx context.Context, id string) (entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE id=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(ctx, q, id).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		dbLogger.Error(err)
		return entity.ShortenedURLInfo{}, err
	}

	return info, nil
}

// GetByOriginalURL получить URL сущность по url
func (repo URLRepositoryPG) GetByOriginalURL(ctx context.Context, url string) (entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE original_url=$1"
	var info entity.ShortenedURLInfo
	err := repo.DB.QueryRow(ctx, q, url).Scan(&info.ID, &info.OriginalURL, &info.Owner, &info.Status)

	if err != nil {
		dbLogger.Error(err)
		return entity.ShortenedURLInfo{}, err
	}

	return info, nil
}

// GetAllByOwner получить все URL по пользователю
func (repo URLRepositoryPG) GetAllByOwner(ctx context.Context, owner string) ([]entity.ShortenedURLInfo, error) {
	q := "SELECT id, original_url, owner, status FROM shortened_urls WHERE owner=$1"

	rows, err := repo.DB.Query(ctx, q, owner)
	if err != nil {
		dbLogger.Error(err)
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

// Ping проверка доступа до БД
func (repo *URLRepositoryPG) Ping(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := repo.DB.Ping(ctx); err != nil {
		return false
	} else {
		return true
	}
}

// DeleteURLsInBatch удалить несколько URL за 1 запрос
func (repo *URLRepositoryPG) DeleteURLsInBatch(_ context.Context, owner string, ids []string) error {
	q := "UPDATE shortened_urls SET status = 1 WHERE owner = $1 AND id = ANY ($2)"
	task := newWorkerTask(q, owner, ids)
	repo.AsyncWorker.executeTask(task)

	return nil
}

// GetAppStatistic получить статистику приложения
func (repo *URLRepositoryPG) GetAppStatistic(ctx context.Context) (dto.AppStatisticItem, error) {
	var res dto.AppStatisticItem

	err := repo.DB.QueryRow(ctx, "SELECT COUNT(*) FROM shortened_urls").Scan(&res.URLs)
	if err != nil {
		return res, err
	}
	err = repo.DB.QueryRow(ctx, "SELECT COUNT(DISTINCT owner) FROM shortened_urls").Scan(&res.Users)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Close для graceful завершения работы БД
func (repo *URLRepositoryPG) Close() {
	repo.DB.Close()
	repo.AsyncWorker.close()
}

func setupPreparedStatements(conn *pgxpool.Pool) error {
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		dbLogger.Error(err)
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Prepare(ctx, preparedSaveBatchStatement, "INSERT INTO shortened_urls (id, original_url, owner, status) VALUES ($1, $2, $3, $4)")
	if err != nil {
		dbLogger.Error(err)
		return err
	}
	_, err = tx.Prepare(ctx, preparedMarkAsDeletedStatement, "UPDATE shortened_urls SET status = 1 WHERE id = $1")
	if err != nil {
		dbLogger.Error(err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		dbLogger.Error(err)
		return err
	}
	return nil
}

func waitWithTimeout(wg *sync.WaitGroup, t time.Duration) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		workerLogger.Info("worker shutdown correctly")
	case <-time.After(t):
		workerLogger.Info("worker shutdown timeout exceeded")
	}
}
