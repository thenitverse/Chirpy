package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header found")
	}
	token := strings.TrimPrefix(authHeader, "Bearer")
	token = strings.TrimSpace(token)
	return token, nil
}
