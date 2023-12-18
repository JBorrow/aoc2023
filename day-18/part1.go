package part1

import (
	"bufio"
	"fmt"
	"log"
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

type Instruction struct {
	direction []int
	steps     int
	color     string
}

var instruction_regex = regexp.MustCompile(`^([UDLR]) (\d+) \(#([a-z0-9]*)\)`)

func line_to_instruction(str string) Instruction {
	matches := instruction_regex.FindStringSubmatch(str)

	if matches == nil {
		panic("No matches")
	}

	direction := directions[matches[1]]
	steps, _ := strconv.Atoi(matches[2])
	color := matches[3]

	return Instruction{direction, steps, color}
}

func dig_trenches(instructions []Instruction) [][]string {
	// Figure out how big our grid needs to be.
	max_x := 0
	max_y := 0
	min_x := 0
	min_y := 0

	x := 0
	y := 0

	for _, instruction := range instructions {
		x += instruction.direction[0] * instruction.steps
		y += instruction.direction[1] * instruction.steps

		max_x = max(max_x, x)
		max_y = max(max_y, y)
		min_x = min(min_x, x)
		min_y = min(min_y, y)
	}

	// Add a little extra space around the edges.
	grid_size_x := max_x - min_x + 3
	grid_size_y := max_y - min_y + 3

	// Starting point is now -(min_x, min_y)
	x = -min_x + 1
	y = -min_y + 1

	grid := make([][]string, grid_size_y)

	for i := 0; i < grid_size_y; i++ {
		grid[i] = make([]string, grid_size_x)
	}

	// First is first instructions colour.
	grid[y][x] = instructions[0].color

	for _, instruction := range instructions {
		for i := 0; i < instruction.steps; i++ {
			x += instruction.direction[0]
			y += instruction.direction[1]

			grid[y][x] = instruction.color
		}
	}

	return grid
}

func print_grid(grid [][]string) {
	for _, row := range grid {
		for _, cell := range row {
			if cell == "" {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

// We can use watershed here.
func watershed_around(grid [][]string, start_x int, start_y int, fill_color string) {
	// First check my points.
	my_node := grid[start_y][start_x]

	if my_node != "" {
		return
	}

	// Fill me in.
	grid[start_y][start_x] = fill_color

	// Now check my neighbors.
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			new_x := start_x + x
			new_y := start_y + y

			if new_x < 0 || new_y < 0 || new_x >= len(grid[0]) || new_y >= len(grid) {
				// We escaped! Panic!!!
				log.Fatal("We escaped!")
			}

			watershed_around(grid, new_x, new_y, fill_color)
		}
	}
}

func fill_grid(grid [][]string, fill_color string) {
	// Start at half way down the grid, iterate until
	// we know we are inside the grid.

	y := len(grid) / 2

	edge := false
	final_x := 0

	for x := 0; x < len(grid[0]); x++ {
		current_node := grid[y][x]

		if edge && current_node == "" {
			final_x = x
			break
		}

		if current_node != "" {
			edge = true
		}
	}

	if !edge {
		log.Fatal("We escaped!")
	}

	start_point_x := final_x
	start_point_y := y

	fmt.Println("Starting watershed around", start_point_x, start_point_y, "with color", fill_color)

	watershed_around(grid, start_point_x, start_point_y, fill_color)
}

func count_filled(grid [][]string) int {
	count := 0

	for _, row := range grid {
		for _, cell := range row {
			if cell != "" {
				count++
			}
		}
	}

	return count
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

	grid := dig_trenches(instructions)

	if DEBUG {
		print_grid(grid)
	}

	fill_grid(grid, "FFFFFF")

	if DEBUG {
		print_grid(grid)
	}

	fmt.Println("Filled:", count_filled(grid))
}
