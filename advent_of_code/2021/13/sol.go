package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Pair struct {
	X int
	Y int
}

type Fold struct {
	AlongY bool
	Index  int
}

func ParseInput(lines []string) (map[Pair]struct{}, []Fold) {
	points := map[Pair]struct{}{}

	lineIndex := 0
	for ; lineIndex < len(lines); lineIndex++ {
		if len(lines[lineIndex]) == 0 {
			break
		}
		numTexts := strings.Split(lines[lineIndex], ",")
		if len(numTexts) != 2 {
			log.Fatalf("Expected two numbers separated by a comma, but got: '%v'\n", lines[lineIndex])
		}
		nums := make([]int, 2)
		for i := range numTexts {
			num, err := strconv.Atoi(numTexts[i])
			if err != nil {
				log.Fatalf("Integer conversion failed on input: '%v'\n", numTexts[i])
			}
			nums[i] = num
		}
		points[Pair{nums[0], nums[1]}] = struct{}{}
	}

	lineIndex++
	folds := []Fold{}
	foldRegexp := regexp.MustCompile(".*(\\w)=(\\d+)")
	for ; lineIndex < len(lines); lineIndex++ {
		submatches := foldRegexp.FindStringSubmatch(lines[lineIndex])
		if len(submatches) != 3 {
			log.Fatalf("Failed to parse fold instructions on input: '%v'\n", lines[lineIndex])
		}
		value, err := strconv.Atoi(submatches[2])
		if err != nil {
			log.Fatalf("Integer conversion failed on input: '%v'\n", submatches[2])
		}
		if submatches[1] == "x" {
			folds = append(folds, Fold{false, value})
		} else if submatches[1] == "y" {
			folds = append(folds, Fold{true, value})
		} else {
			log.Fatalf("Unexpected axis: '%v'\n", submatches[1])
		}
	}

	return points, folds
}

func PartOne(points map[Pair]struct{}, folds []Fold) int {
	afterFold := ExecuteFolds(points, folds, 1)
	return len(afterFold)
}

func ExecuteFolds(points map[Pair]struct{}, folds []Fold, numberOfFolds int) map[Pair]struct{} {
	foldingPoints := map[Pair]struct{}{}
	for k, v := range points {
		foldingPoints[k] = v
	}

	for foldCount := 0; foldCount < numberOfFolds && foldCount < len(folds); foldCount++ {
		fold := folds[foldCount]
		for k, _ := range foldingPoints {
			if fold.AlongY {
				if fold.Index < k.Y {
					foldingPoints[Pair{k.X, fold.Index - (k.Y - fold.Index)}] = struct{}{}
					delete(foldingPoints, k)
				}
			} else {
				if fold.Index < k.X {
					foldingPoints[Pair{fold.Index - (k.X - fold.Index), k.Y}] = struct{}{}
					delete(foldingPoints, k)
				}
			}
		}
	}
	return foldingPoints
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	points, folds := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("Number of dots after the first fold instruction: %v\n", PartOne(points, folds))
	fullyFolded := ExecuteFolds(points, folds, len(folds))
	fmt.Println("Part two")
	xMax := 0
	yMax := 0
	for k := range fullyFolded {
		if xMax < k.X {
			xMax = k.X
		}
		if yMax < k.Y {
			yMax = k.Y
		}
	}
	for y := 0; y <= yMax; y++ {
	for x := 0; x <= xMax; x++ {
			_, ok := fullyFolded[Pair{x, y}]
			if ok {
				fmt.Printf("%c", '#')
				continue
			}
			fmt.Printf("%c", '.')
		}
		fmt.Println()
	}
}
