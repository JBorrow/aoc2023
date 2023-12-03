package main

import "fmt"
import "os"
import "log"
import "regexp"
import "bufio"
import "strconv"
import "slices"

//import "strings"

var number_regex = regexp.MustCompile(`\d*`)
var token_regex = regexp.MustCompile(`\*`)

func parse_line(str string, mapping map[int]int, current_uid int) ([]int, []bool, map[int]int, int) {
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
		// Use a UID not a part number because they may appear twice!
		current_uid++
		converted_value, _ := strconv.Atoi(value)
		mapping[current_uid] = converted_value

		// Loop over n, where n is string length
		for x := 0; x < len(value); x++ {
			ids[index+x] = current_uid
		}
	}

	return ids, tokens, mapping, current_uid
}

func extract_gear_ratios(tokens [][]bool, values [][]int, mapping map[int]int) int {
	total_gear_ratios := 0

	number_of_lines := len(tokens)
	number_of_columns := len(tokens[0])

	for line_index, token_line := range tokens {
		// Valid part numbers are 'around' the token.
		for column_index, token_value := range token_line {
			if !token_value {
				continue
			}

			these_mappings := make([]int, 0)

			start_line := max(min(number_of_lines-1, line_index-1), 0)
			end_line := max(min(number_of_lines-1, line_index+1), 0)

			start_column := max(min(number_of_columns-1, column_index-1), 0)
			end_column := max(min(number_of_columns-1, column_index+1), 0)

			for line := start_line; line <= end_line; line++ {
				for column := start_column; column <= end_column; column++ {
					this_value := values[line][column]
					if this_value > 0 {
						if !(slices.Contains(these_mappings, this_value)) {
							these_mappings = append(these_mappings, this_value)
						}
					}
				}
			}

			if len(these_mappings) == 2 {
				total_gear_ratios += mapping[these_mappings[0]] * mapping[these_mappings[1]]
			}
		}
	}

	return total_gear_ratios
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
	mapping := make(map[int]int)
	current_uid := 1

	for i, l := range raw_input {
		values[i], tokens[i], mapping, current_uid = parse_line(l, mapping, current_uid)
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

	gear_ratio := extract_gear_ratios(tokens, values, mapping)

	fmt.Println("Gear ratio: ", gear_ratio)
}
