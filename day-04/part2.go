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

func calculate_matches(result CardResult) int {
	number_of_matches := 0

	for _, v := range result.card_numbers {
		if slices.Contains(result.winning_numbers, v) {
			number_of_matches++
		}
	}

	return number_of_matches
}

func calculate_score(result []CardResult) ([]int, int) {
	number_of_cards := make([]int, len(result))

	for i := range number_of_cards {
		number_of_cards[i] = 1
	}

	for i, v := range result {
		me_copies := number_of_cards[i]

		if v.matches == 0 {
			continue
		}

		for x := 1; x <= v.matches; x++ {
			number_of_cards[i+x] += me_copies
		}
	}

	total_score := 0

	for _, v := range number_of_cards {
		total_score += v
	}

	return number_of_cards, total_score
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
	}

	result.matches = calculate_matches(result)

	return result
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	cards := make([]CardResult, 0)

	for scanner.Scan() {
		text := scanner.Text()

		card_result := parse_card(text)

		if DEBUG {
			fmt.Println("Output: ", card_result)
			fmt.Println(" Base: ", text)
		}

		cards = append(cards, card_result)
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	_, total_score := calculate_score(cards)

	fmt.Println("Total: ", total_score)

}
