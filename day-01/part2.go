package main

import "fmt"
import "os"
import "log"
import "regexp"
import "bufio"

var string_map = map[string]int{
    "0": 0,
    "1": 1,
    "2": 2,
    "3": 3,
    "4": 4,
    "5": 5,
    "6": 6,
    "7": 7,
    "8": 8,
    "9": 9,
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "seven": 7,
    "eight": 8,
    "nine": 9,
}

func set_regex() *regexp.Regexp {
    regex := ""

    for k := range string_map {
        regex += string(k + "|")
    }

    return regexp.MustCompile(regex[:len(regex)-1])
}

var just_numbers_regex = set_regex()

func single_line_multimatch(str string) int {
    // First match is super easy.
    first_match := string_map[just_numbers_regex.FindString(str)]
    final_match := -1

    for i := len(str); i >= 0; i-- {
        match := just_numbers_regex.FindString(str[i:])

        if (match != "") {
            // First match going backwards!
            final_match = string_map[match];
            break;
        }
    }

    // Combine back together
    return 10 * first_match + final_match
}

func main() {
    DEBUG := false

    scanner := bufio.NewScanner(os.Stdin)

    total := 0;

    for scanner.Scan() {
        text := scanner.Text()

        single := single_line_multimatch(text)

        if (DEBUG) {
            fmt.Println("Output: ", single, " Base: ", text)
        }

        total += single
    }

    if scanner.Err() != nil {
        log.Fatal("No input provided.")
    }

    fmt.Println("Total: ", total)
}
