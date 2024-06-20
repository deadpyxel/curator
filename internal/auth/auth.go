package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No authentication information found.")
	}

	headerContents := strings.Split(authHeader, " ")
	if len(headerContents) != 2 {
		return "", errors.New("Malformed auth header")
	}

	if headerContents[0] != "ApiKey" {
		return "", errors.New("Malformed auth header")
	}

	return headerContents[1], nil
}
