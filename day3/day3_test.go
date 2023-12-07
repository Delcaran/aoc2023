package day3

import (
	_ "embed"
	"testing"
)

//go:embed input_test.txt
var input string

func TestRun(t *testing.T) {
	want_part1 := 4361
	want_part2 := 467835
	part1, part2, err := Run(input)
	if err != nil || part1 != want_part1 || part2 != want_part2 {
		t.Fatalf(`Run = %d, %d, %v, want match for %d, %d, nil`, part1, part2, err, want_part1, want_part2)
	}
}
