package encryption

import (
	"testing"
	"time"
)

// RFC 6238 Appendix B test vector secret ("12345678901234567890").
const rfcSecret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

func TestTOTPCodeRFC6238Vectors(t *testing.T) {
	// RFC 6238 Appendix B — SHA1, 8 digits, 30s step.
	cases := []struct {
		unix int64
		want string
	}{
		{59, "94287082"},
		{1111111109, "07081804"},
		{1111111111, "14050471"},
		{1234567890, "89005924"},
		{2000000000, "69279037"},
		{20000000000, "65353130"},
	}
	for _, c := range cases {
		got, err := TOTPCode(rfcSecret, time.Unix(c.unix, 0).UTC(), 30, 8)
		if err != nil {
			t.Fatalf("TOTPCode(%d): %v", c.unix, err)
		}
		if got != c.want {
			t.Errorf("TOTPCode(%d) = %s, want %s", c.unix, got, c.want)
		}
	}
}

func TestTOTPDefaults(t *testing.T) {
	secret, err := NewTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	code, err := TOTPCode(secret, now, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Fatalf("default digits = %d, want 6", len(code))
	}
	if !VerifyTOTP(secret, code, 1) {
		t.Fatal("verify own code failed")
	}
	if VerifyTOTP(secret, "000000", 1) {
		t.Fatal("wrong code accepted")
	}
	if VerifyTOTP("not-base32!", "000000", 1) {
		t.Fatal("invalid secret accepted")
	}
	// A code 2 steps ahead should fail at window 1 but pass at window 2.
	future, _ := TOTPCode(secret, now.Add(2*DefaultTOTPStep*time.Second), 0, 0)
	if VerifyTOTP(secret, future, 1) {
		t.Fatal("future code accepted at window 1")
	}
	if !VerifyTOTP(secret, future, 2) {
		t.Fatal("future code rejected at window 2")
	}
}

func TestTOTPSecretFormat(t *testing.T) {
	secret, err := NewTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	if len(secret) != 32 {
		t.Fatalf("secret length = %d, want 32", len(secret))
	}
	// Standard authenticator-app alphabet (base32, no padding).
	for _, r := range secret {
		if !(r >= 'A' && r <= 'Z') && !(r >= '2' && r <= '7') {
			t.Fatalf("secret contains non-base32 char %q", r)
		}
	}
}

func TestTOTPURI(t *testing.T) {
	uri := TOTPURI("JBSWY3DPEHPK3PXP", "alice", "ControlServer", 30, 6)
	want := "otpauth://totp/ControlServer:alice?secret=JBSWY3DPEHPK3PXP&issuer=ControlServer&algorithm=SHA1&digits=6&period=30"
	if uri != want {
		t.Errorf("URI = %s, want %s", uri, want)
	}
}
