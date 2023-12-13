package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func horizontal_and_vertical_not_equal(horizontal []int, vertical []int, new_grid []string) bool {
	new_vertical := find_vertical_reflection_points(new_grid)
	new_horizontal := find_horizontal_reflection_points(new_grid)

	if len(new_horizontal) == 0 && len(new_vertical) == 0 {
		return false
	}

	new_line := false

	for _, row := range new_horizontal {
		if !slices.Contains(horizontal, row) {
			fmt.Println("New horizontal reflection")
			new_line = true
		}
	}

	for _, column := range new_vertical {
		if !slices.Contains(vertical, column) {
			fmt.Println("New vertical reflection")
			new_line = true
		}
	}

	return new_line
}

func new_string_replace_at_index(s string, index int) string {
	new_char := "."

	if string(s[index]) == "." {
		new_char = "#"
	} else {
		new_char = "."
	}

	return s[:index] + new_char + s[index+1:]
}

func try_all_replacements(grid []string) ([]int, []int) {
	// First: find the current relfection lines:
	horizontal := find_horizontal_reflection_points(grid)
	vertical := find_vertical_reflection_points(grid)

	new_grid := make([]string, len(grid))

	for i, line := range grid {
		new_grid[i] = line
	}

	character_number := 0
	number_of_characters := len(grid) * len(grid[0])

	for !horizontal_and_vertical_not_equal(horizontal, vertical, new_grid) {
		fmt.Println("Changing character", character_number)
		// Flip that character from "." to "#" or vice versa.
		// First: undo what we did before.
		if character_number != 0 {
			row := (character_number - 1) / len(grid[0])
			column := (character_number - 1) % len(grid[0])

			new_grid[row] = new_string_replace_at_index(new_grid[row], column)
		}

		row := character_number / len(grid[0])
		column := character_number % len(grid[0])

		new_grid[row] = new_string_replace_at_index(new_grid[row], column)

		character_number += 1

		if character_number == number_of_characters {
			log.Fatal("Could not find new line of reflection.")
		}
	}

	new_horizontal := find_horizontal_reflection_points(new_grid)
	new_vertical := find_vertical_reflection_points(new_grid)

	return_horizontal := make([]int, 0)
	return_vertical := make([]int, 0)

	// Only use the _new_ ones
	for _, row := range new_horizontal {
		if !slices.Contains(horizontal, row) {
			return_horizontal = append(return_horizontal, row)
		}
	}

	for _, column := range new_vertical {
		if !slices.Contains(vertical, column) {
			return_vertical = append(return_vertical, column)
		}
	}

	return return_horizontal, return_vertical
}

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

			horizontal, vertical := try_all_replacements(grid)

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
