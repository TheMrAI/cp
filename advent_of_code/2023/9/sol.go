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
		input := make([]int, len(numberSequences[i]))
		copy(input, numberSequences[i])
		total += ExtrapolateNextValue(input)
	}

	return total
}

func PartTwo(numberSequences [][]int) int {
	total := 0

	for i := range numberSequences {
		input := make([]int, len(numberSequences[i]))
		copy(input, numberSequences[i])
		for i := 0; i < len(input)/2; i++ {
			input[i], input[len(input)-1-i] = input[len(input)-1-i], input[i]
		}
		total += ExtrapolateNextValue(input)
	}

	return total
}

// I wonder what is the pattern in Go, but I feel like
// destroying the input may not be such a widespread thing.
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
	fmt.Println("Part two")
	fmt.Printf("Sum of backwards extrapolated values is: %v\n", PartTwo(numberSequences))
}
