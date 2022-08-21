package inmemory

/*
import (
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var dummyURL1 = *entity.NewShortenedURLInfo("0", "owner", "dummyURL1")
var dummyURL2 = *entity.NewShortenedURLInfo("1", "owner", "dummyURL2")
var dummyURL3 = *entity.NewShortenedURLInfo("2", "owner", "dummyURL3")

func TestURLRepositoryInMemoryImpl_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock map[string]entity.ShortenedURLInfo
		urls           []entity.ShortenedURLInfo
		expectedIds    map[string]entity.ShortenedURLInfo
	}{
		{
			name:           "add single item to empty repo",
			repositoryMock: make(map[string]entity.ShortenedURLInfo),
			urls:           []entity.ShortenedURLInfo{dummyURL1},
			expectedIds:    map[string]entity.ShortenedURLInfo{"0": dummyURL1},
		},
		{
			name:           "add multiple item to empty repo",
			repositoryMock: make(map[string]entity.ShortenedURLInfo),
			urls:           []entity.ShortenedURLInfo{dummyURL1, dummyURL2},
			expectedIds:    map[string]entity.ShortenedURLInfo{"0": dummyURL1, "1": dummyURL2},
		},
		{
			name:           "add single item to not empty repo",
			repositoryMock: map[string]entity.ShortenedURLInfo{"0": dummyURL1},
			urls:           []entity.ShortenedURLInfo{dummyURL2},
			expectedIds:    map[string]entity.ShortenedURLInfo{"1": dummyURL2},
		},
	}

	for _, test := range tests {
		repository := &URLRepositoryInMemory{Storage: test.repositoryMock}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[string]int)
			for index, item := range test.urls {
				id, _ := repository.Save(item.GetOriginalURL())
				resultMap[id] = index
			}
			for key, index := range resultMap {
				originalURL := test.urls[index]
				storedURL := test.repositoryMock[key]
				assert.Equal(t, originalURL.GetOriginalURL(), storedURL.GetOriginalURL())
				assert.Equal(t, originalURL.GetOwner(), storedURL.GetOwner())
			}
		})
	}
}

func TestURLRepositoryInMemoryImpl_Save(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock map[string]entity.ShortenedURLInfo
		ids            []int
		expectedIds    map[string]entity.ShortenedURLInfo
	}{
		//TODO: запросить несуществующий идентификатор
		//{
		//	name:           "Get single item from empty repo",
		//	repositoryMock: make(map[int]entity.ShortenedURLInfo),
		//	ids:            []int{0},
		//	expectedIds:    map[int]entity.ShortenedURLInfo{0: dummyURL1},
		//},
		{
			name:           "Get item from not empty repo",
			repositoryMock: map[string]entity.ShortenedURLInfo{"0": dummyURL1},
			ids:            []int{0},
			expectedIds:    map[string]entity.ShortenedURLInfo{"0": dummyURL1},
		},
		{
			name:           "Get multiple items from not empty repo",
			repositoryMock: map[string]entity.ShortenedURLInfo{"0": dummyURL1, "1": dummyURL2},
			ids:            []int{0, 1},
			expectedIds:    map[string]entity.ShortenedURLInfo{"0": dummyURL1, "1": dummyURL2},
		},
	}

	for _, test := range tests {
		repository := &URLRepositoryInMemory{Storage: test.repositoryMock}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[int]entity.ShortenedURLInfo)
			for _, item := range test.ids {
				resultMap[item], _ = repository.GetByID(strconv.Itoa(item))
			}
			assert.Equal(t, test.expectedIds, resultMap)
		})
	}
}

*/
