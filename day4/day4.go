package day4

import (
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

type card struct {
	winning []string
	numbers []string
	winners []int // index of numbers
}

func (c *card) find_winners() {
	for _, w := range c.winning {
		for idx, n := range c.numbers {
			if n == w {
				c.winners = append(c.winners, idx)
			}
		}
	}
}

func (c *card) points() int {
	c.find_winners()
	if len(c.winners) > 0 {
		return int(math.Pow(float64(2), float64(len(c.winners)-1)))
	}
	return 0
}

type lottery struct {
	cards []card
}

func (l *lottery) init(content string) {
	regex := regexp.MustCompile(`Card (?P<card>\d+):\s(?P<winning>[\s+\d+]*)\s\|\s(?P<numbers>[\s+\d+]*)\s`)
	match := regex.FindAllStringSubmatch(content, -1)
	for _, m := range match {
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = m[i]
			}
		}
		c := card{winning: strings.Fields(result["winning"]), numbers: strings.Fields(result["numbers"])}
		l.cards = append(l.cards, c)
	}
}

func (l *lottery) points() int {
	sum := 0
	for _, c := range l.cards {
		sum += c.points()
	}
	return sum
}

func Run(test string) {
	if len(test) > 0 {
		test += "_"
	}
	buffer, err := os.ReadFile("day4/" + test + "input.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(buffer[:])
	var l lottery
	l.init(content)
	log.Println(l.points())
}
