package encryption

import "testing"

func TestClusterToken(t *testing.T) {
	clusterID := "cluster-prod-1"
	secretKey := "super-secret-key-12345"

	token := GenerateClusterToken(clusterID, secretKey)
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	if !ValidateClusterToken(token, clusterID, secretKey) {
		t.Fatal("expected token to be valid")
	}

	// Invalid cluster ID
	if ValidateClusterToken(token, "other-cluster", secretKey) {
		t.Fatal("expected token to be invalid for wrong clusterID")
	}

	// Invalid secret key
	if ValidateClusterToken(token, clusterID, "wrong-key") {
		t.Fatal("expected token to be invalid for wrong key")
	}

	// Malformed token
	if ValidateClusterToken("invalidtokenformat", clusterID, secretKey) {
		t.Fatal("expected malformed token to be invalid")
	}
}
