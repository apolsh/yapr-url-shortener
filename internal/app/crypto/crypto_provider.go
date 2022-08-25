package crypto

type CryptographicProvider interface {
	Encrypt(data []byte) string
	Decrypt(data []byte) (string, error)
}
