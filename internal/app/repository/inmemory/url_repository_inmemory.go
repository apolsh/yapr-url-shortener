package inmemory

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/apolsh/yapr-url-shortener/internal/logger"
)

type void struct{}

var log = logger.LoggerOfComponent("in-memory-repo")

type backupStorage interface {
	read() (entity.ShortenedURLInfo, error)
	write(url entity.ShortenedURLInfo) error
}

type fileBackup struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func (f fileBackup) read() (entity.ShortenedURLInfo, error) {

	event := entity.ShortenedURLInfo{}
	if err := f.decoder.Decode(&event); err != nil {
		log.Error(err)
		return entity.ShortenedURLInfo{}, err
	}
	return event, nil
}

func (f fileBackup) write(url entity.ShortenedURLInfo) error {
	return f.encoder.Encode(&url)
}

func newFileBackup(filename string) (*fileBackup, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	gob.Register(entity.ShortenedURLInfo{})
	return &fileBackup{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

type URLRepositoryInMemory struct {
	Storage       map[string]entity.ShortenedURLInfo
	backupStorage backupStorage
	mu            sync.RWMutex
}

// NewURLRepositoryInMemory создает хранилище URL в памяти, если передан аргумент fileStorage,
// то при создании будут считаны записи ранее сохраненные в этот файл. а так же при работе в этот
// файл будут сохранены новые записи. Если аргумент fileStorage не передан, то все данные будут храниться
// только в памяти, сохранения на диск не происходит.
func NewURLRepositoryInMemory(m map[string]entity.ShortenedURLInfo, fileStorage string) (*URLRepositoryInMemory, error) {
	storage := m
	if fileStorage != "" {
		backupStorage, err := newFileBackup(fileStorage)
		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf(`repository initialization error: %w`, err)
		}
		for {
			url, err := backupStorage.read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error(err)
				return nil, fmt.Errorf(`repository initialization error: %w`, err)
			}
			storage[url.ID] = url
		}
		return &URLRepositoryInMemory{Storage: storage, backupStorage: backupStorage}, nil
	}
	return &URLRepositoryInMemory{Storage: storage, backupStorage: nil}, nil
}

func (r *URLRepositoryInMemory) Save(ctx context.Context, shortenedInfo entity.ShortenedURLInfo) (string, error) {
	_, err := r.GetByOriginalURL(ctx, shortenedInfo.GetOriginalURL())
	if err != nil && errors.Is(err, repository.ErrorItemNotFound) {
		id := repository.NextID()
		shortenedInfo.SetID(id)
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.backupStorage != nil {
			log.Error(err)
			err := r.backupStorage.write(shortenedInfo)
			if err != nil {
				return "", err
			}
		}
		r.Storage[id] = shortenedInfo
		return id, nil
	}
	if err == nil {
		log.Error(err)
		return "", repository.ErrorURLAlreadyStored
	}
	return "", err
}

func (r *URLRepositoryInMemory) update(_ context.Context, shortenedInfo entity.ShortenedURLInfo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.backupStorage != nil {
		err := r.backupStorage.write(shortenedInfo)
		if err != nil {
			log.Error(err)
			return "", err
		}
	}
	r.Storage[shortenedInfo.GetID()] = shortenedInfo
	return shortenedInfo.GetID(), nil
}

func (r *URLRepositoryInMemory) SaveBatch(ctx context.Context, owner string, batch []dto.ShortenInBatchRequestItem) (map[string]string, error) {
	response := make(map[string]string, len(batch))
	for _, item := range batch {
		info := entity.NewUnstoredShortenedURLInfo(owner, item.OriginalURL)
		id, _ := r.Save(ctx, *info)
		response[item.CorrelationID] = id
	}

	return response, nil
}

func (r *URLRepositoryInMemory) DeleteURLsInBatch(ctx context.Context, owner string, ids []string) error {
	for _, id := range ids {
		urlEntity, isFound := r.Storage[id]
		if isFound && urlEntity.GetOwner() == owner {
			urlEntity.SetDeleted()
			_, err := r.update(ctx, urlEntity)
			if err != nil {
				log.Error(fmt.Errorf("failed to save %s, cause: %w", urlEntity.GetID(), err))
			}
		}
	}

	return nil
}

func (r *URLRepositoryInMemory) GetByID(_ context.Context, id string) (entity.ShortenedURLInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, isFound := r.Storage[id]
	if !isFound {
		return entity.ShortenedURLInfo{}, repository.ErrorItemNotFound
	}
	return s, nil
}

func (r *URLRepositoryInMemory) GetByOriginalURL(_ context.Context, url string) (entity.ShortenedURLInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, value := range r.Storage {
		if value.OriginalURL == url {
			return value, nil
		}
	}
	return entity.ShortenedURLInfo{}, repository.ErrorItemNotFound
}

func (r *URLRepositoryInMemory) GetAllByOwner(_ context.Context, owner string) ([]entity.ShortenedURLInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	urls := make([]entity.ShortenedURLInfo, 0)
	for _, value := range r.Storage {
		if value.Owner == owner {
			urls = append(urls, value)
		}
	}
	return urls, nil
}

func (r *URLRepositoryInMemory) Close() {

}

func (r *URLRepositoryInMemory) Ping(_ context.Context) bool {
	return true
}

func (r *URLRepositoryInMemory) GetAppStatistic(_ context.Context) (dto.AppStatisticItem, error) {
	uniqOwners := make(map[string]void, 0)
	for _, urlInfo := range r.Storage {
		uniqOwners[urlInfo.Owner] = void{}
	}

	users := len(uniqOwners)
	urls := len(r.Storage)
	return dto.AppStatisticItem{Users: users, URLs: urls}, nil
}
