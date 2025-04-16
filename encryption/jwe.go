package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/medatechnology/goutil/medaerror"
)

// Standard JWE format
// BASE64URL(UTF8(Protected Header)).BASE64URL(Encrypted Key).BASE64URL(IV).BASE64URL(Ciphertext).BASE64URL(Authentication Tag)
// Protected Header: {"alg":"A128GCM","enc":"A128GCM"}
// Encrypted Key: Encrypted symmetric key (optional, not used in this example)
// IV: Initialization Vector (nonce)
// Ciphertext: Encrypted payload
// Authentication Tag: Used for integrity check (not used in this example)
// NOTE: This is a simplified example. In a real-world application, you should handle errors and edge cases properly.
// Create JWE token
// The payload is the data you want to encrypt
// The key is the symmetric key used for encryption
func CreateJWE(payload []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", medaerror.Errorf("creating cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", medaerror.Errorf("creating GCM: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", medaerror.Errorf("generating nonce: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, payload, nil)

	header := map[string]string{"alg": "A128GCM", "enc": "A128GCM"}
	headerBytes, _ := json.Marshal(header)
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedCiphertext := base64.RawURLEncoding.EncodeToString(ciphertext)
	encodedNonce := base64.RawURLEncoding.EncodeToString(nonce)

	jwe := strings.Join([]string{encodedHeader, encodedNonce, encodedCiphertext}, ".")
	return jwe, nil
}

// Parse JWE token
// The JWE string is the token you want to decrypt
// The key is the symmetric key used for decryption
// The function returns the decrypted payload
// Example usage:
// jweString := "eyJhbGciOiJBMjU2R0NNIiwiZW5jIjoiQTEyOEdDTSJ9.h79q
func ParseJWE(jweString string, key []byte) ([]byte, error) {
	parts := strings.Split(jweString, ".")
	if len(parts) != 3 {
		return nil, medaerror.Errorf("invalid JWE format: %d parts", len(parts))
	}

	encodedHeader := parts[0]
	encodedNonce := parts[1]
	encodedCiphertext := parts[2]

	decodedHeader, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, medaerror.Errorf("decoding header: %w", err)
	}

	var header map[string]string
	err = json.Unmarshal(decodedHeader, &header)
	if err != nil {
		return nil, medaerror.Errorf("unmarshaling header: %w", err)
	}

	decodedCiphertext, err := base64.RawURLEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, medaerror.Errorf("decoding ciphertext: %w", err)
	}

	decodedNonce, err := base64.RawURLEncoding.DecodeString(encodedNonce)
	if err != nil {
		return nil, medaerror.Errorf("decoding nonce: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, medaerror.Errorf("creating cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, medaerror.Errorf("creating GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, decodedNonce, decodedCiphertext, nil)
	if err != nil {
		return nil, medaerror.Errorf("decrypting: %w", err)
	}

	return plaintext, nil
}

// If the cypertext or basically the payload is json of map[string]string then this has the unmarshall
// and return the map[string]string
// NOTE: later if needed use MapToStruct from utils
func ParseJWEToMap(jweString string, key []byte) (map[string]string, error) {
	plainText, err := ParseJWE(jweString, key)
	if err != nil {
		return nil, err
	}
	var tmpMap map[string]string
	err = json.Unmarshal(plainText, &tmpMap)
	if err != nil {
		return nil, medaerror.Errorf("unmarshaling jwe payload: %w", err)
	}
	return tmpMap, nil
}

// IDEA: JWE inside JWE??
// Have handler /createjwe <-- where it needs to be called with valid JWE already (caller/client JWE)
// then pass the payload in json like "{ mydata: mycontent, key: myKey}"
// the API will output those payload into JWE. Then sell this as JWE API service?

// Example for the front-end JS/TS processing JWE

// import * as jose from 'jose'
// async function decryptJWE(jwe, key) {
//   try {
//     const { payload, protectedHeader } = await jose.compactDecrypt(jwe, new Uint8Array(key));
//     console.log('Decrypted payload:', new TextDecoder().decode(payload));
//     console.log('Protected Header', protectedHeader);
//   } catch (error) {
//     console.error('Error decrypting JWE:', error)
//   }
// }

// Example usage (assuming 'jwe' is the string from the Go code)
// const jweString = "eyJhbGciOiJBMjU2R0NNIiwiZW5jIjoiQTEyOEdDTSJ9.h79q794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o794o7
