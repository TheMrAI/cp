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

type Scratchcard struct {
	CardId         int
	WinningNumbers map[int]struct{}
	Numbers        []int
}

func ParseScrachcard(lines []string) []Scratchcard {
	scratchcards := []Scratchcard{}

	cardIdRegexp := regexp.MustCompile("Card +(\\d+)")
	numberRegexp := regexp.MustCompile("\\d+")

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid input received: %v\n", line)
		}
		cardIdText := cardIdRegexp.FindStringSubmatch(parts[0])
		if len(cardIdText) != 2 {
			log.Fatalf("Could not parse game id from: %v\n", parts[0])
		}
		cardId, err := strconv.Atoi(cardIdText[1])
		if err != nil {
			log.Fatalf("Could not convert game id to integer: %v\n", cardId)
		}
		parts = strings.Split(parts[1], "|")
		if len(parts) != 2 {
			log.Fatalf("Could not find | in: %v\n", parts[1])
		}
		winningNumbers := map[int]struct{}{}
		for _, numberText := range numberRegexp.FindAllString(parts[0], -1) {
			number, err := strconv.Atoi(numberText)
			if err != nil {
				log.Fatalf("Failed to convert to integer: %v\n", numberText)
			}
			winningNumbers[number] = struct{}{}
		}
		numbers := []int{}
		for _, numberText := range numberRegexp.FindAllString(parts[1], -1) {
			number, err := strconv.Atoi(numberText)
			if err != nil {
				log.Fatalf("Failed to convert to integer: %v\n", numberText)
			}
			numbers = append(numbers, number)
		}
		scratchcards = append(scratchcards, Scratchcard{cardId, winningNumbers, numbers})
	}

	return scratchcards
}

func PartOne(scratchards []Scratchcard) int {
	totalPoints := 0
	for _, scratchcard := range scratchards {
		matchCount := CountMatches(scratchcard)
		if matchCount > 0 {
			totalPoints += 0b1 << (matchCount - 1)
		}
	}
	return totalPoints
}

func CountMatches(scratchcard Scratchcard) int {
	matchCount := 0
	for _, number := range scratchcard.Numbers {
		_, ok := scratchcard.WinningNumbers[number]
		if ok {
			matchCount++
		}
	}
	return matchCount
}

func PartTwo(scratchcards []Scratchcard) int {
	ticketCounter := make([]int, len(scratchcards))
	for i := range ticketCounter {
		ticketCounter[i] = 1
	}

	for i := len(scratchcards) - 1; i >= 0; i-- {
		matchCount := CountMatches(scratchcards[i])
		for j := i + 1; j < len(scratchcards) && j <= i+matchCount; j++ {
			ticketCounter[i] += ticketCounter[j]
		}
	}

	totalTicketSum := 0
	for _, counter := range ticketCounter {
		totalTicketSum += counter
	}
	return totalTicketSum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	scratchcards := ParseScrachcard(lines)

	fmt.Println("Part one")
	fmt.Printf("Total scratchcard points: %v\n", PartOne(scratchcards))
	fmt.Println("Part two")
	fmt.Printf("Total number of won scratchcards: %v\n", PartTwo(scratchcards))
}
