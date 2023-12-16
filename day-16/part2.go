package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const DEBUG = false

const PART2 = true

var color_red = "\033[31m"
var color_none = "\033[0m"
var color_green = "\033[32m"

type Connection struct {
	input   []int
	outputs [][]int
}

type Node struct {
	symbol      string
	render      string
	used_hashes []string
	connections []Connection
	visited     int
}

var nodes = map[string]Node{
	"|": {
		symbol:  "|",
		render:  "│",
		visited: 0,
		connections: []Connection{
			{
				input: []int{0, -1},
				outputs: [][]int{
					{0, 1},
				},
			},
			{
				input: []int{0, 1},
				outputs: [][]int{
					{0, -1},
				},
			},
			// Now do splitting!
			{
				input: []int{1, 0},
				outputs: [][]int{
					{0, 1},
					{0, -1},
				},
			},
			{
				input: []int{-1, 0},
				outputs: [][]int{
					{0, 1},
					{0, -1},
				},
			},
		},
	},
	"-": {
		symbol:  "-",
		render:  "─",
		visited: 0,
		connections: []Connection{
			{
				input: []int{-1, 0},
				outputs: [][]int{
					{1, 0},
				},
			},
			{
				input: []int{1, 0},
				outputs: [][]int{
					{-1, 0},
				},
			},
			// Now do splitting!
			{
				input: []int{0, 1},
				outputs: [][]int{
					{1, 0},
					{-1, 0},
				},
			},
			{
				input: []int{0, -1},
				outputs: [][]int{
					{1, 0},
					{-1, 0},
				},
			},
		},
	},
	"/": {
		symbol: "/",
		render: "╱",
		connections: []Connection{
			// Reflect up
			{
				input: []int{-1, 0},
				outputs: [][]int{
					{0, -1},
				},
			},
			// Reflect down
			{
				input: []int{1, 0},
				outputs: [][]int{
					{0, 1},
				},
			},
			// Reflect left
			{
				input: []int{0, -1},
				outputs: [][]int{
					{-1, 0},
				},
			},
			// Reflect right
			{
				input: []int{0, 1},
				outputs: [][]int{
					{1, 0},
				},
			},
		},
	},
	"\\": {
		symbol:  "\\",
		render:  "╲",
		visited: 0,
		connections: []Connection{
			// Reflect up
			{
				input: []int{1, 0},
				outputs: [][]int{
					{0, -1},
				},
			},
			// Reflect down
			{
				input: []int{-1, 0},
				outputs: [][]int{
					{0, 1},
				},
			},
			// Reflect left
			{
				input: []int{0, 1},
				outputs: [][]int{
					{-1, 0},
				},
			},
			// Reflect right
			{
				input: []int{0, -1},
				outputs: [][]int{
					{1, 0},
				},
			},
		},
	},
	".": {
		// Basically empty space!
		symbol:  ".",
		render:  " ",
		visited: 0,
		connections: []Connection{
			{
				input: []int{1, 0},
				outputs: [][]int{
					{-1, 0},
				},
			},
			{
				input: []int{-1, 0},
				outputs: [][]int{
					{1, 0},
				},
			},
			{
				input: []int{0, 1},
				outputs: [][]int{
					{0, -1},
				},
			},
			{
				input: []int{0, -1},
				outputs: [][]int{
					{0, 1},
				},
			},
		},
	},
}

func node_from_string(str string) Node {
	return nodes[str]
}

func rows_to_nodes(rows []string) [][]Node {
	nodes := make([][]Node, len(rows))

	for i, row := range rows {
		nodes[i] = make([]Node, len(row))

		for j, char := range row {
			nodes[i][j] = node_from_string(string(char))
		}
	}

	return nodes
}

func print_node_array(nodes [][]Node) {
	for _, row := range nodes {
		for _, node := range row {
			if node.visited > 0 && node.visited < 100 {
				fmt.Print(color_red)
			}
			if node.visited > 100 {
				fmt.Print(color_green)
			}
			if node.visited > 0 && node.symbol == "." {
				fmt.Print("#")
			} else {
				fmt.Print(node.render)
			}
			if node.visited > 0 {
				fmt.Print(color_none)
			}
		}
		fmt.Println()
	}
}

