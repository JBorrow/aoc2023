package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// | is a vertical pipe connecting north and south.
// - is a horizontal pipe connecting east and west.
// L is a 90-degree bend connecting north and east.
// J is a 90-degree bend connecting north and west.
// 7 is a 90-degree bend connecting south and west.
// F is a 90-degree bend connecting south and east.
// . is ground; there is no pipe in this tile.
// S is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.

type Node struct {
	connections [][]int
	symbol      string
	render      string
	distance    int
	visited     bool
}

var color_red = "\033[31m"
var color_none = "\033[0m"

var unicode_map = map[string]string{
	"|": "│",
	"-": "─",
	"L": "└",
	"J": "┘",
	"7": "┐",
	"F": "┌",
	".": " ",
	"S": "╋",
}

var direction_map = map[string][][]int{
	"|": {
		{0, -1},
		{0, 1},
	},
	"-": {
		{-1, 0},
		{1, 0},
	},
	"L": {
		{0, -1},
		{1, 0},
	},
	"J": {
		{0, -1},
		{-1, 0},
	},
	"7": {
		{0, 1},
		{-1, 0},
	},
	"F": {
		{0, 1},
		{1, 0},
	},
	".": {{}},
	"S": {
		{0, -1},
		{0, 1},
		{-1, 0},
		{1, 0},
	},
}

func string_to_node(symbol string) Node {
	return Node{
		connections: direction_map[symbol],
		symbol:      symbol,
		render:      unicode_map[symbol],
		distance:    -1,
		visited:     false,
	}
}

func line_to_nodes(line string) []Node {
	line = strings.TrimSpace(line)

	nodes := make([]Node, len(line))

	for i, symbol := range line {
		nodes[i] = string_to_node(string(symbol))
	}

	return nodes
}

func render_node_grid(nodes [][]Node) {
	for _, row := range nodes {
		for _, node := range row {
			if node.distance > -1 {
				fmt.Print(color_red, node.render, color_none)
			} else {
				fmt.Print(node.render)
			}
		}
		fmt.Println()
	}
}

func render_node_grid_distances(nodes [][]Node) {
	for _, row := range nodes {
		for _, node := range row {
			if node.distance > -1 {
				fmt.Print(color_red)
				if node.distance > 9 {
					fmt.Print("X")
				} else {
					fmt.Print(node.distance)
				}
				fmt.Print()
			} else {
				fmt.Print("N")
			}

			fmt.Print(color_none)
		}
		fmt.Println()
	}
}

func find_connections(nodes [][]Node) {
	// First, find the node with the S symbol.
	start_x := 1
	start_y := 1

	for y, row := range nodes {
		for x, node := range row {
			if node.symbol == "S" {
				start_x = x
				start_y = y
			}
		}
	}

	current_symbol := "NONE"
	nodes[start_y][start_x].distance = 0

	iterations := 0

	for current_symbol != "S" {
		current_node := nodes[start_y][start_x]

		for _, direction := range current_node.connections {
			// fmt.Println("Checking direction: ", direction, "Have symbol: ", current_symbol)
			// Does direction take us out of bounds?
			new_x := start_x + direction[0]
			new_y := start_y + direction[1]

			if new_x < 0 || new_y < 0 || new_x >= len(nodes[0]) || new_y >= len(nodes) {
				continue
			}

			fmt.Println("Iteration: ", iterations)

			if iterations > 50000 {
				fmt.Println("Too many iterations!")
				current_symbol = "S"
				break
			} else {
				iterations++
			}

			// if iterations > 50 {
			// 	fmt.Println("Too many iterations!")
			// 	current_symbol = "S"
			// 	break
			// } else {
			// 	iterations++

			// 	fmt.Println("Currently at: ", start_x, start_y, "Going to: ", new_x, new_y, "Direction: ", direction, "Symbol: ", current_symbol)
			// }

			// Does the node we are going to have a symbol that accepts our connection?
			allowed := false

			for _, connection := range nodes[new_y][new_x].connections {
				// fmt.Println("Checking connection: ", connection, "Against: ", []int{direction[0] * -1, direction[1] * -1})
				if slices.Equal(connection, []int{direction[0] * -1, direction[1] * -1}) {
					allowed = true
					break
				}
			}

			if !allowed {
				continue
			}

			// Is this new node the source and I have a large distance?
			if nodes[new_y][new_x].symbol == "S" && nodes[start_y][start_x].distance > 1 {
				current_symbol = "S"
				break
			}

			// Is this new node unvisited?
			if nodes[new_y][new_x].distance > -1 {
				continue
			}

			symbols := make([][]string, 3)

			for row_id := 0; row_id < 3; row_id++ {
				symbols[row_id] = []string{" ", " ", " "}
			}

			symbols[1][1] = current_node.render
			symbols[1+direction[1]][1+direction[0]] = nodes[new_y][new_x].render

			fmt.Println("We think that the following symbols are connected: ")
			for _, row := range symbols {
				fmt.Println(row)
			}

			// Update our current position.
			start_x = new_x
			start_y = new_y

			nodes[start_y][start_x].distance = current_node.distance + 1

			current_symbol = nodes[start_y][start_x].symbol

			// We found the connection! Onto the next symbol.
			break
		}
	}

	return
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	nodes := make([][]Node, 0)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		node := line_to_nodes(text)

		if DEBUG {
			fmt.Println("Given: ", text)
		}

		nodes = append(nodes, node)
	}

	find_connections(nodes)

	if DEBUG {
		render_node_grid(nodes)
		render_node_grid_distances(nodes)
	}
}
