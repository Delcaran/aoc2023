package day8

import (
	_ "embed"
	"testing"
)

//go:embed input_test1.txt
var input_1 string

//go:embed input_test2.txt
var input_2 string

func TestRun(t *testing.T) {
	want_part1_1 := 2
	part1_1, _ := Run(input_1)
	if part1_1 != want_part1_1 {
		t.Fatalf(`Run -> %d != %d`, part1_1, want_part1_1)
	}

	want_part1_2 := 6
	_, part1_2 := Run(input_2)
	if part1_2 != want_part1_2 {
		t.Fatalf(`Run -> %d != %d`, part1_2, want_part1_2)
	}
}
