package main

import (
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

func part1(input string) int {
	var sum = 0
	for _, line := range strings.Split(strings.Trim(input, "\n"), "\n") {
		first := -1
		last := -1
		for _, char := range strings.Split(line, "") {
			intval, err := strconv.Atoi(char)
			if err == nil {
				// it is an int
				if first == -1 {
					first = intval
				}
				last = intval
			}
		}
		if last == -1 || first == -1 {
			panic("no numbers encountered :(")
		}
		sum = sum + (first*10 + last)
	}
	return sum
}

var strMap = map[string]int{
	"0":     0,
	"zero":  0,
	"1":     1,
	"one":   1,
	"2":     2,
	"two":   2,
	"3":     3,
	"three": 3,
	"4":     4,
	"four":  4,
	"5":     5,
	"five":  5,
	"6":     6,
	"six":   6,
	"7":     7,
	"seven": 7,
	"8":     8,
	"eight": 8,
	"9":     9,
	"nine":  9,
}

func extractNumber(line string) (intval int, remainder string) {
	for key, value := range strMap {
		if strings.HasPrefix(line, key) {
			return value, line[1:]
		}
	}
	return -1, line[1:]
}

func part2(input string) int {
	var sum = 0
	for _, line := range strings.Split(strings.Trim(input, "\n"), "\n") {
		first := -1
		last := -1
		intval := -1
		for len(line) > 0 {
			intval, line = extractNumber(line)
			if intval != -1 {
				if first == -1 {
					first = intval
				}
				last = intval
			}
		}
		if last == -1 || first == -1 {
			panic("no numbers encountered :(")
		}
		calibration := first*10 + last
		sum = sum + calibration
	}
	return sum
}

var TEST_INPUT = `1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
`

var TEST_INPUT2 = `two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
`

func main() {
	// inputs := TEST_INPUT
	input := utils.GetInputs(2023, 1)
	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	// input2 := TEST_INPUT2
	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
