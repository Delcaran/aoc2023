package day10

import (
	"log"

	"github.com/fatih/color"
)

const (
	NONE  = iota
	PATH  = iota
	LEFT  = iota
	RIGHT = iota
)

type vertex struct {
	loc    coord
	kind   rune
	edges  []*edge
	relpos int
	end    bool
}

func (p *vertex) print() string {
	col := color.New(color.FgBlack)
	if p.relpos == PATH {
		if p.kind == 'S' {
			col = col.Add(color.BgGreen)
		} else if p.end {
			col = col.Add(color.BgRed)
		} else {
			col = col.Add(color.BgYellow)
		}
	} else {
		switch p.relpos {
		case LEFT:
			col = col.Add(color.BgMagenta)
		case RIGHT:
			col = col.Add(color.BgCyan)
		default:
			col = col.Add(color.BgWhite)
		}
	}
	return col.Sprintf("%c", p.kind)
}

func (p *vertex) mark(pos int) {
	if p.relpos != PATH {
		p.relpos = pos
	}
}

func (p *vertex) get_next(incoming *vertex) *vertex {
	p.relpos = PATH
	for _, e := range p.edges {
		other_side := e.get_other(p)
		if !other_side.same(incoming) {
			other_side.relpos = PATH
			return other_side
		}
	}
	log.Fatalf("Pipe %d,%d is a dead end!!!", p.loc.row, p.loc.col)
	return nil
}

func (p *vertex) same(o *vertex) bool {
	res := p.loc.row == o.loc.row && p.loc.col == o.loc.col
	return res
}

func (pipe *vertex) linkable_north() bool {
	res := pipe.kind == 'S' || pipe.kind == '|' || pipe.kind == 'J' || pipe.kind == 'L'
	return res
}

func (pipe *vertex) linkable_east() bool {
	res := pipe.kind == 'S' || pipe.kind == '-' || pipe.kind == 'L' || pipe.kind == 'F'
	return res
}

func (pipe *vertex) linkable_west() bool {
	res := pipe.kind == 'S' || pipe.kind == '-' || pipe.kind == 'J' || pipe.kind == '7'
	return res
}

func (pipe *vertex) linkable_south() bool {
	res := pipe.kind == 'S' || pipe.kind == '|' || pipe.kind == 'F' || pipe.kind == '7'
	return res
}

func (a *vertex) make_link(b *vertex) {
	exists := false
	for _, e := range a.edges {
		exists = b.same(e.get_other(a)) || exists
	}
	if !exists {
		link := &edge{a: a, b: b}
		a.edges = append(a.edges, link)
		b.edges = append(b.edges, link)
	}
}

func (south *vertex) make_link_north(north *vertex) int {
	if north != nil && south.linkable_north() && north.linkable_south() {
		south.make_link(north)
		return 1
	}
	return 0
}

func (west *vertex) make_link_east(east *vertex) int {
	if east != nil && west.linkable_east() && east.linkable_west() {
		west.make_link(east)
		return 1
	}
	return 0
}

func (east *vertex) make_link_west(west *vertex) int {
	if west != nil && east.linkable_west() && west.linkable_east() {
		east.make_link(west)
		return 1
	}
	return 0
}

func (north *vertex) make_link_south(south *vertex) int {
	if south != nil && north.linkable_south() && south.linkable_north() {
		north.make_link(south)
		return 1
	}
	return 0
}

func (pipe *vertex) unlink() {
	for _, e := range pipe.edges {
		e.unlink(pipe) // this vertex is not reachable
	}
	pipe.edges = pipe.edges[:0] // this vertex has no knowledge of the graph
}
