package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/scrypt"
)

const (
	DEFAULT_HASH_LENGTH = 32
	// DEFAULT_CPU = 32768
	DEFAULT_CPU = 8192 // CPU cost, usually is 32768, higher takes more resources
	DEFAULT_R   = 8    // Repetition , higher takes more resources, default = 8
	DEFAULT_P   = 1    // Permutation?, higher takes more resources, default = 1
)

// func MD5Hash(password string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(password))
// 	hashBytes := hasher.Sum(nil)
// 	// Mengonversi hasil hash menjadi string heksadesimal
// 	hashStr := hex.EncodeToString(hashBytes)
// 	return hashStr
// }

// MD5Hash generates a 32-character MD5 hash of the input string.
// Usage: hash := MD5Hash("example")
// Output: "1a79a4d60de6718e8e5b326e338ae533"
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// func Init() {
// 	initRandom()
// }

// Need to do this to get a better random generator, called once with Init function
// func initRandom() {
// 	var b [8]byte
// 	_, err := rand.Read(b[:])
// 	if err != nil {
// 		panic("cannot seed math/rand package with cryptographically secure random number generator")
// 	}
// 	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
// }

// CreateHash generates a 128-bit/32-character MD5 hash, useful for AES key generation.
// Usage: hash := CreateHash("my-secret-key")
// Output: "5ebe2294ecd0e0f08eab7690d2a6ee69"
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt Meda password
// Password is in variable pin
// key1 and key2 is the salt. Usually key1 = signature and
// key2 = the additional randomness but that cannot change, usually is the created_at date string
// HashPin hashes a PIN using scrypt with salts (key1 and key2).
// Usage: hash, err := HashPin("1234", "signature", "2025-04-16")
// Output: Base64-encoded hash (e.g., "c29tZS1oYXNoLXZhbHVl")
func HashPin(pin, key1, key2 string) (string, error) {
	salt := []byte(pin)

	if key1 != "" && key2 != "" {
		passwordSalt := CreateHash(key1 + key2)
		dk, err := scrypt.Key([]byte(passwordSalt), salt, DEFAULT_CPU, DEFAULT_R, DEFAULT_P, DEFAULT_HASH_LENGTH)

		if err != nil {
			// utils.LogError( "hashpin", "error on scrypt THIS IS BAD", err)
			return "", err
		}
		// fmt.Println(base64.URLEncoding.EncodeToString(dk))
		return base64.URLEncoding.EncodeToString(dk), nil
	}
	return "", errors.New("error on createHash: empty salt")
}

// SHA256 generates a 64-character SHA-256 hash of the input string.
// Usage: hash := SHA256("example")
// Output: "50d858e8e8c1b6f1c8b8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8"
func SHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func EncryptWithKey(data, key string) (string, error) {
	encryptionKey := []byte(key)
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt login payload
// DecryptWithKey decrypts a Base64-encoded ciphertext using AES in CFB mode with the provided key.
// Usage: decrypted, err := DecryptWithKey(encrypted, "my32byteencryptionkey!")
// Output: "my secret data"
func DecryptWithKey(data, key string) (string, error) {
	encryptionKey := []byte(key)
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// ===== Below is for PGP encryption

// PGPGenerateKey generates a PGP key pair (public and private keys).
// Usage: pubKey, privKey, err := PGPGenerateKey("John Doe", "Comment", "john@example.com")
// Output: Public and private keys as strings
func PGPGenerateKey(name, comment, email string) (string, string, error) {
	config := &packet.Config{
		DefaultHash:            0, // Use default hash algorithm
		DefaultCipher:          0, // Use default cipher
		DefaultCompressionAlgo: 0, // Use default compression
		RSABits:                4096,
		Rand:                   rand.Reader,
	}

	entity, err := openpgp.NewEntity(name, comment, email, config)
	if err != nil {
		return "", "", err
	}

	// Serialize public key
	pubKeyBuf := new(bytes.Buffer)
	err = entity.Serialize(pubKeyBuf)
	if err != nil {
		return "", "", err
	}
	publicKey := pubKeyBuf.String()

	// Serialize private key
	privKeyBuf := new(bytes.Buffer)
	err = entity.SerializePrivate(privKeyBuf, nil)
	if err != nil {
		return "", "", err
	}
	privateKey := privKeyBuf.String()

	return publicKey, privateKey, nil
}

// PGPEncrypt encrypts a message using the recipient's PGP public key.
// Usage: encrypted, err := PGPEncrypt(pubKey, "Hello, World!")
// Output: ASCII-armored encrypted message
func PGPEncrypt(pubKey, message string) (string, error) {
	// Decode the recipient's public key
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(pubKey))
	if err != nil {
		return "", err
	}

	// Create a buffer to store the encrypted message
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}

	// Write the message to the encrypted buffer
	_, err = w.Write([]byte(message))
	if err != nil {
		return "", err
	}

	// Close the writer to finalize the encryption
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode the encrypted message in ASCII armor
	encryptedBuf := new(bytes.Buffer)
	armorWriter, err := armor.Encode(encryptedBuf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	_, err = armorWriter.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	armorWriter.Close()
	return encryptedBuf.String(), nil
}

// PGPDecrypt decrypts an ASCII-armored encrypted message using the recipient's PGP private key.
// Usage: decrypted, err := PGPDecrypt(privKey, encryptedMessage)
// Output: "Hello, World!"
func PGPDecrypt(privKey, encryptedMessage string) (string, error) {
	// Decode the private key
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(privKey))
	if err != nil {
		return "", err
	}

	// Decode the encrypted message
	decoded, err := armor.Decode(bytes.NewBufferString(encryptedMessage))
	if err != nil {
		return "", err
	}

	// Decrypt the message
	md, err := openpgp.ReadMessage(decoded.Body, entityList, nil, nil)
	if err != nil {
		return "", err
	}

	// Read the decrypted message
	decryptedBytes, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
