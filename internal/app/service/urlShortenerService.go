package service

import "github.com/apolsh/yapr-url-shortener/internal/app/repository"

type UrlShortenerService struct {
	repository repository.UrlRepository
}

func NewUrlShortenerService() *UrlShortenerService {
	return &UrlShortenerService{repository: repository.NewUrlRepositoryInMemoryImpl()}
}

func (r UrlShortenerService) AddNewUrl(url string) int {
	return r.repository.Save(url)
}

func (r UrlShortenerService) GetUrlById(id int) string {
	return r.repository.GetById(id)
}
