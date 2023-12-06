package main

import (
	"log"
	"os"
	"strconv"

	"github.com/delcaran/aoc2023/day1"
	"github.com/delcaran/aoc2023/day2"
	"github.com/delcaran/aoc2023/day3"
)

func main() {
	test := 0
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
		test, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Wrong test...")
		}
	default:
		day = 3
		test = 1
	}

	switch day {
	case 1:
		day1.Run(test > 0)
	case 2:
		day2.Run(test > 0)
	case 3:
		day3.Run(test > 0)
	default:
		log.Fatal("Day not done")
	}
}
