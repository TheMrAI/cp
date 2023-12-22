package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Pos struct {
	X int
	Y int
	Z int
}

type Brick struct {
	A           Pos
	B           Pos
	Supports    map[int]struct{}
	SupportedBy map[int]struct{}
}

func ParseInput(lines []string) []Brick {
	bricks := []Brick{}

	numberRegexp := regexp.MustCompile("-?\\d+")

	for _, line := range lines {
		coordinates := strings.Split(line, "~")
		if len(coordinates) != 2 {
			log.Fatalf("Unexpected input received: '%v'\n", line)
		}
		brick := Brick{}
		x, y, z := parseCoordinates(coordinates[0], numberRegexp)
		brick.A = Pos{x, y, z}
		x, y, z = parseCoordinates(coordinates[1], numberRegexp)
		brick.B = Pos{x, y, z}
		brick.SupportedBy = map[int]struct{}{}
		brick.Supports = map[int]struct{}{}
		bricks = append(bricks, brick)
	}

	return bricks
}

func parseCoordinates(data string, numberRegexp *regexp.Regexp) (int, int, int) {
	coordSet := numberRegexp.FindAllString(data, -1)
	if len(coordSet) != 3 {
		log.Fatalf("Expected 3 numbers in input, got: '%v'\n", data)
	}
	return convertNumber(coordSet[0]), convertNumber(coordSet[1]), convertNumber(coordSet[2])
}

func convertNumber(number string) int {
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Fatalf("Integer conversion failed on input: %v\n", number)
	}
	return value
}

func SimulateBrickFall(bricks []Brick) {
	minX := math.MaxInt
	maxX := math.MinInt
	minY := math.MaxInt
	maxY := math.MinInt

	for _, brick := range bricks {
		minX = Min(minX, brick.A.X)
		maxX = Max(maxX, brick.A.X)
		minX = Min(minX, brick.B.X)
		maxX = Max(maxX, brick.B.X)
		minY = Min(minY, brick.A.Y)
		maxY = Max(maxY, brick.A.Y)
		minY = Min(minY, brick.B.Y)
		maxY = Max(maxY, brick.B.Y)
	}
	// shift all bricks close to 0, 0
	for i := range bricks {
		bricks[i].A = Pos{bricks[i].A.X - minX, bricks[i].A.Y - minY, bricks[i].A.Z}
		bricks[i].B = Pos{bricks[i].B.X - minX, bricks[i].B.Y - minY, bricks[i].B.Z}
	}
	// sort the bricks by Z height!!!
	sort.Slice(bricks, func(i, j int) bool {
		return Min(bricks[i].A.Z, bricks[i].B.Z) < Min(bricks[j].A.Z, bricks[j].B.Z)
	})

	// a top-down view of the currently laid brick ids
	zView := make([][]int, 1+maxX-minX)
	for i := range zView {
		zView[i] = make([]int, 1+maxY-minY)
		for j := range zView[i] {
			zView[i][j] = -1
		}
	}

	// lay all bricks, snapping them to the highest z pos
	for i := range bricks {
		projection := collectProjectionBrickIdsAndHeights(zView, bricks, i)
		nextZ := -1
		for _, square := range projection {
			nextZ = Max(nextZ, square.Height)
		}
		brickIds := []int{}
		for _, square := range projection {
			if square.Id == -1 {
				continue
			}
			if bricks[square.Id].A.Z != nextZ && bricks[square.Id].B.Z != nextZ {
				continue
			}
			brickIds = append(brickIds, square.Id)
		}
		// now we increase Z for the brick snapping
		nextZ += 1
		for _, id := range brickIds {
			bricks[i].SupportedBy[id] = struct{}{}
		}
		// snap brick
		zDelta := bricks[i].A.Z - nextZ
		if zDelta > bricks[i].B.Z-nextZ {
			zDelta = bricks[i].B.Z - nextZ
		}
		bricks[i].A.Z -= zDelta
		bricks[i].B.Z -= zDelta

		updateZView(zView, bricks, i)

		// for k := range zView {
		// 	for l := range zView[k] {
		// 		fmt.Printf("%v ", zView[k][l])
		// 	}
		// 	fmt.Println()
		// }
		// fmt.Println()

		for _, bI := range brickIds {
			bricks[bI].Supports[i] = struct{}{}
		}
	}
}

type IdHeight struct {
	Id     int
	Height int
}

