package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

func ending_condition(strings []string) bool {
	for _, s := range strings {
		if string(s[2]) != "Z" {
			return true
		}
	}
	return false
}

func follow_route(route string, nodes map[string]Node, starting_nodes []string) int {
	// Returns the path length.

	path_length := 0
	current_nodes := make([]string, len(starting_nodes))

	for i, starting_node := range starting_nodes {
		current_nodes[i] = starting_node
	}

	for ending_condition(current_nodes) {
		instruction := route[path_length%len(route)]

		if instruction == 'L' {
			for i, node := range current_nodes {
				current_nodes[i] = nodes[node].left
			}
		} else if instruction == 'R' {
			for i, node := range current_nodes {
				current_nodes[i] = nodes[node].right
			}
		} else {
			log.Fatal("Unknown instruction: ", instruction)
		}

		path_length += 1

		if path_length > 1e10 {
			log.Fatal("Path length exceeded 1e6.")
		}

	}

	return path_length
}

func all_nodes_ending_in_A(graph map[string]Node) []string {
	nodes := make([]string, 0)

	for k, _ := range graph {
		if string(k[2]) == "A" {
			nodes = append(nodes, k)
		}
	}

	return nodes
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

	starting_nodes := all_nodes_ending_in_A(nodes)

	if DEBUG {
		fmt.Println("Starting nodes: ", starting_nodes)
	}

	// Now we can follow the route.
	path_length := follow_route(route, nodes, starting_nodes)

	fmt.Println("Path length: ", path_length)

}
