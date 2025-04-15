package encryption

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/lithammer/shortuuid/v4"
)

const (
	MAX_OTP_DIGIT = 10 // max number of otp digit in total (set * digit)
	MIN_OTP_DIGIT = 2  // min number of otp digit per set

	DEFAULT_OTP_DIGIT     = 2 // default otp digit per set
	DEFAULT_OTP_SET       = 2 // default otp set
	DEFAULT_OTP_SEPARATOR = "-"
	DEFAULT_OTP_ERROR     = "7" // use if randomizer is error, then random number is always this const

	DEFAULT_TOKEN_ITERATION = 5 // iterate randomtoken to get long random string. Used in magic link?
)

// GenerateSecureRandomNumber generates a secure random number of the specified number of digits.
func GenerateSecureRandomNumber(numLen int) (string, error) {
	const characters = "0123456789"
	result := make([]byte, numLen)

	for i := 0; i < numLen; i++ {
		// Generate a secure random index
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		result[i] = characters[randIndex.Int64()]
	}

	return string(result), nil
}

func GenerateDefaultOTP() string {
	return GenerateOTP(0, 0, "")
}

// Get random number for OTP (6 digit format with dash)
// Generate random number in form of string usually used for OTP
// The digit is number of digit, set is how many sets are there
// Ex: 57-03   ==> digit=2, set=2
// .   820-587 ==> digit=3, set=2
func GenerateOTP(digit, set int, separator string) string {
	// Generate new seed everytime we call the random function.
	// NOTE: might not need to do this, I think from go 1.20 the seed is automatically
	//       set to be random at startup, but can read the documentation more
	// math_rand.New(math_rand.NewSource(time.Now().UnixNano()))

	// set digit constraint
	if digit < MIN_OTP_DIGIT {
		digit = MIN_OTP_DIGIT
	} else if digit > MAX_OTP_DIGIT {
		digit = MAX_OTP_DIGIT
	}

	// calculate so amount of digit X set is not more than MAX_OTP_DIGIT, adjust set accordingly
	// set cannot be less than 1
	if set < 1 {
		set = DEFAULT_OTP_SET
	}
	if (set * digit) > MAX_OTP_DIGIT {
		// ie digit = 3 , set = 4 == 12 digit, MAX=10 then set = 10/digit
		set = MAX_OTP_DIGIT / digit
		if set < 1 {
			set = 1
		}
	}

	if separator == "" {
		separator = DEFAULT_OTP_SEPARATOR
	}

	var otp []string
	// generate per set
	for i := 0; i < set; i++ {
		str, err := GenerateSecureRandomNumber(digit)
		if err != nil {
			str = strings.Repeat(DEFAULT_OTP_ERROR, digit)
		}
		otp = append(otp, str)
	}

	return strings.Join(otp, DEFAULT_OTP_SEPARATOR)
}

// Generate just random token which is essentially a short-uuid, 22 characters length (this golang implementation)
// But actually short-uuid can be different length on different implementation or different programming language
// While standard UUID format is always 36 characters long.
// Also this is base57 encoding which exclude characters that is simmilar like 0 and O and l and I
func NewRandomToken() string {
	return shortuuid.New() // KwSysDpxcBU9FNhGkn2dCf
}

// Concatenante NewRandomToken (which is short-uuid) x number of times.  This is to be used
// in maybe magic link or public link which takes longer string.
func NewRandomTokenIterate(x int) string {
	if x < 1 {
		x = DEFAULT_TOKEN_ITERATION
	}
	var t string
	for i := 0; i < x; i++ {
		t += NewRandomToken()
	}
	return t
}
