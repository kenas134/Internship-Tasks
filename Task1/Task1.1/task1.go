package main

import (
	"fmt"
	"strings"
)

func WordFrequency(text string) map[string]int {
	freq := make(map[string]int)

	text = strings.ToLower(text)

	word := ""

	for _, ch := range text {
		// If the character is a letter or digit, add it to the current word.
		if isLetter(ch) || isDigit(ch) {
			word += string(ch)
		} else {
			// We reached punctuation or whitespace.
			if word != "" {
				freq[word]++
				word = ""
			}
		}
	}

	// Add the last word if there is one.
	if word != "" {
		freq[word]++
	}

	return freq
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func main() {
	text := "Hello, hello! Go is fun. Go, Go!"

	fmt.Println(WordFrequency(text))
}
