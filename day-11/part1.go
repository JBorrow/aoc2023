package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Galaxy struct {
	x  int
	y  int
	id int
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func minkowski_norm(a Galaxy, b Galaxy) int {
	return intAbs(a.x-b.x) + intAbs(a.y-b.y)
}

func extract_galaxies(grid []string) []Galaxy {
	galaxies := make([]Galaxy, 0)

	for y, row := range grid {
		for x, column := range row {
			if string(column) == "#" {
				galaxy := Galaxy{x, y, len(galaxies)}
				galaxies = append(galaxies, galaxy)
			}
		}
	}

	return galaxies
}

func all_pair_distances(galaxies []Galaxy) []int {
	distances := make([]int, 0)

	max_extract_id := 0

	for i, a := range galaxies {

		max_extract_id = i

		for j, b := range galaxies {
			if i == j {
				continue
			}

			if j > max_extract_id {
				continue
			}

			distances = append(distances, minkowski_norm(a, b))
		}
	}

	return distances
}

func insert_extra_columns(grid []string) []string {
	current_column := 0
	total_columns := len(grid[0])

	for current_column < total_columns {
		// Loop through all the rows and find out if we have a
		// column that is all empty (i.e. contains only ".")

		empty := true

		for _, row := range grid {
			if string(row[current_column]) != "." {
				empty = false
				break
			}
		}

		if !empty {
			current_column++
			continue
		}

		// We have an empty column, so we need to insert a new column
		// at the current position
		for id, row := range grid {
			left := row[:current_column]
			right := row[current_column:]
			new_row := left + "." + right

			grid[id] = new_row
		}

		current_column += 2
		total_columns += 1
	}

	return grid
}

func insert_extra_rows(grid []string) []string {
	current_row := 0
	row_length := len(grid[0])
	total_rows := len(grid)

	for current_row < total_rows {
		// Loop through all the rows and find out if we have a
		// row that is all empty (i.e. contains only ".")

		empty := true

		for _, column := range grid[current_row] {
			if string(column) != "." {
				empty = false
				break
			}
		}

		if !empty {
			current_row++
			continue
		}

		// We have an empty row, so we need to insert a new row
		// at the current position
		left := grid[:current_row+1]
		right := grid[current_row:]
		new_row := strings.Repeat(".", row_length)

		grid = append(left, right...)
		grid[current_row] = new_row

		current_row += 2
		total_rows += 1
	}

	return grid
}

func insert_extra_space(grid []string) []string {
	grid = insert_extra_columns(grid)
	grid = insert_extra_rows(grid)

	return grid
}

func main() {
	DEBUG := true

	scanner := bufio.NewScanner(os.Stdin)

	grid := make([]string, 0)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if DEBUG {
			fmt.Println("Given: ", text)
		}

		grid = append(grid, text)
	}

	if DEBUG {
		fmt.Println("Before adding:")
		for _, row := range grid {
			fmt.Println(row)
		}
	}

	grid = insert_extra_space(grid)

	if DEBUG {
		fmt.Println("After adding:")
		for _, row := range grid {
			fmt.Println(row)
		}
	}

	galaxies := extract_galaxies(grid)

	if DEBUG {
		fmt.Println("Galaxies:")
		for _, galaxy := range galaxies {
			fmt.Println(galaxy)
		}
	}

	distances := all_pair_distances(galaxies)

	if DEBUG {
		fmt.Println("Distances:")
		fmt.Println(distances)
		fmt.Println("Found", len(distances), "pairs")
	}

	total_distances := 0

	for _, distance := range distances {
		total_distances += distance
	}

	fmt.Println(total_distances)

}
