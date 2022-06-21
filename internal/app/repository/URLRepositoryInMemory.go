package repository

import (
	"bufio"
	"fmt"
	"os"
)

type backupStorage interface {
	read() (string, error)
	write(url string) error
}

type fileBackup struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func (f fileBackup) read() (string, error) {
	if !f.scanner.Scan() {
		return "", f.scanner.Err()
	}
	url := f.scanner.Text()
	return url, nil
}

func (f fileBackup) write(url string) error {
	if _, err := f.writer.Write([]byte(url)); err != nil {
		return err
	}
	if err := f.writer.WriteByte('\n'); err != nil {
		return err
	}
	return f.writer.Flush()
}

func NewFileBackup(filename string) (*fileBackup, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &fileBackup{
		file:    file,
		scanner: bufio.NewScanner(file),
		writer:  bufio.NewWriter(file),
	}, nil
}

type URLRepositoryInMemory struct {
	Storage       map[int]string
	backupStorage backupStorage
}

func NewURLRepositoryInMemory(fileStorage string) URLRepository {
	storage := make(map[int]string)
	if fileStorage != "" {
		backupStorage, err := NewFileBackup(fileStorage)
		if err != nil {
			panic(fmt.Sprintf("Repository initialization error: %s", err))
		}
		for {
			url, err := backupStorage.read()
			if err != nil {
				panic(fmt.Sprintf("Backup restore error: %s", err))
			}
			if url == "" {
				break
			}
			storage[len(storage)] = url
		}
		return &URLRepositoryInMemory{Storage: storage, backupStorage: backupStorage}

	}

	return &URLRepositoryInMemory{Storage: storage, backupStorage: nil}
}

func (repository *URLRepositoryInMemory) Save(url string) int {
	id := repository.NextID()
	if repository.backupStorage != nil {
		err := repository.backupStorage.write(url)
		if err != nil {
			return 0
		}
	}
	repository.Storage[id] = url
	return id
}

func (repository URLRepositoryInMemory) GetByID(id int) string {
	s := repository.Storage[id]
	return s
}

func (repository *URLRepositoryInMemory) NextID() int {
	return len(repository.Storage)
}
