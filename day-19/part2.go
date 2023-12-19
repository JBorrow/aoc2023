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

func walk_constraints(lines map[string]Line, test_range map[string][]int, line_name string, accepted_ranges *[]map[string][]int, rejected_ranges *[]map[string][]int) {
	if line_name == "A" {
		// Accept the range.
		fmt.Println("Accepting", test_range)
		*accepted_ranges = append(*accepted_ranges, test_range)
		return
	}

	if line_name == "R" {
		fmt.Println("Rejecting", test_range)
		*rejected_ranges = append(*rejected_ranges, test_range)
		return
	}

	line := lines[line_name]

	if DEBUG {
		fmt.Println("Currently on line name: ", line_name)
	}

	broken := false

	for _, instruction := range line.instructions {
		fmt.Println("Testing range:", test_range)
		if instruction.less {
			value_range := test_range[instruction.condition_on]
			fmt.Println("Testing value range:", value_range)

			if value_range[0] > instruction.value {
				// This will never be true. We always go to the next instruction.
				continue
			} else if value_range[1] < instruction.value {
				// This will always be true. We always go to the next 'instruction'.
				walk_constraints(lines, test_range, instruction.next_true, accepted_ranges, rejected_ranges)
				broken = true
				break
			} else {
				// This might be true. We have to split!
				fmt.Println("Splitting! Test range:", test_range, "at", instruction.value, "(", instruction.condition_on, ")")

				lower_range := make(map[string][]int)
				upper_range := make(map[string][]int)

				for key, value := range test_range {
					lower_range[key] = []int{value[0], value[1]}
					upper_range[key] = []int{value[0], value[1]}
				}

				lower_range[instruction.condition_on][1] = instruction.value - 1
				upper_range[instruction.condition_on][0] = instruction.value

				fmt.Println("Test range:", test_range)
				fmt.Println("Lower range:", lower_range)
				fmt.Println("Upper range:", upper_range)

				walk_constraints(lines, lower_range, instruction.next_true, accepted_ranges, rejected_ranges)

				// Keep the 'false' range and continue
				test_range = upper_range
				continue
			}
		} else {
			value_range := test_range[instruction.condition_on]

			if value_range[0] > instruction.value {
				// This will always be true. We always go to the next 'instruction'.
				walk_constraints(lines, test_range, instruction.next_true, accepted_ranges, rejected_ranges)
				broken = true
				break
			} else if value_range[1] < instruction.value {
				// This will never be true. We always go to the next instruction.
				continue
			} else {
				// This might be true. We have to split!
				fmt.Println("Splitting! Test range:", test_range, "at", instruction.value, "(", instruction.condition_on, ")")

				fmt.Println("Test range:", test_range)
				lower_range := make(map[string][]int)
				upper_range := make(map[string][]int)

				for key, value := range test_range {
					lower_range[key] = []int{value[0], value[1]}
					upper_range[key] = []int{value[0], value[1]}
				}

				lower_range[instruction.condition_on][1] = instruction.value
				upper_range[instruction.condition_on][0] = instruction.value + 1

				fmt.Println("Lower range:", lower_range)
				fmt.Println("Upper range:", upper_range)

				walk_constraints(lines, upper_range, instruction.next_true, accepted_ranges, rejected_ranges)

				// Keep the 'false' range and continue
				test_range = lower_range
				continue
			}
		}
	}

	if !broken {
		walk_constraints(lines, test_range, line.fall_through, accepted_ranges, rejected_ranges)
	}
}

func count_combinations(usable_ranges []map[string][]int) int {
	total_combos := 0

	for _, usable_range := range usable_ranges {
		combos := 1

		combos *= usable_range["x"][1] - usable_range["x"][0] + 1
		combos *= usable_range["m"][1] - usable_range["m"][0] + 1
		combos *= usable_range["a"][1] - usable_range["a"][0] + 1
		combos *= usable_range["s"][1] - usable_range["s"][0] + 1

		total_combos += combos
	}

	return total_combos
}

func check_overlap(usable_ranges []map[string][]int) bool {
	for i, usable_range_i := range usable_ranges {
		for j, usable_range_j := range usable_ranges {
			if i == j {
				continue
			}

			x_overlap := false
			m_overlap := false
			a_overlap := false
			s_overlap := false

			if usable_range_i["x"][0] <= usable_range_j["x"][1] {
				if usable_range_i["x"][1] >= usable_range_j["x"][0] {
					x_overlap = true
				}
			}

			if usable_range_i["m"][0] <= usable_range_j["m"][1] {
				if usable_range_i["m"][1] >= usable_range_j["m"][0] {
					m_overlap = true
				}
			}

			if usable_range_i["a"][0] <= usable_range_j["a"][1] {
				if usable_range_i["a"][1] >= usable_range_j["a"][0] {
					a_overlap = true
				}
			}

			if usable_range_i["s"][0] <= usable_range_j["s"][1] {
				if usable_range_i["s"][1] >= usable_range_j["s"][0] {
					s_overlap = true
				}
			}

			if x_overlap && m_overlap && a_overlap && s_overlap {
				fmt.Println("Overlap between", usable_range_i, "and", usable_range_j)

				return true
			}
		}
	}
	return false
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

	if DEBUG {
		fmt.Println(lines)
	}

	acceptable_range := make(map[string][]int)

	acceptable_range["x"] = []int{1, 4000}
	acceptable_range["m"] = []int{1, 4000}
	acceptable_range["a"] = []int{1, 4000}
	acceptable_range["s"] = []int{1, 4000}

	usable_ranges := make([]map[string][]int, 0)
	unusable_ranges := make([]map[string][]int, 0)

	walk_constraints(lines, acceptable_range, "in", &usable_ranges, &unusable_ranges)

	fmt.Println(usable_ranges)

	check_overlap(usable_ranges)

	fmt.Println("Real answer (for test data):", 167409079868000)
	fmt.Println("Total space:", count_combinations(usable_ranges))
	fmt.Println("Total rejected:", count_combinations(unusable_ranges))
	fmt.Println("Total expected sum:", count_combinations(usable_ranges)+count_combinations(unusable_ranges))
	fmt.Println("Total sum:", 4000*4000*4000*4000)
	fmt.Println("Total not rejected: ", 4000*4000*4000*4000-count_combinations(unusable_ranges))

}
