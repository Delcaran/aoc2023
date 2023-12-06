package main

import (
	"log"
	"os"
	"strconv"

	"github.com/delcaran/aoc2023/day1"
	"github.com/delcaran/aoc2023/day2"
	"github.com/delcaran/aoc2023/day3"
	"github.com/delcaran/aoc2023/day4"
)

func main() {
	test := ""
	day := 0
	var err error
	switch len(os.Args[1:]) {
	case 1:
		day, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Wrong day...")
		}
	case 2:
		day, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Wrong day...")
		}
		test = os.Args[2]
	default:
		day = 3
		test = "test"
	}

	switch day {
	case 1:
		day1.Run(test)
	case 2:
		day2.Run(test)
	case 3:
		day3.Run(test)
	case 4:
		day4.Run(test)
	default:
		log.Fatal("Day not done")
	}
}
