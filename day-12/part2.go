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
	pattern_string := splits[1]

	// Unfold them.
	new_tokens := "" + tokens
	new_patterns := "" + pattern_string

	for i := 0; i < 4; i++ {
		new_tokens = new_tokens + "?" + tokens
		new_patterns = new_patterns + "," + pattern_string
	}

	patterns := strings.Split(new_patterns, ",")

	for _, column := range new_tokens {
		status = append(status, string_to_status[string(column)])
	}

	for _, value := range patterns {
		v, _ := strconv.Atoi(value)
		pattern = append(pattern, uint8(v))
	}

	return status, pattern
}

func to_string(row []uint8) string {
	result := ""

	for _, value := range row {
		if value == UNKNOWN {
			result += "?"
		} else if value == OPERATIONAL {
			result += "."
		} else if value == BROKEN {
			result += "#"
		}
	}

	return result
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

func consume_and_continue(status []uint8, pattern []uint8, starting_index int, cache *map[string]int) int {
	// Now check if we have already computed this.
	hash := fmt.Sprintf("%v-%v", starting_index, pattern)
	if value, ok := (*cache)[hash]; ok {
		// fmt.Println("Found cached value: ", value)
		return value
	}

	// Base case; there are no more patterns left to match. We have matched!
	if len(pattern) == 0 {
		// fmt.Println("Matched end of pattern!")
		(*cache)[hash] = 1
		return 1
	}

	// Base case: there are no more status left to match. We have failed!
	if len(status) <= starting_index {
		// fmt.Println("Failed to match end of status!")
		(*cache)[hash] = 0
		return 0
	}

	// Base case: we have ran out of characters
	left_over_characters := 0
	for _, value := range status[starting_index:] {
		if value == UNKNOWN || value == BROKEN {
			left_over_characters += 1
		}
	}

	need_to_match_total := 0
	for _, value := range pattern {
		need_to_match_total += int(value)
	}

	if left_over_characters < need_to_match_total {
		// fmt.Println("Ran out of characters! Have", left_over_characters, "but need", need_to_match_total)
		// fmt.Println("Status: ", status[starting_index:])
		(*cache)[hash] = 0
		return 0
	}

	need_to_consume := pattern[0]

	for i := starting_index; i < len(status); i++ {
		this_status := status[i]
		// fmt.Println("Considering: ", this_status, "with need_to_consume: ", need_to_consume)

		if need_to_consume == 0 {
			if this_status == BROKEN {
				// Uh, I need to consume nothing, but I have to!
				(*cache)[hash] = 0
				return 0
			}
			// fmt.Println("Have consumed all that I need to consume.")
			// The next one MUST be an OPERATIONAL or we can't match
			// Are we at the end of the array?
			if i+1 == len(status) {
				// We have matched!
				// fmt.Println("Matched up to my point but reached end of status!")
				if len(pattern) == 1 {
					(*cache)[hash] = 1
					return 1
				} else {
					(*cache)[hash] = 0
					return 0
				}
			}

			if status[i] == OPERATIONAL || status[i] == UNKNOWN {
				// Try all possible starting points.
				if len(pattern[1:]) == 0 {
					for _, value := range status[i+1:] {
						if value == BROKEN {
							// We have left overs :(
							(*cache)[hash] = 0
							return 0
						}
					}
					// We don't actually need to match anything else.
					(*cache)[hash] = 1
					return 1
				}

				// fmt.Println("Trying all possible starting points. Still need to match: ", pattern[1:], "with status: ", to_string(status[i+1:]))

				total_matches := 0
				for j := i + 1; j < len(status); j++ {
					if status[j] == OPERATIONAL {
						// No point in launching from here. We will just double count!
						continue
					}
					individual_matches := consume_and_continue(
						status,
						pattern[1:],
						j,
						cache,
					)
					// if individual_matches > 0 {
					// 	// fmt.Println("Found individual matches: ", individual_matches, "with status: ", to_string(status[j:]), "and pattern: ", pattern[1:])
					// }
					total_matches += individual_matches
					if status[j] == BROKEN {
						// Cannot go past this point! Recursion must handle the rest. I MUST match hashes.
						break
					}
				}
				// fmt.Println("Found total matches: ", total_matches)
				(*cache)[hash] = total_matches
				return total_matches
			} else {
				(*cache)[hash] = 0
				return 0
			}
		}

		if this_status == UNKNOWN || this_status == BROKEN {
			// fmt.Println("Matching unknown or broken at position: ", i, " with need_to_consume: ", need_to_consume)
			need_to_consume -= 1
			continue
		}

		if this_status == OPERATIONAL {
			// Fall out
			(*cache)[hash] = 0
			return 0
		}
	}

	if need_to_consume == 0 {
		if len(pattern) == 1 {
			// fmt.Println("Matched up to my point!")
			(*cache)[hash] = 1
			return 1
		}
	}

	(*cache)[hash] = 0
	return 0
}

func consume_all(
	status []uint8,
	pattern []uint8,
) int {
	total_matches := 0
	cache := make(map[string]int)

	for i := 0; i < len(status); i++ {
		if status[i] == OPERATIONAL {
			// No point in launching from here. We will just double count!
			continue
		}

		// fmt.Println("Entering consume_and_continue with status: ", status[i:], "and pattern: ", pattern)
		total_matches += consume_and_continue(
			status,
			pattern,
			i,
			&cache,
		)
		// fmt.Println("Exited consume_and_continue with n total matches", total_matches)

		if status[i] == BROKEN {
			// Cannot go past this point! Recursion must handle the rest.
			break
		}
	}

	return total_matches
}

func parse_and_return(row string) int {
	status, pattern := parse_row(row)
	matches := consume_all(status, pattern)
	return matches
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	all_rows := make([]string, 0)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		all_rows = append(all_rows, text)
	}

	if DEBUG {
		fmt.Println("Rows: ", all_rows)
	}

	// num_matches := make(chan int, len(all_rows))

	all_matches := 0

	for _, row := range all_rows {
		these_matches := parse_and_return(row)
		if DEBUG {
			fmt.Println("Total matches: ", these_matches)
		}
		all_matches += these_matches
	}

	fmt.Println("Total matches: ", all_matches)
}
