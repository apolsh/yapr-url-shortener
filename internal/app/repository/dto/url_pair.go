package dto

// URLPair - пара связанных между собой оригинального URL и его укороченной версии
type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
