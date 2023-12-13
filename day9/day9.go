package day9

import (
	"log"
	"strconv"
	"strings"
)

type sequence struct {
	values     []int
	all_zeroes bool
	diff       *sequence
}

func (s *sequence) init(data []int, all_zeroes bool) {
	s.all_zeroes = all_zeroes
	if s.all_zeroes {
		s.values = nil // don't waste space for a sequence of zeroes
		s.diff = nil
	} else {
		s.values = append(s.values, data...)
		var diffs []int
		all_zeroes := true
		for x := 1; x < len(s.values); x++ {
			diff := s.values[x] - s.values[x-1]
			diffs = append(diffs, diff)
			all_zeroes = all_zeroes && diff == 0
		}
		s.diff = &sequence{}
		s.diff.init(diffs, all_zeroes)
	}
}

func (s *sequence) get_next() int {
	if s.all_zeroes {
		return 0
	} else {
		return s.values[len(s.values)-1] + s.diff.get_next()
	}
}

func (s *sequence) get_prev() int {
	if s.all_zeroes {
		return 0
	} else {
		return s.values[0] - s.diff.get_prev()
	}
}

func Part1(content string) int {
	sum := 0
	for _, line := range strings.Split(content, "\n") {
		if len(line) > 0 {
			sequence_str := strings.Fields(line)
			var numlist []int
			for _, str := range sequence_str {
				num, err := strconv.Atoi(str)
				if err != nil {
					log.Fatal(err)
				}
				numlist = append(numlist, num)
			}
			var sequence sequence
			sequence.init(numlist, false)
			sum += sequence.get_next()
		}
	}
	return sum
}

func Part2(content string) int {
	sum := 0
	for _, line := range strings.Split(content, "\n") {
		if len(line) > 0 {
			sequence_str := strings.Fields(line)
			var numlist []int
			for _, str := range sequence_str {
				num, err := strconv.Atoi(str)
				if err != nil {
					log.Fatal(err)
				}
				numlist = append(numlist, num)
			}
			var sequence sequence
			sequence.init(numlist, false)
			sum += sequence.get_prev()
		}
	}
	return sum
}

func Run(content string) (int, int) {
	return Part1(content), Part2(content)
}
