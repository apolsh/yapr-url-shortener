package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLRepositoryInMemoryImpl_GetByID(t *testing.T) {
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
		repository := &URLRepositoryInMemory{Storage: test.repositoryMock}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[int]string)
			for _, item := range test.urls {
				resultMap[repository.Save(item)] = item
			}
			assert.Equal(t, test.expectedIds, resultMap)
		})
	}
}

func TestURLRepositoryInMemoryImpl_Save(t *testing.T) {
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
		repository := &URLRepositoryInMemory{Storage: test.repositoryMock}

		t.Run(test.name, func(t *testing.T) {
			resultMap := make(map[int]string)
			for _, item := range test.ids {
				resultMap[item] = repository.GetByID(item)
			}
			assert.Equal(t, test.expectedIds, resultMap)
		})
	}
}
