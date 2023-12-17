package main

import (
	"flag"
	"fmt"

	"github.com/skirklin/aoc2023/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var example = flag.Bool("example", false, "Use example input instead of AoC URL")

var TEST_INPUT = `
...
`

func part1(input string) (result int) {
	return result
}

func part2(input string) (result int) {
	return result
}

func main() {
	flag.Parse()

	var input string
	if *example {
		input = TEST_INPUT
	} else {
		input = utils.GetInputs(2023, -1)
	}

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
