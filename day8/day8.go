package day8

import (
	"log"
	"math/big"
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

type position struct {
	key       string
	order_idx int
	move      int
}

type loop struct {
	begin_move int
	end_move   int
	size       int
}

func get_z_ending(full_map map[string][2]string, order read_order, begin string) ([]int, loop) {
	step := true
	key := begin
	var path []position
	var loop loop
	move := 0
	path = append(path, position{key: key, order_idx: order.current, move: move})
	//log.Printf("%d : %s %d", move, key, order.current)
	for step {
		step = true
		key = full_map[key][order.read()]
		move += 1
		curr_post := position{key: key, order_idx: order.current, move: move}
		//log.Printf("%d : %s %d", curr_post.move, curr_post.key, curr_post.order_idx)

		for x := len(path) - 1; x >= 0; x-- {
			oldpos := path[x]
			if oldpos.key == curr_post.key && oldpos.order_idx == curr_post.order_idx {
				loop.begin_move = oldpos.move
				loop.end_move = curr_post.move - 1
				loop.size = curr_post.move - oldpos.move
				//log.Printf("%s %d -> %s %d", curr_post.key, curr_post.order_idx, oldpos.key, oldpos.order_idx)
				step = false
				break
			}
		}

		if step {
			path = append(path, curr_post)
		}
	}
	var z_endings []int
	for _, pos := range path {
		if pos.key[len(pos.key)-1] == 'Z' {
			z_endings = append(z_endings, pos.move)
		}
	}
	return z_endings, loop
}

func Part1(content string) int {
	order, full_map, _ := build_map(content, false)
	return follow_map(full_map, order, "AAA", "ZZZ")
}

func calc_lcm(a *big.Int, b *big.Int) big.Int {
	var gcd, x big.Int
	gcd.GCD(nil, nil, a, b)
	x.Div(a, &gcd).Mul(&x, b)
	return x
}

func Part2(content string) int {
	order, full_map, keys := build_map(content, true)

	z_endings := make(map[string][]int)
	loops := make(map[string]loop)
	var zs []*big.Int
	for _, k := range keys {
		z_endings[k], loops[k] = get_z_ending(full_map, order, k)
		log.Printf("%s : %d (%d->%d) | %d", k, z_endings[k], loops[k].begin_move, loops[k].end_move, loops[k].size)
		for _, z := range z_endings[k] {
			zs = append(zs, big.NewInt(int64(z)))
		}
	}

	// Use LCM, should work with Z-ending all at the end of the loop as it looks like
	// from examples and input
	var lcm big.Int
	for idx, x := range zs {
		if idx == 0 {
			lcm = *x
		} else {
			lcm = calc_lcm(&lcm, x)
		}
	}
	return int(lcm.Int64())
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
