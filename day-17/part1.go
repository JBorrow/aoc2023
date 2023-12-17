package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DEBUG = true
const color_red = "\033[31m"
const color_green = "\033[32m"
const color_yellow = "\033[33m"
const color_none = "\033[0m"

type Node struct {
	cost        int
	visited     bool
	connections []string
	position    []int
	// Current heuristics
	g               int
	h               int
	best_connection string
	on_path         bool
}

type Grid struct {
	nodes  map[string]Node
	height int
	width  int
}

func node_hash(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func connections_from_position(x int, y int, max_x int, max_y int) []string {
	connections := make([]string, 0)

	if x > 0 {
		connections = append(connections, node_hash(x-1, y))
	}

	if x < max_x {
		connections = append(connections, node_hash(x+1, y))
	}

	if y > 0 {
		connections = append(connections, node_hash(x, y-1))
	}

	if y < max_y {
		connections = append(connections, node_hash(x, y+1))
	}

	return connections
}

func rows_to_nodes(rows []string) Grid {
	node_map := make(map[string]Node)

	max_x := len(rows[0]) - 1
	max_y := len(rows) - 1

	for y, row := range rows {
		for x, char := range row {
			// Get character value
			value, _ := strconv.Atoi(string(char))

			hash := node_hash(x, y)

			node_map[hash] = Node{
				cost:            value,
				visited:         false,
				connections:     connections_from_position(x, y, max_x, max_y),
				position:        []int{x, y},
				g:               0,
				h:               0,
				best_connection: "",
			}
		}
	}

	return Grid{
		nodes:  node_map,
		height: max_y + 1,
		width:  max_x + 1,
	}
}

func visualize_grid(grid Grid) {
	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			hash := node_hash(x, y)

			node := grid.nodes[hash]

			if node.on_path {
				fmt.Print(color_yellow, node.cost, color_none)
			} else if node.visited {
				fmt.Print(color_red, node.cost, color_none)
			} else {
				fmt.Print(color_green, node.cost, color_none)
			}
		}

		fmt.Println()
	}
}

func visualize_g(grid Grid) {
	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			hash := node_hash(x, y)

			node := grid.nodes[hash]

			if node.visited {
				fmt.Printf("%s%-4d%s", color_red, node.g, color_none)
			} else {
				fmt.Printf("%s%-4d%s", color_green, node.g, color_none)
			}
		}

		fmt.Println()
	}
}

func visualize_connections(grid Grid) {
	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			hash := node_hash(x, y)

			node := grid.nodes[hash]

			if node.on_path {
				fmt.Print(color_yellow)
			}

			if node.visited {
				// Find connection
				connected_node := grid.nodes[node.best_connection]

				connected_node_x := connected_node.position[0]
				connected_node_y := connected_node.position[1]

				if connected_node_x == x {
					fmt.Print("|")
				} else if connected_node_y == y {
					fmt.Print("-")
				} else {
					log.Fatal("Something went wrong")
				}
			} else {
				fmt.Print(" ")
			}

			if node.on_path {
				fmt.Print(color_none)
			}

		}
		fmt.Println()
	}
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Try with just manhattan distance at first?
func distance_heuristic(node Node, end Node) int {
	return intAbs(end.position[0]-node.position[0]) + intAbs(end.position[1]-node.position[1])
}

func weighted_manhattan_heuristic(node Node, end Node, grid Grid) int {
	start_x := min(node.position[0], end.position[0])
	start_y := min(node.position[1], end.position[1])
	end_x := max(node.position[0], end.position[0])
	end_y := max(node.position[1], end.position[1])

	// L-path cost. This is wrong!
	// total_cost := 0

	// for y := start_y; y <= end_y; y++ {
	// 	hash := node_hash(start_x, y)
	// 	total_cost += grid.nodes[hash].cost
	// }
	// for x := start_x; x <= end_x; x++ {
	// 	hash := node_hash(x, end_y)
	// 	total_cost += grid.nodes[hash].cost
	// }

	// return total_cost

	// For all nodes in the square, find the minimal cost of an item.
	min_cost := 1000000

	for y := start_y; y <= end_y; y++ {
		for x := start_x; x <= end_x; x++ {
			hash := node_hash(x, y)
			if grid.nodes[hash].cost < min_cost {
				min_cost = grid.nodes[hash].cost
			}
		}
	}

	// return min_cost * (end_x - start_x + end_y - start_y)
	return 0
}

// Try with some crazy pathing stuff...
func in_line_heuristic(node Node, prior_node Node, grid Grid) int {
	// If we have three in a row, the cost is massive.
	path := make([][]int, 0)

	path = append(path, node.position)

	current_node := prior_node

	counter := 0

	trace_length := 400

	for current_node.best_connection != "" && counter <= trace_length {
		path = append(path, current_node.position)
		current_node = grid.nodes[current_node.best_connection]
		counter += 1
	}

	if counter <= trace_length {
		// If we have less than three in a row, the cost is simple.
		if DEBUG {
			fmt.Println("Returning reasonable cost for path: ", path)
		}
		return node.cost + prior_node.g
	}

	x := path[0][0]
	y := path[0][1]

	all_xs_equal := true
	all_ys_equal := true

	for _, position := range path {
		if position[0] != x {
			all_xs_equal = false
		}
		if position[1] != y {
			all_ys_equal = false
		}
	}

	if all_xs_equal || all_ys_equal {
		if DEBUG {
			fmt.Println("Returning unreasonable cost for path: ", path)
		}

		return 1000000
	} else {
		if DEBUG {
			fmt.Println("Returning reasonable cost for path: ", path)
		}
		return node.cost + prior_node.g
	}

}

