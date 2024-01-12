package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
)

type Pair struct {
	Lhs byte
	Rhs byte
}

func ParseInput(lines []string) (string, map[Pair]byte) {
	originalPolymer := lines[0]

	polymerMappingRules := map[Pair]byte{}
	polymerRuleRegexp := regexp.MustCompile("(\\w)(\\w) -> (\\w)")
	for _, line := range lines[2:] {
		submatches := polymerRuleRegexp.FindStringSubmatch(line)
		if len(submatches) != 4 {
			log.Fatalf("Failed to parse polymer insertion rule from input: '%v'\n", line)
		}
		polymerMappingRules[Pair{submatches[1][0], submatches[2][0]}] = submatches[3][0]
	}

	return originalPolymer, polymerMappingRules
}

func SimulateInsertions(polymer string, insertionRules map[Pair]byte, insertionCount int) int {
	polymerSegments := map[Pair]int{}
	for i := 1; i < len(polymer); i++ {
		polymerSegments[Pair{polymer[i-1], polymer[i]}] += 1
	}

	polymerSegments = InsertPolymers(polymerSegments, insertionRules, insertionCount)

	charCounter := [26]int{}
	for k, v := range polymerSegments {
		charCounter[k.Lhs-'A'] += v
		charCounter[k.Rhs-'A'] += v
	}
	// correct for counting double
	first := polymer[0]
	last := polymer[len(polymer)-1]
	for i, count := range charCounter {
		charCounter[i] = count / 2
		c := byte(i + 'A')
		if c == first || c == last {
			charCounter[i]++
		}
	}

	min := math.MaxInt
	max := math.MinInt
	for _, count := range charCounter {
		if count > 0 && min > count {
			min = count
		}
		if max < count {
			max = count
		}
	}
	fmt.Println(min, max)
	return max - min
}

func InsertPolymers(polymerSegments map[Pair]int, insertionRules map[Pair]byte, rounds int) map[Pair]int {
	workingPolymerSegments := map[Pair]int{}
	for k, v := range polymerSegments {
		workingPolymerSegments[k] = v
	}

	for insertionRound := 0; insertionRound < rounds; insertionRound++ {
		nextPolymerSegments := map[Pair]int{}
		for k, v := range workingPolymerSegments {
			nextPolymerSegments[Pair{k.Lhs, insertionRules[k]}] += v
			nextPolymerSegments[Pair{insertionRules[k], k.Rhs}] += v
		}
		workingPolymerSegments = nextPolymerSegments
	}
	return workingPolymerSegments
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	polymer, rules := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("The delta between the most and least common elements after 10 insertions: %v\n", SimulateInsertions(polymer, rules, 10))
	fmt.Println("Part two")
	fmt.Printf("The delta between the most and least common elements after 40 insertions: %v\n", SimulateInsertions(polymer, rules, 40))
}
