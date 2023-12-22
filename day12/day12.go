package day12

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
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

func (f *field_row) unfold() {
	var maps []string
	orig_records := f.damaged_records
	for x := 0; x < 5; x++ {
		f.damaged_records = append(f.damaged_records, orig_records...)
		maps = append(maps, f.springs_map)
	}
	f.springs_map = strings.Join(maps, "?")
}

func (f field_row) brute_force(springs_map string) int {
	if strings.Contains(springs_map, "?") {
		arrangements := 0
		for x := SPRING_WORKING; x <= SPRING_DAMAGED; x++ {
			switch x {
			case SPRING_WORKING:
				arrangements += f.brute_force(strings.Replace(springs_map, "?", ".", 1))
			case SPRING_DAMAGED:
				arrangements += f.brute_force(strings.Replace(springs_map, "?", "#", 1))
			}
		}
		return arrangements
	}
	// check current spring
	var damaged_groups []int
	group_count := 0
	for _, ru := range springs_map {
		if ru == '#' {
			group_count += 1
		} else {
			if group_count > 0 {
				damaged_groups = append(damaged_groups, group_count)
				group_count = 0
			}
		}
	}
	// EOL check
	if group_count > 0 {
		damaged_groups = append(damaged_groups, group_count)
		group_count = 0
	}
	if len(damaged_groups) == len(f.damaged_records) {
		same := 0
		for x := 0; x < len(damaged_groups); x++ {
			if damaged_groups[x] == f.damaged_records[x] {
				same += 1
			}
		}
		if same == len(damaged_groups) {
			return 1
		}
	}
	return 0
}

type field struct {
	rows []field_row
}

func (f field) brute_force() int {
	arrangements := 0
	for rown := range f.rows {
		line_arrengements := f.rows[rown].brute_force(f.rows[rown].springs_map)
		arrangements += line_arrengements
	}
	return arrangements
}

func (f *field) unfold() {
	for rown := range f.rows {
		f.rows[rown].unfold()
	}
}

func minimize(springs_map string) string {
	trimmed := strings.Trim(springs_map, ".")
	var compressed string
	count := 0
	last_point := false
	for _, ru := range trimmed {
		switch ru {
		case '#':
			count += 1
			last_point = false
		case '?':
			last_point = false
			if count > 0 {
				compressed = fmt.Sprintf("%s%d", compressed, count)
			}
			compressed = fmt.Sprintf("%s%c", compressed, ru)
			count = 0
		case '.':
			if count > 0 {
				compressed = fmt.Sprintf("%s%d", compressed, count)
			}
			if !last_point {
				last_point = true
				compressed = fmt.Sprintf("%s%c", compressed, ru)
			}
			count = 0
		default:
			// number
			num, err := strconv.Atoi(string(ru))
			if err != nil {
				log.Fatal(err)
			}
			count += num
		}
	}
	if count > 0 {
		compressed = fmt.Sprintf("%s%d", compressed, count)
		count = 0
	}
	return compressed
}

func stringify(damaged_map []int) string {
	stringified := make([]string, len(damaged_map))
	for x, v := range damaged_map {
		stringified[x] = strconv.Itoa(v)
	}
	return strings.Join(stringified, ".")
}

func check_valid_and_recurr(left_to_parse int, damaged_map []int, damaged_index int, count int) (int, bool) {
	// check if we have finished
	if count < damaged_map[damaged_index] {
		// not enough
		return 0, false
	} else if count > damaged_map[damaged_index] {
		// too many
		return 0, false
	} else {
		if damaged_index == len(damaged_map)-1 {
			// finished, return results
			if left_to_parse > 0 {
				// map left to parse but no more groups... error!
				return 0, false
			} else {
				// nothing to parse and no more groups, it's working!
				return 1, false
			}
		} else {
			if left_to_parse > 0 {
				// no results here, recursion needed
				return 0, true
			} else {
				// no results and nothing left to do
				return 0, false
			}
		}
	}
}

func find_groups(springs_map string, damaged_map []int, damaged_index int, already_found map[string]int) int {
	found, ok := already_found[springs_map]
	if ok {
		return found
	}
	minimized := minimize(springs_map)
	foundmin, ok := already_found[minimized]
	if ok {
		already_found[springs_map] = foundmin
		return foundmin
	}
	count := 0
	last_taken := -1
	arrangements_broken := 0
	recurr_broken := false
	arrangements_working := 0
	recurr_working := false
	for ch_x, ch := range minimized {
		if unicode.IsDigit(ch) {
			// number, adding it to count
			num, err := strconv.Atoi(string(ch))
			if err != nil {
				log.Fatal(err)
			}
			count += num
			last_taken = ch_x
		} else {
			if ch == '?' {
				last_taken = ch_x
				// simulate #
				arrangements_broken, recurr_broken = check_valid_and_recurr(len(minimized[last_taken+1:]), damaged_map, damaged_index, count+1)
				// simulate .
				arrangements_working, recurr_working = check_valid_and_recurr(len(minimized[last_taken+1:]), damaged_map, damaged_index, count)
				count += 1
				break
			} else {
				// '.'
				break
			}
		}
	}

	arrangements := 0
	if recurr_broken || recurr_working {
		next_idx := last_taken + 1
		if recurr_broken {
			fmt.Printf("%s -> %s\n", minimized, "#"+minimized[next_idx:])
			arrangements += find_groups("#"+minimized[next_idx:], damaged_map, damaged_index+1, already_found)
		}
		if recurr_working {
			fmt.Printf("%s -> %s\n", minimized, "."+minimized[next_idx:])
			arrangements += find_groups("."+minimized[next_idx:], damaged_map, damaged_index+1, already_found)
		}
	} else {
		arrangements = arrangements_broken + arrangements_working
		fmt.Printf("%s -> %d\n", minimized, arrangements)
		already_found[springs_map] = arrangements
		already_found[minimized] = arrangements
	}
	return arrangements
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

func Part1(f field) int {
	//return f.brute_force()
	arrangements := 0
	memoization := make(map[string]int)
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
