package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

// type Node struct {
// 	path_length int
// 	visited     bool
// 	links       []string
// }

type Node struct {
	left  string
	right string
}

var node_regex = regexp.MustCompile(`([0-9A-Z]*) = \(([0-9A-Z]*), ([0-9A-Z]*)\)`)

func process_node(to_parse string) (string, Node) {
	matches := node_regex.FindStringSubmatch(to_parse)

	if len(matches) != 4 {
		log.Fatal("Could not parse: ", to_parse)
	}

	name := matches[1]
	left := matches[2]
	right := matches[3]

	return name, Node{left, right}
}

func follow_route(route string, nodes map[string]Node, starting_node string, ending_node string) []int {
	// Returns the path length.

	path_length := 0
	current_node := starting_node

	find_n_overlaps := 10

	overlaps := make([]int, find_n_overlaps)

	for i := 0; i < find_n_overlaps; i++ {
		for string(current_node[2]) != "Z" {
			instruction := route[path_length%len(route)]

			if instruction == 'L' {
				current_node = nodes[current_node].left
			} else if instruction == 'R' {
				current_node = nodes[current_node].right
			} else {
				log.Fatal("Unknown instruction: ", instruction)
			}

			path_length += 1
		}

		overlaps[i] = path_length

		instruction := route[path_length%len(route)]

		if instruction == 'L' {
			current_node = nodes[current_node].left
		} else if instruction == 'R' {
			current_node = nodes[current_node].right
		} else {
			log.Fatal("Unknown instruction: ", instruction)
		}

		path_length += 1
	}

	return overlaps
}

func follow_routes(route string, nodes map[string]Node, starting_nodes []string, ending_ndoes []string) [][]int {
	// Returns the path length.

	path_lengths := make([][]int, len(starting_nodes))

	for i, starting_node := range starting_nodes {
		path_lengths[i] = follow_route(route, nodes, starting_node, "ZZZ")
	}

	return path_lengths
}

func all_nodes_ending_in(graph map[string]Node, char string) []string {
	nodes := make([]string, 0)

	for k := range graph {
		if string(k[2]) == char {
			nodes = append(nodes, k)
		}
	}

	return nodes
}

// Get all prime factors of a given number n
func PrimeFactors(n int) (pfs []int) {
	// Get the number of 2s that divide n
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n = n / 2
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	for i := 3; i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			pfs = append(pfs, i)
			n = n / i
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n > 2 {
		pfs = append(pfs, n)
	}

	return
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	number_of_lines := 0
	route := ""
	nodes := make(map[string]Node)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if number_of_lines == 0 {
			route = text
			number_of_lines += 1
			continue
		}

		if len(text) < 2 {
			continue
		}

		name, node := process_node(text)

		if DEBUG {
			fmt.Println("Given: ", text)
			fmt.Println("Parsed to: ", node)
		}

		nodes[name] = node

		number_of_lines += 1
	}

	if DEBUG {
		fmt.Println("Route: ", route)
	}

	if scanner.Err() != nil {
		log.Fatal("No input provided.")
	}

	starting_nodes := all_nodes_ending_in(nodes, "A")
	ending_nodes := all_nodes_ending_in(nodes, "Z")

	if DEBUG {
		fmt.Println("Starting nodes: ", starting_nodes)
		fmt.Println("Ending nodes: ", ending_nodes)
	}

	path_length := follow_routes(route, nodes, starting_nodes, ending_nodes)

	for _, path := range path_length {
		fmt.Println("Path lengths: ", path)
	}

	unique_prime_factors := make([]int, 0)
	lcm := 1

	for _, path := range path_length {
		loop_lengths := make([]int, len(path))

		for i, v := range path {
			loop_lengths[i] = v - path[max(0, i-1)]
		}

		fmt.Println("Loop lengths: ", loop_lengths)

		this_loop_length := loop_lengths[1]

		prime_factors := PrimeFactors(this_loop_length)

		fmt.Println("Prime factors: ", prime_factors)

		for _, v := range prime_factors {
			if !slices.Contains(unique_prime_factors, v) {
				unique_prime_factors = append(unique_prime_factors, v)
				lcm *= v
			}
		}
	}

	fmt.Println("Unique prime factors: ", unique_prime_factors)
	fmt.Println("Smallest factorization: ", lcm)
}
