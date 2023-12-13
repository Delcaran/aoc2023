package day9

import (
	_ "embed"
	"testing"
)

//go:embed input_test.txt
var input string

func TestRun(t *testing.T) {
	want_part1 := 114
	want_part2 := 2
	part1, part2 := Run(input)
	if part1 != want_part1 || part2 != want_part2 {
		t.Fatalf(`Run = (%d, %d) != (%d, %d)`, part1, part2, want_part1, want_part2)
	}
}
