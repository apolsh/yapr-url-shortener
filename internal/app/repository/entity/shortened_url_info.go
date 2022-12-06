package entity

import (
	"fmt"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
)

const (
	Active = iota
	Deleted
)

// ShortenedURLInfo - представляет собой хранимый объект информации о сохраненном URL
type ShortenedURLInfo struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OriginalURL string `json:"originalURL"`
	Status      int    `json:"status"`
}

func (s *ShortenedURLInfo) ToURLPair(hostURL string) dto.URLPair {
	return dto.URLPair{OriginalURL: s.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", hostURL, s.ID)}
}

func (s *ShortenedURLInfo) GetOwner() string {
	return s.Owner
}

func (s *ShortenedURLInfo) GetOriginalURL() string {
	return s.OriginalURL
}

func (s *ShortenedURLInfo) GetID() string {
	return s.ID
}

func (s *ShortenedURLInfo) SetID(id string) {
	s.ID = id
}

func (s *ShortenedURLInfo) GetStatus() int {
	return s.Status
}

func (s *ShortenedURLInfo) SetDeleted() {
	s.Status = Deleted
}

func (s *ShortenedURLInfo) IsDeleted() bool {
	return s.Status == Deleted
}

func NewUnstoredShortenedURLInfo(owner, originalURL string) *ShortenedURLInfo {
	return &ShortenedURLInfo{ID: "", Owner: owner, OriginalURL: originalURL, Status: Active}
}
