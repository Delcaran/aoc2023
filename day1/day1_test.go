package day1

import (
	_ "embed"
	"testing"
)

//go:embed input1_test.txt
var input1 string

//go:embed input2_test.txt
var input2 string

func TestRun(t *testing.T) {
	want_part1 := 142
	part1 := Part1(input1)
	if part1 != want_part1 {
		t.Fatalf(`Output 1 = %d, want match for %d`, part1, want_part1)
	}
	want_part2 := 281
	part2 := Part2(input2)
	if part2 != want_part2 {
		t.Fatalf(`Output 2 = %d, want match for %d`, part2, want_part2)
	}
}
