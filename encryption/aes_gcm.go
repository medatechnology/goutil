package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const EncryptedPrefix = "enc:v1:"

// DeriveKey derives a deterministic 32-byte AES key from a project secret and salt.
func DeriveKey(secret, salt string) []byte {
	h := sha256.New()
	h.Write([]byte(secret))
	h.Write([]byte(":"))
	h.Write([]byte(salt))
	return h.Sum(nil)
}

// EncryptField encrypts plaintext using AES-256-GCM and returns standard envelope "enc:v1:<base64>".
func EncryptField(plaintext string, key []byte) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	if len(key) != 32 {
		return "", errors.New("key must be exactly 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce read error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return EncryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptField decrypts an "enc:v1:<base64>" envelope or returns the original string if not encrypted.
func DecryptField(cipherEnvelope string, key []byte) (string, error) {
	if cipherEnvelope == "" || !strings.HasPrefix(cipherEnvelope, EncryptedPrefix) {
		return cipherEnvelope, nil
	}
	if len(key) != 32 {
		return "", errors.New("key must be exactly 32 bytes for AES-256")
	}

	rawB64 := strings.TrimPrefix(cipherEnvelope, EncryptedPrefix)
	ciphertext, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext shorter than nonce size")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt error: %w", err)
	}

	return string(plaintext), nil
}
