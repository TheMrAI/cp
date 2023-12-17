package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func ParseInput(lines []string) [][]int {
	matrix := make([][]int, 0, len(lines))

	numberRegexp := regexp.MustCompile("\\d")

	for _, line := range lines {
		numberTexts := numberRegexp.FindAllString(line, -1)
		matrixLine := make([]int, 0, len(line))
		for _, text := range numberTexts {
			val, err := strconv.Atoi(text)
			if err != nil {
				log.Fatalf("Failed to convert input to string: '%v'\n", text)
			}
			matrixLine = append(matrixLine, val)
		}
		matrix = append(matrix, matrixLine)
	}

	return matrix
}

type Pair struct {
	I int
	J int
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func RotateLeft(dir Direction) Direction {
	return (dir + 3) % 4
}

func RotateRight(dir Direction) Direction {
	return (dir + 1) % 4
}

func DirToOffsets(dir Direction) Pair {
	dirToOffsets := [4]Pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	return dirToOffsets[dir]
}

type BlockId struct {
	Pos        Pair
	Dir        Direction
	HeatLoss   int
	StepsTaken int // how many steps did we take in a straight line
}

type Block struct {
	Id       BlockId
	priority int
	index    int
}

type PriorityQueue []*Block

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want the lowest priorities first, because we
	// use the manhattan distance as priority
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Block)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// First solution time:
// real    2m17.350s
// user    2m5.124s
// sys     0m12.837s
// and it eats up a whopping 30 GiB of memory
func PartOne(matrix [][]int) int {
	target := Pair{len(matrix) - 1, len(matrix[0]) - 1}
	return FindMinimalPath(matrix, target, 1, 3)
}

// First solution didn't finish for the big input in 10 minutes.
func PartTwo(matrix [][]int) int {
	target := Pair{len(matrix) - 1, len(matrix[0]) - 1}
	return FindMinimalPath(matrix, target, 4, 10)
}

type MinimalHeatLoss [][]int // indexed by direction, then by steps taken - minimumStepCount

// no support for handling the start point, for that we would have to programmatically find
// the starting directions and it isn't necessary here
func FindMinimalPath(matrix [][]int, target Pair, minStepCount, maxStepCount int) int {
	minimalHeatLosses := make([][]MinimalHeatLoss, len(matrix))
	for i := range minimalHeatLosses {
		row := make([]MinimalHeatLoss, len(matrix[0]))
		for j := range row {
			row[j] = make(MinimalHeatLoss, 4)
			for dir := range row[j] {
				row[j][dir] = make([]int, 1+maxStepCount-minStepCount)
				for step := range row[j][dir] {
					row[j][dir][step] = math.MaxInt
				}
			}
		}
		minimalHeatLosses[i] = row
	}
	minimalHeatLoss := math.MaxInt

	queued := map[BlockId]struct{}{}
	toBeChecked := PriorityQueue{}
	heap.Init(&toBeChecked)
	// We add two search points immediately. Since we can never go forward first without turning, to cover all
	// starting possibilities we pretend we arrived there from 2 directions.
	starterBlockIdOne := BlockId{Pair{0, 0}, Right, 0, 1}
	starterBlockIdTwo := BlockId{Pair{0, 0}, Down, 0, 1}
	insertSteps(matrix, &toBeChecked, queued, target, minStepCount, maxStepCount, Block{starterBlockIdOne, ManhattanDistance(Pair{0, 0}, target), 0}, true)
	insertSteps(matrix, &toBeChecked, queued, target, minStepCount, maxStepCount, Block{starterBlockIdTwo, ManhattanDistance(Pair{0, 0}, target), 0}, true)

	for len(toBeChecked) != 0 {
		checking := heap.Pop(&toBeChecked).(*Block)
		if checking.Id.Pos.I == target.I && checking.Id.Pos.J == target.J {
			fmt.Printf("Target reached with heat loss: %v\n", checking.Id.HeatLoss)
			if checking.Id.HeatLoss < minimalHeatLoss {
				minimalHeatLoss = checking.Id.HeatLoss
			}
			continue
		}

		if checking.Id.HeatLoss >= minimalHeatLosses[checking.Id.Pos.I][checking.Id.Pos.J][checking.Id.Dir][checking.Id.StepsTaken-minStepCount] {
			continue
		}
		minimalHeatLosses[checking.Id.Pos.I][checking.Id.Pos.J][checking.Id.Dir][checking.Id.StepsTaken-minStepCount] = checking.Id.HeatLoss
		// If we assume just 1 heat lost per block, based on the current heat loss and manhattan distance
		// we can be absolutely sure, this branch can't bring an improvement.
		if checking.priority+checking.Id.HeatLoss >= minimalHeatLoss {
			continue
		}

		leftBlock, leftValid := getBlockOnLeft(matrix, target, *checking)
		insertSteps(matrix, &toBeChecked, queued, target, minStepCount, maxStepCount, leftBlock, leftValid)

		rightBlock, rightValid := getBlockOnRight(matrix, target, *checking)
		insertSteps(matrix, &toBeChecked, queued, target, minStepCount, maxStepCount, rightBlock, rightValid)
	}

	return minimalHeatLoss
}

