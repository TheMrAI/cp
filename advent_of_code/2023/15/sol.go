package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func PartOne(inputs []string) int {
	sumOfHashes := 0

	for _, input := range inputs {
		sumOfHashes += ASCIIHashString(input, 0)
	}
	return sumOfHashes
}

func ASCIIHashStrings(inputs []string) int {
	hash := 0

	for _, input := range inputs {
		hash = ASCIIHashString(input, hash)
	}
	return hash
}

func ASCIIHashString(input string, hash int) int {
	for i := range input {
		hash += int(input[i])
		hash *= 17
		hash = hash % 256
	}
	return hash
}

type Lens struct {
	Label       string
	FocalLength int
}

func PartTwo(inputs []string) int {
	boxes := make([][]Lens, 256)

	for _, input := range inputs {
		label := ""
		add := true
		inputParts := strings.Split(input, "=")
		if len(inputParts) == 1 {
			inputParts = strings.Split(input, "-")
			add = false
		}
		label = inputParts[0]

		boxId := ASCIIHashString(label, 0)
		lensIndex := -1
		for i, lens := range boxes[boxId] {
			if lens.Label == label {
				lensIndex = i
				break
			}
		}

		if add {
			focalLength, err := strconv.Atoi(inputParts[1])
			if err != nil {
				log.Fatalf("Integer conversion failed for input: '%v'\n", inputParts[1])
			}
			if lensIndex != -1 {
				boxes[boxId][lensIndex] = Lens{label, focalLength}
			} else {
				boxes[boxId] = append(boxes[boxId], Lens{label, focalLength})
			}
			continue
		}
		if lensIndex != -1 {
			boxes[boxId] = append(boxes[boxId][:lensIndex], boxes[boxId][lensIndex+1:]...)
		}
	}

	totalFocusingPower := 0
	for i, box := range boxes {
		for slot, lens := range box {
			totalFocusingPower += (i + 1) * (slot + 1) * lens.FocalLength
		}
	}
	return totalFocusingPower
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	initializationSequence := strings.Split(lines[0], ",")

	fmt.Println("Part one")
	fmt.Printf("Sum of all hash results: %v\n", PartOne(initializationSequence))
	fmt.Println("Part two")
	fmt.Printf("Total focusing power of final lens configuration: %v\n", PartTwo(initializationSequence))
}
