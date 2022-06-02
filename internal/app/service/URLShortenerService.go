package service

type URLShortenerService interface {
	AddNewURL(url string) int

	GetURLByID(id int) string
}
