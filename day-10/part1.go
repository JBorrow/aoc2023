package main

import (
	"bufio"
	"fmt"
	"log"
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
	connections     [][]int
	connection_used []bool
	symbol          string
	render          string
	distance        int
}

var color_red = "\033[31m"

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
		{0, 1},
		{1, 0},
	},
	"J": {
		{0, 1},
		{-1, 0},
	},
	"7": {
		{0, -1},
		{-1, 0},
	},
	"F": {
		{0, -1},
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
		connections:     direction_map[symbol],
		symbol:          symbol,
		render:          unicode_map[symbol],
		distance:        -1,
		connection_used: make([]bool, len(direction_map[symbol])),
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
				fmt.Print(color_red, node.render)
			} else {
				fmt.Print(node.render)
			}
		}
		fmt.Println()
	}
}

func find_connections(nodes [][]Node) {
	// First, find the node with the S symbol.
	var start_x, start_y int

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

		// Look around our possible directions.
		for direction_id, direction := range current_node.connections {
			// Check if we have already used this connection.
			if current_node.connection_used[direction_id] {
				continue
			}

			// Check if we can go in this direction.
			if len(direction) == 0 {
				log.Fatal("Could not find a direction to go in.")
			}

			if iterations > 100 {
				current_symbol = "S"
				break
			} else {
				iterations++
			}

			fmt.Println("Iteration: ", iterations, "Current position: ", start_x, " ", start_y, " Symbol: ", current_symbol, " Current distance: ", current_node.distance)

			// Calculate the new coordinates.
			new_x := start_x + direction[0]
			new_y := start_y + direction[1]

			// Check if we are still in bounds.
			if new_x < 0 || new_y < 0 || new_x >= len(nodes[0]) || new_y >= len(nodes) {
				continue
			}

			potential_node := nodes[new_y][new_x]
			required_direction := []int{direction[0] * -1, direction[1] * -1}

			// Does the node we are going to have a symbol that accepts our connection?
			// Is it unvisited?
			allowed := false
			connection_id := 10

			for id, connection := range potential_node.connections {
				if slices.Equal(connection, required_direction) {
					allowed = true
					connection_id = id
					break
				}
			}

			if !allowed {
				continue
			}

			// Mark my and their connection as used.
			nodes[start_y][start_x].connection_used[direction_id] = true
			nodes[new_y][new_x].connection_used[connection_id] = true

			// We can move in this direction.
			nodes[new_y][new_x].distance = current_node.distance + 1

			// Update our current position.
			start_x = new_x
			start_y = new_y

			// Update our current symbol.
			current_symbol = nodes[new_y][new_x].symbol
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
	}
}
