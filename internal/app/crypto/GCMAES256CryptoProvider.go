package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

type CCMAES256CryptoProvider struct {
	aesBlock cipher.Block
}

func NewCCMAES256CryptoProvider(keyString string) Provider {
	keyHash := sha256.Sum256([]byte(keyString))

	aesBlock, err := aes.NewCipher(keyHash[:])
	if err != nil {
		panic(err)
	}
	return CCMAES256CryptoProvider{aesBlock: aesBlock}
}

func (c CCMAES256CryptoProvider) Encrypt(data []byte) string {
	encrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Encrypt(encrypted, data)
	return hex.EncodeToString(encrypted)
}

func (c CCMAES256CryptoProvider) Decrypt(data []byte) (string, error) {
	decrypted := make([]byte, aes.BlockSize)
	c.aesBlock.Decrypt(decrypted, data)
	decryptedUUID, err := uuid.FromBytes(decrypted)
	if err != nil {
		panic(err)
	}
	return decryptedUUID.String(), nil
}
