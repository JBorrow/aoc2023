package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func find_vertical_reflection_points(grid []string) []int {
	// Find all points in the array where there is a vertical reflection
	// (i.e. all rows above mirror those below.

	// First step: find all possible points where two rows are identical.

	identical_rows := make([]int, 0)
	number_of_rows := len(grid)

	for i := 0; i < number_of_rows-1; i++ {
		if grid[i] == grid[i+1] {
			identical_rows = append(identical_rows, i)
		}
	}

	// fmt.Println("Found identical rows:", identical_rows)

	reflection_points := make([]int, 0)

	for _, row := range identical_rows {
		// Check if all the rows above and below are identical.
		for i := 0; i < number_of_rows; i++ {
			// fmt.Println("Checking rows", row-i, row+i+1)
			if ((row + i + 1) == number_of_rows) || (row-i) == -1 {
				// We found a match.
				// fmt.Println("Found a match at row", row)
				reflection_points = append(reflection_points, row)
				break
			}

			// Check if the rows are identical.
			// AAAAAAA
			// BBBBBBB <- row
			// BBBBBBB
			// AAAAAAA

			above := grid[row-i]
			below := grid[row+i+1]

			if above != below {
				// fmt.Println("Found that rows", above, below, "are not identical.")
				break
			}
		}
	}

	return reflection_points
}

func find_horizontal_reflection_points(grid []string) []int {
	// First step: re-format the grid into a slice of columns.
	new_grid := make([]string, len(grid[0]))

	for i := 0; i < len(grid[0]); i++ {
		line := ""
		for j := 0; j < len(grid); j++ {
			line += string(grid[j][i])
		}
		new_grid[i] = line
	}

	return find_vertical_reflection_points(new_grid)
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	grid := make([]string, 0)

	summary := 0

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			// Process grid and move on.

			if len(grid) < 2 {
				grid = make([]string, 0)
				continue
			}

			horizontal := find_horizontal_reflection_points(grid)
			vertical := find_vertical_reflection_points(grid)

			if DEBUG {
				fmt.Println("Grid:")
				for i, line := range grid {
					fmt.Println(line, i)
				}
				fmt.Println("Horizontal reflection points: ", horizontal)
				fmt.Println("Verical reflection points: ", vertical)
			}

			for _, row := range horizontal {
				summary += 1 * (row + 1)
			}

			for _, column := range vertical {
				summary += 100 * (column + 1)
			}

			// Reset grid
			grid = make([]string, 0)
		} else {
			// Add to grid.
			grid = append(grid, text)
		}
	}

	fmt.Println("Summary:", summary)
}
