package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	UNKNOWN     = uint8(0)
	OPERATIONAL = uint8(1)
	BROKEN      = uint8(2)
)

var string_to_status = map[string]uint8{
	"?": UNKNOWN,
	".": OPERATIONAL,
	"#": BROKEN,
}

func parse_row(row string) ([]uint8, []uint8) {
	status := make([]uint8, 0)
	pattern := make([]uint8, 0)

	splits := strings.Split(strings.TrimSpace(row), " ")
	tokens := splits[0]
	patterns := strings.Split(splits[1], ",")

	for _, column := range tokens {
		status = append(status, string_to_status[string(column)])
	}

	for _, value := range patterns {
		v, _ := strconv.Atoi(value)
		pattern = append(pattern, uint8(v))
	}

	return status, pattern
}

func matches_pattern(row []uint8, pattern []uint8) bool {
	start := 0

	for i, value := range row {
		if value != OPERATIONAL {
			start = i
			break
		}
	}

	current_pattern_index := 0
	matched_in_pattern := uint8(0)
	current_pattern_to_match := pattern[current_pattern_index]

	matched := make([]bool, len(pattern))

	for i := start; i < len(row); i++ {
		if row[i] == BROKEN {
			if current_pattern_to_match > matched_in_pattern {
				matched_in_pattern += 1

				if matched_in_pattern == current_pattern_to_match {
					// We have matched the current pattern
					matched[current_pattern_index] = true
				}
			} else {
				// We have broken the current pattern
				return false
			}
		} else if row[i] == OPERATIONAL {
			// Skip to next pattern.
			if matched_in_pattern == 0 {
				continue
			}

			current_pattern_index += 1
			matched_in_pattern = 0

			if current_pattern_index >= len(pattern) {
				// We have matched all the patterns
				// Check if there are any more BROKENs
				for j := i; j < len(row); j++ {
					if row[j] == BROKEN {
						return false
					}
				}

				break
			}

			current_pattern_to_match = pattern[current_pattern_index]
		}
	}

	// Did we match all the patterns?
	for _, match := range matched {
		if !match {
			return false
		}
	}

	return true
}

func print(row []uint8) {
	for _, value := range row {
		if value == UNKNOWN {
			fmt.Print("?")
		} else if value == OPERATIONAL {
			fmt.Print(".")
		} else if value == BROKEN {
			fmt.Print("#")
		}
	}

	fmt.Println()
}

func replace_and_continue(
	row []uint8,
	index int,
	pattern []uint8,
	total_matches *int,
) {
	// Base case
	if index == len(row) {
		// Can actually check our row for matches
		if matches_pattern(row, pattern) {
			*total_matches += 1
			return
		} else {
			return
		}
	}

	// Recursive case
	if row[index] == UNKNOWN {
		// Replace with BROKEN
		new_arr := make([]uint8, len(row))
		copy(new_arr, row)
		new_arr[index] = BROKEN

		replace_and_continue(
			new_arr,
			index+1,
			pattern,
			total_matches,
		)

		// Replace with OPERATIONAL
		copy(new_arr, row)
		new_arr[index] = OPERATIONAL

		replace_and_continue(
			new_arr,
			index+1,
			pattern,
			total_matches,
		)
	} else {
		// Just continue
		replace_and_continue(
			row,
			index+1,
			pattern,
			total_matches,
		)
	}

	return
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	all_matches := 0

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		status, pattern := parse_row(text)

		var total_matches *int
		total_matches = new(int)
		*total_matches = 0
		replace_and_continue(status, 0, pattern, total_matches)

		if DEBUG {
			fmt.Println("Given: ", text)
			fmt.Println("Parsed to: ", status, pattern)
			fmt.Println("Total matches: ", *total_matches)
		}

		all_matches += *total_matches

	}

	fmt.Println("Total matches: ", all_matches)
}
