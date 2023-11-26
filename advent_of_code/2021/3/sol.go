package main

import (
	"bufio"
	"fmt"
	"os"
)

func binaryStringToUint64(binaryString string) uint64 {
	result := uint64(0)
	for _, r := range binaryString {
		result = result << 1
		if r == '1' {
			result = result | 0b1
		}
		// well we ignore all invalid inputs
	}
	return result
}

func partOne(oneCounter []uint64, maxBitsSet int, numberCount int) (uint64, uint64) {
	gammaRate := uint64(0)
	for _, val := range oneCounter {
		gammaRate = gammaRate << 1
		if val > uint64(numberCount)/2 {
			gammaRate = gammaRate | 1
		}
	}
	epsilonRate := (^gammaRate) & ((0b1 << maxBitsSet) - 1)
	return gammaRate, epsilonRate
}

func countSetBits(numbers []uint64, maxBitsSet int) []uint64 {
	oneCounter := make([]uint64, 64)
	for _, number := range numbers {
		for i := 0; i < maxBitsSet; i++ {
			if (number & (0b1 << i)) != 0 {
				oneCounter[63-i] += 1
			}
		}
	}
	return oneCounter
}

func partTwo(numbers []uint64, maxBitsSet int) (uint64, uint64) {
	forOxygenRating := make([]uint64, len(numbers))
	copy(forOxygenRating, numbers)
	for i := maxBitsSet - 1; i >= 0 && len(forOxygenRating) > 1; i-- {
		mostCommon := mostCommonIsSet(forOxygenRating, i)
		forOxygenRating = separate(forOxygenRating, i, mostCommon)
	}
	oxygenRating := forOxygenRating[0]

	forScrubberRating := make([]uint64, len(numbers))
	copy(forScrubberRating, numbers)
	for i := maxBitsSet - 1; i >= 0 && len(forScrubberRating) > 1; i-- {
		leastCommon := !mostCommonIsSet(forScrubberRating, i)
		forScrubberRating = separate(forScrubberRating, i, leastCommon)
	}
	scrubberRating := forScrubberRating[0]

	return oxygenRating, scrubberRating
}

func separate(numbers []uint64, onBit int, filterSetBits bool) []uint64 {
	preferredUntil := 0
	for i := range numbers {
		bitIsSet := (numbers[i] & (0b1 << onBit)) != 0
		if bitIsSet == filterSetBits {
			numbers[preferredUntil], numbers[i] = numbers[i], numbers[preferredUntil]
			preferredUntil++
		}
	}
	return numbers[:preferredUntil]
}

func mostCommonIsSet(numbers []uint64, onBit int) bool {
	setCounter := 0
	for _, num := range numbers {
		if (num & (0b1 << onBit)) != 0 {
			setCounter++
		}
	}

	if setCounter >= (len(numbers) - setCounter) {
		return true
	}
	return false
}

func main() {
	maxBitsSet := int(0)
	numbers := []uint64{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		lineLength := len(line)
		if maxBitsSet < lineLength {
			maxBitsSet = lineLength
		}
		numbers = append(numbers, binaryStringToUint64(line))
	}

	setBitsCount := countSetBits(numbers, maxBitsSet)

	gammaRate, epsilonRate := partOne(setBitsCount, maxBitsSet, len(numbers))
	fmt.Printf("Gamma rate: %5b, Epsilon rate: %5b, power consumption: %v\n", gammaRate, epsilonRate, gammaRate*epsilonRate)

	oxygenRating, scrubberRating := partTwo(numbers, maxBitsSet)
	fmt.Printf("Oxygen rating: %v, CO2 scrubber rating: %v, Life support rating: %v\n", oxygenRating, scrubberRating, oxygenRating*scrubberRating)
}
