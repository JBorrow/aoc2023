package main

import "fmt"
import "os"
import "log"
import "strconv"
import "strings"
import "regexp"
import "bufio"

var just_numbers_regex = regexp.MustCompile(`[0-9]`)

func single_line(str string) int {
	sub_string := just_numbers_regex.FindAllString(str, -1)
	n_strings := len(sub_string)

	// Deal with the fact we may have more than one character
	first := string(sub_string[0])
	final := first

	if n_strings > 1 {
		final = string(sub_string[n_strings-1])
	}

	first_and_last := strings.TrimSpace(
		first + final,
	)

	real_int, err := strconv.Atoi(first_and_last)

	if err != nil {
		log.Fatal(err)
	}

	return real_int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		single := single_line(text)
		fmt.Println("Output: ", single, " Base: ", text)

		total += single
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	fmt.Println("Total: ", total)
}
