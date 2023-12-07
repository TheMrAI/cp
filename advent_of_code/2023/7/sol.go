package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func ClassifyHand(hand string) HandType {
	counter := map[rune]int{}
	for _, r := range hand {
		counter[r] += 1
	}
	return ClassifyMap(counter)	
}

func ClassifyMap(counter map[rune]int) HandType {
	oneCount := 0
	for _, v := range counter {
		if v == 1 {
			oneCount++
		}
	}
	if len(counter) == 1 {
		return FiveOfAKind
	} else if len(counter) == 2 {
		if oneCount == 1 {
			return FourOfAKind
		}
		return FullHouse
	} else if len(counter) == 3 {
		if oneCount == 2 {
			return ThreeOfAKind
		}
		return TwoPair
	} else if len(counter) == 4 {
		return OnePair
	}
	return HighCard
}

func ClassifyHandWithJoker(hand string) HandType {
	counter := map[rune]int{}
	for _, r := range hand {
		counter[r] += 1
	}
	if (counter['J'] > 0) {
		mostCommonRune := 'A'
		mostCommonCount := 0
		for r, count := range counter {
			if r == 'J' {
				continue
			}
			if count > mostCommonCount {
				mostCommonRune = r
				mostCommonCount = count
			}
		}
		counter[mostCommonRune] += counter['J']
		delete(counter, 'J')
	}
	return ClassifyMap(counter)
}

type Hand struct {
	Cards string
	Class HandType
}

type HandBid struct {
	Hand Hand
	Bid  int
}

func CardValue(r rune) (int, bool) {
	conversion := map[rune]int{
		'A': 14,
		'K': 13,
		'Q': 12,
		'J': 11,
		'T': 10,
		'9': 9,
		'8': 8,
		'7': 7,
		'6': 6,
		'5': 5,
		'4': 4,
		'3': 3,
		'2': 2,
	}
	val, ok := conversion[r]
	return val, ok
}

func CardValueWithJoker(r rune) (int, bool) {
	conversion := map[rune]int{
		'A': 13,
		'K': 12,
		'Q': 11,
		'T': 10,
		'9': 9,
		'8': 8,
		'7': 7,
		'6': 6,
		'5': 5,
		'4': 4,
		'3': 3,
		'2': 2,
		'J': 1,
	}
	val, ok := conversion[r]
	return val, ok
}

func ParseInput(lines []string) []HandBid {
	handBids := []HandBid{}

	handRegexp := regexp.MustCompile("(\\S{5}) (\\d+)")

	for _, line := range lines {
		matches := handRegexp.FindAllStringSubmatch(line, -1)
		if len(matches) != 1 && len(matches[0]) != 3 {
			log.Fatalf("Unexpected input received: %v\n", line)
		}
		number, err := strconv.Atoi(matches[0][2])
		if err != nil {
			log.Fatalf("Number parsing failed: %v\n", matches[0][2])
		}
		handBids = append(handBids, HandBid{Hand{matches[0][1], ClassifyHand(matches[0][1])}, number})
	}

	return handBids
}

func PartOne(handBids []HandBid) int {
	sort.Slice(handBids, func(i, j int) bool {
		if handBids[i].Hand.Class == handBids[j].Hand.Class {
			for k := range handBids[i].Hand.Cards {
				lhs, _ := CardValue(rune(handBids[i].Hand.Cards[k]))
				rhs, _ := CardValue(rune(handBids[j].Hand.Cards[k]))
				if lhs == rhs {
					continue
				}
				return lhs < rhs
			}
		}
		return handBids[i].Hand.Class < handBids[j].Hand.Class
	})

	totalWinnings := 0
	for i, handBid := range handBids {
		rank := i + 1
		totalWinnings += rank * handBid.Bid
	}
	return totalWinnings
}

func PartTwo(handBids []HandBid) int {
	for i := range handBids {
		handBids[i].Hand.Class = ClassifyHandWithJoker(handBids[i].Hand.Cards)
	}

	sort.Slice(handBids, func(i, j int) bool {
		if handBids[i].Hand.Class == handBids[j].Hand.Class {
			for k := range handBids[i].Hand.Cards {
				lhs, _ := CardValueWithJoker(rune(handBids[i].Hand.Cards[k]))
				rhs, _ := CardValueWithJoker(rune(handBids[j].Hand.Cards[k]))
				if lhs == rhs {
					continue
				}
				return lhs < rhs
			}
		}
		return handBids[i].Hand.Class < handBids[j].Hand.Class
	})

	totalWinnings := 0
	for i, handBid := range handBids {
		rank := i + 1
		totalWinnings += rank * handBid.Bid
	}
	return totalWinnings
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	handBids := ParseInput(lines)
	
	fmt.Println("Part one")
	fmt.Printf("Total winnings: %v\n", PartOne(handBids))
	fmt.Println("Part two")
	fmt.Printf("Total winnings: %v\n", PartTwo(handBids))
}
