package main

import (
	"bufio"
	"fmt"
	"os"
)

func ParseInput(lines []string) [][]string {
	notes := [][]string{}

	note := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			notes = append(notes, note)
			note = []string{}
			continue
		}
		note = append(note, line)
	}
	if len(note) != 0 {
		notes = append(notes, note)
	}

	return notes
}

func PartOne(matrices [][]string) int {
	sum := 0
	for _, matrix := range matrices {
		horizontalMirrorPositions := FindMirrorPositions(matrix, 0)
		if len(horizontalMirrorPositions) > 0 {
			sum += 100 * horizontalMirrorPositions[0]
		}
		verticalMirrorPositions := FindMirrorPositions(Transpose(matrix), 0)
		if len(verticalMirrorPositions) > 0 {
			sum += verticalMirrorPositions[0]
		}
	}
	return sum
}

func PartTwo(matrices [][]string) int {
	sum := 0
	for _, matrix := range matrices {
		horizontalReflection := FindMirrorPositions(matrix, 0)
		horizontalReflectionWithCorrection := FindMirrorPositions(matrix, 1)
		horizontalDifference := setDifference(horizontalReflectionWithCorrection, horizontalReflection)
		if len(horizontalDifference) > 0 {
			sum += 100 * horizontalDifference[0]
		}

		verticalReflection := FindMirrorPositions(Transpose(matrix), 0)
		verticalReflectionWithCorrection := FindMirrorPositions(Transpose(matrix), 1)
		verticalDifferences := setDifference(verticalReflectionWithCorrection, verticalReflection)
		if len(verticalDifferences) > 0 {
			sum += verticalDifferences[0]
		}
	}
	return sum
}

// A - B
func setDifference(A, B []int) []int {
	i := 0
	j := 0
	difference := []int{}
	for ; i < len(A); i++ {
		for ; j < len(B); j++ {
			if A[i] <= B[j] {
				break
			}
		}
		if j < len(B) && A[i] == B[j] {
			continue
		}
		difference = append(difference, A[i])
	}
	return difference
}

// meh not checking for errors, life
func Transpose(matrix []string) []string {
	transposed := make([]string, len(matrix[0]))
	for j := 0; j < len(matrix[0]); j++ {
		line := make([]byte, len(matrix))
		for i := 0; i < len(matrix); i++ {
			line[i] = matrix[i][j]
		}
		transposed[j] = string(line)
	}
	return transposed
}

func FindMirrorPositions(input []string, correctionBudget int) []int {
	mirrorPositions := []int{}
	mirrorPosition := 1
	for ; mirrorPosition < len(input); mirrorPosition++ {
		sideLength := mirrorPosition
		onRight := len(input) - mirrorPosition
		if onRight < sideLength {
			sideLength = onRight
		}
		if IsSymmetric(input[mirrorPosition-sideLength:mirrorPosition+sideLength], correctionBudget) {
			mirrorPositions = append(mirrorPositions, mirrorPosition)
		}
	}
	return mirrorPositions
}

func IsSymmetric(input []string, correctionBudget int) bool {
	front := 0
	back := len(input) - 1
	correctionApplied := false
	for front < back {
		if input[front] != input[back] {
			if !correctionApplied {
				differenceCount := 0
				for i := range input[front] {
					if input[front][i] != input[back][i] {
						differenceCount++
					}
				}
				if differenceCount == correctionBudget {
					correctionApplied = true
					front++
					back--
					continue
				}
			}
			return false
		}
		front++
		back--
	}
	return true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	notes := ParseInput(lines)
	
	fmt.Println("Part one")
	fmt.Printf("Note summary result: %v\n", PartOne(notes))
	fmt.Println("Two one")
	fmt.Printf("Note summary result, with smudge correction: %v\n", PartTwo(notes))
}
