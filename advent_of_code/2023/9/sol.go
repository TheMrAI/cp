package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(lines []string) [][]int {
	numberSequences := [][]int{}

	numberRegexp := regexp.MustCompile("-?\\d+")
	for _, line := range lines {
		numbersText := numberRegexp.FindAllString(line, -1)
		numbers := []int{}
		for _, numberText := range numbersText {
			number, err := strconv.Atoi(numberText)
			if err != nil {
				log.Fatalf("Number conversion failed on input: %v\n", numberText)
			}
			numbers = append(numbers, number)
		}
		numberSequences = append(numberSequences, numbers)
	}

	return numberSequences
}

func PartOne(numberSequences [][]int) int {
	total := 0

	for i := range numberSequences {
		total += ExtrapolateNextValue(numberSequences[i])
	}

	return total
}

func ExtrapolateNextValue(numbers []int) int {
	length := len(numbers)
	for {
		allZeroes := true
		for i := 1; i < length; i++ {
			numbers[i-1] = numbers[i] - numbers[i-1]
			if numbers[i-1] != 0 {
				allZeroes = false
			}
		}
		length--
		if allZeroes {
			break
		}
	}

	// edge case that doesn't even appear in the inputs
	// but I was getting paranoid
	if length == 0 && numbers[0] != 0 {
		return -1
	}
	prediction := 0
	for ; length < len(numbers); length++ {
		prediction += numbers[length]
	}
	return prediction
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	numberSequences := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Sum of extrapolated values is: %v\n", PartOne(numberSequences))
}
