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

var directions = map[string][]int{
	"U": {0, -1},
	"D": {0, 1},
	"L": {-1, 0},
	"R": {1, 0},
}

var new_directions = map[string]string{
	"0": "R",
	"1": "D",
	"2": "L",
	"3": "U",
}

type Instruction struct {
	direction []int
	steps     int
	color     string
	vertex    []int
}

var instruction_regex = regexp.MustCompile(`^([UDLR]) (\d+) \(#([a-z0-9A-Z]*)\)`)

func line_to_instruction(str string) Instruction {
	matches := instruction_regex.FindStringSubmatch(str)

	if matches == nil {
		panic("No matches")
	}

	// direction := directions[matches[1]]
	// steps, _ := strconv.Atoi(matches[2])
	color := matches[3]

	// Unpack steps
	steps, _ := strconv.ParseInt(string(color[:len(color)-1]), 16, 32)
	direction := directions[new_directions[string(color[len(color)-1])]]

	vertex := []int{0, 0}

	return Instruction{direction, int(steps), color, vertex}
}

func dig_trenches(instructions []Instruction) {
	// Actually figure out verticies
	x := 0
	y := 0

	for i, instruction := range instructions {
		x += instruction.direction[0] * instruction.steps
		y += instruction.direction[1] * instruction.steps
		instructions[i].vertex = []int{x, y}
	}
}

func instructions_to_vertex(instructions []Instruction) [][]int {
	verticies := make([][]int, len(instructions)+2)

	verticies[0] = []int{0, 0}

	for i, instruction := range instructions {
		verticies[i+1] = instruction.vertex
	}

	verticies[len(verticies)-1] = verticies[0]

	return verticies
}

func shoelace(instructions []Instruction) int {
	verticies := instructions_to_vertex(instructions)

	area := 0

	for i := 0; i <= len(verticies)-2; i++ {
		current_x := verticies[i][0]
		current_y := verticies[i][1]

		next_x := verticies[i+1][0]
		next_y := verticies[i+1][1]

		area += current_x*next_y - current_y*next_x
	}

	area = area / 2

	if area < 0 {
		area = -area
	}

	edges := 0

	for _, instruction := range instructions {
		edges += instruction.steps
	}

	fmt.Println("Edges:", edges)

	return area + edges/2 + 1
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	all_rows := make([]string, 0)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if DEBUG {
			fmt.Println("Given: ", text)
		}

		all_rows = append(all_rows, text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}

	instructions := make([]Instruction, len(all_rows))

	for i, row := range all_rows {
		instructions[i] = line_to_instruction(row)
	}

	if DEBUG {
		fmt.Println(instructions)
	}

	dig_trenches(instructions)

	if DEBUG {
		for _, instruction := range instructions {
			fmt.Println(instruction.vertex)
		}
	}

	fmt.Println("Filled:", shoelace(instructions))
}
