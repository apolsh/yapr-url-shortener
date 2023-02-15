package entity

import (
	"fmt"

	"github.com/apolsh/yapr-url-shortener/internal/app/repository/dto"
)

const (
	active = iota
	deleted
)

// ShortenedURLInfo - представляет собой хранимый объект информации о сохраненном URL
type ShortenedURLInfo struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OriginalURL string `json:"originalURL"`
	Status      int    `json:"status"`
}

// ToURLPair ShortenedURLInfo трансформирует в dto.URLPair
func (s *ShortenedURLInfo) ToURLPair(hostURL string) dto.URLPair {
	return dto.URLPair{OriginalURL: s.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", hostURL, s.ID)}
}

// GetOwner геттер
func (s *ShortenedURLInfo) GetOwner() string {
	return s.Owner
}

// GetOriginalURL геттер
func (s *ShortenedURLInfo) GetOriginalURL() string {
	return s.OriginalURL
}

// GetID геттер
func (s *ShortenedURLInfo) GetID() string {
	return s.ID
}

// SetID сеттер
func (s *ShortenedURLInfo) SetID(id string) {
	s.ID = id
}

// GetStatus геттер
func (s *ShortenedURLInfo) GetStatus() int {
	return s.Status
}

// SetDeleted сеттер
func (s *ShortenedURLInfo) SetDeleted() {
	s.Status = deleted
}

// IsDeleted возвращает статус URL == deleted
func (s *ShortenedURLInfo) IsDeleted() bool {
	return s.Status == deleted
}

// NewUnstoredShortenedURLInfo конструктор
func NewUnstoredShortenedURLInfo(owner, originalURL string) *ShortenedURLInfo {
	return &ShortenedURLInfo{ID: "", Owner: owner, OriginalURL: originalURL, Status: active}
}
