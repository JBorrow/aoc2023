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
	used        []bool
	visited     bool
	outside     bool
	inside      bool
}

var color_red = "\033[31m"
var color_none = "\033[0m"
var color_green = "\033[32m"
var color_cyan = "\033[36m"

var unicode_map = map[string]string{
	"|": "│",
	"-": "─",
	"L": "└",
	"J": "┘",
	"7": "┐",
	"F": "┌",
	".": "█",
	"S": "╋",
}

var render_map = map[string][][]uint8{
	"|": {
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	},
	"-": {
		{0, 0, 0},
		{1, 1, 1},
		{0, 0, 0},
	},
	"L": {
		{0, 1, 0},
		{0, 1, 1},
		{0, 0, 0},
	},
	"J": {
		{0, 1, 0},
		{1, 1, 0},
		{0, 0, 0},
	},
	"7": {
		{0, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	},
	"F": {
		{0, 0, 0},
		{0, 1, 1},
		{0, 1, 0},
	},
	".": {
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	},
	"S": {
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	},
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
		used:        make([]bool, len(direction_map[symbol])),
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
			if node.visited {
				fmt.Print(color_red, node.render, color_none)
			} else if node.outside {
				fmt.Print(color_green, node.render, color_none)
			} else if node.inside {
				fmt.Print(color_cyan, node.render, color_none)
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
			if node.visited {
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
	beginning_x := 1
	beginning_y := 1

	for y, row := range nodes {
		for x, node := range row {
			if node.symbol == "S" {
				beginning_x = x
				beginning_y = y
			}
		}
	}

	for core_iteration := 0; core_iteration < 2; core_iteration++ {
		start_x := beginning_x
		start_y := beginning_y

		prev_x := beginning_x
		prev_y := beginning_y

		// Reset the visits
		for y, row := range nodes {
			for x := range row {
				nodes[y][x].visited = false
			}
		}

		current_symbol := "NONE"
		nodes[start_y][start_x].visited = true
		nodes[start_y][start_x].distance = 0

		iterations := 0

		// fmt.Println("Starting node: ", nodes[start_y][start_x])

		for current_symbol != "S" {
			// fmt.Println("Current symbol: ", current_symbol)
			current_node := nodes[start_y][start_x]

			for direction_id, direction := range current_node.connections {
				// fmt.Println("Iteration: ", iterations)

				if iterations > 250000 {
					fmt.Println("Too many iterations!")
					current_symbol = "S"
					break
				} else {
					iterations++
				}

				if current_node.used[direction_id] && current_node.symbol == "S" {
					continue
				}
				// fmt.Println("Checking direction: ", direction, "Have symbol: ", current_symbol)
				// Does direction take us out of bounds?
				new_x := start_x + direction[0]
				new_y := start_y + direction[1]

				if new_x < 0 || new_y < 0 || new_x >= len(nodes[0]) || new_y >= len(nodes) {
					// fmt.Println("Out of bounds")
					continue
				}

				if new_x == prev_x && new_y == prev_y {
					continue
				}

				// Does the node we are going to have a symbol that accepts our connection?
				allowed := false

				for _, connection := range nodes[new_y][new_x].connections {
					// fmt.Println("Checking connection: ", connection, "Against: ", []int{direction[0] * -1, direction[1] * -1})
					if slices.Equal(connection, []int{direction[0] * -1, direction[1] * -1}) {
						// fmt.Println("Found connection!")
						allowed = true
						break
					}
				}

				if !allowed {
					continue
				}

				// Is this new node the source and I have a large distance?
				if nodes[new_y][new_x].symbol == "S" && nodes[start_y][start_x].distance >= 1 {
					current_symbol = "S"
					break
				}

				// Is this new node unvisited?
				if nodes[new_y][new_x].visited {
					// fmt.Println("Already visited...")
					continue
				}

				// symbols := make([][]string, 3)

				// for row_id := 0; row_id < 3; row_id++ {
				// 	symbols[row_id] = []string{" ", " ", " "}
				// }

				// symbols[1][1] = current_node.render
				// symbols[1+direction[1]][1+direction[0]] = nodes[new_y][new_x].render

				// fmt.Println("We think that the following symbols are connected: ")
				// for _, row := range symbols {
				// 	fmt.Println(row)
				// }

				// We used this connection!
				nodes[start_y][start_x].used[direction_id] = true

				// Update our current position.
				prev_x = start_x
				prev_y = start_y

				start_x = new_x
				start_y = new_y

				if nodes[start_y][start_x].distance == -1 {
					nodes[start_y][start_x].distance = current_node.distance + 1
				} else {
					// fmt.Println("Found a shorter distance!", nodes[start_y][start_x].distance, current_node.distance+1)
					nodes[start_y][start_x].distance = min(current_node.distance+1, nodes[start_y][start_x].distance)
				}

				current_symbol = nodes[start_y][start_x].symbol
				nodes[start_y][start_x].visited = true

				// We found the connection! Onto the next symbol.
				break
			}
		}
	}

	return
}

// Watershed around the edges of the map.
func watershed_around(nodes [][]Node, start_x int, start_y int) {
	// First check my points.
	my_node := nodes[start_y][start_x]

	if my_node.visited || my_node.outside {
		return
	}

	nodes[start_y][start_x].outside = true
	nodes[start_y][start_x].render = "█"

	// Now check my neighbors.
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			new_x := start_x + x
			new_y := start_y + y

			if new_x < 0 || new_y < 0 || new_x >= len(nodes[0]) || new_y >= len(nodes) {
				continue
			}

			watershed_around(nodes, new_x, new_y)
		}
	}
}

func watershed_all_edges(nodes [][]Node) {
	// Watershed around all possible edge points.
	for x := 0; x < len(nodes[0]); x++ {
		watershed_around(nodes, x, 0)
		watershed_around(nodes, x, len(nodes)-1)
	}

	for y := 0; y < len(nodes); y++ {
		watershed_around(nodes, 0, y)
		watershed_around(nodes, len(nodes[0])-1, y)
	}
}

func watershed_around_integers(nodes [][]uint8, start_x int, start_y int) {
	// First check my points.
	my_node := nodes[start_y][start_x]

	if my_node > 0 {
		return
	}

	nodes[start_y][start_x] = 2

	// Now check my neighbors.
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			new_x := start_x + x
			new_y := start_y + y

			if new_x < 0 || new_y < 0 || new_x >= len(nodes[0]) || new_y >= len(nodes) {
				continue
			}

			watershed_around_integers(nodes, new_x, new_y)
		}
	}
}

