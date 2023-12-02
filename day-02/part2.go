package main

import "fmt"
import "os"
import "log"
import "regexp"
import "bufio"
import "strconv"
import "strings"

var blue_regex = regexp.MustCompile(`(\d*) blue`)
var red_regex = regexp.MustCompile(`(\d*) red`)
var green_regex = regexp.MustCompile(`(\d*) green`)
var game_regex = regexp.MustCompile(`Game (\d*)`)

type ReplacementResult struct {
	red   int
	green int
	blue  int
}

type GameResult struct {
	replacements []ReplacementResult
	id           int
}

func parse_game(str string) GameResult {
	game_string := game_regex.FindStringSubmatch(str)
	id, _ := strconv.Atoi(game_string[1])

	replacement_strings := strings.Split(str, ";")
	replacements := make([]ReplacementResult, len(replacement_strings))

	for i := range replacement_strings {
		replacements[i] = parse_replacement(replacement_strings[i])
	}

	return GameResult{replacements: replacements, id: id}
}

func parse_replacement(str string) ReplacementResult {
	red_string := red_regex.FindStringSubmatch(str)
	blue_string := blue_regex.FindStringSubmatch(str)
	green_string := green_regex.FindStringSubmatch(str)

	result := ReplacementResult{red: 0, green: 0, blue: 0}

	if len(red_string) > 0 {
		result.red, _ = strconv.Atoi(red_string[1])
	}

	if len(blue_string) > 0 {
		result.blue, _ = strconv.Atoi(blue_string[1])
	}

	if len(green_string) > 0 {
		result.green, _ = strconv.Atoi(green_string[1])
	}

	return result
}

func minimum_possible(results []ReplacementResult) ReplacementResult {
	balls_in_bag := ReplacementResult{red: 0, green: 0, blue: 0}

	for i := range results {
		current_result := results[i]

		if current_result.red > balls_in_bag.red {
			balls_in_bag.red = current_result.red
		}

		if current_result.green > balls_in_bag.green {
			balls_in_bag.green = current_result.green
		}

		if current_result.blue > balls_in_bag.blue {
			balls_in_bag.blue = current_result.blue
		}
	}

	return balls_in_bag
}

func bag_power(result ReplacementResult) int {
	return result.red * result.green * result.blue
}

func minimum_game_power(result GameResult) int {
	balls_in_bag := minimum_possible(result.replacements)

	return bag_power(balls_in_bag)
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		game_result := parse_game(text)
		game_power := minimum_game_power(game_result)

		if DEBUG {
			fmt.Println("Output: ", game_result, " Base: ", text)
			fmt.Println("Game minimum power: ", game_power)
		}

		total += game_power
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	fmt.Println("Total: ", total)
}
