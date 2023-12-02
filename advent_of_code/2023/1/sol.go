package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func partOne(lines []string) int {
	numberRegexp := regexp.MustCompile("\\d")

	numbers := [][]int{}
	for _, line := range lines {
		numbersInLine := []int{}
		numberStrings := numberRegexp.FindAllString(line, -1)
		for _, numberString := range numberStrings {
			number, err := strconv.Atoi(numberString)
			if err != nil {
				log.Fatalf("Failed to convert: %v\n", numberString)
			}
			numbersInLine = append(numbersInLine, number)
		}
		numbers = append(numbers, numbersInLine)
	}
	return addNumbers(numbers)
}

func partTwo(lines []string) int {
	numberRegexp := regexp.MustCompile("\\d|one|two|three|four|five|six|seven|eight|nine")

	numbers := [][]int{}
	for _, line := range lines {
		numbersInLine := []int{}
		leftMatch := numberRegexp.FindString(line)
		numbersInLine = append(numbersInLine, stringToNumber(leftMatch))
		for i := len(line)-1; i >= 0; i-- {
			rightMatch :=  numberRegexp.FindString(line[i:])
			if len(rightMatch) != 0 {
				numbersInLine = append(numbersInLine, stringToNumber(rightMatch))
				break;
			}
		}

		numbers = append(numbers, numbersInLine)
	}
	return addNumbers(numbers)
}

func stringToNumber(numberString string) int {
	if len(numberString) == 1 {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			log.Fatalf("Failed to convert: %v\n", numberString)
		}
		return number
	}
	switch numberString {
	case "one":
		{
			return 1
		}
	case "two":
		{
			return 2
		}
	case "three":
		{
			return 3
		}
	case "four":
		{
			return 4
		}
	case "five":
		{
			return 5
		}
	case "six":
		{
			return 6
		}
	case "seven":
		{
			return 7
		}
	case "eight":
		{
			return 8
		}
	case "nine":
		{
			return 9
		}
	default:
		{
			return -1
		}
	}
}

func addNumbers(numbers [][]int) int {
	sum := 0
	for _, numbersLine := range numbers {
		sum += (numbersLine[0] * 10) + numbersLine[len(numbersLine)-1]
	}
	return sum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// fmt.Println("Part one")
	// fmt.Printf("Calibration value: %v\n", partOne(lines))
	fmt.Println("Part two")
	fmt.Printf("Calibration value with extended parsing: %v\n", partTwo(lines))
}
