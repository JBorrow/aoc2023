package main

import "fmt"
import "os"
import "log"
import "strconv"
import "strings"
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

var just_numbers_regex = regexp.MustCompile(
    `[0-9]|one|two|three|four|five|six|seven|eight|nine`,
)

func single_line_multimatch(str string) int {
    first_match := -1
    final_match := -1

    for i := 0; i < len(str); i++ {
        match := just_numbers_regex.FindString(str[i:])

        if (match == "") {
            // There can be no more matches!
            break;
        }
        
        mapped_integer := string_map[match]

        if (first_match < 0) {
            first_match = mapped_integer;
        }

        final_match = mapped_integer
    }

    // Combine back together

    first_string := strconv.Itoa(first_match)
    final_string := strconv.Itoa(final_match)
    combined := first_string + final_string

    final, _ := strconv.Atoi(combined)

    return final
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    total := 0;

    for scanner.Scan() {
        text := scanner.Text()

        single := single_line_multimatch(text)
        fmt.Println("Output: ", single, " Base: ", text)

        total += single
    }

    if scanner.Err() != nil {
        log.Fatal("No input provided.")
    }

    fmt.Println("Total: ", total)
}
