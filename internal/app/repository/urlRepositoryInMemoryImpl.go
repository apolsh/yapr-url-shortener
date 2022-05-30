package repository

type UrlRepositoryInMemoryImpl struct {
	storage map[int]string
}

func NewUrlRepositoryInMemoryImpl() UrlRepository {
	return &UrlRepositoryInMemoryImpl{storage: make(map[int]string)}
}

func (receiver *UrlRepositoryInMemoryImpl) Save(url string) int {
	id := len(receiver.storage)
	receiver.storage[id] = url
	return id
}

func (receiver UrlRepositoryInMemoryImpl) GetById(id int) string {
	s := receiver.storage[id]
	return s
}
