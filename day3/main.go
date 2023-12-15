package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/skirklin/aoc2023/utils"
)

// Cell represents a cell in the 2D array.
type Cell struct {
	Row, Col int
}

// NumberInfo represents information about a multi-digit number.
type NumberInfo struct {
	Start  Cell
	Length int
}

// Schematic represents the 2D array of runes.
type Schematic struct {
	Grid [][]rune
}

// NewSchematic creates a new Schematic from a string.
func NewSchematic(input string) *Schematic {
	var grid [][]rune

	for _, line := range strings.Fields(input) {
		lineRunes := []rune(line)
		grid = append(grid, lineRunes)
	}

	return &Schematic{Grid: grid}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// FindNumbers finds the locations and lengths of multi-digit numbers in the schematic.
func (s *Schematic) GetNumberValue(numInfo NumberInfo) int {
	currValue := 0
	for i := 0; i < numInfo.Length; i++ {
		char := s.Grid[numInfo.Start.Row][numInfo.Start.Col+i]
		if !unicode.IsDigit(char) {
			panic(fmt.Sprintf("bad conversion of %s", string(char)))
		}
		currValue = currValue*10 + int(char-'0')
	}
	return currValue
}

// FindNumbers finds the locations and lengths of multi-digit numbers in the schematic.
func (s *Schematic) FindNumbers() []NumberInfo {
	var numbers []NumberInfo
	var currNumber *NumberInfo

	for i, row := range s.Grid {
		for j, char := range row {
			if unicode.IsDigit(char) {
				switch {
				case currNumber == nil:
					// Start of a number, find its length
					start := Cell{Row: i, Col: j}
					currNumber = &NumberInfo{start, 1}
				case currNumber != nil:
					currNumber.Length += 1
				}
			} else {
				if currNumber != nil {
					numbers = append(numbers, *currNumber)
					currNumber = nil
				}
			}
		}
		if currNumber != nil {
			numbers = append(numbers, *currNumber)
			currNumber = nil
		}
	}

	return numbers
}

// CheckNeighbors checks neighboring characters for each number.
func (s *Schematic) CheckNeighbors(number NumberInfo) []Cell {
	result := []Cell{}
	// Coordinates of the number
	row, col := number.Start.Row, number.Start.Col

	// Check all possible adjacent positions
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= number.Length; dj++ {
			if di == 0 && dj == 0 {
				// Skip the current position
				continue
			}

			ni, nj := row+di, col+dj
			if ni >= 0 && ni < len(s.Grid) && nj >= 0 && nj < len(s.Grid[ni]) {
				// Check if the character in the adjacent position is a symbol of interest
				adjacentChar := s.Grid[ni][nj]
				if !(unicode.IsDigit(adjacentChar) || adjacentChar == '.') {
					// If a symbol is found, mark the number as having adjacent symbols
					result = append(result, Cell{ni, nj})
				}
			}
		}
	}
	return result
}

var TEST_INPUT = `467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`

func part1(input string) int {
	schematic := NewSchematic(input)

	// Find the locations and lengths of multi-digit numbers
	numbers := schematic.FindNumbers()

	total := 0
	for _, number := range numbers {
		hasAdjacentSymbol := schematic.CheckNeighbors(number)

		if len(hasAdjacentSymbol) > 0 {
			numValue := schematic.GetNumberValue(number)
			total += numValue
		}
	}

	return total
}

func part2(input string) int {
	schematic := NewSchematic(input)

	// Find the locations and lengths of multi-digit numbers
	numbers := schematic.FindNumbers()

	total := 0
	gearMap := map[Cell][]NumberInfo{}
	for _, number := range numbers {
		for _, neighbor := range schematic.CheckNeighbors(number) {
			if schematic.Grid[neighbor.Row][neighbor.Col] == '*' {
				if gearMap[neighbor] == nil {
					gearMap[neighbor] = []NumberInfo{}
				}
				gearMap[neighbor] = append(gearMap[neighbor], number)
			}
		}
	}

	for _, numbers := range gearMap {
		if len(numbers) == 2 {
			ratio := schematic.GetNumberValue(numbers[0]) * schematic.GetNumberValue(numbers[1])
			total += ratio
		}
	}

	return total
}

func main() {
	// input := TEST_INPUT
	input := utils.GetInputs(2023, 3)

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
