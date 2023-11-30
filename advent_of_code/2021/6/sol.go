package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func partOne(numbers []int, simulationRounds int) uint64 {
	lifecycle := make([]uint64, 9)
	for _, number := range numbers {
		lifecycle[number]++
	}

	simulationDay := 1
	for ;simulationDay <= simulationRounds; simulationDay++ {
		updatedLifecycle := make([]uint64, 9)
		for age, count := range lifecycle {
			if age == 0 {
				updatedLifecycle[8] += count
				updatedLifecycle[6] += count
				continue
			}
			updatedLifecycle[age-1] += count
		}
		lifecycle = updatedLifecycle
	}
	sum := uint64(0)
	for _, count := range lifecycle {
		sum += count
	}
	return sum
}

func main() {
	numberRegexp := regexp.MustCompile("\\d+")

	scanner := bufio.NewScanner(os.Stdin)
	numbers := []int{}
	for scanner.Scan() {
		line := scanner.Text()

		numberStrings := numberRegexp.FindAllString(line, -1)
		for _, numberString := range numberStrings {
			number, err := strconv.Atoi(numberString)
			if err != nil {
				log.Fatalf("Failed to convert: %v\n", numberString)
			}
			numbers = append(numbers, number)
		}
	}

	fmt.Println("Part one")
	fmt.Printf("Lanternfish count after 80 days: %v\n", partOne(numbers, 80))
	fmt.Println("Part two")
	fmt.Printf("Lanternfish count after 256 days: %v\n", partOne(numbers, 256))
}
