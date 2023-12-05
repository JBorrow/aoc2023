// To run this, yuo need to add a few blank lines to the input.

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

var numbers_regex = regexp.MustCompile(`\d+`)

// Capture groups give mapping ingredients (0 to 1)
var mapping_regex = regexp.MustCompile(`([a-z]+)-to-([a-z]+) map`)

// var MAXIMUM_NUMBER_OF_SEEDS = 100

type SeedMapping struct {
	from    string
	to      string
	mapping map[int]int
}

func parse_name_string(str string) (string, string) {
	matches := mapping_regex.FindStringSubmatch(str)

	return matches[1], matches[2]
}

func parse_mapping(strs []string) SeedMapping {
	result := SeedMapping{
		mapping: make(map[int]int),
	}

	// // Default is direct mapping to self
	// for i := 0; i < MAXIMUM_NUMBER_OF_SEEDS; i++ {
	// 	result.mapping[i] = i
	// }

	for _, v := range strs {
		if strings.Contains(v, ":") {
			result.from, result.to = parse_name_string(v)
		} else {
			match_numbers := numbers_regex.FindAllString(strings.TrimSpace(v), -1)

			to_start, _ := strconv.Atoi(match_numbers[0])
			from_start, _ := strconv.Atoi(match_numbers[1])
			length, _ := strconv.Atoi(match_numbers[2])

			for i := 0; i < length; i++ {
				result.mapping[from_start+i] = to_start + i
			}
		}
	}

	return result
}

func make_hops(have string, want string, mappings []SeedMapping, number int) int {
	// fmt.Println("Hopping: have ", have, " want ", want, " number ", number)
	if have == want {
		return number
	}

	// Find the correct mapping...
	for _, v := range mappings {
		if v.from == have {
			new_number, ok := v.mapping[number]

			if !ok {
				new_number = number
			}

			return make_hops(
				v.to, want, mappings, new_number,
			)
		}
	}

	fmt.Println("FAILED TO FIND MAPPING FOR: ", have, want, number)
	os.Exit(1)

	return 1
}

func parse_seeds(str string) []int {
	result := make([]int, 0)

	for _, v := range numbers_regex.FindAllString(str, -1) {
		integer, _ := strconv.Atoi(v)

		result = append(result, integer)
	}

	return result
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	mappings := make([]SeedMapping, 0)
	seeds := make([]int, 0)

	buffer := make([]string, 0)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "seeds") {
			seeds = parse_seeds(text)

			if DEBUG {
				fmt.Println("Seeds: ", seeds)
			}
		} else {
			if len(text) == 0 {
				this_mapping := parse_mapping(buffer)

				if DEBUG {
					fmt.Println("Base: ", buffer)
					fmt.Println("Output: ", this_mapping)
				}

				mappings = append(mappings, this_mapping)
				buffer = make([]string, 0)
			} else {
				buffer = append(buffer, text)
			}
		}

	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	smallest_location := 1000000000000

	for _, v := range seeds {
		new_location := make_hops("seed", "location", mappings, v)
		smallest_location = min(new_location, smallest_location)

		if DEBUG {
			fmt.Println("Seed: ", v)
			fmt.Println("Maps to location: ", new_location)
			fmt.Println("Smallest location so far: ", smallest_location)
		}
	}

	fmt.Println("Smallest location: ", smallest_location)

}
