package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	buffer, err := os.ReadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	content := string(buffer[:])
	sum := 0

	for _, line := range strings.Split(content, "\n") {
		first := 0
		last := 0
		for _, letter := range line {
			n := int(letter - '0')
			if n > 0 && n <= 9 {
				if first == 0 {
					first = n * 10
					last = n
				} else {
					last = n
				}
			}
		}
		sum += (first + last)
	}

	log.Println(sum)
}
