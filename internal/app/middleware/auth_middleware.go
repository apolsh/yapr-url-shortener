package middleware

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/google/uuid"
)

const authCookieName = "sessionId"

type ContextKey string

var OwnerID ContextKey = "userId"

func generateNewUserIDCookie(uuid []byte, cryptoProvider crypto.CryptographicProvider) *http.Cookie {

	encryptedUserID := cryptoProvider.Encrypt(uuid)
	return &http.Cookie{Name: authCookieName, Value: encryptedUserID}
}

// AuthMiddleware middleware, который присваивает пользователям идентификаторы и
// идентифицирует пользователя по cookie sessionId, если у пользователя
// отсутствует этот cookie, то присваивает новый идентификатор и добавляет SetCookie заголовок с sessionId
func AuthMiddleware(cryptoProvider crypto.CryptographicProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFunction := func(w http.ResponseWriter, r *http.Request) {
			var userID string
			cookie, err := r.Cookie(authCookieName)
			if err != nil {
				userUUID := uuid.New()
				http.SetCookie(w, generateNewUserIDCookie(userUUID[:], cryptoProvider))
				userID = userUUID.String()
			} else {
				sessionIDBytes, err := hex.DecodeString(cookie.Value)
				if err != nil {
					userUUID := uuid.New()
					http.SetCookie(w, generateNewUserIDCookie(userUUID[:], cryptoProvider))
					userID = userUUID.String()
				} else {
					userID, err = cryptoProvider.Decrypt(sessionIDBytes)
					if err != nil {
						userUUID := uuid.New()
						http.SetCookie(w, generateNewUserIDCookie(userUUID[:], cryptoProvider))
						userID = userUUID.String()
					}
				}
			}

			ctxWithUser := context.WithValue(r.Context(), OwnerID, userID)
			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		}
		return http.HandlerFunc(handlerFunction)
	}
}
