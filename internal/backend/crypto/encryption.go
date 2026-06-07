package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	KeySize  = 32
	SaltSize = 16

	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
)

var ErrInvalidKey = errors.New("key must be 32 bytes")

// GenerateVaultKey creates a random 256-bit key used to encrypt vault entries.
func GenerateVaultKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateSalt creates a random salt for Argon2id master-key derivation.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveMasterKey derives a 256-bit key from a user-specified master password.
func DeriveMasterKey(masterPassword string, salt []byte) []byte {
	return argon2.IDKey([]byte(masterPassword), salt, argonTime, argonMemory, argonThreads, KeySize)
}

// Encrypt encrypts plaintext using AES-256-GCM and a 32-byte key.
func Encrypt(plaintext string, key []byte) (string, error) {
	ciphertext, err := EncryptBytes([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES-256-GCM ciphertext using a 32-byte key.
func Decrypt(ciphertext string, key []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := DecryptBytes(decoded, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptBytes encrypts bytes using AES-256-GCM and prefixes the nonce.
func EncryptBytes(plaintext, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptBytes decrypts bytes produced by EncryptBytes.
func DecryptBytes(ciphertext, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, encryptedText := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, encryptedText, nil)
}

// WrapVaultKey encrypts a vault key with a master key.
func WrapVaultKey(vaultKey, masterKey []byte) (string, error) {
	wrapped, err := EncryptBytes(vaultKey, masterKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(wrapped), nil
}

// UnwrapVaultKey decrypts a vault key with a master key.
func UnwrapVaultKey(wrappedVaultKey string, masterKey []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(wrappedVaultKey)
	if err != nil {
		return nil, err
	}
	return DecryptBytes(decoded, masterKey)
}
