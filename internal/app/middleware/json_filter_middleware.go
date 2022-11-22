package middleware

import (
	"net/http"
	"strings"
)

var allowedTypes = []string{"application/json", "application/x-gzip"}

const invalidContentTypeError = "invalid content type"

func JsonFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualContentType := r.Header.Get("Content-Type")

		haveRightContentType := false

		for _, allowed := range allowedTypes {
			if strings.Contains(actualContentType, allowed) {
				haveRightContentType = true
				break
			}
		}
		if !haveRightContentType {
			http.Error(w, invalidContentTypeError, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
