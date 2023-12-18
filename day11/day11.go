package day11

import (
	"fmt"
	"strings"
)

type sector struct {
	row          int
	col          int
	expanded_row int
	expanded_col int
	galaxy       bool
}

func (s *sector) print() string {
	if s.galaxy {
		return fmt.Sprintf(" %dX%d ", s.row, s.col)
	}
	return fmt.Sprintf(" %d.%d ", s.row, s.col)
}

func (start *sector) distance_from(end *sector) int {
	// taxicab distance
	col_distance := max(start.expanded_col, end.expanded_col) - min(start.expanded_col, end.expanded_col)
	row_distance := max(start.expanded_row, end.expanded_row) - min(start.expanded_row, end.expanded_row)
	return col_distance + row_distance
}

type skymap struct {
	galaxies         []*sector
	sectors          [][]*sector
	expanded_sectors [][]*sector
	distance_matrix  map[int]map[int]int
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

func (s *skymap) print_expanded() {
	for row := 0; row < len(s.expanded_sectors); row++ {
		for col := 0; col < len(s.expanded_sectors[row]); col++ {
			sect := s.expanded_sectors[row][col]
			fmt.Print(sect.print())
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func build_sky(content string) skymap {
	var sky skymap
	var empty_cols []int
	for row, l := range strings.Split(content, "\n") {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			sky.sectors = append(sky.sectors, make([]*sector, len(line)))
			sky.expanded_sectors = append(sky.expanded_sectors, make([]*sector, len(line)))
			galaxy_in_row := false
			for col, ru := range line {
				if row == 0 {
					empty_cols = append(empty_cols, col)
				}
				s := &sector{row: row, col: col, galaxy: ru == '#'}
				sky.sectors[len(sky.sectors)-1][col] = s
				sky.expanded_sectors[len(sky.expanded_sectors)-1][col] = s
				if s.galaxy {
					galaxy_in_row = true
					for idx := 0; idx < len(empty_cols); idx++ {
						if empty_cols[idx] == col {
							empty_cols = append(empty_cols[:idx], empty_cols[idx+1:]...)
						}
					}
					sky.galaxies = append(sky.galaxies, s)
				}
			}

			if !galaxy_in_row {
				sky.expanded_sectors = append(sky.expanded_sectors, sky.expanded_sectors[len(sky.expanded_sectors)-1])
			}
		}
	}
	// check empty cols
	expansions := 0
	for _, empty_col := range empty_cols {
		for row := 0; row < len(sky.expanded_sectors); row++ {
			idx := empty_col + expansions
			tmp := sky.expanded_sectors[row][idx]
			sky.expanded_sectors[row] = append(sky.expanded_sectors[row][:idx+1], sky.expanded_sectors[row][idx:]...)
			sky.expanded_sectors[row][idx] = tmp
		}
		expansions += 1
	}

	// update values
	for row := 0; row < len(sky.expanded_sectors); row++ {
		for col := 0; col < len(sky.expanded_sectors[row]); col++ {
			sky.expanded_sectors[row][col].expanded_row = row
			sky.expanded_sectors[row][col].expanded_col = col
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
					distance := start.distance_from(end)
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
	sky := build_sky(content)
	sky.print()
	sky.print_expanded()
	return sky.galaxy_distance()
}

func Part2(content string) int {
	return 0
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
