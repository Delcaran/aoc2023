package day5

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type generic_map struct {
	src_type  string
	dest_type string
	mapping   map[int]int
}

func (m *generic_map) init(src string, dest string) {
	m.src_type = src
	m.dest_type = dest
	m.mapping = make(map[int]int)
}

func (m *generic_map) convert(input []int) ([]int, string) {
	var output []int
	for i := range input {
		in := input[i]
		out, ok := m.mapping[in]
		if !ok {
			out = in
		}
		output = append(output, out)
	}
	return output, m.dest_type
}

type almanac struct {
	seeds []int
	maps  map[string]generic_map // src_type as key
}

func initialize(content string) almanac {
	var a almanac
	regex_seeds := regexp.MustCompile(`seeds:\s(?P<seeds>.*)`)
	regex_header := regexp.MustCompile(`(?P<src>\w+)-to-(?P<dest>\w+)\smap:`)
	a.maps = make(map[string]generic_map)
	var current_map generic_map
	for _, line := range strings.Split(content, "\n") {
		// looking for seeds (once)
		if len(a.seeds) == 0 {
			match := regex_seeds.FindStringSubmatch(line)
			if len(match) > 0 {
				result := make(map[string]string)
				for i, name := range regex_seeds.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = match[i]
					}
				}
				for _, s := range strings.Fields(result["seeds"]) {
					n, err := strconv.Atoi(s)
					if err != nil {
						log.Fatalln(err)
					}
					a.seeds = append(a.seeds, n)
				}
			}
		}
		// look for maps
		if len(current_map.src_type) == 0 {
			match := regex_header.FindStringSubmatch(line)
			if len(match) > 0 {
				result := make(map[string]string)
				for i, name := range regex_header.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = match[i]
					}
				}
				current_map.init(result["src"], result["dest"])
			}
		} else {
			// parsing map
			fields := strings.Fields(line)
			if len(fields) > 0 {
				dest_begin, err := strconv.Atoi(fields[0])
				if err != nil {
					log.Fatalln(err)
				}
				src_begin, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Fatalln(err)
				}
				count, err := strconv.Atoi(fields[2])
				if err != nil {
					log.Fatalln(err)
				}
				for i := 0; i < count; i++ {
					current_map.mapping[src_begin+i] = dest_begin + i
				}
			} else {
				// finisced map
				a.maps[current_map.src_type] = current_map
				current_map.src_type = ""
			}
		}
	}
	// check if missing last empty line
	if current_map.src_type != "" {
		a.maps[current_map.src_type] = current_map
		current_map.src_type = ""
	}
	// we got all maps, now build the conversion chain
	return a
}

func (a *almanac) decode(destination string) []int {
	data := a.seeds

	ok := true
	for src := "seed"; ok; {
		m, have_map := a.maps[src]
		if have_map {
			data, src = m.convert(data)
			ok = have_map
		}
		if src == destination {
			// reached my destination
			ok = false
		}
	}
	return data
}

func part1(a *almanac) int {
	locations := a.decode("location")
	lowest_location := int(^uint(0) >> 1) // initialized at max int
	for _, x := range locations {
		lowest_location = min(lowest_location, x)
	}
	return lowest_location
}

func Run(content string) (int, int, error) {
	almanac := initialize(content)

	return part1(&almanac), 0, nil
}
