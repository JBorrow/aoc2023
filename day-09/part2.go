package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var extract_sequence_regex = regexp.MustCompile(`-?[0-9]*`)

func extract_sequence(to_parse string) []int {
	matches := extract_sequence_regex.FindAllString(strings.TrimSpace(to_parse), -1)

	if len(matches) == 0 {
		log.Fatal("Could not parse: ", to_parse)
	}

	sequence := make([]int, len(matches))

	for i, match := range matches {
		sequence[i], _ = strconv.Atoi(match)
	}

	return sequence
}

func find_differences(sequence []int) []int {
	differences := make([]int, len(sequence)-1)

	for i := 0; i < len(sequence)-1; i++ {
		differences[i] = sequence[i+1] - sequence[i]
	}

	return differences
}

func are_all_values_equal(values []int) bool {
	for i := 1; i < len(values); i++ {
		if values[i] != values[0] {
			return false
		}
	}

	return true
}

func are_all_values_zero(values []int) bool {
	for i := 0; i < len(values); i++ {
		if values[i] != 0 {
			return false
		}
	}

	return true
}

func find_next_value_in_sequence(values []int) int {
	fmt.Println("Finding next value in sequence: ", values)
	if are_all_values_zero(values) {
		return 0
	}

	differences := find_differences(values)

	return values[0] - find_next_value_in_sequence(differences)
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	sequences := make([]string, 0)
	next_values := make([]int, 0)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		next_value := find_next_value_in_sequence(extract_sequence(text))

		sequences = append(sequences, text)
		next_values = append(next_values, next_value)

		if DEBUG {
			fmt.Println("Given: ", text)
			fmt.Println("Next value: ", next_value)
		}
	}

	total_next_values := 0

	for _, v := range next_values {
		total_next_values += v
	}

	fmt.Println("Total of next values: ", total_next_values)

}
