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

type MappingRange struct {
	start_from int
	stop_from  int
	start_to   int
	stop_to    int
	width      int
}

func (mapping *SeedMapping) FindMapping(number int) int {
	for _, v := range mapping.ranges {
		if number >= v.start_from && number <= v.stop_from {
			return v.start_to + (number - v.start_from)
		}
	}

	return number
}

type SeedMapping struct {
	from   string
	to     string
	ranges []MappingRange
}

func parse_name_string(str string) (string, string) {
	matches := mapping_regex.FindStringSubmatch(str)

	return matches[1], matches[2]
}

func parse_mapping(strs []string) SeedMapping {
	result := SeedMapping{
		ranges: make([]MappingRange, 0),
	}

	for _, v := range strs {
		if strings.Contains(v, ":") {
			result.from, result.to = parse_name_string(v)
		} else {
			match_numbers := numbers_regex.FindAllString(strings.TrimSpace(v), -1)

			to_start, _ := strconv.Atoi(match_numbers[0])
			from_start, _ := strconv.Atoi(match_numbers[1])
			length, _ := strconv.Atoi(match_numbers[2])

			new_mapping_range := MappingRange{
				start_from: from_start,
				stop_from:  from_start + length - 1,
				start_to:   to_start,
				stop_to:    to_start + length - 1,
				width:      length,
			}

			result.ranges = append(result.ranges, new_mapping_range)
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
			new_number := v.FindMapping(number)

			return make_hops(
				v.to, want, mappings, new_number,
			)
		}
	}

	fmt.Println("FAILED TO FIND MAPPING FOR: ", have, want, number)
	os.Exit(1)

	return 1
}

func parse_seeds(str string) [][]int {
	result := make([][]int, 0)

	// Consider seeds in batches.

	string_seeds := numbers_regex.FindAllString(strings.TrimSpace(str), -1)

	for i := 0; i < len(string_seeds)/2; i++ {
		from, _ := strconv.Atoi(string_seeds[i*2])
		length, _ := strconv.Atoi(string_seeds[i*2+1])

		range_array := make([]int, 2)
		range_array[0] = from
		range_array[1] = from + length

		result = append(result, range_array)
	}

	return result
}

func main() {
	DEBUG := false

	scanner := bufio.NewScanner(os.Stdin)

	mappings := make([]SeedMapping, 0)
	seeds := make([][]int, 0)

	buffer := make([]string, 0)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "seeds") {
			seeds = parse_seeds(text)

			if DEBUG {
				fmt.Println("Seeds input: ", text)
				fmt.Println("Seeds: ", seeds)
			}
		} else {
			if len(text) == 0 {
				if len(buffer) == 0 {
					continue
				}

				if DEBUG {
					fmt.Println("Buffer: ")
					for _, v := range buffer {
						fmt.Println(v)
					}
				}

				this_mapping := parse_mapping(buffer)

				if DEBUG {
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
		fmt.Println("Starting seed block: ", v)
		for i := v[0]; i < v[1]; i++ {
			new_location := make_hops("seed", "location", mappings, i)
			smallest_location = min(new_location, smallest_location)

			if DEBUG {
				fmt.Println("Seed: ", i)
				fmt.Println("Maps to location: ", new_location)
				fmt.Println("Smallest location so far: ", smallest_location)
			}

		}
	}

	fmt.Println("Smallest location: ", smallest_location)

}
