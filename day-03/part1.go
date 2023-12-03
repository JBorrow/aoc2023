package main

import "fmt"
import "os"
import "log"
import "regexp"
import "bufio"
import "strconv"
//import "strings"

var number_regex = regexp.MustCompile(`\d*`)
var token_regex = regexp.MustCompile(`[^\.|\d|\n]`)

func parse_line(str string) ([]int, []bool) {
	number_of_characters := len(str)
	
	ids := make([]int, number_of_characters)
	tokens := make([]bool, number_of_characters)

	// First parse tokens, much easier.
	indicies := token_regex.FindAllStringIndex(str, -1)

	for _, index := range indicies {
		tokens[index[0]] = true
	}

	// Ids a bit more complex, requires both indicies and values
	values := number_regex.FindAllString(str, -1)
	indicies = number_regex.FindAllStringIndex(str, -1)

	for i := 0; i < len(values); i++ {
		value := values[i]
		index := indicies[i][0]

		// Loop over n, where n is string length
		for x := 0; x < len(value); x++ {
			ids[index + x], _ = strconv.Atoi(value)
		}
	}

	return ids, tokens
}

func extract_part_numbers(tokens [][]bool, values [][]int) []int {
	// Go does not have a set so we use this instead.
	m := make(map[int]bool)

	number_of_lines := len(tokens)
	number_of_columns := len(tokens[0])

	for line_index, token_line := range tokens {
		// Valid part numbers are 'around' the token.
		for column_index, token_value := range token_line {
			if (!token_value) {
				continue
			}

			start_line := max(min(number_of_lines - 1, line_index - 1), 0)
			end_line := max(min(number_of_lines - 1, line_index + 1), 0)

			start_column := max(min(number_of_columns - 1, column_index - 1), 0)
			end_column := max(min(number_of_columns - 1, column_index + 1), 0)

			for line := start_line; line <= end_line; line++ {
				for column := start_column; column <= end_column; column++ {
					this_value := values[line][column]
					if (this_value > 0) {
						m[this_value] = true
					}
				}
			}
		}
	}

	// Now have full map; extract
	MAX_PART_NUMBER := 100000

	part_numbers := make([]int, 0)

	for i := 0; i <= MAX_PART_NUMBER; i++ {
		if (m[i]) {
			part_numbers = append(part_numbers, i)
		}
	}

	return part_numbers
}



func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	raw_input := make([]string, 0)

	for scanner.Scan() {
		text := scanner.Text()

		raw_input = append(raw_input, text)
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	number_of_lines := len(raw_input)

	tokens := make([][]bool, number_of_lines)
	values := make([][]int, number_of_lines)

	for i, l := range raw_input {
		v, t := parse_line(l)

		tokens[i] = t
		values[i] = v
	}

	if DEBUG {
		fmt.Println("Raw:")
		for _, v := range raw_input {
			fmt.Println(v)
		}

		fmt.Println("Tokens:")
		for _, v := range tokens {
			fmt.Println(v)
		}

		fmt.Println("Values:")
		for _, v := range values {
			fmt.Println(v)
		}
	}

	part_numbers := extract_part_numbers(tokens, values)

	if DEBUG {
		fmt.Println("Part numbers: ")
		fmt.Println(part_numbers)
	}

	total_part_numbers := 0

	for _, v := range part_numbers {
		total_part_numbers += v
	}

	fmt.Println("Total of part numbers: ", total_part_numbers)
}
