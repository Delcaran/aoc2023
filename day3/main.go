package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Gear struct {
	col     int
	row     int
	numbers []int
	zone    [][]bool
}

func (g *Gear) init(s *Schematic, row int, col int) {
	g.zone = make([][]bool, s.rows)
	for i := range g.zone {
		g.zone[i] = make([]bool, s.cols)
		for j := range g.zone[i] {
			g.zone[i][j] = false
		}
	}

	g.col = col
	g.row = row

	up := max(0, row-1)
	down := min(s.rows-1, row+1)
	left := max(0, col-1)
	right := min(s.cols-1, col+1)

	g.zone[up][left] = true
	g.zone[up][col] = true
	g.zone[up][right] = true
	g.zone[row][left] = true
	g.zone[row][col] = true
	g.zone[row][right] = true
	g.zone[down][left] = true
	g.zone[down][col] = true
	g.zone[down][right] = true
}

func (g *Gear) ratio() int {
	if len(g.numbers) != 2 {
		return 0
	}
	prod := 1
	for _, value := range g.numbers {
		prod = prod * value
	}
	return prod
}

func (g *Gear) equal(o *Gear) bool {
	return g.col == o.col && g.row == o.row
}

func (g *Gear) add_part_number(num int) bool {
	if len(g.numbers) != 2 {
		g.numbers = append(g.numbers, num)
		return true
	}
	return false
}

func (g *Gear) in_range(nrow int, ncol int) bool {
	for r, row := range g.zone {
		for c := range row {
			if row[c] && nrow == r && ncol == c {
				return true
			}
		}
	}
	return false
}

type PartNumber struct {
	value int
	gears []*Gear
}

func (pn *PartNumber) equal(pno *PartNumber) bool {
	return pn.value == pno.value
}

type Schematic struct {
	schema       []string
	cols         int
	rows         int
	valid        [][]bool
	gears        []Gear
	part_numbers []PartNumber
}

func (s *Schematic) find_part_numbers() {
	for row, line := range s.schema {
		//fmt.Printf("----- %d -----\n", row+1)
		tmp_num := []rune{}
		valid := false
		var possible_gears []*Gear
		for col, char := range line {
			if unicode.IsDigit(char) {
				tmp_num = append(tmp_num, char)
				//fmt.Printf("%d,%d = %v\n", row, col, s.valid[row][col])
				if s.valid[row][col] {
					valid = true
				}
				// check possibile gear
				for i := range s.gears {
					g := &s.gears[i]
					if g.in_range(row, col) {
						add := true
						for _, eg := range possible_gears {
							if eg.equal(g) {
								add = false
							}
						}
						if add {
							possible_gears = append(possible_gears, g)
						}
					}
				}
			} else {
				// finished parsing number
				var pn PartNumber
				pn.value = check_num(string(tmp_num), valid)
				pn.gears = possible_gears
				if pn.value != 0 {
					s.part_numbers = append(s.part_numbers, pn)
				}
				valid = false
				possible_gears = possible_gears[:0]
				tmp_num = []rune{}
			}
		}
		// check right border numbers
		var pn PartNumber
		pn.value = check_num(string(tmp_num), valid)
		pn.gears = possible_gears
		if pn.value != 0 {
			s.part_numbers = append(s.part_numbers, pn)
		}
		valid = false
		possible_gears = possible_gears[:0]
		tmp_num = []rune{}
	}
}

func (s *Schematic) find_gears() {
	for _, pn := range s.part_numbers {
		for _, g := range pn.gears {
			// find global
			to_delete_global := -1
			for sidx := range s.gears {
				sg := &s.gears[sidx]
				if sg.equal(g) {
					if sg.add_part_number(pn.value) == false {
						to_delete_global = sidx
					}
				}
			}
			if to_delete_global != -1 {
				if to_delete_global == len(s.gears)-1 {
					s.gears = s.gears[:to_delete_global]
				} else {
					s.gears = append(s.gears[:to_delete_global], s.gears[to_delete_global+1:]...)
				}
			}
		}
	}
}

func (s *Schematic) parse(content string) {
	s.schema = strings.Fields(strings.TrimSpace(content))
	s.rows = len(s.schema)
	s.cols = len(s.schema[0])
	s.valid = make([][]bool, s.rows)
	for i := range s.valid {
		s.valid[i] = make([]bool, s.cols)
		for j := range s.valid[i] {
			s.valid[i][j] = false
		}
	}
	// check valid area
	for row, line := range s.schema {
		for col, char := range line {
			if !unicode.IsDigit(char) && char != '.' {
				up := max(0, row-1)
				down := min(s.rows-1, row+1)
				left := max(0, col-1)
				right := min(s.cols-1, col+1)
				s.valid[up][left] = true
				s.valid[up][col] = true
				s.valid[up][right] = true
				s.valid[row][left] = true
				s.valid[row][col] = true
				s.valid[row][right] = true
				s.valid[down][left] = true
				s.valid[down][col] = true
				s.valid[down][right] = true

				if char == '*' {
					var new_g Gear
					new_g.init(s, row, col)
					s.gears = append(s.gears, new_g)
				}
			}
		}
	}
	s.find_part_numbers()
	s.find_gears()
}

func check_num(num_str string, valid bool) int {
	if len(num_str) > 0 {
		//fmt.Printf("%s ", num_str)
		if valid {
			num, err := strconv.Atoi(num_str)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("part number\n")
			return num
		} else {
			//fmt.Printf("not a part number\n")
			return 0
		}
	}
	return 0
}

func part1(s *Schematic) int {
	sum := 0
	for _, pn := range s.part_numbers {
		sum += pn.value
	}
	return sum
}

func part2(s *Schematic) int {
	sum := 0
	for _, g := range s.gears {
		ratio := g.ratio()
		sum += ratio
	}
	return sum
}

func main() {
	buffer, err := os.ReadFile("test_input.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(buffer[:])
	var schema Schematic
	schema.parse(content)

	fmt.Printf("Part 1: %d\n", part1(&schema))
	fmt.Printf("Part 1: %d\n", part2(&schema))
}