func insertSteps(matrix [][]int, toBeChecked *PriorityQueue, queued map[BlockId]struct{}, target Pair, minStepCount, maxStepCount int, block Block, blockValid bool) {
	for k := 1; k <= maxStepCount; k++ {
		if !blockValid {
			break
		}
		if block.Id.StepsTaken < minStepCount {
			block, blockValid = nextBlock(matrix, target, block)
			continue
		}
		insert(toBeChecked, queued, block)
		block, blockValid = nextBlock(matrix, target, block)
	}
}

func insert(toBeChecked *PriorityQueue, queued map[BlockId]struct{}, block Block) {
	_, alreadyQueued := queued[block.Id]
	if !alreadyQueued {
		heap.Push(toBeChecked, &block)
		queued[block.Id] = struct{}{}
	}
}

func getBlockOnLeft(matrix [][]int, target Pair, block Block) (Block, bool) {
	leftDir := RotateLeft(block.Id.Dir)
	leftOffset := DirToOffsets(leftDir)
	leftPos := Pair{block.Id.Pos.I + leftOffset.I, block.Id.Pos.J + leftOffset.J}
	if leftPos.I < 0 || leftPos.I >= len(matrix) || leftPos.J < 0 || leftPos.J >= len(matrix[0]) {
		return Block{}, false
	}
	leftHeatLoss := block.Id.HeatLoss + matrix[leftPos.I][leftPos.J]
	leftBlockId := BlockId{leftPos, leftDir, leftHeatLoss, 1} // we took one step
	return Block{leftBlockId, ManhattanDistance(leftPos, target), 0}, true
}

func getBlockOnRight(matrix [][]int, target Pair, block Block) (Block, bool) {
	rightDir := RotateRight(block.Id.Dir)
	rightOffset := DirToOffsets(rightDir)
	rightPos := Pair{block.Id.Pos.I + rightOffset.I, block.Id.Pos.J + rightOffset.J}
	if rightPos.I < 0 || rightPos.I >= len(matrix) || rightPos.J < 0 || rightPos.J >= len(matrix[0]) {
		return Block{}, false
	}
	rightHeatLoss := block.Id.HeatLoss + matrix[rightPos.I][rightPos.J]
	rightBlockId := BlockId{rightPos, rightDir, rightHeatLoss, 1} // we took one step
	return Block{rightBlockId, ManhattanDistance(rightPos, target), 0}, true
}

func nextBlock(matrix [][]int, target Pair, block Block) (Block, bool) {
	frontOffset := DirToOffsets(block.Id.Dir)
	frontPos := Pair{block.Id.Pos.I + frontOffset.I, block.Id.Pos.J + frontOffset.J}
	if frontPos.I < 0 || frontPos.I >= len(matrix) || frontPos.J < 0 || frontPos.J >= len(matrix[0]) {
		return Block{}, false
	}
	frontHeatLoss := block.Id.HeatLoss + matrix[frontPos.I][frontPos.J]
	frontBlockId := BlockId{frontPos, block.Id.Dir, frontHeatLoss, block.Id.StepsTaken + 1}
	return Block{frontBlockId, ManhattanDistance(frontPos, target), 0}, true
}

func ManhattanDistance(from, to Pair) int {
	deltaI := to.I - from.I
	if deltaI < 0 {
		deltaI = 0 - deltaI
	}
	deltaJ := to.J - from.J
	if deltaJ < 0 {
		deltaJ = 0 - deltaJ
	}
	return deltaI + deltaJ
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	data := ParseInput(lines)
	for _, d := range data {
		fmt.Println(d)
	}

	// fmt.Println("Part one")
	// fmt.Printf("Path length with minimal heat loss: %v\n", PartOne(data))
	fmt.Println("Part two")
	fmt.Printf("Path length with minimal heat loss using ULTRA crucibles: %v\n", PartTwo(data))
}
