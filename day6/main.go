package main

import (
	"flag"
	"fmt"
	"math"
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

var TEST_INPUT = `Time:      7  15   30
Distance:  9  40  200
`

/*
Formula:
d = x * t
t = (N - x)
d = x * (N - x) = x*N - x^2

where x is the number of ms holding the button (and therefor also the speed) d
is distance and N is the total time for the race. To maximize d, take the
derivative and solve for 0

dd/dx = N - 2*x = 0
x = N/2

The problem statement is to find the range of x such that d > D, so, no longer just looking to
maximize, we instead want:

D < x*N - x^2
0 < x*N - x^2 - D

Solving for x here will have two solutions:
-b +/- (b^2 - 4ac)^0.5 / 2a

Using our values of:

a = -1
b = N
c = -D

Plugging in:
(-N +- (N^2 - 4*-1*-D)^0.5) / (2*-1)

And simplifying:
N/2 - (N^2 - 4D)^0.5 / 2
N/2 + (N^2 - 4D)^0.5 / 2

I.e. the middle amount of time, with a window that is larger proportional to the
total amount of time available and threshold you need to exceed.


*/

type Race struct {
	time     int
	distance int
}

func parseInputs(input string) []Race {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	times := strings.Fields(lines[0])[1:]
	distances := strings.Fields(lines[1])[1:]
	races := make([]Race, len(times))
	for i, tstr := range times {
		dstr := distances[i]
		t, ok := strconv.Atoi(tstr)
		check(ok)
		d, ok := strconv.Atoi(dstr)
		check(ok)
		races[i] = Race{time: t, distance: d}
	}
	return races
}

func computeWindow(r Race) int {
	// the formula we need is:
	// low end = N/2 - (N^2 - 4D)^0.5 / 2
	// high end = N/2 + (N^2 - 4D)^0.5 / 2
	// but rounding towards the middle
	N := float64(r.time)
	// add epsilon to D because it is a strict inequality
	D := float64(r.distance) + 1e-8
	lowEnd := N/2. - math.Pow(N*N-4*D, 0.5)/2.
	highEnd := N/2. + math.Pow(N*N-4*D, 0.5)/2.
	window := int(math.Floor(highEnd)-math.Ceil(lowEnd)) + 1
	return window
}

func part1(input string) (result int) {
	races := parseInputs(input)
	result = 1
	for _, r := range races {
		window := computeWindow(r)
		result *= window
	}
	return result
}

func parseInputs2(input string) Race {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	tstr := strings.ReplaceAll(strings.SplitN(lines[0], ":", 2)[1], " ", "")
	dstr := strings.ReplaceAll(strings.SplitN(lines[1], ":", 2)[1], " ", "")
	t, ok := strconv.Atoi(tstr)
	check(ok)
	d, ok := strconv.Atoi(dstr)
	check(ok)
	return Race{distance: d, time: t}
}

func part2(input string) (result int) {
	race := parseInputs2(input)
	window := computeWindow(race)
	return window
}

func main() {
	flag.Parse()

	var input string
	if *example {
		input = TEST_INPUT
	} else {
		input = utils.GetInputs(2023, 6)
	}

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
