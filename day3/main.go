package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Number struct {
	num string
	row int
	col int
	len int
	pn  bool
}

type Schematic struct {
	cols    int
	rows    int
	data    []string
	numbers []Number
}

func (s *Schematic) parse(content string) {
	regex := regexp.MustCompile(`(\d+)`)
	for row, line := range strings.Split(content, "\n") {
		if len(line) > 0 {
			s.cols = len(line) - 1
			s.data = append(s.data, line)
			match := regex.FindAllStringSubmatch(line, -1)
			for _, m := range match {
				col := strings.Index(line, m[0])
				num := Number{num: m[0], row: row, col: col, len: len(m[0])}
				s.numbers = append(s.numbers, num)
			}
			s.rows = row + 1
		}
	}
}

func (s *Schematic) isPartNumber(n Number) bool {
	left_idx := max(n.col-1, 0)
	right_idx := min(n.col+n.len, s.cols-1)
	top_idx := max(n.row-1, 0)
	bottom_idx := min(n.row+1, s.rows-1)

	left_dot := 1
	right_dot := 1
	if left_idx == n.col {
		left_dot = 0
	}
	if right_idx == n.col+n.len-1 {
		right_dot = 0
	}
	str_len := n.len + left_dot + right_dot

	up_down_regex := regexp.MustCompile(`^\.{` + strconv.Itoa(str_len) + `}$`)
	central_regex := regexp.MustCompile(`^\.{` + strconv.Itoa(left_dot) + `}` + n.num + `\.{` + strconv.Itoa(right_dot) + `}$`)

	top_row := s.data[top_idx][left_idx : right_idx+1]
	bottom_row := s.data[bottom_idx][left_idx : right_idx+1]
	number_row := s.data[n.row][left_idx : right_idx+1]

	match_up := true
	if top_idx != n.row {
		log.Printf("\t%s\n", top_row)
		match_up = up_down_regex.MatchString(top_row)
	}
	match_central := central_regex.MatchString(number_row)
	log.Printf("\t%s\n", number_row)
	match_down := true
	if bottom_idx != n.row {
		match_down = up_down_regex.MatchString(bottom_row)
		log.Printf("\t%s\n", bottom_row)
	}
	not_pn := match_up && match_down && match_central
	if not_pn {
		log.Printf("\tnot a part number")
	} else {
		log.Printf("\tis part number")
	}
	return !not_pn
}

func (s *Schematic) sumPartNumbers() int {
	sum := 0
	for _, n := range s.numbers {
		log.Printf("Number %s", n.num)
		n.pn = s.isPartNumber(n)
		if n.pn {

			num, err := strconv.Atoi(n.num)
			if err != nil {
				log.Fatal(err)
			}
			sum += num
		} else {

		}

	}
	return sum
}

func part1(content string) int {
	var schema Schematic
	schema.parse(content)
	return schema.sumPartNumbers()
}

func main() {
	buffer, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(buffer[:])
	log.Printf("Part 1: %d\n", part1(content))
}
