package inmemory

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/rs/xid"
	"io"
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
	//
	//if !f.scanner.Scan() {
	//	return "", f.scanner.Err()
	//}
	//url := f.scanner.Text()
	//return url, nil
}

func (f fileBackup) write(url *entity.ShortenedURLInfo) error {
	return f.encoder.Encode(&url)
	//if _, err := f.writer.Write([]byte(url)); err != nil {
	//	return err
	//}
	//if err := f.writer.WriteByte('\n'); err != nil {
	//	return err
	//}
	//return f.writer.Flush()
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
}

func NewURLRepositoryInMemory(fileStorage string) repository.URLRepository {
	storage := make(map[string]entity.ShortenedURLInfo)
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
			storage[url.ID] = *url
		}
		return &URLRepositoryInMemory{Storage: storage, backupStorage: backupStorage}
	}
	return &URLRepositoryInMemory{Storage: storage, backupStorage: nil}
}

func (repository *URLRepositoryInMemory) Save(shortenedURLEntity entity.ShortenedURLInfo) (string, error) {
	id := repository.nextID()
	shortenedURLEntity.SetID(id)
	if repository.backupStorage != nil {
		err := repository.backupStorage.write(&shortenedURLEntity)
		if err != nil {
			return "", err
		}
	}
	repository.Storage[id] = shortenedURLEntity
	return id, nil
}

func (repository URLRepositoryInMemory) GetByID(id string) (entity.ShortenedURLInfo, error) {
	s, isFound := repository.Storage[id]
	if !isFound {
		return s, repository.ItemNotFoundError
	}
	return s, nil
}

func (repository URLRepositoryInMemory) GetAllByOwner(owner string) ([]entity.ShortenedURLInfo, error) {
	urls := make([]entity.ShortenedURLInfo, 0)
	for _, value := range repository.Storage {
		if value.Owner == owner {
			urls = append(urls, value)
		}
	}
	return urls, nil
}

func (repository *URLRepositoryInMemory) nextID() string {
	return xid.New().String()
	//sum256 := sha256.Sum256([]byte(url.GetOriginalURL() + url.GetOwner()))
	//return hex.EncodeToString(sum256[:])
}
