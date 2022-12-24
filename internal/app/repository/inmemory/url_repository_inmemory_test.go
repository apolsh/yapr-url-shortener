package inmemory

import (
	"testing"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
	"github.com/apolsh/yapr-url-shortener/internal/app/repository/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	ownerID   = "ownderID"
	dummyURL1 = "http://dummyURL1.com"
	dummyURL2 = "http://dummyURL2.com"
	dummyURL3 = "http://dummyURL3.com"
	dummyID1  = "1"
	dummyID2  = "2"
	dummyID3  = "3"
)

type URLRepositoryInMemorySuite struct {
	suite.Suite
	rep *URLRepositoryInMemory
	m   map[string]entity.ShortenedURLInfo
}

func TestURLRepositoryInMemorySuite(t *testing.T) {
	suite.Run(t, new(URLRepositoryInMemorySuite))
}

func (s *URLRepositoryInMemorySuite) SetupTest() {
	m := make(map[string]entity.ShortenedURLInfo)
	s.m = m
	s.rep, _ = NewURLRepositoryInMemory(m, "")
}

func (s *URLRepositoryInMemorySuite) TestSaveAlreadyStored() {
	_, _ = s.rep.Save(*entity.NewUnstoredShortenedURLInfo(ownerID, dummyURL1))
	_, err := s.rep.Save(*entity.NewUnstoredShortenedURLInfo(ownerID, dummyURL1))
	assert.ErrorIs(s.T(), repository.ErrorURLAlreadyStored, err)
}

func (s *URLRepositoryInMemorySuite) TestSaveNewURL() {
	id, err := s.rep.Save(*entity.NewUnstoredShortenedURLInfo(ownerID, dummyURL1))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.m[id].OriginalURL, dummyURL1)
}

func (s *URLRepositoryInMemorySuite) TestSaveBatch() {
	batch := make([]dto.ShortenInBatchRequestItem, 2)
	batch = append(batch, dto.ShortenInBatchRequestItem{OriginalURL: dummyURL1, CorrelationID: "1"})
	batch = append(batch, dto.ShortenInBatchRequestItem{OriginalURL: dummyURL2, CorrelationID: "2"})
	response, err := s.rep.SaveBatch(ownerID, batch)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), dummyURL1, s.m[response["1"]].OriginalURL)
	assert.Equal(s.T(), dummyURL2, s.m[response["2"]].OriginalURL)
}

func (s *URLRepositoryInMemorySuite) TestGetByID() {
	urlInfo := *entity.NewUnstoredShortenedURLInfo(ownerID, dummyURL1)
	urlInfo.ID = dummyID1
	s.m[dummyID1] = urlInfo
	info, err := s.rep.GetByID(dummyID1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), info.Owner, ownerID)
	assert.Equal(s.T(), info.OriginalURL, dummyURL1)
	assert.Equal(s.T(), info.ID, dummyID1)
}

func (s *URLRepositoryInMemorySuite) TestGetByIDNotFound() {
	_, err := s.rep.GetByID(dummyID1)
	assert.ErrorIs(s.T(), repository.ErrorItemNotFound, err)
}

func (s *URLRepositoryInMemorySuite) TestGetByOriginalURL() {
	urlInfo := *entity.NewUnstoredShortenedURLInfo(ownerID, dummyURL1)
	urlInfo.ID = dummyID1
	s.m[dummyID1] = urlInfo
	info, err := s.rep.GetByOriginalURL(dummyURL1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), info.Owner, ownerID)
	assert.Equal(s.T(), info.OriginalURL, dummyURL1)
	assert.Equal(s.T(), info.ID, dummyID1)
}

func (s *URLRepositoryInMemorySuite) TestGetByOriginalURLItemNotFound() {
	_, err := s.rep.GetByOriginalURL(dummyURL1)
	assert.ErrorIs(s.T(), repository.ErrorItemNotFound, err)
}

func (s *URLRepositoryInMemorySuite) TestGetAllByOwner() {
	info1 := createURLInfo(dummyID1, ownerID, dummyURL1)
	info2 := createURLInfo(dummyID2, ownerID, dummyURL2)
	s.m[dummyID1] = info1
	s.m[dummyID2] = info2
	infos, err := s.rep.GetAllByOwner(ownerID)
	assert.NoError(s.T(), err)
	for _, infoItem := range infos {
		if infoItem.ID == dummyID1 {
			assert.Equal(s.T(), info1.Owner, infoItem.Owner)
			assert.Equal(s.T(), info1.OriginalURL, infoItem.OriginalURL)
		} else if infoItem.ID == dummyID2 {
			assert.Equal(s.T(), info2.Owner, infoItem.Owner)
			assert.Equal(s.T(), info2.OriginalURL, infoItem.OriginalURL)
		} else {
			assert.True(s.T(), false)
		}
	}
}

func (s *URLRepositoryInMemorySuite) TestDeleteURLsInBatch() {
	info1 := createURLInfo(dummyID1, ownerID, dummyURL1)
	info2 := createURLInfo(dummyID2, ownerID, dummyURL2)
	info3 := createURLInfo(dummyID3, ownerID, dummyURL3)
	s.m[dummyID1] = info1
	s.m[dummyID2] = info2
	s.m[dummyID3] = info3
	err := s.rep.DeleteURLsInBatch(ownerID, []string{dummyID1, dummyID2})
	assert.NoError(s.T(), err)
	info1 = s.m[dummyID1]
	info2 = s.m[dummyID2]
	info3 = s.m[dummyID3]
	assert.True(s.T(), info1.IsDeleted())
	assert.True(s.T(), info2.IsDeleted())
	assert.False(s.T(), info3.IsDeleted())
}

func createURLInfo(id, ownerID, url string) entity.ShortenedURLInfo {
	urlInfo := *entity.NewUnstoredShortenedURLInfo(ownerID, url)
	urlInfo.ID = id
	return urlInfo
}
