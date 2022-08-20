package entity

type ShortenedURLInfo struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OriginalURL string `json:"originalURL"`
}

func NewShortenedURLInfo(id, owner, originalURL string) *ShortenedURLInfo {
	return &ShortenedURLInfo{ID: id, Owner: owner, OriginalURL: originalURL}
}

func (shortenedURL ShortenedURLInfo) GetOwner() string {
	return shortenedURL.Owner
}

func (shortenedURL ShortenedURLInfo) GetOriginalURL() string {
	return shortenedURL.OriginalURL
}

func (shortenedURL ShortenedURLInfo) GetID() string {
	return shortenedURL.ID
}

func (shortenedURL *ShortenedURLInfo) SetID(id string) {
	shortenedURL.ID = id
}

func NewUnstoredShortenedURLInfo(owner, originalURL string) *ShortenedURLInfo {
	return &ShortenedURLInfo{ID: "", Owner: owner, OriginalURL: originalURL}
}
