package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func partOne(measurements []int) int {
	increaseCount := 0
	for index := range measurements {
		if index == 0 {
			continue
		}
		if measurements[index-1] < measurements[index] {
			increaseCount++
		}
	}
	return increaseCount
}

func partTwo(measurements []int) int {
	fourSum := 0
	increaseCount := 0
	for index := range measurements {
		fourSum += measurements[index]
		if index < 3 {
			continue
		}
		previousWindow := fourSum - measurements[index]
		currentWindow := fourSum - measurements[index-3]
		if previousWindow < currentWindow {
			increaseCount++
		}
		fourSum = currentWindow
	}
	return increaseCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	measurements := []int{}
	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Failed conversion")
		}
		measurements = append(measurements, value)
	}
	fmt.Printf("Depth increases: %v\n", partOne(measurements))
	fmt.Printf("Three-measurement sliding window increases: %v\n", partTwo(measurements))
}
