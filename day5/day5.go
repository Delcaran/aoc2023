package day5

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type map_entry struct {
	src        int
	dest       int
	data_range int
}

func (m *map_entry) convert(input int) (int, error) {
	if input >= m.src && input < m.src+m.data_range {
		offset := input - m.src
		return m.dest + offset, nil
	}
	return 0, fmt.Errorf("out of range")
}

type generic_map struct {
	src_type  string
	dest_type string
	mapping   []map_entry
}

func (m *generic_map) convert(input []int) ([]int, string) {
	var output []int
	for _, in := range input {
		found := false
		for _, m := range m.mapping {
			out, err := m.convert(in)
			if err == nil {
				output = append(output, out)
				found = true
				break
			}
		}
		if !found {
			output = append(output, in)
		}
	}
	return output, m.dest_type
}

type almanac struct {
	seeds         []int
	seed_list_idx int
	maps          map[string]generic_map // src_type as key
}

func initialize(content string) almanac {
	var a almanac
	regex_seeds := regexp.MustCompile(`seeds:\s(?P<seeds>.*)`)
	regex_header := regexp.MustCompile(`(?P<src>\w+)-to-(?P<dest>\w+)\smap:`)
	a.maps = make(map[string]generic_map)
	a.seed_list_idx = 0
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
				current_map.src_type = result["src"]
				current_map.dest_type = result["dest"]
				current_map.mapping = make([]map_entry, 0)
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
				current_map.mapping = append(current_map.mapping,
					map_entry{src: src_begin, dest: dest_begin, data_range: count})
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

func (a *almanac) decode(destination string, data []int) []int {
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
	data := a.seeds
	locations := a.decode("location", data)
	lowest_location := int(^uint(0) >> 1) // initialized at max int
	for _, x := range locations {
		lowest_location = min(lowest_location, x)
	}
	return lowest_location
}

func part2(a *almanac) int {
	lowest_location := int(^uint(0) >> 1) // initialized at max int
	for a.seed_list_idx < len(a.seeds)-2 {
		log.Printf("a.seed_idx %d\n", a.seed_list_idx)
		first_seed := a.seeds[a.seed_list_idx]
		seed := first_seed
		seeds_count := a.seeds[a.seed_list_idx+1]
		last_seed := first_seed + seeds_count
		const seed_limit = 100
		for seed < last_seed {
			var seed_portion []int
			for ; seed < min(seed+seed_limit, last_seed); seed++ {
				seed_portion = append(seed_portion, seed)
			}
			for _, x := range a.decode("location", seed_portion) {
				lowest_location = min(lowest_location, x)
			}
		}
		a.seed_list_idx += 2
	}
	return lowest_location
}

func Run(content string) (int, int, error) {
	almanac := initialize(content)

	return part1(&almanac), part2(&almanac), nil
}
