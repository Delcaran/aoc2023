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

func (h hand) map_cards() map[rune]int {
	labels := make(map[rune]int)
	for _, r := range h.cards {
		_, found := labels[r]
		if found {
			labels[r] += 1
		} else {
			labels[r] = 1
		}
	}

	return labels
}

func (h hand) find_type_simple() int {
	var hand_type int
	labels := h.map_cards()
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

func (h hand) find_type(jokers bool) int {
	if jokers {
		// find most numerous card
		labels := h.map_cards()
		max_same_cards := 0
		max_card := 'J'
		for x := range labels {
			if x != 'J' {
				if max_same_cards < labels[x] {
					max_same_cards = labels[x]
					max_card = x
				} else {
					if max_same_cards == labels[x] {
						val_current := get_card_value(max_card, jokers)
						val_new := get_card_value(x, jokers)
						if val_new > val_current {
							max_card = x
						}
					}
				}
			}
		}
		if max_card == 'J' {
			max_card = 'A'
		}
		// replace jokers with that card
		h.cards = strings.ReplaceAll(h.cards, "J", string(max_card))
	}
	return h.find_type_simple()
}

func get_card_value(card rune, jokers bool) int {
	card_value := "23456789TJQKA"
	if jokers {
		card_value = "J23456789TQKA"
	}
	return strings.IndexRune(card_value, card)
}

func (h *hand) calc_value(jokers bool) {
	powers := []int{8, 6, 4, 2, 0}
	h.value = h.find_type(jokers)
	for i, r := range h.cards {
		val := get_card_value(r, jokers)
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
			h.calc_value(false)
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
			h.calc_value(true)
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

func Run(content string) (int, int, error) {
	return part1(content), part2(content), nil
}
