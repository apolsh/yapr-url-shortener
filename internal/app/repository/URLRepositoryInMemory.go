package repository

type URLRepositoryInMemory struct {
	Storage map[int]string
}

func NewURLRepositoryInMemory() URLRepository {
	return &URLRepositoryInMemory{Storage: make(map[int]string)}
}

func (repository *URLRepositoryInMemory) Save(url string) int {
	id := repository.NextID()
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
