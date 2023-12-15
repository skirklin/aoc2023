package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/skirklin/aoc2023/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type draw struct {
	red   int
	blue  int
	green int
}

type game struct {
	id    int
	hands []draw
}

var LINE_REGEX = regexp.MustCompile(`^Game (?P<gameId>\d+): (.*)`)
var BLOCK_REGEX = regexp.MustCompile(`^(\d+) (red|green|blue)`)

func parseLine(line string) game {
	lineMatch := LINE_REGEX.FindStringSubmatch(line)
	gameId, err := strconv.Atoi(lineMatch[1])
	check(err)
	draws := []draw{}
	for _, drawString := range strings.Split(lineMatch[2], "; ") {
		d := draw{}
		for _, blockString := range strings.Split(drawString, ", ") {
			blockMatch := BLOCK_REGEX.FindStringSubmatch(blockString)
			num, err := strconv.Atoi(blockMatch[1])
			check(err)
			switch blockMatch[2] {
			case "red":
				d.red = num
			case "blue":
				d.blue = num
			case "green":
				d.green = num
			}
		}
		draws = append(draws, d)
	}
	g := game{gameId, draws}
	// fmt.Println(g, line)
	return g
}

func parseInput(input string) []game {
	var games []game
	for _, line := range strings.Split(strings.Trim(input, "\n"), "\n") {
		game := parseLine(line)
		games = append(games, game)
	}
	return games
}

func AllValid(draws []draw) bool {
	for _, hand := range draws {
		if hand.red > 12 || hand.green > 13 || hand.blue > 14 {
			return false
		}
	}
	return true
}

func part1(input string) int {
	games := parseInput(input)
	count := 0
	for _, game := range games {
		if AllValid(game.hands) {
			count += game.id
		}
	}
	return count
}

func part2(input string) int {
	games := parseInput(input)
	total := 0
	for _, game := range games {
		minBag := draw{0, 0, 0}
		for _, hand := range game.hands {
			minBag.red = max(minBag.red, hand.red)
			minBag.green = max(minBag.green, hand.green)
			minBag.blue = max(minBag.blue, hand.blue)
		}
		product := minBag.red * minBag.green * minBag.blue
		total += product
	}
	return total
}

var TEST_INPUT = `Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green`

func main() {
	// input := TEST_INPUT
	input := utils.GetInputs(2023, 2)
	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
