package day6

import (
	"log"
	"strconv"
	"strings"
)

type method struct {
	release_ms int
	speed      int
	distance   int
}

type race struct {
	time    int
	record  int
	methods []method
}

func (r *race) estimate_brute_force() {
	for release_ms := 1; release_ms < r.time; release_ms++ {
		m := method{speed: release_ms, release_ms: release_ms, distance: release_ms * (r.time - release_ms)}
		if m.distance > r.record {
			r.methods = append(r.methods, m)
		}
	}
}

func (r *race) margin_of_error() int {
	if len(r.methods) == 0 {
		r.estimate_brute_force()
	}
	return len(r.methods)
}

func part1(content string) int {
	races := make([]race, 0)
	lines := strings.Split(content, "\n")
	times := strings.Fields(lines[0])
	distances := strings.Fields(lines[1])
	for r := 1; r < len(times); r++ {
		time, err := strconv.Atoi(times[r])
		if err != nil {
			log.Fatal(err)
		}
		distance, err := strconv.Atoi(distances[r])
		if err != nil {
			log.Fatal(err)
		}
		races = append(races, race{time: time, record: distance})
	}

	margin := 1
	for _, r := range races {
		margin *= r.margin_of_error()
	}
	return margin
}

func part2(content string) int {
	lines := strings.Split(content, "\n")
	time, err := strconv.Atoi(strings.Join(strings.Fields(strings.Split(lines[0], ":")[1]), ""))
	if err != nil {
		log.Fatal(err)
	}
	distance, err := strconv.Atoi(strings.Join(strings.Fields(strings.Split(lines[1], ":")[1]), ""))
	if err != nil {
		log.Fatal(err)
	}
	race := race{time: time, record: distance}
	margin := race.margin_of_error()

	return margin
}

func Run(content string) (int, int, error) {
	return part1(content), part2(content), nil
}