func collectProjectionBrickIdsAndHeights(zView [][]int, bricks []Brick, i int) []IdHeight {
	projection := []IdHeight{}

	startX := bricks[i].A.X
	lastX := bricks[i].B.X
	deltaX := lastX - startX
	if deltaX != 0 {
		deltaX /= Abs(lastX - startX)
	}

	startY := bricks[i].A.Y
	lastY := bricks[i].B.Y
	deltaY := lastY - startY
	if deltaY != 0 {
		deltaY /= Abs(lastY - startY)
	}

	for x := startX; x <= lastX; x += deltaX {
		for y := startY; y <= lastY; y += deltaY {
			id := zView[x][y]

			height := -1
			if id != -1 {
				height = Max(bricks[id].A.Z, bricks[id].B.Z)
			}
			projection = append(projection, IdHeight{id, height})
			if deltaY == 0 {
				break
			}
		}
		if deltaX == 0 {
			break
		}
	}

	return projection
}

func updateZView(zView [][]int, bricks []Brick, i int) {
	startX := bricks[i].A.X
	lastX := bricks[i].B.X
	deltaX := lastX - startX
	if deltaX != 0 {
		deltaX /= Abs(lastX - startX)
	}

	startY := bricks[i].A.Y
	lastY := bricks[i].B.Y
	deltaY := lastY - startY
	if deltaY != 0 {
		deltaY /= Abs(lastY - startY)
	}

	for x := startX; x <= lastX; x += deltaX {
		for y := startY; y <= lastY; y += deltaY {
			zView[x][y] = i
			if deltaY == 0 {
				break
			}
		}
		if deltaX == 0 {
			break
		}
	}
}

func Min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func Max(lhs, rhs int) int {
	if lhs < rhs {
		return rhs
	}
	return lhs
}

func Abs(val int) int {
	if val < 0 {
		return 0 - val
	}
	return val
}

func PartOne(bricks []Brick) int {
	disintegrateCount := 0
	for _, brick := range bricks {
		if len(brick.Supports) == 0 {
			disintegrateCount++
			continue
		}
		supportedAllHaveMultiSupport := true
		for supportedId := range brick.Supports {
			if len(bricks[supportedId].SupportedBy) == 1 {
				supportedAllHaveMultiSupport = false
				break
			}
		}
		if supportedAllHaveMultiSupport {
			disintegrateCount++
		}
	}
	return disintegrateCount
}

func PartTwo(bricks []Brick) int {
	// For each brick we will count how many times do we reach it with from
	// individual bricks
	// If the id counts for a given brick match it's len(SupportedBy) we
	// can be sure that through all support paths the given id affects the brick
	// and thus would cause it to eventually fall if removed.
	brickIdPropagationCounts := make([]map[int]int, len(bricks))
	for i := range brickIdPropagationCounts {
		brickIdPropagationCounts[i] = map[int]int{}
	}

	bricksOnTheFloor := []int{}
	for i := range bricks {
		if len(bricks[i].SupportedBy) == 0 {
			bricksOnTheFloor = append(bricksOnTheFloor, i)
		}
	}

	propagateDependencyCountersBFT(bricks, bricksOnTheFloor, brickIdPropagationCounts)

	propagatedDisintegrationCount := 0
	for id := range bricks {
		for _, count := range brickIdPropagationCounts[id] {
			if count == len(bricks[id].SupportedBy) {
				propagatedDisintegrationCount++
			}
		}
	}

	return propagatedDisintegrationCount
}

func propagateDependencyCountersBFT(bricks []Brick, bricksOnTheFloor []int, brickIdPropagationCounts []map[int]int) {
	prepared := make([]bool, len(bricks))
	toProcess := PriorityQueue{}
	heap.Init(&toProcess)
	for _, id := range bricksOnTheFloor {
		heap.Push(&toProcess, &BrickByHeight{id, Min(bricks[id].A.Z, bricks[id].B.Z), 0})
	}

	for len(toProcess) != 0 {
		checking := heap.Pop(&toProcess).(*BrickByHeight)
		id := checking.Id

		for supporterId := range bricks[id].SupportedBy {
			// add immediate dependencies as one
			brickIdPropagationCounts[id][supporterId] += 1
			// increase dependency counts transitively
			for transId, count := range brickIdPropagationCounts[supporterId] {
				if count != len(bricks[supporterId].SupportedBy) {
					continue
				}
				brickIdPropagationCounts[id][transId] += 1
			}
		}

		for supportedId := range bricks[id].Supports {
			if prepared[supportedId] {
				continue
			}
			heap.Push(&toProcess, &BrickByHeight{supportedId, Min(bricks[supportedId].A.Z, bricks[supportedId].B.Z), 0})
			prepared[supportedId] = true
		}
	}
}

type BrickByHeight struct {
	Id       int
	priority int
	index    int
}

type PriorityQueue []*BrickByHeight

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want the lowest elevation first.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*BrickByHeight)
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	bricks := ParseInput(lines)
	SimulateBrickFall(bricks)

	// 468
	fmt.Println("Part one")
	fmt.Printf("The number of safely removable single bricks: %v\n", PartOne(bricks))
	fmt.Println("Part two")
	// too low 24047, 25576
	fmt.Printf("The total number of cascade falling bricks: %v\n", PartTwo(bricks))
}
