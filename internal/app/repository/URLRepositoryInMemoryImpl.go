package repository

type URLRepositoryInMemoryImpl struct {
	Storage map[int]string
}

func NewURLRepositoryInMemoryImpl() URLRepository {
	return &URLRepositoryInMemoryImpl{Storage: make(map[int]string)}
}

func (receiver *URLRepositoryInMemoryImpl) Save(url string) int {
	id := len(receiver.Storage)
	receiver.Storage[id] = url
	return id
}

func (receiver URLRepositoryInMemoryImpl) GetByID(id int) string {
	s := receiver.Storage[id]
	return s
}
