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

func propagate_all_balls_east(grid [][]int) [][]int {
	new_grid := make([][]int, len(grid))

	for i, line := range grid {
		new_grid[i] = make([]int, len(line))
		for j, char := range line {
			new_grid[i][j] = char
		}
	}

	for i, line := range grid {
		for j := len(line) - 1; j >= 0; j-- {
			char := line[j]
			if char == 2 {
				// Try to move it as far eastwards as it will go
				for k := j; k < len(line)-1; k++ {
					if new_grid[i][k+1] == 0 {
						new_grid[i][k+1] = 2
						new_grid[i][k] = 0
					} else {
						break
					}
				}
			}
		}
	}

	return new_grid
}

func propagate_all_balls_west(grid [][]int) [][]int {
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
				// Try to move it as far westwards as it will go
				for k := j; k > 0; k-- {
					if new_grid[i][k-1] == 0 {
						new_grid[i][k-1] = 2
						new_grid[i][k] = 0
					} else {
						break
					}
				}
			}
		}
	}

	return new_grid
}

func propagate_all_balls_south(grid [][]int) [][]int {
	new_grid := make([][]int, len(grid))

	for i, line := range grid {
		new_grid[i] = make([]int, len(line))
		for j, char := range line {
			new_grid[i][j] = char
		}
	}

	for i := len(grid) - 1; i >= 0; i-- {
		line := grid[i]
		for j, char := range line {
			if char == 2 {
				// Try to move it as far southwards as it will go
				for k := i; k < len(grid)-1; k++ {
					if new_grid[k+1][j] == 0 {
						new_grid[k+1][j] = 2
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

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			// Process grid and move on.

			if len(grid) < 2 {
				grid = make([]string, 0)
				continue
			}

			parsed_grid := parse_grid(grid)
			propagated_grid := parsed_grid

			if DEBUG {
				fmt.Println("Original Grid:")
				vis_grid(parsed_grid)
			}

			n_iters := 1000

			scores := make([]int, n_iters)

			for i := 0; i < n_iters; i++ {
				propagated_grid = propagate_all_balls_north(propagated_grid)
				propagated_grid = propagate_all_balls_west(propagated_grid)
				propagated_grid = propagate_all_balls_south(propagated_grid)
				propagated_grid = propagate_all_balls_east(propagated_grid)

				grid_score := score_grid(propagated_grid)

				if DEBUG {
					fmt.Println("Iteration: ", i+1)
					fmt.Println("Score: ", grid_score)
					// vis_grid(propagated_grid)
				}

				scores[i] = grid_score
			}

			// Now need to find the periodicity of the sequence.
			up_to_now := make([]int, 0)

			period := 0

			for i := n_iters - 1; i > 0; i-- {
				up_to_now = append(up_to_now, scores[i])

				// Go from i through up to now and see if it repeats.
				// If it does, then we have found the period.

				if len(up_to_now) < 2 {
					continue
				}

				// Double up up to now... (caveman style)

				double_up_to_now := make([]int, 0)
				for _, val := range up_to_now {
					double_up_to_now = append(double_up_to_now, val)
				}
				for _, val := range up_to_now {
					double_up_to_now = append(double_up_to_now, val)
				}

				// fmt.Println("Double up to now: ", double_up_to_now)

				repeats := true

				for j := 0; j < len(double_up_to_now)-1; j++ {
					// fmt.Println("Comparinig: ", scores[n_iters-j-1], double_up_to_now[j])
					if scores[n_iters-j-1] != double_up_to_now[j] {
						repeats = false
						break
					}
				}

				if !repeats {
					continue
				} else {
					period = len(up_to_now)
					fmt.Println("Found period: ", period)
					fmt.Println("Periodic sequence: ", up_to_now)
					break
				}
			}

			n_cycles := 1000000000

			// Now we can calculate the score at the end.
			n_cycles_maps_to := (n_cycles - n_iters) % period
			fmt.Println("n_cycles_maps_to: ", n_cycles_maps_to)
			// Into the sequence
			final_score := scores[n_iters-period+n_cycles_maps_to-1]

			fmt.Println("Final Score: ", final_score)

			// Reset grid
			grid = make([]string, 0)
		} else {
			// Add to grid.
			grid = append(grid, text)
		}
	}
}
