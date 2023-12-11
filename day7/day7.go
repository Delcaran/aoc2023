package day7

import (
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	HIGH_CARD  = iota
	PAIR       = iota
	TWO_PAIR   = iota
	THREE_KIND = iota
	FULL       = iota
	FOUR_KIND  = iota
	FIVE_KIND  = iota
)

type hand struct {
	cards    string
	bet      int
	value    int
	rank     int
	winnings int
}

func (h hand) find_type() int {
	var hand_type int
	labels := make(map[rune]int)
	for _, r := range h.cards {
		_, found := labels[r]
		if found {
			labels[r] += 1
		} else {
			labels[r] = 1
		}
	}
	max_same_cards := 0
	for x := range labels {
		max_same_cards = max(max_same_cards, labels[x])
	}
	switch max_same_cards {
	case 5:
		hand_type = FIVE_KIND
	case 4:
		hand_type = FOUR_KIND
	case 3:
		if len(labels) == 2 {
			hand_type = FULL
		} else {
			hand_type = THREE_KIND
		}
	case 2:
		if len(labels) == 3 {
			hand_type = TWO_PAIR
		} else {
			hand_type = PAIR
		}
	case 1:
		hand_type = HIGH_CARD
	}
	power := int(math.Pow10(10))
	return hand_type * power
}

func (h *hand) calc_value() {
	const card_value = "23456789TJQKA"
	powers := []int{8, 6, 4, 2, 0}
	h.value = h.find_type()
	for i, r := range h.cards {
		val := strings.IndexRune(card_value, r)
		value := val * int(math.Pow10(powers[i]))
		h.value += value
	}
}

func (h *hand) calc_win(rank int) {
	h.rank = rank
	h.winnings = h.rank * h.bet
}

func mergeSort(items []hand) []hand {
	if len(items) < 2 {
		return items
	}
	first := mergeSort(items[:len(items)/2])
	second := mergeSort(items[len(items)/2:])
	return merge(first, second)
}

func merge(a []hand, b []hand) []hand {
	final := []hand{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].value < b[j].value {
			final = append(final, a[i])
			i++
		} else {
			final = append(final, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		final = append(final, a[i])
	}
	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}

func part1(content string) int {
	var hands []hand
	for _, line := range strings.Split(content, "\n") {
		if len(line) > 0 {
			var h hand
			data := strings.Fields(line)
			h.cards = data[0]
			bet, err := strconv.Atoi(data[1])
			if err != nil {
				log.Fatal(err)
			}
			h.bet = bet
			h.calc_value()
			hands = append(hands, h)
		}
	}
	sorted := mergeSort(hands)
	winnings := 0
	for rank := range sorted {
		sorted[rank].calc_win(rank + 1)
		winnings += sorted[rank].winnings
	}

	return winnings
}

func part2(content string) int {
	return 0
}

func Run(content string) (int, int, error) {
	return part1(content), part2(content), nil
}
