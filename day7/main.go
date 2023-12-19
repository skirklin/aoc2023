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

var TEST_INPUT = `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

type Card rune

type Hand struct {
	raw      string
	cards    []Card
	sortKey  string
	handType HandType
	bid      int
}

type HandType string

type Count struct {
	card Card
	n    int
}

type RuleSet interface {
	getHandType(count []Count) HandType
	cardOrder() []Card
}

type BasicRules struct {
}

func (b BasicRules) cardOrder() []Card {
	return []Card("23456789TJQKA")
}

func (b BasicRules) getHandType(counts []Count) HandType {
	return basicHandType(counts)
}

func basicHandType(counts []Count) HandType {
	switch {
	case counts[0].n == 5:
		return "five"
	case counts[0].n == 4:
		return "four"
	case counts[0].n == 3 && counts[1].n == 2:
		return "full"
	case counts[0].n == 3:
		return "three"
	case counts[0].n == 2 && counts[1].n == 2:
		return "two"
	case counts[0].n == 2:
		return "one"
	case counts[0].n == 1:
		return "high"
	default:
		panic(fmt.Sprintf("what is %s?", counts))
	}
}

type AdvancedRules struct {
}

func (r AdvancedRules) cardOrder() []Card {
	return []Card("J23456789TQKA")
}

func (r AdvancedRules) getHandType(counts []Count) HandType {
	// I think just piling jokers into the top count category is good enough to
	// get the best possible hand. Specifically,
	// if I have one pair + one joker, it is better to make three of a kind than two pair
	// if I have three of a kind + one joker it is better to make four of a kind than a full house
	// if I have one pair + two jokers it is better to make four of a kind than a full house
	// ....
	modifiedCounts := []Count{}
	jokers := 0
	for _, count := range counts {
		if count.card == Card('J') {
			jokers = count.n
		} else {
			modifiedCounts = append(modifiedCounts, count)
		}
	}
	if len(modifiedCounts) == 0 {
		// means it was all jokers
		return basicHandType(counts)
	} else {
		modifiedCounts[0].n += jokers
		return basicHandType(modifiedCounts)
	}
}

func (c Count) String() string {
	return fmt.Sprintf("%s -> %d", string(c.card), c.n)
}
func (h Hand) String() string {
	return fmt.Sprintf("%s %d (%s)", h.raw, h.bid, h.handType)
}

func parseHand(line string, ruleset RuleSet) Hand {
	parts := strings.Fields(line)
	bid, ok := strconv.Atoi(parts[1])
	check(ok)
	cards := []Card{}
	for _, card := range parts[0] {
		cards = append(cards, Card(card))
	}

	countMap := map[Card]int{}
	// returns sorted counts
	for _, c := range cards {
		countMap[c] += 1
	}
	counts := []Count{}
	for k, v := range countMap {
		counts = append(counts, Count{k, v})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].n > counts[j].n
	})

	handType := ruleset.getHandType(counts)

	cardOrder := ruleset.cardOrder()
	// sorting is a lexicographic sorting of the tuples:
	// (handType, card1, card2, card3, card4, card5)
	runes := make([]rune, 6)
	for i, t := range handOrder {
		if t == handType {
			runes[0] = 'a' + rune(i)
			break
		}
	}
	for i, r := range cards {
		for j, t := range cardOrder {
			if t == r {
				runes[1+i] = 'a' + rune(j)
				break
			}
		}
	}

	sortKey := string(runes)
	hand := Hand{parts[0], cards, sortKey, handType, bid}
	// fmt.Printf("Parsed %s to %s\n", line, hand)
	return hand
}

func parseHands(input string, ruleset RuleSet) []Hand {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	hands := make([]Hand, len(lines))
	for i, line := range lines {
		hands[i] = parseHand(line, ruleset)
	}

	// to get this effect, I'll convert each hand to a 6 rune string containing the scores, then sort those.
	sort.Slice(hands, func(i, j int) bool {
		return hands[i].sortKey < hands[j].sortKey
	})

	return hands
}

var handOrder = []HandType{
	HandType("high"),
	HandType("one"),
	HandType("two"),
	HandType("three"),
	HandType("full"),
	HandType("four"),
	HandType("five"),
}

func part1(input string) (result int) {
	ruleset := BasicRules{}
	hands := parseHands(input, ruleset)

	result = 0
	fmt.Println("sorted order:")
	for i, hand := range hands {
		fmt.Println(" ", hand)
		result += (i + 1) * hand.bid
	}
	return result
}

func part2(input string) (result int) {
	ruleset := AdvancedRules{}
	hands := parseHands(input, ruleset)

	result = 0
	fmt.Println("sorted order:")
	for i, hand := range hands {
		fmt.Println(" ", hand)
		result += (i + 1) * hand.bid
	}
	return result
}

func main() {
	flag.Parse()

	var input string
	if *example {
		input = TEST_INPUT
	} else {
		input = utils.GetInputs(2023, 7)
	}

	ans1 := part1(input)
	fmt.Println("Part 1 answer:", ans1)

	ans2 := part2(input)
	fmt.Println("Part 2 answer:", ans2)
}
