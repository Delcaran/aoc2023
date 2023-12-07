package day1

import (
	"strings"
)

func Part1(content string) int {
	sum := 0

	for _, line := range strings.Split(content, "\n") {
		first := 0
		last := 0
		for _, letter := range line {
			n := int(letter - '0')
			if n > 0 && n <= 9 {
				if first == 0 {
					first = n * 10
					last = n
				} else {
					last = n
				}
			}
		}
		sum += (first + last)
	}

	return sum
}

func Part2(content string) int {
	sum := 0
	valid_spelled_numbers := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	}

	for _, line := range strings.Split(content, "\n") {
		founds := make(map[int]int)
		for n := 0; n < len(line); n++ {
			// look for regular digits
			digit := int(line[n] - '0')
			if digit > 0 && digit <= 9 {
				founds[n] = digit
			} else {
				// look for spelled-out digits
				for v, w := range valid_spelled_numbers {
					index := strings.Index(line[n:], w)
					if index >= 0 {
						founds[n+index] = v + 1
					}
				}
			}
		}

		// summing
		min_idx := len(line) + 1
		max_idx := 0
		for index := range founds {
			if index < min_idx {
				min_idx = index
			}
			if index > max_idx {
				max_idx = index
			}
		}
		value := founds[min_idx]*10 + founds[max_idx]
		sum += value
	}

	return sum
}

func Run(content string) (int, int, error) {
	return Part1(content), Part2(content), nil
}
