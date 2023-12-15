package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	REMOVE  = "-"
	REPLACE = "="
)

type Token struct {
	base         string
	hash         uint8
	label        string
	instruction  string
	focal_length int
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
		if strings.Contains(char, REMOVE) {
			instruction := strings.Split(char, REMOVE)[0]
			tokens = append(tokens, Token{
				base:         char,
				hash:         calculate_hash(instruction),
				label:        instruction,
				focal_length: -1,
				instruction:  REMOVE,
			})
		} else if strings.Contains(char, REPLACE) {
			token_split := strings.Split(char, REPLACE)
			focal_length, _ := strconv.Atoi(token_split[1])
			tokens = append(tokens, Token{
				base:         char,
				hash:         calculate_hash(token_split[0]),
				label:        token_split[0],
				focal_length: focal_length,
				instruction:  REPLACE,
			})
		} else {
			log.Fatal("Invalid token: ", char)
		}
	}

	return tokens
}

func print_box(box map[uint8][]Token) {
	for hash, tokens := range box {
		fmt.Print(hash, ": [")
		for _, token := range tokens {
			fmt.Print(token.base, " ")
		}
		fmt.Println("]")
	}
}

func focusing_power(tokens []Token) int {
	if len(tokens) == 0 {
		return 0
	}

	box_factor := (int(tokens[0].hash) + 1)

	total_focusing_power := 0

	for i, token := range tokens {
		total_focusing_power += box_factor * (i + 1) * token.focal_length
	}

	return total_focusing_power
}

func tokens_to_box(tokens []Token) map[uint8][]Token {
	box := make(map[uint8][]Token)

	for _, token := range tokens {
		if token.instruction == REPLACE {
			// Current box state
			container := box[token.hash]

			if len(container) == 0 {
				// If box is empty, add the token
				box[token.hash] = []Token{token}
			} else {
				// Check if we have a token with the same label
				found := false
				for i, box_token := range container {
					if box_token.label == token.label {
						// If we do, replace it
						container[i] = token
						found = true
					}
				}

				if !found {
					// If we don't, append it
					box[token.hash] = append(container, token)
				}
			}
		} else if token.instruction == REMOVE {
			if len(box[token.hash]) > 0 {
				// If the box is not empty, remove the token
				container := box[token.hash]

				for i, box_token := range container {
					if box_token.label == token.label {
						box[token.hash] = append(container[:i], container[i+1:]...)
					}
				}
			}
		} else {
			log.Fatal("Invalid instruction: ", token.instruction)
		}
		fmt.Println("Current token:", token.base)
		fmt.Println("Current box:")
		print_box(box)

	}

	return box
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

	box := tokens_to_box(tokens)

	if DEBUG {
		fmt.Println("Box: ", box)
	}

	total_power := 0

	for _, tokens := range box {
		total_power += focusing_power(tokens)
	}

	fmt.Println("Total focusing power:", total_power)
}
