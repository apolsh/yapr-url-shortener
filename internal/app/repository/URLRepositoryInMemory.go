package repository

type URLRepositoryInMemory struct {
	Storage map[int]string
}

func NewURLRepositoryInMemoryImpl() URLRepository {
	return &URLRepositoryInMemory{Storage: make(map[int]string)}
}

func (receiver *URLRepositoryInMemory) Save(url string) int {
	id := len(receiver.Storage)
	receiver.Storage[id] = url
	return id
}

func (receiver URLRepositoryInMemory) GetByID(id int) string {
	s := receiver.Storage[id]
	return s
}
