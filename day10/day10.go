package day10

import (
	"log"
	"strings"
)

type coord struct {
	row int
	col int
}

type pipe struct {
	loc   coord
	kind  rune
	links []*pipe
}

func (p *pipe) next_pipe(incoming *pipe) *pipe {
	for _, outgoing := range p.links {
		if outgoing != nil && !outgoing.same(incoming) {
			log.Printf("\t%c -> [%c] -> %c", incoming.kind, p.kind, outgoing.kind)
			return outgoing
		}
	}
	log.Fatalf("[%d, %d] %c coming from [%d, %d] %c leads to nowhere!\n", p.loc.row, p.loc.col, p.kind, incoming.loc.row, incoming.loc.col, incoming.kind)
	return nil // should never be there
}

func (p *pipe) int(c coord, k rune) {
	p.kind = k
	p.loc = c
}

/*
func (p *pipe) north() *pipe {
	r, ok := p.links['n']
	if !ok {
		return nil
	}
	return r
}
func (p *pipe) south() *pipe {
	r, ok := p.links['s']
	if !ok {
		return nil
	}
	return r
}
func (p *pipe) east() *pipe {
	r, ok := p.links['e']
	if !ok {
		return nil
	}
	return r
}
func (p *pipe) west() *pipe {
	r, ok := p.links['w']
	if !ok {
		return nil
	}
	return r
}
*/

func (p *pipe) same(o *pipe) bool {
	return p.loc.row == o.loc.row && p.loc.col == o.loc.col
}

func (p *pipe) link(l *pipe) {
	p.links = append(p.links, l)
}

type island struct {
	start coord
	pipes map[int]map[int]*pipe // [row][col]->pipe
}

func (i *island) get_pipe(r int, c int) *pipe {
	row_map, ok := i.pipes[r]
	if !ok {
		return nil
	}
	p, ok := row_map[c]
	if !ok {
		return nil
	}
	return p
}

func (isl *island) make_links() []coord {
	var to_delete []coord
	for row := range isl.pipes {
		for col := range isl.pipes[row] {
			p := isl.pipes[row][col]
			north := isl.get_pipe(p.loc.row-1, p.loc.col)
			south := isl.get_pipe(p.loc.row+1, p.loc.col)
			east := isl.get_pipe(p.loc.row, p.loc.col+1)
			west := isl.get_pipe(p.loc.row, p.loc.col-1)
			switch p.kind {
			case '|':
				if north != nil && south != nil {
					p.link(north)
					p.link(south)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case '-':
				if east != nil && west != nil {
					p.link(east)
					p.link(west)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case 'L':
				if north != nil && east != nil {
					p.link(north)
					p.link(east)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case 'J':
				if north != nil && west != nil {
					p.link(north)
					p.link(west)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case '7':
				if west != nil && south != nil {
					p.link(south)
					p.link(west)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case 'F':
				if east != nil && south != nil {
					p.link(south)
					p.link(east)
				} else {
					to_delete = append(to_delete, p.loc)
				}
			case 'S':
				if north != nil {
					p.link(north)
				}
				if east != nil {
					p.link(east)
				}
				if south != nil {
					p.link(south)
				}
				if west != nil {
					p.link(west)
				}
				if len(p.links) != 2 {
					log.Fatal("Invalid S pipe")
				}
			}
		}
	}
	return to_delete
}

func build_island(content string) island {
	var isl island
	isl.pipes = make(map[int]map[int]*pipe)
	for row, l := range strings.Split(content, "\n") {
		isl.pipes[row] = make(map[int]*pipe)
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			for col, ru := range line {
				if ru != '.' {
					loc := coord{row: row, col: col}
					var p pipe
					p.int(loc, ru)
					isl.pipes[row][col] = &p
					if ru == 'S' {
						isl.start = loc
					}
				}
			}
		}
	}

	// creating links and removing wrong nodes
	for _, c := range isl.make_links() {
		delete(isl.pipes[c.row], c.col)
	}

	// deleting all dead-ends and unreachable nodes
	cleanup := true
	for cleanup {
		var unuseful_nodes []coord
		for _, rowmap := range isl.pipes {
			for _, p := range rowmap {
				if len(p.links) != 2 {
					unuseful_nodes = append(unuseful_nodes, p.loc)
				}
			}
		}
		cleanup = len(unuseful_nodes) > 0
		if cleanup {
			for _, c := range unuseful_nodes {
				delete(isl.pipes[c.row], c.col)
			}
		}

	}
	return isl
}

func Part1(content string) int {
	isl := build_island(content)
	s_pipe := isl.get_pipe(isl.start.row, isl.start.col)
	var current_1, current_2 *pipe
	for _, p := range s_pipe.links {
		if p != nil {
			if current_1 == nil {
				current_1 = p
			} else {
				current_2 = p
			}
		}
	}
	distance := 1
	prev_1 := s_pipe
	prev_2 := s_pipe
	if current_1 != nil && current_2 != nil {
		for !current_1.same(current_2) {
			log.Println(distance)
			next_1 := current_1.next_pipe(prev_1)
			prev_1 = current_1
			current_1 = next_1
			next_2 := current_2.next_pipe(prev_2)
			prev_2 = current_2
			current_2 = next_2
			distance += 1
		}
	} else {
		log.Fatal("S pipe is floating all alone")
	}

	return distance
}

func Part2(content string) int {
	sum := 0
	return sum
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
