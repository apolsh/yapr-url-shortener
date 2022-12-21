package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

// AESCryptoProvider реализует интерфейс CryptographicProvider, используемый для кодирования и декодирования
// токенов пользователя в Cookies, данная реализация использует алгоритм блочного шифрования AES на базе хэшсуммы
// алгоритмом SHA256
type AESCryptoProvider struct {
	aesBlock cipher.Block
}

// NewAESCryptoProvider в качестве аргумента принимает строку, которая будет использован как AES ключ
func NewAESCryptoProvider(keyString string) AESCryptoProvider {
	keyHash := sha256.Sum256([]byte(keyString))

	aesBlock, err := aes.NewCipher(keyHash[:])
	if err != nil {
		panic(err)
	}
	return AESCryptoProvider{aesBlock: aesBlock}
}

// Encrypt шифрует данные в base64
func (c AESCryptoProvider) Encrypt(data []byte) string {
	encrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Encrypt(encrypted, data)
	return hex.EncodeToString(encrypted)
}

// Decrypt дешифрует UUID пользователя
func (c AESCryptoProvider) Decrypt(data []byte) (string, error) {
	decrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Decrypt(decrypted, data)
	decryptedUUID, err := uuid.FromBytes(decrypted)
	if err != nil {
		panic(err)
	}
	return decryptedUUID.String(), nil
}
