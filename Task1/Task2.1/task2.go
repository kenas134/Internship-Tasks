package main

import (
	"fmt"
	"strings"
)

func ispal(text string) bool {
	text = strings.ToLower(text)
	word := ""

	for _, ch := range text {
		if isLetter(ch) {
			word += string(ch)
		}
	}

	i, j := 0, len(word)-1
	for i < j {
		if word[i] != word[j] {
			return false
		}
		i += 1
		j -= 1
	}
	return true

}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z')
}

func main() {

	text := "abcddcba"
	text2 := "asjhbajhsbjh"
	fmt.Println(ispal(text))
	fmt.Println(ispal(text2))

}
