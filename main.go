package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/delcaran/aoc2023/day1"
	"github.com/delcaran/aoc2023/day2"
	"github.com/delcaran/aoc2023/day3"
	"github.com/delcaran/aoc2023/day4"
	"github.com/delcaran/aoc2023/day5"
	"github.com/delcaran/aoc2023/day6"
	"github.com/delcaran/aoc2023/day7"
	"github.com/delcaran/aoc2023/day8"
	"github.com/delcaran/aoc2023/day9"
)

//go:embed day1/input.txt
var day1_input string

//go:embed day2/input.txt
var day2_input string

//go:embed day3/input.txt
var day3_input string

//go:embed day4/input.txt
var day4_input string

//go:embed day5/input.txt
var day5_input string

//go:embed day6/input.txt
var day6_input string

//go:embed day7/input.txt
var day7_input string

//go:embed day8/input.txt
var day8_input string

//go:embed day9/input.txt
var day9_input string

func main() {
	day := 0
	var err error
	switch len(os.Args[1:]) {
	case 1:
		day, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Wrong day...")
		}
	}

	switch day {
	case 1:
		fmt.Println(day1.Run(day1_input))
	case 2:
		fmt.Println(day2.Run(day2_input))
	case 3:
		fmt.Println(day3.Run(day3_input))
	case 4:
		fmt.Println(day4.Run(day4_input))
	case 5:
		fmt.Println(day5.Run(day5_input))
	case 6:
		fmt.Println(day6.Run(day6_input))
	case 7:
		fmt.Println(day7.Run(day7_input))
	case 8:
		fmt.Println(day8.Run(day8_input))
	case 9:
		fmt.Println(day9.Run(day9_input))
	default:
		log.Fatal("Day not done")
	}
}
