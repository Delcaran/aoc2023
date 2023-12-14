package day10

import (
	_ "embed"
	"testing"
)

//go:embed input_test_1_1.txt
var input_1_1 string

//go:embed input_test_1_2.txt
var input_1_2 string

//go:embed input_test_2_1.txt
var input_2_1 string

//go:embed input_test_2_2.txt
var input_2_2 string

//go:embed input_test_2_3.txt
var input_2_3 string

//go:embed input_test_2_4.txt
var input_2_4 string

func TestRun(t *testing.T) {
	/*
		want_part1_1 := 4
		part1_1 := Part1(input_1_1)
		if part1_1 != want_part1_1 {
			t.Fatalf(`Part1 1 -> %d != %d`, part1_1, want_part1_1)
		}

		want_part1_2 := 8
		part1_2 := Part1(input_1_2)
		if part1_2 != want_part1_2 {
			t.Fatalf(`Part1 2 -> %d != %d`, part1_2, want_part1_2)
		}
	*/

	want_part2_1 := 4
	part2_1 := Part2(input_2_1)
	if part2_1 != want_part2_1 {
		t.Fatalf(`Part2 1 -> %d != %d`, part2_1, want_part2_1)
	}

	want_part2_2 := 4
	part2_2 := Part2(input_2_2)
	if part2_2 != want_part2_2 {
		t.Fatalf(`Part2 2 -> %d != %d`, part2_2, want_part2_2)
	}

	want_part2_3 := 8
	part2_3 := Part2(input_2_3)
	if part2_3 != want_part2_3 {
		t.Fatalf(`Part2 1 -> %d != %d`, part2_3, want_part2_3)
	}

	want_part2_4 := 10
	part2_4 := Part2(input_2_4)
	if part2_4 != want_part2_4 {
		t.Fatalf(`Part1 2 -> %d != %d`, part2_4, want_part2_4)
	}
}