func watershed_all_edges_integers(nodes [][]uint8) {
	// Watershed around all possible edge points.
	for x := 0; x < len(nodes[0]); x++ {
		watershed_around_integers(nodes, x, 0)
		watershed_around_integers(nodes, x, len(nodes)-1)
	}

	for y := 0; y < len(nodes); y++ {
		watershed_around_integers(nodes, 0, y)
		watershed_around_integers(nodes, len(nodes[0])-1, y)
	}
}

func label_unvisited(nodes [][]Node) int {
	number_unvisited := 0

	for y, row := range nodes {
		for x := range row {
			if !(nodes[y][x].visited || nodes[y][x].outside) {
				nodes[y][x].inside = true
				nodes[y][x].render = "█"
				number_unvisited += 1
			}
		}
	}

	return number_unvisited
}

func render_to_2d_array(nodes [][]Node) [][]uint8 {
	render := make([][]uint8, 3*len(nodes))

	for y, _ := range nodes {
		for i := 0; i < 3; i++ {
			render[y*3+i] = make([]uint8, 3*len(nodes[0]))
		}
	}

	for y, row := range nodes {
		for x := range row {
			if nodes[y][x].visited {
				for i, row := range render_map[nodes[y][x].symbol] {
					for j, value := range row {
						render[y*3+i][x*3+j] = value
					}
				}
			}
		}
	}

	return render
}

func extract_watershed_status(nodes [][]Node, render_array [][]uint8) int {
	total_number_inside := 0

	for y, row := range nodes {
		for x := range row {
			rendered_value := render_array[y*3+1][x*3+1]
			if rendered_value == 2 {
				nodes[y][x].outside = true
			} else if rendered_value == 1 {
				nodes[y][x].visited = true
			} else if rendered_value == 0 {
				total_number_inside += 1
				nodes[y][x].inside = true
			}
		}
	}

	return total_number_inside
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

	// Now watershed

	watershed_all_edges(nodes)

	// Now label all unvisited nodes.

	number_unvisited := label_unvisited(nodes)

	if DEBUG {
		render_node_grid(nodes)
	}

	fmt.Println("Number of unvisited nodes: ", number_unvisited)
	fmt.Println("Warning - do not use this number - it does not watershed correctly into literal corner cases.")

	rendered_array := render_to_2d_array(nodes)
	watershed_all_edges_integers(rendered_array)

	number_unvisited = extract_watershed_status(nodes, rendered_array)

	if DEBUG {
		fmt.Println("Full rendered grid: ")
		for _, row := range rendered_array {
			fmt.Println(row)
		}
	}

	if DEBUG {
		render_node_grid(nodes)
	}

	fmt.Println("Number of unvisited nodes (true): ", number_unvisited)

}
