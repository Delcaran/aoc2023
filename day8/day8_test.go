package day8

import (
	_ "embed"
	"testing"
)

//go:embed input_test_1_1.txt
var input_1_1 string

//go:embed input_test_1_2.txt
var input_1_2 string

//go:embed input_test_2.txt
var input_2 string

func TestRun(t *testing.T) {
	want_part1_1 := 2
	part1_1, _ := Run(input_1_1)
	if part1_1 != want_part1_1 {
		t.Fatalf(`Run 1 -> %d != %d`, part1_1, want_part1_1)
	}

	want_part1_2 := 6
	part1_2, _ := Run(input_1_2)
	if part1_2 != want_part1_2 {
		t.Fatalf(`Run 2 -> %d != %d`, part1_2, want_part1_2)
	}

	want_part2 := 6
	_, part2 := Run(input_2)
	if part2 != want_part2 {
		t.Fatalf(`Run 2 -> %d != %d`, part2, want_part2)
	}
}
