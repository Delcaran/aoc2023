package day2

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type extraction_t struct {
	red   int
	green int
	blue  int
}

func (c *extraction_t) isPossible(e extraction_t) bool {
	red_ok := c.red <= e.red
	green_ok := c.green <= e.green
	blue_ok := c.blue <= e.blue
	return red_ok && green_ok && blue_ok
}

func (c *extraction_t) power() int {
	return c.red * c.blue * c.green
}

func (c *extraction_t) parse(s string) {
	c.red = 0
	c.green = 0
	c.blue = 0
	regex := regexp.MustCompile(`\s*(?P<num>\d+)\s+(?P<color>red|green|blue)`)
	match := regex.FindAllStringSubmatch(s, -1)
	var err error
	for _, m := range match {
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = m[i]
			}
		}
		switch result["color"] {
		case "red":
			c.red, err = strconv.Atoi(result["num"])
			if err != nil {
				log.Fatal(err)
			}
		case "green":
			c.green, err = strconv.Atoi(result["num"])
			if err != nil {
				log.Fatal(err)
			}
		case "blue":
			c.blue, err = strconv.Atoi(result["num"])
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("Unknown dice")
		}
	}
}

type game_t struct {
	id          int
	extractions []extraction_t
}

func (g *game_t) minimal() extraction_t {
	min := extraction_t{0, 0, 0}
	for _, ext := range g.extractions {
		min.red = max(min.red, ext.red)
		min.green = max(min.green, ext.green)
		min.blue = max(min.blue, ext.blue)
	}
	return min
}

func (g *game_t) power() int {
	min := g.minimal()
	return min.power()
}

func (g *game_t) isPossible(e extraction_t) bool {
	for _, ext := range g.extractions {
		if !ext.isPossible(e) {
			return false
		}
	}
	return true
}

func (g *game_t) parse(s string) bool {
	regex := regexp.MustCompile(`Game (?P<game>\d+):(?P<extractions>.*)`)
	match := regex.FindStringSubmatch(s)
	if len(match) > 0 {
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		var err error
		g.id, err = strconv.Atoi(result["game"])
		if err != nil {
			log.Fatal(err)
		}
		extractions_str := result["extractions"]
		for _, extraction_str := range strings.Split(extractions_str, ";") {
			var e extraction_t
			e.parse(extraction_str)
			g.extractions = append(g.extractions, e)
		}
		return true
	}
	return false
}

func getGames(content string) []game_t {
	var games []game_t
	for _, line := range strings.Split(content, "\n") {
		var g game_t
		if g.parse(line) {
			games = append(games, g)
		}
	}
	return games
}

func part1(content string) int {
	given_cubes := extraction_t{red: 12, green: 13, blue: 14}
	sum := 0
	for _, g := range getGames(content) {
		if g.isPossible(given_cubes) {
			sum += g.id
		}
	}
	return sum
}

func part2(content string) int {
	sum := 0
	for _, g := range getGames(content) {
		sum += g.power()
	}
	return sum
}

func Run(content string) (int, int) {
	return part1(content), part2(content)
}
