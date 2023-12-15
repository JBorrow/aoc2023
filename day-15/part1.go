package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Token struct {
	base string
	hash uint8
}

func calculate_hash(str string) uint8 {
	hash := uint8(0)

	for _, char := range str {
		numerical_value := uint8(char)

		hash += uint8(numerical_value)
		hash *= uint8(17)
	}

	return hash
}

func parse_line(str string) []Token {
	tokens := []Token{}

	split := strings.Split(str, ",")

	for _, char := range split {
		tokens = append(tokens, Token{string(char), calculate_hash(string(char))})
	}

	return tokens
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	tokens := make([]Token, 0)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if DEBUG {
			fmt.Println("Input: ", text)
		}

		new_tokens := parse_line(text)

		if DEBUG {
			fmt.Println("Tokens: ", new_tokens)
		}

		tokens = append(tokens, new_tokens...)
	}

	total := 0

	for _, token := range tokens {
		total += int(token.hash)
	}

	fmt.Println("Total score: ", total)
}
