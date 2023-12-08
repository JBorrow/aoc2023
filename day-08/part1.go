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

var node_regex = regexp.MustCompile(`([A-Z]*) = \(([A-Z]*), ([A-Z]*)\)`)

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

func follow_route(route string, nodes map[string]Node, starting_node string, ending_node string) int {
	// Returns the path length.

	path_length := 0
	current_node := starting_node

	for current_node != ending_node {
		instruction := route[path_length%len(route)]

		if instruction == 'L' {
			current_node = nodes[current_node].left
		} else if instruction == 'R' {
			current_node = nodes[current_node].right
		} else {
			log.Fatal("Unknown instruction: ", instruction)
		}

		path_length += 1

		if path_length > 1e6 {
			log.Fatal("Path length exceeded 1e6.")
		}

	}

	if ending_node != current_node {
		log.Fatal("Route did not end at the correct node. Expected: ", ending_node, " Got: ", current_node)
	}

	return path_length
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

	// Now we can follow the route.
	path_length := follow_route(route, nodes, "AAA", "ZZZ")

	fmt.Println("Path length: ", path_length)

}
