package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Generate the `access token` and `refresh token` for the secret key.
func GenerateJWTTokenPair(hmacSecret []byte, sessionUuid string, d time.Duration) (string, string, error) {
	//
	// Generate token.
	//
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_uuid"] = sessionUuid
	claims["exp"] = time.Now().Add(d).Unix()

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", "", err
	}

	//
	// Generate refresh token.
	//
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["session_uuid"] = sessionUuid
	rtClaims["exp"] = time.Now().Add(d + time.Hour*72).Unix()

	refreshTokenString, err := refreshToken.SignedString(hmacSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// Validates either the `access token` or `refresh token` and returns either the
// `uuid` if success or error on failure.
func ProcessBearerToken(hmacSecret []byte, reqToken string) (string, error) {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sessionUuid := claims["session_uuid"].(string)
			// m["exp"] := string(claims["exp"].(float64))
			return sessionUuid, nil
		} else {
			return "", err
		}

	} else {
		return "", err
	}
	return "", nil
}
