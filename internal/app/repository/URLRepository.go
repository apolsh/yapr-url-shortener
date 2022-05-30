package repository

type URLRepository interface {
	Save(url string) int

	GetByID(id int) string
}
