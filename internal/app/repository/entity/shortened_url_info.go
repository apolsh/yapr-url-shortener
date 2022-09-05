package entity

const (
	Active = iota
	Deleted
)

type ShortenedURLInfo struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OriginalURL string `json:"originalURL"`
	Status      int    `json:"status"`
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

func (shortenedURL *ShortenedURLInfo) GetStatus() int {
	return shortenedURL.Status
}

func (shortenedURL *ShortenedURLInfo) SetDeleted() {
	shortenedURL.Status = Deleted
}

func (shortenedURL *ShortenedURLInfo) IsDeleted() bool {
	return shortenedURL.Status == Deleted
}

func NewUnstoredShortenedURLInfo(owner, originalURL string) *ShortenedURLInfo {
	return &ShortenedURLInfo{ID: "", Owner: owner, OriginalURL: originalURL, Status: Active}
}
