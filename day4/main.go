package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/skirklin/aoc2023/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var example = flag.Bool("example", false, "Use example input instead of AoC URL")

var TEST_INPUT = `Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11`

// Card represents a parsed card
type Card struct {
	Number      int
	Left, Right []int
}

// CardParser provides a channel to iterator over cards
type CardParser struct {
	Input string
	Cards chan Card
}

// Parse async parses the input and sends cards to channel
func (cp *CardParser) Parse() {
	// Start goroutine to parse cards in background
	go func() {
		lines := strings.Split(strings.Trim(cp.Input, "\n"), "\n")
		for _, line := range lines {
			// Parse each line into Card struct
			card := parseLine(line)

			// Send to channel
			cp.Cards <- card
		}

		// Close channel when done
		close(cp.Cards)
	}()
}

func parseLine(line string) Card {
	parts := strings.SplitN(line, ": ", 2)

	numString := parts[0][5:]
	number, ok := strconv.Atoi(strings.Trim(numString, " "))
	check(ok)

	values := strings.SplitN(parts[1], " | ", 2)
	left := stringToInts(values[0])
	right := stringToInts(values[1])

	return Card{
		Number: number,
		Left:   left,
		Right:  right,
	}
}

func stringToInts(values string) []int {
	var ints []int
	for _, s := range strings.Fields(values) {
		v, _ := strconv.Atoi(s)
		ints = append(ints, v)
	}
	return ints
}

func intersection(a, b []int) []int {
	var result []int

	// Create a map to track elements in b
	bMap := make(map[int]bool)
	for _, v := range b {
		bMap[v] = true
	}

	// Iterate a, add common elements
	for _, v := range a {
		if bMap[v] {
			result = append(result, v)
		}
	}
	return result
}

func part1(input string) (result int) {
	parser := &CardParser{Input: input, Cards: make(chan Card, 100)}
	parser.Parse()
	for card := range parser.Cards {
		// Use card
		common := intersection(card.Left, card.Right)
		if len(common) > 0 {
			score := 1 << (len(common) - 1)
			result += score
		}
	}

	return result
}

func part2(input string) (result int) {
	parser := &CardParser{Input: input, Cards: make(chan Card, 100)}
	parser.Parse()

	counts := map[int]int{}
	cards := map[int]Card{}
	for card := range parser.Cards {
		cards[card.Number] = card
		counts[card.Number] = 1
	}

	result = 0
	for i := 1; i <= len(cards); i++ {
		card := cards[i]
		count := counts[i]
		result += count
		matches := len(intersection(card.Left, card.Right))
		fmt.Printf("card %d (%d) -> %d\n", i, card.Number, count)
		for j := 1; j <= matches; j++ {
			counts[i+j] += count
		}
	}

	return result
}

func main() {
	flag.Parse()

	var input string
	if *example {
		input = TEST_INPUT
	} else {
		input = utils.GetInputs(2023, 4)
	}

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
