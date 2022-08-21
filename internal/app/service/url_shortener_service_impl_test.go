package service

import (
	"github.com/apolsh/yapr-url-shortener/internal/app/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLShortenerService_AddNewURL(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock map[int]string
		urls           []string
		expectedIds    map[int]string
	}{
		{
			name:           "add single item to empty repo",
			repositoryMock: make(map[int]string),
			urls:           []string{"dummyUrl1"},
			expectedIds:    map[int]string{0: "dummyUrl1"},
		},
		{
			name:           "add multiple item to empty repo",
			repositoryMock: make(map[int]string),
			urls:           []string{"dummyUrl1", "dummyUrl2"},
			expectedIds:    map[int]string{0: "dummyUrl1", 1: "dummyUrl2"},
		},
		{
			name:           "add single item to not empty repo",
			repositoryMock: map[int]string{0: "dummyUrl1"},
			urls:           []string{"dummyUrl2"},
			expectedIds:    map[int]string{1: "dummyUrl2"},
		},
	}

	for _, test := range tests {
		service := &URLShortenerServiceImpl{repository: &mock.URLRepositoryMock{Storage: test.repositoryMock}}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[int]string)
			for _, item := range test.urls {
				resultMap[service.AddNewURL()] = item
			}
			assert.Equal(t, test.expectedIds, resultMap)
		})
	}
}

func TestURLShortenerService_GetURLByID(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock map[int]string
		ids            []int
		expectedIds    map[int]string
	}{
		{
			name:           "Get single item from empty repo",
			repositoryMock: make(map[int]string),
			ids:            []int{0},
			expectedIds:    map[int]string{0: ""},
		},
		{
			name:           "Get item from not empty repo",
			repositoryMock: map[int]string{0: "dummyUrl1"},
			ids:            []int{0},
			expectedIds:    map[int]string{0: "dummyUrl1"},
		},
		{
			name:           "Get multiple items from not empty repo",
			repositoryMock: map[int]string{0: "dummyUrl1", 1: "dummyUrl2"},
			ids:            []int{0, 1},
			expectedIds:    map[int]string{0: "dummyUrl1", 1: "dummyUrl2"},
		},
	}

	for _, test := range tests {
		service := &URLShortenerServiceImpl{repository: &mock.URLRepositoryMock{Storage: test.repositoryMock}}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[int]string)
			for _, item := range test.ids {
				resultMap[item], _ = service.GetURLByID(item)
			}
			assert.Equal(t, test.expectedIds, resultMap)
		})
	}
}
