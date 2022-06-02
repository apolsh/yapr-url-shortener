package service

import "github.com/apolsh/yapr-url-shortener/internal/app/repository"

type URLShortenerServiceImpl struct {
	repository repository.URLRepository
}

func NewURLShortenerService(repo repository.URLRepository) URLShortenerService {
	return &URLShortenerServiceImpl{repository: repo}
}

func (r URLShortenerServiceImpl) AddNewURL(url string) int {
	return r.repository.Save(url)
}

func (r URLShortenerServiceImpl) GetURLByID(id int) string {
	return r.repository.GetByID(id)
}
