package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Extraction struct {
	red   int
	green int
	blue  int
}

func (c *Extraction) isPossible(e Extraction) bool {
	red_ok := c.red <= e.red
	green_ok := c.green <= e.green
	blue_ok := c.blue <= e.blue
	return red_ok && green_ok && blue_ok
}

func (c *Extraction) power() int {
	return c.red * c.blue * c.green
}

func (c *Extraction) parse(s string) {
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

type Game struct {
	id          int
	extractions []Extraction
}

func (g *Game) minimal() Extraction {
	min := Extraction{0, 0, 0}
	for _, ext := range g.extractions {
		min.red = max(min.red, ext.red)
		min.green = max(min.green, ext.green)
		min.blue = max(min.blue, ext.blue)
	}
	return min
}

func (g *Game) power() int {
	min := g.minimal()
	return min.power()
}

func (g *Game) isPossible(e Extraction) bool {
	for _, ext := range g.extractions {
		if ext.isPossible(e) == false {
			return false
		}
	}
	return true
}

func (g *Game) parse(s string) bool {
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
			var e Extraction
			e.parse(extraction_str)
			g.extractions = append(g.extractions, e)
		}
		return true
	}
	return false
}

func getGames(content string) []Game {
	var games []Game
	for _, line := range strings.Split(content, "\n") {
		var g Game
		if g.parse(line) {
			games = append(games, g)
		}
	}
	return games
}

func part1(content string) int {
	given_cubes := Extraction{red: 12, green: 13, blue: 14}
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

func main() {
	buffer, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	content := string(buffer[:])
	log.Printf("Part 1: %d\n", part1(content))
	log.Printf("Part 2: %d\n", part2(content))
}
