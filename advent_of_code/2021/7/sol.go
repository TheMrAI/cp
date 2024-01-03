package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(lines []string) []int {
	numbers := []int{}

	numberRegexp := regexp.MustCompile("\\d+")

	for _, line := range lines {
		numbersTexts := numberRegexp.FindAllString(line, -1)
		for _, text := range numbersTexts {
			number, err := strconv.Atoi(text)
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", text)
			}
			numbers = append(numbers, number)
		}
	}

	return numbers
}

func PartOne(numbers []int) int {
	minNum, maxNum := MinMax(numbers)

	minimalFuelUsage := math.MaxInt
	for candidate := minNum; candidate <= maxNum; candidate++ {
		fuelUsage := 0
		for _, position := range numbers {
			fuelUsage += Abs(position - candidate)
		}
		if fuelUsage < minimalFuelUsage {
			minimalFuelUsage = fuelUsage
		}
	}
	return minimalFuelUsage
}

func PartTwo(numbers []int) int {
	minNum, maxNum := MinMax(numbers)

	minimalFuelUsage := math.MaxInt
	for candidate := minNum; candidate <= maxNum; candidate++ {
		fuelUsage := 0
		for _, position := range numbers {
			// we need to add 1 so that we are summing on the range [1, n]
			// not [0, n-1]
			fuel := Abs(position-candidate) + 1
			fuelUsage += (fuel * (fuel - 1)) / 2
		}
		if fuelUsage < minimalFuelUsage {
			minimalFuelUsage = fuelUsage
		}
	}
	return minimalFuelUsage
}

func MinMax(numbers []int) (int, int) {
	minNum := math.MaxInt
	maxNum := math.MinInt
	for i := range numbers {
		if numbers[i] < minNum {
			minNum = numbers[i]
		}
		if maxNum < numbers[i] {
			maxNum = numbers[i]
		}
	}
	return minNum, maxNum
}

func Abs(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	numbers := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Total fuel necessary to align all crabs horizontally: %v\n", PartOne(numbers))
	fmt.Println("Part two")
	fmt.Printf("Total fuel necessary to align all crabs horizontally, with crab engines: %v\n", PartTwo(numbers))
}
