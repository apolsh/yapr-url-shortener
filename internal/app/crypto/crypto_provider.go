//go:generate mockgen -destination=../mocks/crypto_provider_mock.go -package=mocks github.com/apolsh/yapr-url-shortener/internal/app/crypto CryptographicProvider
package crypto

// CryptographicProvider интерфейс для криптопровайдера, шифрующего и дешифрующего данные
type CryptographicProvider interface {
	Encrypt(data []byte) string
	Decrypt(data []byte) (string, error)
}
