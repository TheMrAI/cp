package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func parseNumbers(numbersText []string) ([]int, error) {
	numbers := make([]int, 0, len(numbersText))
	for _, numText := range numbersText {
		num, err := strconv.Atoi(numText)
		if err != nil {
			return numbers, fmt.Errorf("Invalid input: %v", numText)
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

type Coordinate struct {
	I int
	J int
}

type BingoBoard struct {
	Matrix          [][]int
	VisitedInRow    []uint64
	VisitedInColumn []uint64
	Done bool
	valueAt         map[int]Coordinate
}

func NewBingoBoard(matrix [][]int) BingoBoard {
	valueAt := map[int]Coordinate{}
	for i, line := range matrix {
		for j, num := range line {
			valueAt[num] = Coordinate{i, j}
		}
	}
	board := BingoBoard{matrix, make([]uint64, len(matrix)), make([]uint64, len(matrix)), false, valueAt}
	return board
}

func (board *BingoBoard) Mark(number int) bool {
	atLocation, ok := board.valueAt[number]
	if !ok {
		return false
	}
	fullyMarkedMask := (uint64(1) << len(board.Matrix)) - 1
	board.VisitedInRow[atLocation.I] = board.VisitedInRow[atLocation.I] | (0b1 << atLocation.J)
	board.VisitedInColumn[atLocation.J] = board.VisitedInColumn[atLocation.J] | (0b1 << atLocation.I)
	return (board.VisitedInRow[atLocation.I] == fullyMarkedMask) || (board.VisitedInColumn[atLocation.J] == fullyMarkedMask)
}

type Winner struct {
	LastNumber int
	Board BingoBoard
}

func partOneAndTwo(bingoBoards []BingoBoard, numbers []int) (int, int) {
	winningOrder := [](Winner){}
	for _, num := range numbers {
		for i := range bingoBoards {
			if bingoBoards[i].Done {
				continue
			}
			if bingoBoards[i].Mark(num) {
				bingoBoards[i].Done = true
				winningOrder = append(winningOrder, Winner{num, bingoBoards[i]})
			}
		}
	}

	firstWinnerScore := unmarkedSum(winningOrder[0].Board) * winningOrder[0].LastNumber
	lastWinnerScore := unmarkedSum(winningOrder[len(winningOrder)-1].Board) * winningOrder[len(winningOrder)-1].LastNumber
	return firstWinnerScore, lastWinnerScore
}

func unmarkedSum(board BingoBoard) int {
	sum := 0
	for i, line := range board.Matrix {
		for j, num := range line {
			if (board.VisitedInRow[i] & (0b1 << j)) == 0 {
				sum += num
			}
		}
	}
	return sum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		log.Fatalln("Empty input received.")
	}

	numbersText := strings.FieldsFunc(scanner.Text(), func(c rune) bool {
		return c == ','
	})
	numbers, err := parseNumbers(numbersText)
	if err != nil {
		log.Fatal(err)
	}

	bingoBoards := []BingoBoard{}
	matrix := [][]int{}
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		numbersText = strings.FieldsFunc(text, func(c rune) bool {
			return unicode.IsSpace(c)
		})
		matrixRow, err := parseNumbers(numbersText)
		if err != nil {
			log.Fatal(err)
		}
		matrix = append(matrix, matrixRow)
		if len(matrix) == len(matrixRow) {
			bingoBoards = append(bingoBoards, NewBingoBoard(matrix))
			matrix = [][]int{}
		}
	}

	partOneResult, partTwoResult := partOneAndTwo(bingoBoards, numbers)
	fmt.Println("Part one")
	fmt.Printf("Score of winning board: %v\n", partOneResult)
	fmt.Println("Part two")
	fmt.Printf("Score of last winning board: %v\n", partTwoResult)
}
