package mock

import "github.com/apolsh/yapr-url-shortener/internal/app/repository"

type URLRepositoryMock struct {
	Storage map[int]string
}

func NewURLRepositoryMock() repository.URLRepository {
	return &URLRepositoryMock{Storage: make(map[int]string)}
}

func (receiver *URLRepositoryMock) Save(url string) int {
	id := len(receiver.Storage)
	receiver.Storage[id] = url
	return id
}

func (receiver URLRepositoryMock) GetByID(id int) string {
	s := receiver.Storage[id]
	return s
}
