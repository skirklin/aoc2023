package main

import (
	"fmt"

	"github.com/skirklin/aoc2023/utils"
)

func accumulate(input string) int {
	total := 0
	for _, c := range input {
		switch c {
		case '(':
			total++
		case ')':
			total--
		default:
			fmt.Printf("Fuuuck %c", c)
		}
	}
	return total
}

func firstdrop(input string) int {
	total := 0
	for i, c := range input {
		switch c {
		case '(':
			total++
		case ')':
			total--
		default:
			fmt.Printf("Fuuuck %c", c)
		}
		if total == -1 {
			return i + 1
		}
	}
	return -1
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	utils.GetInputs(2023, 1)
}
