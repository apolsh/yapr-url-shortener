package dto

// ShortenInBatchRequestItem используется для запроса в сохранении нескольких URL за один запрос.
type ShortenInBatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortenInBatchResponseItem используется для ответа в сохранении нескольких URL за один запрос.
type ShortenInBatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// AppStatisticItem статистические данные приложения.
type AppStatisticItem struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// URLPair - пара связанных между собой оригинального URL и его укороченной версии.
type URLPair struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
