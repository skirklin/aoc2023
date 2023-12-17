package main

import (
	"flag"
	"fmt"
	"sort"
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

var TEST_INPUT = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

func stringsToInts(input []string) (out []int64) {
	out = make([]int64, len(input))
	for i, numstr := range input {
		num, ok := strconv.ParseInt(numstr, 10, 64)
		check(ok)
		out[i] = int64(num)
	}
	return out
}

type Range struct {
	Start int64
	End   int64
}

type RangePair struct {
	source Range
	dest   int64
}

type RangeMap struct {
	ranges []RangePair
}

func (r *Range) contains(n int64) bool {
	return n >= r.Start && n <= r.End
}

func (r1 *Range) overlaps(r2 Range) bool {
	return (r2.Start <= r1.Start && r1.Start <= r2.End) ||
		(r2.Start <= r1.End && r1.End <= r2.End)
}

func (rmap *RangeMap) apply(n int64) (result int64) {
	result = n
	hits := 0
	for _, pair := range rmap.ranges {
		if pair.source.contains(n) {
			hits = 1
			result = pair.dest + n - pair.source.Start
		}
	}
	if hits > 1 {
		panic("too many matches")
	}
	return result
}

func (rmap *RangeMap) applyOn(r1 Range) []Range {
	fmt.Println("applying map:", rmap.ranges)
	// all inputs in r1 must be covered by the return, defaulting to their existing value
	for _, r2 := range rmap.ranges {
		if r1.overlaps(r2.source) {
			fmt.Println("r1:", r1, "r2", r2)
			newStart := max(r1.Start, r2.source.Start)
			newEnd := min(r1.End, r2.source.End)
			shift := r2.dest - r2.source.Start
			result := []Range{}

			if newStart > r1.Start {
				result = append(result, Range{r1.Start, newStart})
			}
			result = append(result, Range{Start: newStart + shift, End: newEnd + shift})
			if newEnd < r1.End {
				result = append(result, Range{newEnd, r1.End})
			}
			fmt.Println("mapped newStart:", newStart, "newend:", newEnd, "shift:", shift)
			fmt.Println("after mapping:", result)
			sort.Slice(result, func(i, j int) bool {
				return result[i].Start < result[j].Start
			})
			return result
		}
	}
	return []Range{r1}
}

func parseBlock(input string) (result RangeMap) {
	ranges := []RangePair{}
	for _, line := range strings.Split(strings.Trim(input, "\n"), "\n")[1:] {
		parts := stringsToInts(strings.Fields(line))
		fmt.Println(parts)
		if len(parts) != 3 {
			panic(fmt.Sprintf("invalid line: %s", line))
		}

		r := Range{Start: parts[1], End: parts[1] + parts[2] - 1}
		rangePair := RangePair{r, parts[0]}
		ranges = append(ranges, rangePair)
	}
	result = RangeMap{ranges}
	return result
}

func part1(input string) (result int64) {
	chunks := strings.Split(input, "\n\n")
	seeds := stringsToInts(strings.Fields(strings.Split(chunks[0], ":")[1]))
	curr := seeds
	fmt.Println(curr)
	for _, chunk := range chunks[1:] {
		parsed := parseBlock(chunk)
		for i, val := range curr {
			curr[i] = parsed.apply(val)
		}
		fmt.Println(curr)
	}
	result = curr[0]
	for _, val := range curr {
		if val < result {
			result = val
		}
	}
	return result
}

func part2(input string) (result int64) {
	chunks := strings.Split(input, "\n\n")
	// make ranges instead of a single array
	seeds := stringsToInts(strings.Fields(strings.Split(chunks[0], ":")[1]))
	fmt.Println(seeds)
	curr := []Range{}
	initialSeeds := int64(0)
	for i := 0; i < len(seeds)-1; i = i + 2 {
		r := Range{seeds[i], seeds[i] + seeds[i+1]}
		initialSeeds += r.End - r.Start
		curr = append(curr, r)
	}
	// fmt.Println("initial seeds:", initialSeeds)

	fmt.Println(curr)
	for _, chunk := range chunks[1:] {
		tmp := []Range{}
		parsed := parseBlock(chunk)
		for _, val := range curr {
			tmp = append(tmp, parsed.applyOn(val)...)
		}
		curr = tmp
		sort.Slice(curr, func(i, j int) bool {
			return curr[i].Start < curr[j].Start
		})
		numSeeds := int64(0)
		for _, val := range curr {
			numSeeds += val.End - val.Start
		}
		// fmt.Println("number of seeds:", numSeeds)
		fmt.Println(curr)
	}
	span := int64(0)
	result = curr[0].Start
	for _, val := range curr {
		span += val.End - val.Start
		if val.Start < result {
			result = val.Start
		}
	}
	// fmt.Println("number of seeds:", span)
	return result
}

func main() {
	flag.Parse()

	var input string
	if *example {
		input = TEST_INPUT
	} else {
		input = utils.GetInputs(2023, 5)
	}

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