func hash_position(position []int, direction []int) string {
	// Hash the position and direction into a string.
	return fmt.Sprintf("%d,%d,%d,%d", position[0], position[1], direction[0], direction[1])
}

func follow_path(nodes [][]Node, position []int, direction []int) {
	// Base case: our new position is out of bounds!
	next_position := []int{position[0] + direction[0], position[1] + direction[1]}

	x := next_position[0]
	y := next_position[1]

	if x < 0 || x >= len(nodes[0]) || y < 0 || y >= len(nodes) {
		if DEBUG {
			fmt.Println("Terminating at: ", next_position)
		}
		return
	}

	// Base case: we are caught in a trap!
	next := nodes[y][x]
	if slices.Contains(next.used_hashes, hash_position(position, direction)) {
		return
	}

	if DEBUG {
		fmt.Println("Visiting: ", next_position, "with symbol", nodes[y][x].symbol)
		fmt.Println("Should terminate at: ", len(nodes[0]), len(nodes))
	}

	for _, connection := range next.connections {
		if connection.input[0] == -direction[0] && connection.input[1] == -direction[1] {
			// We found the input connection.
			nodes[y][x].visited++

			// Add the path to the history.
			nodes[y][x].used_hashes = append(nodes[y][x].used_hashes, hash_position(position, direction))

			// Need to do something different if we're splititng... Maybe?
			for _, output := range connection.outputs {
				follow_path(nodes, next_position, output)
			}
		}
	}
}

func count_energized(nodes [][]Node) int {
	energized := 0

	for _, row := range nodes {
		for _, node := range row {
			if node.visited > 0 {
				energized++
			}
		}
	}

	return energized
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

	if !PART2 {
		nodes := rows_to_nodes(all_rows)

		// We visited our first position!
		nodes[0][0].visited++

		follow_path(nodes, []int{-1, 0}, []int{1, 0})

		print_node_array(nodes)

		fmt.Println("Energized nodes: ", count_energized(nodes))
	}

	if PART2 {
		maximal_energized_nodes := 0
		best_starting_location := []int{0, 0}

		for i := 0; i < len(all_rows); i++ {
			// Launch from every possible point.
			launch_point := []int{-1, i}
			launch_direction := []int{1, 0}

			fmt.Println("Launching horizontally: ", i)

			// Reset the grid
			nodes := rows_to_nodes(all_rows)
			nodes[i][0].visited++

			follow_path(nodes, launch_point, launch_direction)
			number_energized := count_energized(nodes)

			if number_energized > maximal_energized_nodes {
				maximal_energized_nodes = number_energized
				best_starting_location = launch_point
			}

			// Launch from the right hand side
			launch_point = []int{len(all_rows[0]), i}
			launch_direction = []int{-1, 0}

			fmt.Println("Launching horizontally: ", -i)

			// Reset the grid
			nodes = rows_to_nodes(all_rows)
			nodes[i][len(all_rows[0])-1].visited++

			follow_path(nodes, launch_point, launch_direction)
			number_energized = count_energized(nodes)

			if number_energized > maximal_energized_nodes {
				maximal_energized_nodes = number_energized
				best_starting_location = launch_point
			}
		}

		// Now vertically
		for i := 0; i < len(all_rows[0]); i++ {
			// Launch from every possible point.
			launch_point := []int{i, -1}
			launch_direction := []int{0, 1}

			fmt.Println("Launching vertically: ", i)

			// Reset the grid
			nodes := rows_to_nodes(all_rows)
			nodes[0][i].visited++

			follow_path(nodes, launch_point, launch_direction)
			number_energized := count_energized(nodes)

			if number_energized > maximal_energized_nodes {
				maximal_energized_nodes = number_energized
				best_starting_location = launch_point
			}

			// Launch from the right hand side
			launch_point = []int{i, len(all_rows)}
			launch_direction = []int{0, -1}

			fmt.Println("Launching vertically: ", -i)

			// Reset the grid
			nodes = rows_to_nodes(all_rows)
			nodes[len(all_rows)-1][i].visited++

			follow_path(nodes, launch_point, launch_direction)
			number_energized = count_energized(nodes)

			if number_energized > maximal_energized_nodes {
				maximal_energized_nodes = number_energized
				best_starting_location = launch_point
			}
		}

		fmt.Println("Maximal energized nodes: ", maximal_energized_nodes, "from starting location:", best_starting_location)
	}

}
