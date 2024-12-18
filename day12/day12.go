package day12

import (
	"log"
	"strconv"
	"strings"
)

const (
	SPRING_UNKNOWN = iota
	SPRING_WORKING = iota
	SPRING_DAMAGED = iota
)

type field_row struct {
	springs_map     string
	damaged_records []int
}

type field struct {
	rows []field_row
}

func parse_input(content string) field {
	var f field
	for _, l := range strings.Split(content, "\n") {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			var fr field_row
			line_data := strings.Split(line, " ")
			fr.springs_map = line_data[0]
			for _, ru := range strings.Split(line_data[1], ",") {
				num, err := strconv.Atoi(ru)
				if err != nil {
					log.Fatal(err)
				}
				fr.damaged_records = append(fr.damaged_records, num)
			}
			f.rows = append(f.rows, fr)
		}
	}
	return f
}

func minify(input string) string {
	var minified string
	for pos, char := range input {
		if pos == 0 {
			minified = string(char)
		} else {
			if char == '.' && minified[len(minified)-1] == '.' {
				continue
			} else {
				minified = minified + string(char)
			}
		}
	}
	return minified
}

func (f *field_row) count_damaged_groups() {
	count := 0
	for _, char := range f.springs_map {
		if char == '#' {
			count = count + 1
		}
		if char == '.' {
			if count > 0 {
				f.damaged_records = append(f.damaged_records, count)
				count = 0
			}
		}
		if char == '?' {
			panic("count_damaged_groups() MUST be called on fully reconstructed strings")
		}
	}
	if count > 0 {
		f.damaged_records = append(f.damaged_records, count)
		count = 0
	}
}

func solve_brute_force(input string, solved *map[string]field_row) {
	minified := minify(input)
	first_unknown := strings.Index(input, "?")
	if first_unknown == -1 {
		// end recursion: memorize string
		count_damaged_groups()
	} else {
		// recursion
		solve_brute_force(string('#')+input[first_unknown+1:], solved)
		solve_brute_force(string('.')+input[first_unknown+1:], solved)
	}
}

func Part1(f field) int {
	arrangements := 0
	strings_solved := make(map[string]int)
	for rown := range f.rows {
		row_arrangements := find_groups(f.rows[rown].springs_map, f.rows[rown].damaged_records, 0, memoization)
		arrangements += row_arrangements
	}
	return arrangements
}

func Part2(f field) int {
	f.unfold()
	arrangements := 0
	memoization := make(map[string]int)
	for rown := range f.rows {
		row_arrangements := find_groups(f.rows[rown].springs_map, f.rows[rown].damaged_records, 0, memoization)
		arrangements += row_arrangements
	}
	return arrangements
}

func Run(content string) (int, int) {
	field_info := parse_input(content)
	return Part1(field_info), 0 //Part2(field_info)
}
