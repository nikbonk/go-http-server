package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", errors.New("No bearer token in request")
	}

	return parts[1], nil
}

func GetBearerRefreshToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", errors.New("No bearer token in request")
	}

	return parts[1], nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)

	return hex.EncodeToString(token)
}