func consider_adjacent_nodes(grid Grid, node Node, end Node) []string {
	updated_nodes := make([]string, 0)

	for _, connection := range node.connections {
		// Get node
		new_node := grid.nodes[connection]

		if node.best_connection == connection {
			// This is the node that got us here. We don't want to consider it.
			continue
		}

		if !new_node.visited {
			// This guy is gonna get me whether it likes it or not.
			new_node.g = in_line_heuristic(new_node, node, grid)

			if new_node.g == 1000000 {
				// CANNOT VISIT
				continue
			}

			new_node.visited = true
			new_node.best_connection = node_hash(node.position[0], node.position[1])
			new_node.h = weighted_manhattan_heuristic(new_node, end, grid)

			grid.nodes[connection] = new_node
			updated_nodes = append(updated_nodes, connection)

			continue
		}

		// If not visited, calculate g-cost and h-cost
		new_g_cost := in_line_heuristic(new_node, node, grid)
		new_h_cost := weighted_manhattan_heuristic(new_node, end, grid)

		if new_g_cost == 1000000 {
			// CANNOT VISIT
			continue
		}

		if (new_g_cost + new_h_cost) > (new_node.g + new_node.h) {
			// We lose this time boys, sorry.
			continue
		}

		// I take over 😈
		new_node.g = new_g_cost
		new_node.h = new_h_cost
		new_node.best_connection = node_hash(node.position[0], node.position[1])

		// Set visited
		grid.nodes[connection] = new_node
		updated_nodes = append(updated_nodes, connection)
	}

	return updated_nodes
}

// We're gonna need a priority queue here
type QueueItem struct {
	hash  string
	cost  int
	index int
}

type NodeQueue []*QueueItem

func (nq NodeQueue) Len() int {
	return len(nq)
}

func (nq NodeQueue) Less(i, j int) bool {
	// This is the wrong way around on purpose. We want the items with the
	// SMALLEST f-cost to be at the top of the queue.
	return (nq[i].cost) < (nq[j].cost)
}

func (nq NodeQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
	nq[i].index = i
	nq[j].index = j
}

func (nq *NodeQueue) Push(x any) {
	item := x.(*QueueItem)
	item.index = len(*nq)
	*nq = append(*nq, item)
	heap.Fix(nq, item.index)
}

func (nq *NodeQueue) Pop() any {
	old := *nq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*nq = old[0 : n-1]
	item.index = -1
	return item
}

func (nq *NodeQueue) update(item *QueueItem) {
	// In our implemetation, the node is modified already. We just need to put it in
	// the right place in the queue.
	heap.Fix(nq, item.index)
}

func print_queue(nq NodeQueue) {
	fmt.Println("Current queue state: ")
	for _, item := range nq {
		fmt.Print("{hash:", item.hash, " cost:", item.cost, "}, ")
	}
	fmt.Println()
}

func traverse_graph(grid Grid, start string, end string) int {
	// Get start and end nodes
	start_node := grid.nodes[start]
	end_node := grid.nodes[end]

	// Priority queue
	to_process := make(NodeQueue, 1)

	// Add starting node to queue.
	if DEBUG {
		fmt.Println("Starting node: ", start_node.position)
		fmt.Println("Ending node: ", end_node.position)
	}
	to_process[0] = &QueueItem{hash: node_hash(start_node.position[0], start_node.position[1]), cost: 0, index: 0}
	if DEBUG {
		fmt.Println("Initiating priority queue")
	}
	heap.Init(&to_process)
	if DEBUG {
		fmt.Println("Priority queue initiated")
	}

	// While there are nodes to process
	for to_process.Len() > 0 {
		// Get the node with the lowest f-cost
		current_item := heap.Pop(&to_process).(*QueueItem)
		current_node := grid.nodes[current_item.hash]

		if DEBUG {
			fmt.Println("Current node: ", current_node.position)
		}

		// // If we've reached the end, we're done.
		// if slices.Equal(current_node.position, end_node.position) {
		// 	break
		// }

		// Consider adjacent nodes
		updated_nodes := consider_adjacent_nodes(grid, current_node, end_node)

		// Add updated nodes to queue
		for _, node := range updated_nodes {
			to_update := grid.nodes[node]
			to_process.Push(
				&QueueItem{
					hash: node_hash(to_update.position[0], to_update.position[1]),
					cost: to_update.g + to_update.h,
				},
			)
		}
	}

	if DEBUG {
		fmt.Println("Finished processing")
		print_queue(to_process)
	}

	// Now we're done, we need to highlight the best path by walking backwards
	// from the end node.
	start_node_hash := node_hash(start_node.position[0], start_node.position[1])
	current_node := grid.nodes[end]

	for node_hash(current_node.position[0], current_node.position[1]) != start_node_hash {
		current_node.on_path = true
		grid.nodes[node_hash(current_node.position[0], current_node.position[1])] = current_node
		current_node = grid.nodes[current_node.best_connection]
	}

	// Do not highlight the start node

	total_cost := 0
	starting_hash := node_hash(start_node.position[0], start_node.position[1])
	for hash, node := range grid.nodes {
		if node.on_path {
			if starting_hash != hash {
				total_cost += node.cost
			}
		}
	}

	return total_cost
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

	nodes := rows_to_nodes(all_rows)

	cost := traverse_graph(nodes, node_hash(0, 0), node_hash(nodes.width-1, nodes.height-1))

	fmt.Println("Total cost: ", cost)

	if DEBUG {
		visualize_grid(nodes)
		visualize_g(nodes)
		visualize_connections(nodes)
	}
}
