package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type CipherLine struct {
	SignalPatterns []string
	OutputSignals  []string
}

func ParseInput(lines []string) []CipherLine {
	cipherLines := []CipherLine{}

	cipher := regexp.MustCompile("\\w+")

	for _, line := range lines {
		cipherLine := CipherLine{}

		sections := strings.Split(line, "|")
		if len(sections) != 2 {
			log.Fatalf("Input expected to contain |, instead: '%v'\n", line)
		}
		cipherLine.SignalPatterns = cipher.FindAllString(sections[0], -1)
		if len(cipherLine.SignalPatterns) == 0 {
			log.Fatalf("Expected to find signal patterns before the |, got: '%v'\n", sections[0])
		}
		cipherLine.OutputSignals = cipher.FindAllString(sections[1], -1)
		if len(cipherLine.OutputSignals) == 0 {
			log.Fatalf("Expected to find output patterns after the |, got: '%v'\n", sections[1])
		}
		cipherLines = append(cipherLines, cipherLine)
	}

	return cipherLines
}

func PartOne(cipherLines []CipherLine) int {
	digitOfInterestCount := 0 // 1, 4, 7, 8
	for _, cipherLines := range cipherLines {
		for j := range cipherLines.OutputSignals {
			length := len(cipherLines.OutputSignals[j])
			if length == 2 || length == 3 || length == 4 || length == 7 {
				digitOfInterestCount++
			}
		}
	}
	return digitOfInterestCount
}

func PartTwo(cipherLines []CipherLine) int {
	sum := 0
	for i := range cipherLines {
		segmentMapping := DecodeCipher(cipherLines[i].SignalPatterns)
		value := 0
		for j := range cipherLines[i].OutputSignals {
			value *= 10
			value += DecodeValue(segmentMapping, cipherLines[i].OutputSignals[j])
		}
		sum += value
	}
	return sum
}

func DecodeValue(segmentMapping []byte, outputSignal string) int {
	// we invert the mask
	segmentMask := uint(0)
	for i := range outputSignal {
		for j := range segmentMapping {
			if outputSignal[i] == segmentMapping[j] {
				segmentMask |= 1 << j
				break
			}
		}
	}
	return SegmentMaskToNumber(segmentMask)
}

func DecodeCipher(signalPatterns []string) []byte {
	// group the codes by their lengths
	codeForOne := map[byte]int{}
	codeForFour := map[byte]int{}
	codeForSeven := map[byte]int{}
	codeForEight := map[byte]int{}
	codesForTwoThreeFive := map[byte]int{}
	codesForZeroSixNine := map[byte]int{}
	for _, signal := range signalPatterns {
		length := len(signal)
		if length == 2 {
			codeForOne = ToByteMap(signal)
		}
		if length == 3 {
			codeForSeven = ToByteMap(signal)
		}
		if length == 4 {
			codeForFour = ToByteMap(signal)
		}
		if length == 7 {
			codeForEight = ToByteMap(signal)
		}
		if length == 5 {
			codesForTwoThreeFive = MapUnion(codesForTwoThreeFive, ToByteMap(signal))
		}
		if length == 6 {
			codesForZeroSixNine = MapUnion(codesForZeroSixNine, ToByteMap(signal))
		}
	}

	segmentCodes := make([]byte, 7)
	// set top segment key
	for k := range codeForSeven {
		_, inOne := codeForOne[k]
		if !inOne {
			segmentCodes[0] = k
			delete(codeForEight, k)
			break
		}
	}
	for k, v := range codesForTwoThreeFive {
		_, inFour := codeForFour[k]
		// set top left segment key
		if v == 1 && inFour {
			segmentCodes[5] = k
			delete(codeForEight, k)
		}
		// set bottom left segment key
		if v == 1 && !inFour {
			segmentCodes[4] = k
			delete(codeForEight, k)
		}
		// set center segment key
		if v == 3 && inFour {
			segmentCodes[6] = k
			delete(codeForEight, k)
		}
	}
	for k, v := range codesForZeroSixNine {
		_, inOne := codeForOne[k]
		// set top right segment key
		if v == 2 && inOne {
			segmentCodes[1] = k
			delete(codeForEight, k)
		}
		// set bottom right segment key
		if v == 3 && inOne {
			segmentCodes[2] = k
			delete(codeForEight, k)
		}
	}
	// set bottom segment key
	segmentCodes[3] = MapKeys(codeForEight)[0]

	return segmentCodes
}

func ToByteMap(input string) map[byte]int {
	counter := map[byte]int{}
	for i := range input {
		counter[input[i]]++
	}
	return counter
}

func MapKeys(data map[byte]int) []byte {
	keys := []byte{}
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func MapIntersect(lhs map[byte]int, rhs map[byte]int) map[byte]int {
	result := map[byte]int{}
	for k, v := range lhs {
		rv, in := rhs[k]
		if in {
			result[k] = v + rv
		}
	}
	return result
}

func MapUnion(lhs map[byte]int, rhs map[byte]int) map[byte]int {
	result := map[byte]int{}
	for k, v := range lhs {
		result[k] = v
	}
	for k, v := range rhs {
		result[k] += v
	}
	return result
}

// Encoding (sorry gofmt scrambles it):
//
//	    0
//	  5   1
//		   6
//	  4   2
//		   3
func SegmentMaskToNumber(segmentMask uint) int {
	switch segmentMask {
	case 0b0111111:
		return 0
	case 0b0000110:
		return 1
	case 0b1011011:
		return 2
	case 0b1001111:
		return 3
	case 0b1100110:
		return 4
	case 0b1101101:
		return 5
	case 0b1111101:
		return 6
	case 0b0000111:
		return 7
	case 0b1111111:
		return 8
	case 0b1101111:
		return 9
	}
	return -1
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	cipherLines := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Digits 1, 4, 7, 8 appear in the output values %v times.\n", PartOne(cipherLines))
	fmt.Println("Part two")
	fmt.Printf("The sum of the deciphered values: %v\n", PartTwo(cipherLines))
}
