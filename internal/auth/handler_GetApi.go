package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apikey := headers.Get("Authorization")
	if apikey == "" {
		return "", errors.New("no authorization header found")
	}
	token := strings.TrimPrefix(apikey, "ApiKey ")
	token = strings.TrimSpace(token)
	return token, nil
}
