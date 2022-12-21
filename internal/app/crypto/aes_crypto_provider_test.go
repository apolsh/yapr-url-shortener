package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	aesKey       = "secret_key"
	testUUID     = "96b2af59-c817-4751-8ac2-aa1a9d530948"
	encryptedB64 = "e9391d2972e90322f31f2dbef0a74a67"
)

type AESCryptoProviderTestSuite struct {
	suite.Suite
	provider CryptographicProvider
}

func (s *AESCryptoProviderTestSuite) SetupSuite() {
	provider := NewAESCryptoProvider(aesKey)
	s.provider = provider
}

func TestAESCryptoProviderTestSuite(t *testing.T) {
	suite.Run(t, new(AESCryptoProviderTestSuite))
}

func (s *AESCryptoProviderTestSuite) TestEncrypt() {
	userUUID, _ := uuid.Parse(testUUID)
	encrypt := s.provider.Encrypt(userUUID[:])
	assert.Equal(s.T(), encryptedB64, encrypt)
}

func (s *AESCryptoProviderTestSuite) TestDecrypt() {
	sessionIDBytes, _ := hex.DecodeString(encryptedB64)
	decrypted, _ := s.provider.Decrypt(sessionIDBytes)
	assert.Equal(s.T(), testUUID, decrypted)
}
