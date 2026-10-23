package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

// GenerateClusterToken generates a cryptographically signed cluster token using HMAC-SHA256.
// Format: clusterID.hex(HMAC-SHA256(clusterID, secretKey))
func GenerateClusterToken(clusterID, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(clusterID))
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s.%s", clusterID, sig)
}

// ValidateClusterToken verifies that the provided token was signed with the expected secretKey
// and matches the clusterID. Constant-time comparison prevents timing attacks.
func ValidateClusterToken(token, clusterID, secretKey string) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 || parts[0] != clusterID {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(clusterID))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(parts[1]), []byte(expectedSig)) == 1
}
