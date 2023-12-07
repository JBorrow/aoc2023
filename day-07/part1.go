package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var NUMBER_OF_CARDS_PER_HAND = 5

var card_ranking = map[string]int{
	"2": 1,
	"3": 2,
	"4": 3,
	"5": 4,
	"6": 5,
	"7": 6,
	"8": 7,
	"9": 8,
	"T": 9,
	"J": 10,
	"Q": 11,
	"K": 12,
	"A": 13,
}

// We are using this as an enum as that doesn't exist in go
var hand_strengths = map[string]int{
	"high_card":       1,
	"one_pair":        2,
	"two_pair":        3,
	"three_of_a_kind": 4,
	"full_house":      5,
	"four_of_a_kind":  6,
	"five_of_a_kind":  7,
}

type Hand struct {
	cards         []int
	hand_strength int
	bid           int
}

type ByHandStrength []Hand

func (a ByHandStrength) Len() int { return len(a) }
func (a ByHandStrength) Less(i, j int) bool {
	// Only complication is if the two hand strengths are the same
	if a[i].hand_strength == a[j].hand_strength {
		// If they are the same, compare the cards in turn.
		for k := 0; k < NUMBER_OF_CARDS_PER_HAND; k += 1 {
			if a[i].cards[k] == a[j].cards[k] {
				continue
			}

			return a[i].cards[k] < a[j].cards[k]
		}

		// If we get here, the hands are identical.
		return false
	}

	return a[i].hand_strength < a[j].hand_strength
}
func (a ByHandStrength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func calculate_hand_strength(hand Hand) int {
	number_of_cards := map[int]int{}

	for _, v := range hand.cards {
		number_of_cards[v] += 1
	}

	switch unique := len(number_of_cards); unique {
	case 5:
		return hand_strengths["high_card"]
	case 4:
		// We have one pair and only one pair.
		return hand_strengths["one_pair"]
	case 3:
		// We have either two pair or three of a kind.
		for _, v := range number_of_cards {
			if v == 3 {
				return hand_strengths["three_of_a_kind"]
			}
		}
		return hand_strengths["two_pair"]
	case 2:
		// We either have three of a kind, or four of a kind.
		for _, v := range number_of_cards {
			if v == 4 {
				return hand_strengths["four_of_a_kind"]
			}
		}
		return hand_strengths["full_house"]
	case 1:
		return hand_strengths["five_of_a_kind"]
	default:
		fmt.Printf("Warning: %d is not a valid number of unique cards.\n", unique)
	}

	fmt.Println("Warning: fell through case statement in calculate_hand_strength.")

	return 0
}

func process_hand(to_parse string) Hand {
	cards := make([]int, NUMBER_OF_CARDS_PER_HAND)

	split_string := strings.Split(strings.TrimSpace(to_parse), " ")

	if len(split_string) != 2 {
		fmt.Println("Warning: invalid input.")
		fmt.Println("Was given: ", to_parse)
	}

	hand_string := split_string[0]
	bid_string := split_string[1]

	bid, _ := strconv.Atoi(bid_string)

	for i := 0; i < len(hand_string); i += 1 {
		cards[i] = card_ranking[string(hand_string[i])]
	}

	hand := Hand{
		cards: cards,
		bid:   bid,
	}

	hand.hand_strength = calculate_hand_strength(hand)

	return hand
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	hands := make([]Hand, 0)

	for scanner.Scan() {
		text := scanner.Text()

		hand := process_hand(text)

		if DEBUG {
			fmt.Println("Given: ", text)
			fmt.Println("Parsed to: ", hand)
		}

		hands = append(hands, hand)
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	if DEBUG {
		fmt.Println("Unsorted hands: ", hands)
	}

	// Now we can sort by hand strength.
	sort.Sort(ByHandStrength(hands))

	if DEBUG {
		fmt.Println("Sorted hands: ", hands)
	}

	total_winnings := 0

	for rank, hand := range hands {
		// We are using 1-indexing here.
		total_winnings += (rank + 1) * hand.bid
	}

	fmt.Println("Total winnings: ", total_winnings)

}
