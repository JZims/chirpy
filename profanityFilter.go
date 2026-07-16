package main

import (
	"slices"
	"strings"
)

func profanityFilter(text string) string {

	badWords := make([]string, 0, 3)
	badWords = append(badWords, "kerfuffle", "sharbert", "fornax")

	words := strings.Split(text, " ")

	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
