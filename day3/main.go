package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Number struct {
	num string
	row int
	col int
	len int
}

type Schematic struct {
	cols  int
	rows  int
	valid [][]bool
}

func (s *Schematic) parse(content string) int {
	lines := strings.Fields(strings.TrimSpace(content))
	s.rows = len(lines)
	s.cols = len(lines[0])
	s.valid = make([][]bool, s.rows)
	for i := range s.valid {
		s.valid[i] = make([]bool, s.cols)
		for j := range s.valid[i] {
			s.valid[i][j] = false
		}
	}
	// check valid data
	for row, line := range lines {
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
			}
		}
	}

	// get part numbers
	sum := 0
	for row, line := range lines {
		//fmt.Printf("----- %d -----\n", row+1)
		tmp_num := []rune{}
		valid := false
		for col, char := range line {
			if unicode.IsDigit(char) {
				tmp_num = append(tmp_num, char)
				//fmt.Printf("%d,%d = %v\n", row, col, s.valid[row][col])
				if s.valid[row][col] {
					valid = true
				}
			} else {
				// finished parsing number
				sum += check_num(string(tmp_num), valid)
				valid = false
				tmp_num = []rune{}
			}
		}
		// check right border numbers
		sum += check_num(string(tmp_num), valid)
		valid = false
		tmp_num = []rune{}
	}
	return sum
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

func part1(content string) int {
	var schema Schematic
	return schema.parse(content)
}

func main() {
	buffer, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(buffer[:])
	fmt.Printf("Part 1: %d\n", part1(content))
}
