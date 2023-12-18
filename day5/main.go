package main

import (
	"flag"
	"fmt"
	"math"
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

type Segment struct {
	source Range
	dest   Range
}

type LinearPiecewise struct {
	ranges []Segment
}

func (r *Range) contains(n int64) bool {
	return n >= r.Start && n < r.End
}

func (s *Segment) apply(n int64) (result int64) {
	if !(s.source.contains(n) || s.source.End == n) {
		panic("out of valid range")
	}
	return s.dest.Start + n - s.source.Start
}

func (s *Segment) unapply(x int64) (result int64) {
	return x - s.dest.Start + s.source.Start
}

func (rmap *LinearPiecewise) apply(n int64) (result int64) {
	result = n
	hits := 0
	for _, segment := range rmap.ranges {
		if segment.source.contains(n) {
			hits = 1
			result = segment.apply(n)
		}
	}
	if hits > 1 {
		panic("too many matches")
	}
	return result
}

func sortRanges(ranges []Segment) {
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].source.Start < ranges[j].source.Start
	})
}

func (r1 Range) overlaps(r2 Range) bool {
	//  r1 |-----|
	//  r2    |------|
	overlapStart := max(r1.Start, r2.Start)
	overlapEnd := min(r1.End, r2.End)
	return overlapStart < overlapEnd
}

func (p1 LinearPiecewise) compose(p2 LinearPiecewise) LinearPiecewise {
	ranges := []Segment{}
	for _, pwi := range p1.ranges {
		for _, pwj := range p2.ranges {
			if pwi.dest.overlaps(pwj.source) {
				// in the "j" basis
				overlapStart := max(pwi.dest.Start, pwj.source.Start)
				overlapEnd := min(pwi.dest.End, pwj.source.End)

				// convert to the "i" basis
				x := Range{pwi.unapply(overlapStart), pwi.unapply(overlapEnd)}
				ranges = append(ranges, Segment{x, Range{pwj.apply(overlapStart), pwj.apply(overlapEnd)}})
			}
		}
	}

	span := int64(0)
	for _, seg := range ranges {
		span += seg.source.End - seg.source.Start
	}
	if span != math.MaxInt64 {
		fmt.Println(ranges)
		panic(fmt.Sprintf("function doesn't cover all of int64. %d != %d (diff = %d)", span, math.MaxInt64, math.MaxInt64-span))
	}
	return LinearPiecewise{ranges}
}

func parseBlock(input string) (result LinearPiecewise) {
	ranges := []Segment{}
	for _, line := range strings.Split(strings.Trim(input, "\n"), "\n")[1:] {
		parts := stringsToInts(strings.Fields(line))
		if len(parts) != 3 {
			panic(fmt.Sprintf("invalid line: %s", line))
		}

		source := Range{Start: parts[1], End: parts[1] + parts[2]}
		dest := Range{Start: parts[0], End: parts[0] + parts[2]}
		segment := Segment{source, dest}
		ranges = append(ranges, segment)
	}
	sortRanges(ranges)

	fillers := []Segment{}

	if ranges[0].source.Start > 0 {
		r := Range{int64(0), ranges[0].source.Start}
		fillers = append(fillers, Segment{r, r})
	}

	for i := range ranges[:len(ranges)-1] {
		r := Range{ranges[i].source.End, ranges[i+1].source.Start}
		if r.Start > r.End {
			fmt.Println(r)
			panic("invalid overlap")
		}
		if r.End-r.Start > 0 {
			// insert a segment to represent the multiplication with unity
			fillers = append(fillers, Segment{r, r})
		}
	}
	ranges = append(ranges, fillers...)

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].source.Start < ranges[j].source.Start
	})

	last := Range{ranges[len(ranges)-1].source.End, math.MaxInt64}
	ranges = append(ranges, Segment{last, last})

	result = LinearPiecewise{ranges}
	return result
}

func part1(input string) (result int64) {
	chunks := strings.Split(input, "\n\n")
	seeds := stringsToInts(strings.Fields(strings.Split(chunks[0], ":")[1]))
	curr := seeds
	for _, chunk := range chunks[1:] {
		parsed := parseBlock(chunk)
		for i, val := range curr {
			curr[i] = parsed.apply(val)
		}
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
	inputRanges := []Range{}
	for i := 0; i < len(seeds)-1; i = i + 2 {
		r := Range{seeds[i], seeds[i] + seeds[i+1]}
		inputRanges = append(inputRanges, r)
	}
	unit := Segment{Range{int64(0), math.MaxInt64}, Range{int64(0), math.MaxInt64}}
	mapping := LinearPiecewise{[]Segment{unit}}

	for _, chunk := range chunks[1:] {
		pwfunc := parseBlock(chunk)
		mapping = mapping.compose(pwfunc)
	}

	bestX := inputRanges[0].Start
	bestY := mapping.apply(bestX)
	check := func(n int64) {
		// fmt.Printf("checking %d -> %d (curr best %d)\n", n, mapping.apply(n), bestY)
		if val := mapping.apply(n); val < bestY {
			bestX = n
			bestY = val
		}
	}
	for _, inRange := range inputRanges {
		check(inRange.Start)
	}
	for _, segment := range mapping.ranges {
		for _, inRange := range inputRanges {
			if inRange.contains(segment.source.Start) {
				check(segment.source.Start)
			}
		}
	}
	result = bestY

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
