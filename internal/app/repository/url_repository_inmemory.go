package repository

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
)

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
		return entity.ShortenedURLInfo{}, err
	}
	return event, nil
}

func (f fileBackup) write(url entity.ShortenedURLInfo) error {
	return f.encoder.Encode(&url)
}

func NewFileBackup(filename string) (*fileBackup, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
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

func NewURLRepositoryInMemory(fileStorage string) (URLRepository, error) {
	storage := make(map[string]entity.ShortenedURLInfo)
	if fileStorage != "" {
		backupStorage, err := NewFileBackup(fileStorage)
		if err != nil {
			return nil, fmt.Errorf(`repository initialization error: %w`, err)
		}
		for {
			url, err := backupStorage.read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf(`repository initialization error: %w`, err)
			}
			storage[url.ID] = url
		}
		return &URLRepositoryInMemory{Storage: storage, backupStorage: backupStorage}, nil
	}
	return &URLRepositoryInMemory{Storage: storage, backupStorage: nil}, nil
}

func (r *URLRepositoryInMemory) Save(shortenedInfo entity.ShortenedURLInfo) (string, error) {
	_, err := r.GetByOriginalURL(shortenedInfo.GetOriginalURL())
	if err != nil && errors.Is(err, ErrorItemNotFound) {
		id := nextID()
		shortenedInfo.SetID(id)
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.backupStorage != nil {
			err := r.backupStorage.write(shortenedInfo)
			if err != nil {
				return "", err
			}
		}
		r.Storage[id] = shortenedInfo
		return id, nil
	}
	return "", err
}

func (r *URLRepositoryInMemory) update(shortenedInfo entity.ShortenedURLInfo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.backupStorage != nil {
		err := r.backupStorage.write(shortenedInfo)
		if err != nil {
			return "", err
		}
	}
	r.Storage[shortenedInfo.GetID()] = shortenedInfo
	return shortenedInfo.GetID(), nil
}

func (r *URLRepositoryInMemory) SaveBatch(owner string, batch []dto.ShortenInBatchRequestItem) (map[string]string, error) {
	response := make(map[string]string, len(batch))
	for _, item := range batch {
		info := entity.NewUnstoredShortenedURLInfo(owner, item.OriginalURL)
		id, _ := r.Save(*info)
		response[item.CorrelationID] = id
	}

	return response, nil
}

func (r *URLRepositoryInMemory) DeleteURLsInBatch(owner string, ids []string) error {
	for _, id := range ids {
		urlEntity, isFound := r.Storage[id]
		if isFound && urlEntity.GetOwner() == owner {
			urlEntity.SetDeleted()
			_, err := r.update(urlEntity)
			if err != nil {
				log.Println("failed to save ", urlEntity.GetID(), err.Error())
			}
		}
	}

	return nil
}

func (r *URLRepositoryInMemory) GetByID(id string) (entity.ShortenedURLInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, isFound := r.Storage[id]
	if !isFound {
		return entity.ShortenedURLInfo{}, ErrorItemNotFound
	}
	return s, nil
}

func (r *URLRepositoryInMemory) GetByOriginalURL(url string) (entity.ShortenedURLInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, value := range r.Storage {
		if value.OriginalURL == url {
			return value, nil
		}
	}
	return entity.ShortenedURLInfo{}, ErrorItemNotFound
}

func (r *URLRepositoryInMemory) GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error) {
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

func (r *URLRepositoryInMemory) Ping() bool {
	return true
}
