package repository

type URLRepositoryInMemoryImpl struct {
	storage map[int]string
}

func NewURLRepositoryInMemoryImpl() URLRepository {
	return &URLRepositoryInMemoryImpl{storage: make(map[int]string)}
}

func (receiver *URLRepositoryInMemoryImpl) Save(url string) int {
	id := len(receiver.storage)
	receiver.storage[id] = url
	return id
}

func (receiver URLRepositoryInMemoryImpl) GetByID(id int) string {
	s := receiver.storage[id]
	return s
}
