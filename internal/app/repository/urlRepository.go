package repository

type UrlRepository interface {
	Save(url string) int

	GetById(id int) string
}
