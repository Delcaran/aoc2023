package day12

import (
	"fmt"
	"strings"
)

const (
	SPRING_UNKNOWN = iota
	SPRING_WORKING = iota
	SPRING_DAMAGED = iota
)

type field_row struct {
	springs_map     string
	damaged_records string
	arrangements    int
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
			fr.damaged_records = line_data[1]
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

func count_damaged_groups(springs_map string) string {
	damaged_records := make([]int, 0)
	count := 0
	for _, char := range springs_map {
		if char == '#' {
			count = count + 1
		}
		if char == '.' {
			if count > 0 {
				damaged_records = append(damaged_records, count)
				count = 0
			}
		}
		if char == '?' {
			panic("count_damaged_groups() MUST be called with fully reconstructed strings")
		}
	}
	if count > 0 {
		damaged_records = append(damaged_records, count)
		count = 0
	}
	var recordstring string
	for _, n := range damaged_records {
		if len(recordstring) > 0 {
			recordstring = recordstring + ","
		}
		recordstring = recordstring + fmt.Sprint(n)
	}
	return recordstring
}

func (f *field) unfold() {
	unfolding_size := 5
	for rown := range f.rows {
		solving_map := f.rows[rown].springs_map
		solving_records := f.rows[rown].damaged_records
		unfolded_map := make([]string, unfolding_size)
		unfolded_records := make([]string, unfolding_size)
		for x := 0; x < unfolding_size; x++ {
			unfolded_map = append(unfolded_map, solving_map)
			unfolded_records = append(unfolded_map, solving_records)
		}
		f.rows[rown].springs_map = strings.Join(unfolded_map, "?")
		f.rows[rown].damaged_records = strings.Join(unfolded_records, ",")
	}
}

func add_to_solved(solved map[string][]string, solution string) {
	_, exists := solved[solution]
	if !exists {
		solved[solution] = make([]string, 0)
	}
	solved_records, ok := solved[solution]
	if ok {
		found := false
		damaged_records := count_damaged_groups(solution)
		for _, x := range solved_records {
			if x == damaged_records {
				found = true
				break
			}
		}
		if !found {
			solved[solution] = append(solved[solution], damaged_records)
		}
	}
}

func build_group(size int) string {
	var group string
	for x := 0; x < size; x++ {
		group += "#"
	}
	return group
}

func expand_map(basic_map string) []string {
	expanded_maps := make([]string, 0)
	for idx, ch := range basic_map {
		if ch == '.' {
			expanded_maps = append(expanded_maps, basic_map[:idx]+"."+basic_map[idx+1:])
		} else {
			if idx == 0 {
				expanded_maps = append(expanded_maps, "."+basic_map)
			}
			if idx == len(basic_map)-1 {
				expanded_maps = append(expanded_maps, basic_map+".")
			}
		}
	}
	return expanded_maps
}

func build_maps(damage_groups []int, length int) []string {
	dot_positions := len(damage_groups) + 1
	groups := make([]string, 0)
	for _, size := range damage_groups {
		groups = append(groups, build_group(size))
	}
	basic_map := strings.Join(groups, ".")
	missing_dots := length - len(basic_map)
	maps := make([]string, 0)
	for x := 0; x < missing_dots; x++ {

	}
	return maps
}

func solve_by_build(row_map string, damage_groups []int, solved map[string][]string) {
	test_strin
}

func solve_brute_force(row_map string, solved map[string][]string) {
	_, ok := solved[row_map]
	if !ok {
		first_unknown := strings.Index(row_map, "?")
		if first_unknown == -1 {
			// end recursion: memorize string and all substrings
			for idx := range row_map {
				substring := row_map[idx:]
				add_to_solved(solved, substring)
				add_to_solved(solved, minify(substring))
			}
		} else {
			// recursion
			solve_brute_force(row_map[:first_unknown]+string('#')+row_map[first_unknown+1:], solved)
			solve_brute_force(row_map[:first_unknown]+string('.')+row_map[first_unknown+1:], solved)
		}
	}
}

func map_match(original string, solved string) bool {
	if len(original) == len(solved) {
		for idx := 0; idx < len(original); idx++ {
			if original[idx] != '?' {
				if original[idx] != solved[idx] {
					return false
				}
			}
		}
		return true
	}
	return false
}

func Part1(f field) int {
	strings_solved := make(map[string][]string)
	for rown := range f.rows {
		solving_map := f.rows[rown].springs_map
		solving_records := f.rows[rown].damaged_records
		solve_brute_force(solving_map, strings_solved)
		// all possible combinations of the row are found. Now check if something matches
		for solved_map, solved_records := range strings_solved {
			if map_match(solving_map, solved_map) {
				for _, solved_record := range solved_records {
					if solving_records == solved_record {
						f.rows[rown].arrangements += 1
					}
				}
			}
		}
	}
	arrangements := 0
	for rown := range f.rows {
		arrangements += f.rows[rown].arrangements
	}
	return arrangements
}

func Part2(f field) int {
	f.unfold()
	return Part1(f)
}

func Run(content string) (int, int) {
	field_info := parse_input(content)
	return Part1(field_info), Part2(field_info)
}
