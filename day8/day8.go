package day8

import (
	"regexp"
	"strings"
)

type read_order struct {
	order   string
	current int
}

func (r *read_order) read() int {
	index := 0
	if r.order[r.current] == 'R' {
		index = 1
	}
	r.current += 1
	if r.current == len(r.order) {
		r.current = 0
	}
	return index
}

func build_map(content string, ghostmode bool) (read_order, map[string][2]string, []string) {
	var order read_order
	var start_keys []string
	full_map := make(map[string][2]string)
	tract_regex := regexp.MustCompile(`(?P<idx>\w{3})\s+=\s+\((?P<l>\w{3})\,\s+(?P<r>\w{3})\)`)

	for _, line := range strings.Split(content, "\n") {
		if len(line) > 0 {
			if len(order.order) == 0 {
				order.order = strings.TrimSpace(line)
				order.current = 0
			} else {
				//single tract of map
				match := tract_regex.FindStringSubmatch(line)
				if len(match) > 0 {
					result := make(map[string]string)
					for i, name := range tract_regex.SubexpNames() {
						if i != 0 && name != "" {
							result[name] = match[i]
						}
					}
					key := result["idx"]
					full_map[key] = [2]string{result["l"], result["r"]}
					if ghostmode && key[len(key)-1] == 'A' {
						start_keys = append(start_keys, key)
					}
				}
			}
		}
	}
	return order, full_map, start_keys
}

func follow_map(full_map map[string][2]string, order read_order, begin string, end string) int {
	moves := 0
	current_position := begin
	for current_position != end {
		current_position = full_map[current_position][order.read()]
		moves += 1
	}

	return moves
}

func follow_ghost_map(order read_order, full_map map[string][2]string, start_keys []string) int {
	moves := 0
	keys := start_keys
	step := true
	for step {
		step = false
		moves += 1
		read_index := order.read()
		for idx, k := range keys {
			//log.Printf("Move %d step %d", moves, idx)
			keys[idx] = full_map[k][read_index]
			step = step || keys[idx][len(keys[idx])-1] != 'Z'
		}
	}
	return moves
}

func Part1(content string) int {
	order, full_map, _ := build_map(content, false)
	return follow_map(full_map, order, "AAA", "ZZZ")
}

func Part2(content string) int {
	return follow_ghost_map(build_map(content, true))
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
