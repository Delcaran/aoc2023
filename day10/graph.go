package day10

import (
	"fmt"
	"log"
)

type graph struct {
	start    *vertex
	vertexes map[int]map[int]*vertex
}

func (island *graph) get_pipe(row int, col int) *vertex {
	pipe, ok := island.vertexes[row][col]
	if ok {
		return pipe
	}
	return nil
}

func (island *graph) make_pipe_links(pipe *vertex) bool {
	row := pipe.loc.row
	col := pipe.loc.col
	count := 0
	count += pipe.make_link_north(island.get_pipe(row-1, col))
	count += pipe.make_link_east(island.get_pipe(row, col+1))
	count += pipe.make_link_west(island.get_pipe(row, col-1))
	count += pipe.make_link_south(island.get_pipe(row+1, col))
	if pipe.kind == 'S' {
		return false // never delete
	}
	return count != 2
}

func (island *graph) make_island_links() {
	var deletable []coord
	for row := range island.vertexes {
		for col := range island.vertexes[row] {
			pipe := island.vertexes[row][col]
			wrong := island.make_pipe_links(pipe)
			if wrong {
				deletable = append(deletable, pipe.loc)
			}
		}
	}
	island.remove_useless_links(deletable)
}

func (island *graph) cleanup() {
	do_delete := true
	for do_delete {
		var deletable []coord
		for row := range island.vertexes {
			for col := range island.vertexes[row] {
				pipe := island.vertexes[row][col]
				num_edges := len(pipe.edges)
				if pipe.kind == 'S' {
					// remove dead links from S
					for i := num_edges - 1; i >= 0; i-- {
						edge := pipe.edges[i]
						if edge.get_other(pipe) == nil {
							pipe.edges = append(pipe.edges[:i], pipe.edges[i+1:]...)
						}
					}
				} else {
					if num_edges > 0 && num_edges != 2 {
						deletable = append(deletable, pipe.loc)
					}
				}
			}
		}
		island.remove_useless_links(deletable)
		do_delete = len(deletable) > 0
	}
}

func (island *graph) remove_useless_links(useless []coord) {
	for _, c := range useless {
		pipe := island.vertexes[c.row][c.col]
		pipe.unlink()
	}
}

func (island *graph) print() {
	fmt.Print("\n")
	for row := 0; row < len(island.vertexes); row++ {
		for col := 0; col < len(island.vertexes[row]); col++ {
			pipe := island.vertexes[row][col]
			fmt.Print(pipe.print())
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (island *graph) double_walk() int {
	if len(island.start.edges) != 2 {
		log.Fatalf("S is not a valid starting point with %d links", len(island.start.edges))
	} else {
		s_pipe := island.start
		s_pipe.relpos = PATH
		var current_1, current_2 *vertex
		for _, p := range s_pipe.edges {
			if p != nil {
				if current_1 == nil {
					current_1 = p.get_other(s_pipe)
					current_1.relpos = PATH
				} else {
					current_2 = p.get_other(s_pipe)
					current_1.relpos = PATH
				}
			}
		}
		distance := 1
		prev_1 := s_pipe
		prev_2 := s_pipe
		if current_1 != nil && current_2 != nil {
			for !current_1.same(current_2) {
				next_1 := current_1.get_next(prev_1)
				prev_1 = current_1
				current_1 = next_1
				next_2 := current_2.get_next(prev_2)
				prev_2 = current_2
				current_2 = next_2
				distance += 1
			}
			current_1.end = true
			current_2.end = true
		} else {
			log.Fatal("S pipe is floating all alone")
		}
		return distance
	}
	return -1
}

func (island *graph) mark_pipes(prev *vertex, current *vertex) {
	prev_col := prev.loc.col
	prev_row := prev.loc.row
	curr_col := current.loc.col
	curr_row := current.loc.row

	if prev_col == curr_col {
		e := prev_col + 1
		w := prev_col - 1
		// vertical
		if prev_row < curr_row {
			// north -> south
			island.get_pipe(prev_row, e).mark(LEFT)
			island.get_pipe(curr_row, e).mark(LEFT)
			island.get_pipe(prev_row, w).mark(RIGHT)
			island.get_pipe(curr_row, w).mark(RIGHT)
		} else {
			// south -> north
			island.get_pipe(prev_row, w).mark(LEFT)
			island.get_pipe(curr_row, w).mark(LEFT)
			island.get_pipe(prev_row, e).mark(RIGHT)
			island.get_pipe(curr_row, e).mark(RIGHT)
		}
	} else {
		n := prev_row - 1
		s := prev_row + 1
		// horizontal
		if prev_col < curr_col {
			// west -> east
			island.get_pipe(n, prev_col).mark(LEFT)
			island.get_pipe(n, curr_col).mark(LEFT)
			island.get_pipe(s, prev_col).mark(RIGHT)
			island.get_pipe(s, curr_col).mark(RIGHT)
		} else {
			// east -> west
			island.get_pipe(s, prev_col).mark(LEFT)
			island.get_pipe(s, curr_col).mark(LEFT)
			island.get_pipe(n, prev_col).mark(RIGHT)
			island.get_pipe(n, curr_col).mark(RIGHT)
		}
	}
}

func (island *graph) walk() {
	if len(island.start.edges) != 2 {
		log.Fatalf("S is not a valid starting point with %d links", len(island.start.edges))
	} else {
		s_pipe := island.start
		s_pipe.relpos = PATH
		current := s_pipe.edges[0].get_other(s_pipe)
		current.relpos = PATH
		prev := s_pipe
		for !current.same(s_pipe) {
			next := current.get_next(prev)
			prev = current
			current = next
			island.mark_pipes(prev, current)
		}
		current.end = true
		// walk has been walked, tiles have been marked
		// now marked tiles will be spread
	}
}
