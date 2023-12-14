package day10

import (
	"log"
	"strings"
)

type coord struct {
	row int
	col int
}

type edge struct {
	a *vertex
	b *vertex
}

func (e *edge) get_other(x *vertex) *vertex {
	if x.same(e.a) {
		return e.b
	} else {
		return e.a
	}
}

func (e *edge) unlink(x *vertex) {
	if e.a != nil && x.same(e.a) {
		e.a = nil
		return
	}
	if e.b != nil && x.same(e.b) {
		e.b = nil
		return
	}
}

func build_island(content string) graph {
	var island graph
	island.vertexes = make(map[int]map[int]*vertex)
	for row, l := range strings.Split(content, "\n") {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			island.vertexes[row] = make(map[int]*vertex)
			for col, ru := range line {
				if ru != '.' {
					location := coord{row: row, col: col}
					pipe := &vertex{kind: ru, loc: location}
					island.vertexes[row][col] = pipe
					if ru == 'S' {
						island.start = pipe
					}
				}
			}
		}
	}
	island.make_island_links()
	island.cleanup()
	return island
}

func Part1(content string) int {
	island := build_island(content)
	if len(island.start.edges) != 2 {
		log.Fatalf("S is not a valid starting point with %d links", len(island.start.edges))
	} else {
		s_pipe := island.start
		var current_1, current_2 *vertex
		for _, p := range s_pipe.edges {
			if p != nil {
				if current_1 == nil {
					current_1 = p.get_other(s_pipe)
				} else {
					current_2 = p.get_other(s_pipe)
				}
			}
		}
		distance := 1
		prev_1 := s_pipe
		prev_2 := s_pipe
		if current_1 != nil && current_2 != nil {
			for !current_1.same(current_2) {
				log.Println(distance)
				next_1 := current_1.get_next(prev_1)
				prev_1 = current_1
				current_1 = next_1
				next_2 := current_2.get_next(prev_2)
				prev_2 = current_2
				current_2 = next_2
				distance += 1
			}
		} else {
			log.Fatal("S pipe is floating all alone")
		}
		return distance
	}
	return -1
}

func Part2(content string) int {
	sum := 0
	return sum
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
