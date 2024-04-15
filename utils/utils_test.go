package utils

import (
	"fmt"
	"testing"
	"unicode"
)

func TestSanitizeWords(t *testing.T) {
	words := []string{
		"Can't",
		"Don't",
		"Isn't",
		"Won't",
		"I'm",
		"You'r",
		"They'r",
		"It's",
		"He's",
		"She's",
		"Cafe360",
		"Web2.0",
		"Room101",
		"3DPrinter",
		"IPv6",
	}

	for _, wd := range words {
		for _, r := range sanitizeWord(wd) {
			if !unicode.IsLetter(r) {
				t.Fatalf("Invalid Rune: %v\n", r)
			}
		}
		fmt.Println(sanitizeWord(wd))
	}

}
