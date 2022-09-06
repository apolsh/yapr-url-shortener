package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

type AESCryptoProvider struct {
	aesBlock cipher.Block
}

func NewAESCryptoProvider(keyString string) CryptographicProvider {
	keyHash := sha256.Sum256([]byte(keyString))

	aesBlock, err := aes.NewCipher(keyHash[:])
	if err != nil {
		panic(err)
	}
	return AESCryptoProvider{aesBlock: aesBlock}
}

func (c AESCryptoProvider) Encrypt(data []byte) string {
	encrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Encrypt(encrypted, data)
	return hex.EncodeToString(encrypted)
}

func (c AESCryptoProvider) Decrypt(data []byte) (string, error) {
	decrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Decrypt(decrypted, data)
	decryptedUUID, err := uuid.FromBytes(decrypted)
	if err != nil {
		panic(err)
	}
	return decryptedUUID.String(), nil
}
