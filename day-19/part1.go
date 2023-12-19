package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DEBUG = true

// px{a<2006:qkq,m>2090:A,rfg}
// pv{a>1716:R,A}
// lnx{m>1548:A,A}
// rfg{s<537:gd,x>2440:R,A}
// qs{s>3448:A,lnx}
// qkq{x<1416:A,crn}
// crn{x>2662:A,R}
// in{s<1351:px,qqz}
// qqz{s>2770:qs,m<1801:hdj,R}
// gd{a>3333:R,R}
// hdj{m>838:A,pv}

// {x=787,m=2655,a=1222,s=2876}
// {x=1679,m=44,a=2067,s=496}
// {x=2036,m=264,a=79,s=2244}
// {x=2461,m=1339,a=466,s=291}
// {x=2127,m=1623,a=2188,s=1013}

type Instruction struct {
	condition_on string
	less         bool
	value        int
	next_true    string
}

type Line struct {
	instructions []Instruction
	fall_through string
}

var unpack_data_regex = regexp.MustCompile(`{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}`)
var unpack_line_regex = regexp.MustCompile(`^([a-z]+){(.*)}`)

func unpack_instruction(str string) Instruction {
	if strings.Contains(str, ">") {
		split_by_greater := strings.Split(str, ">")
		split_by_colon := strings.Split(split_by_greater[1], ":")

		name := split_by_greater[0]
		condition_on, _ := strconv.Atoi(split_by_colon[0])
		next_true := split_by_colon[1]

		return Instruction{
			condition_on: name,
			less:         false,
			value:        condition_on,
			next_true:    next_true,
		}
	} else if strings.Contains(str, "<") {
		split_by_less := strings.Split(str, "<")
		split_by_colon := strings.Split(split_by_less[1], ":")

		name := split_by_less[0]
		condition_on, _ := strconv.Atoi(split_by_colon[0])
		next_true := split_by_colon[1]

		return Instruction{
			condition_on: name,
			less:         true,
			value:        condition_on,
			next_true:    next_true,
		}
	} else {
		panic("No comparison")
	}
}

func unpack_line(str string) (string, Line) {
	result := unpack_line_regex.FindStringSubmatch(str)

	if result == nil {
		panic("No matches")
	}

	name := result[1]
	instruction_strings := strings.Split(result[2], ",")

	instructions := make([]Instruction, len(instruction_strings)-1)

	for i, instruction_string := range instruction_strings[:len(instruction_strings)-1] {
		instructions[i] = unpack_instruction(instruction_string)
	}

	fall_through := instruction_strings[len(instruction_strings)-1]

	return name, Line{instructions, fall_through}
}

func unpack_data(str string) map[string]int {
	result := unpack_data_regex.FindStringSubmatch(str)

	if result == nil {
		panic("No matches")
	}

	data := make(map[string]int)

	data["x"], _ = strconv.Atoi(result[1])
	data["m"], _ = strconv.Atoi(result[2])
	data["a"], _ = strconv.Atoi(result[3])
	data["s"], _ = strconv.Atoi(result[4])

	return data
}

func walk(data map[string]int, lines map[string]Line) int {
	line_name := "in"

	if DEBUG {
		fmt.Println("Using data: ", data)
	}

	for {
		if DEBUG {
			fmt.Println("Currently on line name: ", line_name)
		}
		if line_name == "R" {
			if DEBUG {
				fmt.Println("Rejecting part from R.")
			}
			return 0
		} else if line_name == "A" {
			if DEBUG {
				fmt.Println("Accepting part from A.")
			}
			return data["a"] + data["m"] + data["s"] + data["x"]
		}

		line := lines[line_name]
		broken := false

		for _, instruction := range line.instructions {
			if DEBUG {
				fmt.Println("Currently on instruction: ", instruction)
			}
			if instruction.less {
				if DEBUG {
					fmt.Println("Checking", data[instruction.condition_on], "Less than", instruction.value)
				}
				if data[instruction.condition_on] < instruction.value {
					line_name = instruction.next_true
					broken = true
					break
				}
			} else {
				if DEBUG {
					fmt.Println("Checking", data[instruction.condition_on], "Greater than", instruction.value)
				}
				if data[instruction.condition_on] > instruction.value {
					line_name = instruction.next_true
					broken = true
					break
				}
			}
		}

		if broken {
			continue
		}

		if line.fall_through == "A" {
			// Accept the part.
			if DEBUG {
				fmt.Println("Accepting part from fall through.")
			}
			return data["a"] + data["m"] + data["s"] + data["x"]
		} else if line.fall_through == "R" {
			if DEBUG {
				fmt.Println("Rejecting part from fall through.")
			}
			return 0
		} else {
			line_name = line.fall_through
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	instructions := make([]string, 0)
	inputs := make([]string, 0)

	done_instructions := false

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if done_instructions {
			inputs = append(inputs, input)
		} else {
			if input == "" {
				done_instructions = true
			} else {
				instructions = append(instructions, input)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}

	lines := make(map[string]Line)

	for _, instruction := range instructions {
		name, line := unpack_line(instruction)
		lines[name] = line
	}

	data := make([]map[string]int, len(inputs))

	for i, input := range inputs {
		data[i] = unpack_data(input)
	}

	if DEBUG {
		fmt.Println(lines)
		fmt.Println(data)
	}

	n_accepted := 0
	total := 0

	for _, d := range data {
		this := walk(d, lines)
		if this > 0 {
			n_accepted++
			total += this
		}
	}

	fmt.Println("Accepted parts:", n_accepted)
	fmt.Println("Total value:", total)

}
