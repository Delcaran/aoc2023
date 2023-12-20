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

func (f *field_row) unfold() {
	var maps []string
	orig_records := f.damaged_records
	for x := 0; x < 5; x++ {
		f.damaged_records = append(f.damaged_records, orig_records...)
		maps = append(maps, f.springs_map)
	}
	f.springs_map = strings.Join(maps, "?")
}

func (f field_row) count_arrangements(springs_map string) int {
	if strings.Contains(springs_map, "?") {
		arrangements := 0
		for x := SPRING_WORKING; x <= SPRING_DAMAGED; x++ {
			switch x {
			case SPRING_WORKING:
				arrangements += f.count_arrangements(strings.Replace(springs_map, "?", ".", 1))
			case SPRING_DAMAGED:
				arrangements += f.count_arrangements(strings.Replace(springs_map, "?", "#", 1))
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

func (f field) count_arrangements() int {
	arrangements := 0
	for rown := range f.rows {
		line_arrengements := f.rows[rown].count_arrangements(f.rows[rown].springs_map)
		arrangements += line_arrengements
	}
	return arrangements
}

func (f *field) unfold() {
	for rown := range f.rows {
		f.rows[rown].unfold()
	}
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
	return f.count_arrangements()
}

func Part2(f field) int {
	f.unfold()
	return f.count_arrangements()
}

func Run(content string) (int, int) {
	field_info := parse_input(content)
	return Part1(field_info), Part2(field_info)
}
