package main

import "fmt"
import "os"
import "log"
import "regexp"
import "bufio"
import "strconv"
import "slices"
import "strings"

var numbers_regex = regexp.MustCompile(`\d+`)

type CardResult struct {
	winning_numbers []int
	card_numbers    []int
	id              int
	matches         int
	score           int
}

func parse_number_string(str string) []int {
	parsed_numbers := numbers_regex.FindAllString(strings.TrimSpace(str), -1)
	integers := make([]int, len(parsed_numbers))

	for i, v := range parsed_numbers {
		integers[i], _ = strconv.Atoi(v)
	}

	return integers
}

func parse_card_string(str string) int {
	parsed_number := numbers_regex.FindString(strings.TrimSpace(str))
	integer, _ := strconv.Atoi(parsed_number)

	return integer
}

func calculate_score(result CardResult) (int, int) {
	number_of_matches := 0

	for _, v := range result.card_numbers {
		if slices.Contains(result.winning_numbers, v) {
			number_of_matches++
		}
	}

	score := 0

	if number_of_matches > 0 {
		// Go does not have exponentiation for integers...
		score = 1 << (number_of_matches - 1)
	}

	return number_of_matches, score
}

func parse_card(str string) CardResult {
	split_first := strings.Split(str, ":")
	card_identifier := split_first[0]
	winning_and_card := split_first[1]

	split_second := strings.Split(winning_and_card, "|")
	winning_numbers := split_second[0]
	card_numbers := split_second[1]

	result := CardResult{
		winning_numbers: parse_number_string(winning_numbers),
		card_numbers:    parse_number_string(card_numbers),
		id:              parse_card_string(card_identifier),
		score:           0,
	}

	result.matches, result.score = calculate_score(result)

	return result
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		card_result := parse_card(text)

		if DEBUG {
			fmt.Println("Output: ", card_result)
			fmt.Println(" Base: ", text)
		}

		total += card_result.score
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	fmt.Println("Total: ", total)
}
