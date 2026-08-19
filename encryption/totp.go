package encryption

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// TOTP implements RFC 6238 time-based one-time passwords (HMAC-SHA1).
// Generic and reusable across every medatechnology service: secrets are
// standard base32 (authenticator-app compatible) and codes verify with a
// configurable ±window of steps.
//
// Usage:
//
//	secret, _ := encryption.NewTOTPSecret()            // e.g. "JBSWY3DPEHPK3PXP"
//	code, _ := encryption.TOTPCode(secret, time.Now(), 30, 6)
//	ok := encryption.VerifyTOTP(secret, code, 1)
//	uri := encryption.TOTPURI(secret, "alice", "MyApp", 30, 6) // QR / manual entry

const (
	// DefaultTOTPStep is the standard 30-second TOTP time step.
	DefaultTOTPStep = 30
	// DefaultTOTPDigits is the standard 6-digit code length.
	DefaultTOTPDigits = 6
)

// NewTOTPSecret generates a random base32 secret (20 bytes → 32 chars, no padding).
func NewTOTPSecret() (string, error) {
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), nil
}

// TOTPCode computes the RFC 6238 code for the given secret and time.
// periodSec and digits default to 30 and 6 when <= 0.
func TOTPCode(secret string, t time.Time, periodSec, digits int) (string, error) {
	if periodSec <= 0 {
		periodSec = DefaultTOTPStep
	}
	if digits <= 0 {
		digits = DefaultTOTPDigits
	}
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
	if err != nil {
		return "", err
	}
	counter := t.Unix() / int64(periodSec)
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(counter))
	mac := hmac.New(sha1.New, key)
	mac.Write(msg[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	code := (binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff) % uint32(pow10(digits))
	return fmt.Sprintf("%0*d", digits, code), nil
}

// VerifyTOTP checks a user-provided code within ±window steps of the current
// time (defaults: 30s step, 6 digits, window 1). window <= 0 means 1.
func VerifyTOTP(secret, code string, window int) bool {
	return VerifyTOTPAt(secret, code, window, time.Now(), DefaultTOTPStep, DefaultTOTPDigits)
}

// VerifyTOTPAt is VerifyTOTP with an explicit reference time and TOTP
// parameters (used by tests and services with custom step/digits).
func VerifyTOTPAt(secret, code string, window int, at time.Time, periodSec, digits int) bool {
	if window <= 0 {
		window = 1
	}
	code = strings.TrimSpace(code)
	for i := -window; i <= window; i++ {
		want, err := TOTPCode(secret, at.Add(time.Duration(i*periodSec)*time.Second), periodSec, digits)
		if err == nil && hmac.Equal([]byte(want), []byte(code)) {
			return true
		}
	}
	return false
}

// TOTPURI builds the otpauth:// enrollment URI for QR codes / manual entry.
// periodSec/digits default to 30/6; issuer defaults to "Meda" when empty.
func TOTPURI(secret, account, issuer string, periodSec, digits int) string {
	if periodSec <= 0 {
		periodSec = DefaultTOTPStep
	}
	if digits <= 0 {
		digits = DefaultTOTPDigits
	}
	if issuer == "" {
		issuer = "Meda"
	}
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=%d&period=%d",
		issuer, account, secret, issuer, digits, periodSec)
}

func pow10(n int) int {
	p := 1
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}
