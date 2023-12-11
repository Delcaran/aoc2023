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

func part1(content string) int {
	var order read_order
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
					full_map[result["idx"]] = [2]string{result["l"], result["r"]}
				}
			}
		}
	}

	moves := 0
	current_position := "AAA"
	for current_position != "ZZZ" {
		current_position = full_map[current_position][order.read()]
		moves += 1
	}

	return moves
}

func part2(content string) int {
	return 0
}

func Run(content string) (int, int) {
	return part1(content), part2(content)
}
