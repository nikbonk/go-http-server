package json

import (
	"errors"

	profanityCheck "github.com/nikbonk/go-http-server/internal/profanityCheck"
)

func ValidaChirp(body string) (bool, string, error) {
	const maxChars = 140
	if len(body) > maxChars {
		return false, "", errors.New("Max character limit (140) exceeded")
	}

	cleanedBody := profanityCheck.ContainsProfanity(body)
	return true, cleanedBody, nil
}
