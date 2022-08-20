package crypto

type Provider interface {
	Encrypt(data []byte) string
	Decrypt(data []byte) (string, error)
}
