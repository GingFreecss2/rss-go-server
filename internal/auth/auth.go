package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("No Auth info found in headers")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("Incorrect auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Incorrect first part of auth header")
	}

	return vals[1], nil
}
