package day4

import (
	"fmt"
	"log"
	"os"
)

func Run(test string) {
	if len(test) > 0 {
		test += "_"
	}
	buffer, err := os.ReadFile("day4/" + test + "input.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(buffer[:])
	fmt.Println(content)
}
