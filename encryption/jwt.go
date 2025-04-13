package encryption

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Header: Authorization bearer [token]
// the authstring passed is the "bearer [token]"
// This is basically splitting the string only.
// Get authorization from header, return: Bearer 'token'
func GetAuthorizationFromHeader(authstring string) (string, string) {
	// NOTE: somehow if Authorize header is supplied, though it's empty (from postman/rested)
	// it will make authstring == "null" <-- as a string. So we flag this by returning below
	if authstring == "null" {
		return "empty", "empty"
	}
	onlyToken := strings.Split(authstring, " ")
	if len(onlyToken) == 2 {
		return onlyToken[0], onlyToken[1]
	} else {
		// this is anything but 2 strings separated by space
		return "", ""
	}
}

// Get JWT Claim manually (without using the JWT middleware)
// Parse the header manually then get the JWT. This function is needed to check
// if JWT is valid but expired, then we use it to renew/extends the expiration
func GetJWTClaimMapFromTokenString(t, JWTKey string) (jwt.MapClaims, error) {
	// p := new(jwt.Parser)
	// p.SkipClaimsValidation = true
	// DebugLog.Println("GetJWTClaimFromTokenString: Start, t = ", t)
	oldToken, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKey), nil
	})
	if err != nil {
		return nil, errors.New("JWT is not valid here 2")
	}

	// this is claims in form of map[string]
	claims, ok := oldToken.Claims.(jwt.MapClaims)
	// Do we need to check if token.Valid? For now skipping this, because maybe
	// token.Valid will be false if expiration date is expired, and that's what we
	// need to check manually
	if !ok {
		return nil, errors.New("JWT is not valid 3")
	}
	return claims, nil
}

// From Authorization : Basic [This part is JWTKey]
// Format accepted is JWTKey == base64(ID:SECRET)
// Authorization Basic JWTKey
func GetClientIDSecretFromTokenString(jwtKey string) (string, string, error) {
	authByte, err := base64.StdEncoding.DecodeString(jwtKey)
	if err != nil {
		return "", "", err
	}
	// MAYBE: Shouldn't it be split first then decode?
	splitted := strings.Split(string(authByte), ":")
	if len(splitted) != 2 {
		return "", "", errors.New("basic auth token is not valid")
	}
	return splitted[0], splitted[1], nil
}

// This is needed for JWT library Parser function
// func jwtHQKeyFunc(token *jwt.Token) (interface{}, error) {
// 	return []byte(JWTKey), nil
// }
