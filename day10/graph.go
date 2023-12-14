package day10

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
