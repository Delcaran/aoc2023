package day10

import (
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
				location := coord{row: row, col: col}
				pipe := &vertex{kind: ru, loc: location}
				island.vertexes[row][col] = pipe
				if ru == 'S' {
					island.start = pipe
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
	distance := island.double_walk()
	island.print()
	return distance
}

func Part2(content string) int {
	island := build_island(content)
	island.walk()
	island.print()
	/*
		IDEA: run the walk and mark left/right
		- run the path in one direction
		- mark left cells with L
		- mark right cells with R
		- "diffuse" L and R
		- if L or R finds the borders, the other mark is on the tiles I need
	*/
	return island.count_inner()
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
