package repository

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"io"
	"log"
	"os"
)

type backupStorage interface {
	read() (*entity.ShortenedURLInfo, error)
	write(url *entity.ShortenedURLInfo) error
}

type fileBackup struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func (f fileBackup) read() (*entity.ShortenedURLInfo, error) {

	event := &entity.ShortenedURLInfo{}
	if err := f.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (f fileBackup) write(url *entity.ShortenedURLInfo) error {
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
	Storage       map[string]*entity.ShortenedURLInfo
	backupStorage backupStorage
}

func NewURLRepositoryInMemory(fileStorage string) URLRepository {
	storage := make(map[string]*entity.ShortenedURLInfo)
	if fileStorage != "" {
		backupStorage, err := NewFileBackup(fileStorage)
		if err != nil {
			panic(fmt.Sprintf("Repository initialization error: %s", err))
		}
		for {
			url, err := backupStorage.read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(fmt.Sprintf("Backup restore error: %s", err))
			}
			storage[url.ID] = url
		}
		return &URLRepositoryInMemory{Storage: storage, backupStorage: backupStorage}
	}
	return &URLRepositoryInMemory{Storage: storage, backupStorage: nil}
}

func (r *URLRepositoryInMemory) Save(shortenedInfo *entity.ShortenedURLInfo) (string, error) {
	_, err := r.GetByOriginalURL(shortenedInfo.GetOriginalURL())
	if err != nil && errors.Is(err, ErrorItemNotFound) {
		if shortenedInfo.GetID() == "" {
			id := nextID()
			shortenedInfo.SetID(id)
		}
		if r.backupStorage != nil {
			err := r.backupStorage.write(shortenedInfo)
			if err != nil {
				return "", err
			}
		}
		r.Storage[shortenedInfo.GetID()] = shortenedInfo
		return shortenedInfo.GetID(), nil
	}
	return "", err
}

func (r *URLRepositoryInMemory) SaveBatch(owner string, batch []*dto.ShortenInBatchRequestItem) ([]*dto.ShortenInBatchResponseItem, error) {
	response := make([]*dto.ShortenInBatchResponseItem, 0, len(batch))
	for _, item := range batch {
		info := entity.NewUnstoredShortenedURLInfo(owner, item.OriginalURL)
		id, _ := r.Save(info)
		response = append(response, &dto.ShortenInBatchResponseItem{CorrelationID: item.CorrelationID, ShortURL: id})
	}

	return response, nil
}

func (r *URLRepositoryInMemory) DeleteURLsInBatch(owner string, ids []*string) error {
	for _, id := range ids {
		urlEntity, isFound := r.Storage[*id]
		if isFound && urlEntity.GetOwner() == owner {
			urlEntity.SetDeleted()
			delete(r.Storage, urlEntity.GetID())
			_, err := r.Save(urlEntity)
			if err != nil {
				log.Println("failed to save ", urlEntity.GetID(), err.Error())
			}
		}
	}

	return nil
}

func (r URLRepositoryInMemory) GetByID(id string) (*entity.ShortenedURLInfo, error) {
	s, isFound := r.Storage[id]
	if !isFound {
		return nil, ErrorItemNotFound
	}
	return s, nil
}

func (r URLRepositoryInMemory) GetByOriginalURL(url string) (*entity.ShortenedURLInfo, error) {
	for _, value := range r.Storage {
		if value.OriginalURL == url {
			return value, nil
		}
	}
	return nil, ErrorItemNotFound
}

func (r URLRepositoryInMemory) GetAllByOwner(owner string) ([]*entity.ShortenedURLInfo, error) {
	urls := make([]*entity.ShortenedURLInfo, 0)
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
