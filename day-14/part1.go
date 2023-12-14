package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var mapping = map[string]int{
	".": 0,
	"#": 1,
	"O": 2,
}

var vis_mapping = map[int]string{
	0: " ",
	1: "█",
	2: "◯",
}

func parse_grid(lines []string) [][]int {
	grid := make([][]int, len(lines))

	for i, line := range lines {
		grid[i] = make([]int, len(line))

		for j, char := range line {
			grid[i][j] = mapping[string(char)]
		}
	}

	return grid
}

func vis_grid(grid [][]int) {
	fmt.Print("┏")
	for i := 0; i < len(grid[0]); i++ {
		fmt.Print("━")
	}
	fmt.Println("┓")
	for _, line := range grid {
		fmt.Print("┃")
		for _, char := range line {
			fmt.Print(vis_mapping[char])
		}
		fmt.Println("┃")
	}
	fmt.Print("┗")
	for i := 0; i < len(grid[0]); i++ {
		fmt.Print("━")
	}
	fmt.Println("┛")
}

func propagate_all_balls_north(grid [][]int) [][]int {
	new_grid := make([][]int, len(grid))

	for i, line := range grid {
		new_grid[i] = make([]int, len(line))
		for j, char := range line {
			new_grid[i][j] = char
		}
	}

	for i, line := range grid {
		for j, char := range line {
			if char == 2 {
				// Try to move it as far northwards as it will go
				for k := i; k > 0; k-- {
					if new_grid[k-1][j] == 0 {
						new_grid[k-1][j] = 2
						new_grid[k][j] = 0
					} else {
						break
					}
				}
			}
		}
	}

	return new_grid
}

func score_grid(grid [][]int) int {
	line_score := len(grid)

	score := 0

	for _, line := range grid {
		for _, char := range line {
			if char == 2 {
				score += line_score
			}
		}
		line_score -= 1
	}

	return score
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	grid := make([]string, 0)

	total_score := 0

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			// Process grid and move on.

			if len(grid) < 2 {
				grid = make([]string, 0)
				continue
			}

			parsed_grid := parse_grid(grid)
			propagated_grid := propagate_all_balls_north(parsed_grid)
			grid_score := score_grid(propagated_grid)

			if DEBUG {
				fmt.Println("Original Grid:")
				vis_grid(parsed_grid)
				fmt.Println("Propagated Grid:")
				vis_grid(propagated_grid)
				fmt.Println("Score: ", grid_score)
			}

			total_score += grid_score

			// Reset grid
			grid = make([]string, 0)
		} else {
			// Add to grid.
			grid = append(grid, text)
		}
	}

	fmt.Println("Total score: ", total_score)
}
