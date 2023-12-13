package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Record struct {
	Pattern       string
	DamagedGroups []int
}

func ParseInput(lines []string) []Record {
	data := []Record{}

	numberRegexp := regexp.MustCompile("\\d+")

	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			log.Fatalf("Input doesn't contain exactly one space: '%v'\n", line)
		}

		numberTexts := numberRegexp.FindAllString(parts[1], -1)
		if len(numberTexts) < 1 {
			log.Fatalf("Missing damaged groupings information in input: '%v'\n", line)
		}
		numbers := []int{}
		for _, numberText := range numberTexts {
			number, err := strconv.Atoi(numberText)
			if err != nil {
				log.Fatalf("Failed to convert string to number: '%v'\n", numberText)
			}
			numbers = append(numbers, number)
		}

		data = append(data, Record{parts[0], numbers})
	}

	return data
}

func PartOne(records []Record) int {
	return SumAllArrangements(records)
}

func PartTwo(records []Record, multiplier int) int {
	if multiplier <= 0 {
		return 0
	}
	multipliedRecords := []Record{}
	for _, record := range records {
		multipliedPattern := record.Pattern + strings.Repeat("?"+record.Pattern, multiplier-1)
		multipliedDamageGroups := make([]int, len(record.DamagedGroups)*multiplier)
		for i := 0; i < multiplier; i++ {
			copy(multipliedDamageGroups[i*len(record.DamagedGroups):], record.DamagedGroups)
		}
		multipliedRecords = append(multipliedRecords, Record{multipliedPattern, multipliedDamageGroups})
	}
	return SumAllArrangements(multipliedRecords)
}

func SumAllArrangements(records []Record) int {
	sum := 0
	for _, record := range records {
		memo := map[Pair]int{}
		sum += CountPossibleArrangements(record.Pattern, record.DamagedGroups, memo)
	}
	return sum
}

type Pair struct {
	RemainingPatternLength int
	GroupIdFromBack        int
}

func CountPossibleArrangements(pattern string, damagedGroups []int, memo map[Pair]int) int {
	memoKey := Pair{len(pattern), len(damagedGroups)}
	value, seen := memo[memoKey]
	if seen {
		return value
	}
	if len(damagedGroups) == 0 {
		numberOfHashmarks := 0
		for _, r := range pattern {
			if r == '#' {
				numberOfHashmarks++
			}
		}
		if numberOfHashmarks != 0 {
			return 0
		}
		return 1
	}

	// fmt.Printf("%v\n", pattern)
	windowSize := damagedGroups[0]
	if len(pattern) < windowSize {
		return 0
	}
	goodCountInWindow := 0
	i := 0
	for ; i < windowSize; i++ {
		if pattern[i] == '.' {
			goodCountInWindow++
		}
	}

	totalArrangementCount := 0
	for ; i <= len(pattern); i++ {
		nextIsNotHash := true
		if i < len(pattern) {
			nextIsNotHash = pattern[i] != '#'
		}
		if goodCountInWindow == 0 && nextIsNotHash {
			nextPattern := ""
			if i < len(pattern)-1 {
				nextPattern = pattern[i+1:]
			}
			totalArrangementCount += CountPossibleArrangements(nextPattern, damagedGroups[1:], memo)
		}
		// # would get out of the window => we matched everything we could
		if pattern[i-windowSize] == '#' {
			break
		}
		if pattern[i-windowSize] == '.' {
			goodCountInWindow--
		}
		if i < len(pattern) && pattern[i] == '.' {
			goodCountInWindow++
		}
	}
	memo[memoKey] = totalArrangementCount
	return totalArrangementCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Sum of possible arrangements: %v\n", PartOne(data))
	fmt.Println("Part two")
	fmt.Printf("Sum of possible arrangements after repeating inputs 5 time(s): %v\n", PartTwo(data, 5))
}
