package day11

import (
	"fmt"
	"strings"
)

type sector struct {
	row    int
	col    int
	galaxy bool
}

func (s *sector) print() string {
	if s.galaxy {
		return fmt.Sprintf(" %dX%d ", s.row, s.col)
	}
	return fmt.Sprintf(" %d.%d ", s.row, s.col)
}

func (start *sector) distance_from(end *sector, sky *skymap) int {
	// taxicab distance
	x_distance := max(start.col, end.col) - min(start.col, end.col)
	y_distance := max(start.row, end.row) - min(start.row, end.row)

	// factor in age of the universe
	for _, ex := range sky.empty_cols {
		if min(start.col, end.col) <= ex && ex <= max(start.col, end.col) {
			x_distance += sky.expansion - 1
		}
	}
	for _, ey := range sky.empty_rows {
		if min(start.row, end.row) <= ey && ey <= max(start.row, end.row) {
			y_distance += sky.expansion - 1
		}
	}

	return x_distance + y_distance
}

type skymap struct {
	expansion       int
	galaxies        []*sector
	sectors         [][]*sector
	empty_rows      []int
	empty_cols      []int
	distance_matrix map[int]map[int]int
}

func (s *skymap) print() {
	for row := 0; row < len(s.sectors); row++ {
		for col := 0; col < len(s.sectors[row]); col++ {
			sect := s.sectors[row][col]
			fmt.Print(sect.print())
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func build_sky(content string, expansion int) skymap {
	var sky skymap
	sky.expansion = expansion
	for row, l := range strings.Split(content, "\n") {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			sky.sectors = append(sky.sectors, make([]*sector, len(line)))
			galaxy_in_row := false
			for col, ru := range line {
				if row == 0 {
					sky.empty_cols = append(sky.empty_cols, col)
				}
				s := &sector{row: row, col: col, galaxy: ru == '#'}
				sky.sectors[len(sky.sectors)-1][col] = s
				if s.galaxy {
					galaxy_in_row = true
					for idx := 0; idx < len(sky.empty_cols); idx++ {
						if sky.empty_cols[idx] == col {
							sky.empty_cols = append(sky.empty_cols[:idx], sky.empty_cols[idx+1:]...)
						}
					}
					sky.galaxies = append(sky.galaxies, s)
				}
			}

			if !galaxy_in_row {
				sky.empty_rows = append(sky.empty_rows, row)
			}
		}
	}

	return sky
}

func (sky *skymap) galaxy_distance() int {
	sum := 0
	sky.distance_matrix = make(map[int]map[int]int)
	for x := 0; x < len(sky.galaxies); x++ {
		_, ok := sky.distance_matrix[x]
		if !ok {
			sky.distance_matrix[x] = make(map[int]int, len(sky.galaxies))
		}
		for y := 0; y <= x; y++ {
			_, ok := sky.distance_matrix[x][y]
			if !ok {
				if x == y {
					sky.distance_matrix[x][y] = 0
				} else {
					start := sky.galaxies[x]
					end := sky.galaxies[y]
					distance := start.distance_from(end, sky)
					sky.distance_matrix[x][y] = distance
					sum += distance
				}
				sky.distance_matrix[y][x] = sky.distance_matrix[x][y]
			}
		}
	}
	return sum
}

func Part1(content string) int {
	sky := build_sky(content, 2)
	//sky.print()
	//sky.print_expanded()
	return sky.galaxy_distance()
}

func Part2(content string) int {
	sky := build_sky(content, 1000000)
	//sky.print()
	//sky.print_expanded()
	return sky.galaxy_distance()
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
