package middleware

import (
	"context"
	"encoding/hex"
	"github.com/apolsh/yapr-url-shortener/internal/app/crypto"
	"github.com/google/uuid"
	"net/http"
)

const authCookieName = "sessionId"

type ContextKey = string

var OwnerID ContextKey = "userId"

func generateNewUserIDCookie(uuid []byte, cryptoProvider crypto.Provider) *http.Cookie {

	encryptedUserID := cryptoProvider.Encrypt(uuid)
	return &http.Cookie{Name: authCookieName, Value: encryptedUserID}
}

func AuthMiddleware(cryptoProvider crypto.Provider) func(http.Handler) http.Handler {
	//TODO: rewrite without else ?
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

	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	cookie, err := r.Cookie(authCookieName)
	//	if err != nil {
	//		if errors.Is(err, http.ErrNoCookie) {
	//			userID := uuid.New()
	//			userUUID := sha256.Sum256(userID[:])
	//
	//			src := []byte("Этюд в розовых тонах") // данные, которые хотим зашифровать
	//
	//			// будем использовать AES256, создав ключ длиной 32 байта
	//			key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//				return
	//			}
	//
	//			aesblock, err := aes.NewCipher(key)
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//				return
	//			}
	//
	//			aesgcm, err := cipher.NewGCM(aesblock)
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//				return
	//			}
	//
	//			// создаём вектор инициализации
	//			nonce, err := generateRandom(aesgcm.NonceSize())
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//				return
	//			}
	//
	//			dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
	//			fmt.Printf("encrypted: %x\n", dst)
	//
	//			src2, err := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем
	//			if err != nil {
	//				fmt.Printf("error: %v\n", err)
	//				return
	//			}
	//			fmt.Printf("decrypted: %s\n", src2)
	//
	//		}
	//	}
	//
	//	//if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
	//	//	next.ServeHTTP(w, r)
	//	//	return
	//	//}
	//	//
	//	//gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	//	//if err != nil {
	//	//	http.Error(w, "Error while body reading", http.StatusInternalServerError)
	//	//}
	//	//defer gz.Close()
	//	//
	//	//w.Header().Set("Content-Encoding", "gzip")
	//	sessionIdCookie := &http.Cookie{Name: "sessionId", Value: "abc"}
	//	http.SetCookie(w, sessionIdCookie)
	//
	//	ctxWithUser := context.WithValue(r.Context(), OwnerID, "user")
	//	rWithUser := r.WithContext(ctxWithUser)
	//	next.ServeHTTP(w, rWithUser)
	//}

}
