package profanity

import "strings"

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func ContainsProfanity(body string) (bool, string) {
	newBody := []string{}
	badWordCount := 0

	for _, word := range strings.Split(body, " ") {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			newBody = append(newBody, filterProfanity(word))
			badWordCount++
		} else {
			newBody = append(newBody, word)
		}
	}
	return badWordCount > 0, strings.Join(newBody, " ")
}

func filterProfanity(body string) string {
	return strings.Repeat("*", 4) // magic number, could instead also be replaced by len(body)
}
