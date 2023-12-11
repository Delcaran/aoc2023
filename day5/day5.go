package day5

// go.dev/blog/pipelines

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

func chunk_gen(seeds []int) <-chan []int {
	out := make(chan []int)

	go func() {
		const chunk_size = 1000000
		blocks := len(seeds) / 2
		index := 0
		for index+1 < len(seeds) {
			current_seed := seeds[index]
			seeds_left := seeds[index+1]
			chunk_num := 0
			chunks := seeds_left/chunk_size + 1
			for seeds_left > 0 {
				limit := min(chunk_size, seeds_left)
				seeds_left -= limit
				chunk := make([]int, 0)
				for count := 0; count < limit; count++ {
					chunk = append(chunk, current_seed)
					current_seed += 1
				}
				chunk_num += 1
				block := int(math.Round((float64(index + 1)) / 2))
				fmt.Printf("%d/%d : %d/%d : %d -> %d\n", block, blocks, chunk_num, chunks, chunk[0], chunk[len(chunk)-1])
				out <- chunk // blocks if out channel is full
			}
			index += 2
		}
		close(out)
	}()
	return out
}

func parse_chunk(input <-chan []int, a *almanac) <-chan int {
	lowest := int(^uint(0) >> 1) // initialized at max int
	out := make(chan int)
	go func() {
		for chunk := range input {
			for _, x := range a.decode("location", chunk) {
				lowest = min(lowest, x)
			}
			out <- lowest
		}
		close(out)
	}()
	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// start output goroutine for each input channel
	// copies values from c to out until c is closed
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// start goroutine to close out once all the output goroutines are done
	// must start after Add call
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func part2(a *almanac) int {
	outs := make([]<-chan int, 0)
	input := chunk_gen(a.seeds)

	for x := 0; x < runtime.NumCPU(); x++ {
		outs = append(outs, parse_chunk(input, a))
	}

	lowest := int(^uint(0) >> 1) // initialized at max int
	for n := range merge(outs...) {
		lowest = min(lowest, n)
	}

	return lowest
}

func Run(content string) (int, int) {
	almanac := initialize(content)

	return part1(&almanac), part2(&almanac)
}
