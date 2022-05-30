package service

import "github.com/apolsh/yapr-url-shortener/internal/app/repository"

type URLShortenerService struct {
	repository repository.URLRepository
}

func NewURLShortenerService() *URLShortenerService {
	return &URLShortenerService{repository: repository.NewURLRepositoryInMemoryImpl()}
}

func (r URLShortenerService) AddNewURL(url string) int {
	return r.repository.Save(url)
}

func (r URLShortenerService) GetURLByID(id int) string {
	return r.repository.GetByID(id)
}
