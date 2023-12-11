package day4

import (
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type card struct {
	id      int
	copy    int
	winning []string
	numbers []string
	winners int
	next    *card
}

func (c *card) find_winners() {
	c.winners = 0
	for _, w := range c.winning {
		found := false
		for _, n := range c.numbers {
			if n == w {
				c.winners += 1
				found = true
				break
			}
		}
		if found {
			continue
		}
	}
}

func (c *card) make_copy() card {
	var new card
	new.id = c.id
	new.winning = append(new.winning, c.winning...)
	new.numbers = append(new.numbers, c.numbers...)
	new.winners = c.winners
	new.next = c.next
	new.copy = c.copy + 1
	// caller must link c.next->new
	return new
}

func (c *card) set_next(n *card) {
	c.next = n
}

func (c *card) points() int {
	c.find_winners()
	points := 0
	if c.winners > 0 {
		points = int(math.Pow(float64(2), float64(c.winners-1)))
	}
	return points
}

type lottery struct {
	played  int
	cards   map[int][]card
	current *card
}

func (l *lottery) init(content string) {
	regex := regexp.MustCompile(`Card\s+(?P<card>\d+):\s(?P<winning>[\s+\d+]*)\s\|\s(?P<numbers>[\s+\d+]*)\s`)
	match := regex.FindAllStringSubmatch(content, -1)
	l.cards = make(map[int][]card)
	var prev_card *card
	for _, m := range match {
		result := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = m[i]
			}
		}
		id, err := strconv.Atoi(result["card"])
		if err != nil {
			log.Fatal(err)
		}
		new_card := card{id: id, copy: 0, winning: strings.Fields(result["winning"]), numbers: strings.Fields(result["numbers"])}
		l.cards[new_card.id] = append(l.cards[new_card.id], new_card)
		last := &l.cards[new_card.id][len(l.cards[new_card.id])-1]
		if prev_card != nil {
			prev_card.set_next(last)
		} else {
			// first
			l.current = last
		}
		prev_card = last
	}
}

func (l *lottery) points() int {
	sum := 0
	for _, cards := range l.cards {
		for idx := range cards {
			c := &cards[idx]
			sum += c.points()
		}
	}
	return sum
}

func (l *lottery) play() int {
	l.played = 0
	for l.current != nil {
		c := l.current
		//log.Printf("Playing card %d copy %d", c.id, c.copy)
		// adding newly won cards
		for count := 1; count <= c.winners; count++ {
			id := c.id + count
			//log.Printf("\tCard %d copy %d won copy of card %d", c.id, c.copy, id)
			old_card := &l.cards[id][len(l.cards[id])-1]
			l.cards[id] = append(l.cards[id], old_card.make_copy())
			old_card.set_next(&l.cards[id][len(l.cards[id])-1])
		}
		//log.Printf("Played card %d copy %d", c.id, c.copy)
		l.played += 1
		l.current = c.next
	}
	return l.played
}

func Run(content string) (int, int) {
	var l lottery
	l.init(content)
	return l.points(), l.play()
}
