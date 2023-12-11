package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Galaxy struct {
	x       int
	y       int
	id      int
	shift_x int
	shift_y int
}

// This is 2 for part 1, 1000000 for part 2
var SHIFT_FACTOR = 1000000

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
				galaxy := Galaxy{x, y, len(galaxies), 0, 0}
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

func insert_extra_space_x(galaxies []Galaxy) {
	inject_space_at := make([]int, 0)

	maximal_x := 0

	for _, galaxy := range galaxies {
		if galaxy.x > maximal_x {
			maximal_x = galaxy.x
		}
	}

	for x := 0; x <= maximal_x; x++ {
		no_galaxies := true

		for _, galaxy := range galaxies {
			if galaxy.x == x {
				no_galaxies = false
				break
			}
		}

		if !no_galaxies {
			continue
		} else {
			inject_space_at = append(inject_space_at, x)
		}
	}

	for _, x := range inject_space_at {
		for i, galaxy := range galaxies {
			if galaxy.x > x {
				galaxies[i].shift_x += 1
			}
		}
	}

	// Apply the shift
	for i, galaxy := range galaxies {
		galaxies[i].x += galaxy.shift_x * (SHIFT_FACTOR - 1)
	}

	return
}

func insert_extra_space_y(galaxies []Galaxy) {
	inject_space_at := make([]int, 0)

	maximal_y := 0

	for _, galaxy := range galaxies {
		if galaxy.x > maximal_y {
			maximal_y = galaxy.y
		}
	}

	for y := 0; y <= maximal_y; y++ {
		no_galaxies := true

		for _, galaxy := range galaxies {
			if galaxy.y == y {
				no_galaxies = false
				break
			}
		}

		if !no_galaxies {
			continue
		} else {
			inject_space_at = append(inject_space_at, y)
		}
	}

	for _, y := range inject_space_at {
		for i, galaxy := range galaxies {
			if galaxy.y > y {
				galaxies[i].shift_y += 1
			}
		}
	}

	// Apply the shift
	for i, galaxy := range galaxies {
		galaxies[i].y += galaxy.shift_y * (SHIFT_FACTOR - 1)
	}

	return
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

	galaxies := extract_galaxies(grid)

	if DEBUG {
		fmt.Println("Galaxies:")
		for _, galaxy := range galaxies {
			fmt.Println(galaxy)
		}
	}

	// Shift the galaxies
	insert_extra_space_x(galaxies)
	insert_extra_space_y(galaxies)

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
