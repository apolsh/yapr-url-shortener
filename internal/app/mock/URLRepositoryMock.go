package mock

import "github.com/apolsh/yapr-url-shortener/internal/app/repository"

type URLRepositoryMock struct {
	Storage map[int]string
}

func NewURLRepositoryMock(predefined map[int]string) repository.URLRepository {
	return &URLRepositoryMock{Storage: predefined}
}

func (repository *URLRepositoryMock) Save(url string) int {
	id := repository.nextID()
	repository.Storage[id] = url
	return id
}

func (repository URLRepositoryMock) GetByID(id int) string {
	s := repository.Storage[id]
	return s
}

func (repository *URLRepositoryMock) nextID() int {
	return len(repository.Storage)
}
