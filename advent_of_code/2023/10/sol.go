package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func PartOne(maze []string) int {
	i, j := FindS(maze)
	if i == -1 || j == -1 {
		log.Fatalf("No 'S' was found in the input.")
	}

	path := TrackPipePath(maze, i, j)
	return len(path) / 2
}

// It may look a little hacky I admit, but the algorithm is very simple.
// We know we have a closed loop, so one side of the loop has to be the outside
// the other the inside. Initially we just go through the loop marking the empty space
// on the left as the left group and on the right as the right. After that we just
// grow each group. The one that reaches the edge of the matrix is the outer group.
//
// This very nicely sidesteps the issue where there is no empty space between two pipe
// sections but it acts as a connection. The side of the pipe in a closed loop never changes.
func PartTwo(maze []string) int {
	i, j := FindS(maze)
	if i == -1 || j == -1 {
		log.Fatalf("No 'S' was found in the input.")
	}

	path := TrackPipePath(maze, i, j)
	markerMaze := make([][]rune, len(maze))
	emptyRow := strings.Repeat("_", len(maze[0]))
	for i := range markerMaze {
		markerMaze[i] = []rune(emptyRow)
	}

	// first mark all pipe locations
	for _, pathPoint := range path {
		// pipe is marked with O for visibility
		markerMaze[pathPoint.Location.I][pathPoint.Location.J] = 'O'
	}
	// trace pipe and mark sides, initializing the partitions
	onLeft := map[Pair]struct{}{}
	onRight := map[Pair]struct{}{}
	for i := range path {
		if i > 0 {
			// we must mark the same location with the previous direction to handle the bends
			left, right := getSides(markerMaze, path[i].Location, path[i-1].Direction, 'O')
			if insertIfUnmarked(markerMaze, onLeft, left) {
				markerMaze[left.I][left.J] = '.'
			}
			if insertIfUnmarked(markerMaze, onRight, right) {
				markerMaze[right.I][right.J] = ':'
			}
		}
		left, right := getSides(markerMaze, path[i].Location, path[i].Direction, 'O')
		if insertIfUnmarked(markerMaze, onLeft, left) {
			markerMaze[left.I][left.J] = '.'
		}
		if insertIfUnmarked(markerMaze, onRight, right) {
			markerMaze[right.I][right.J] = ':'
		}
	}

	// grow the partitions until the and, if one finds the outer bound of the matrix
	// that will be the outside
	leftIsOuter := growGroup(markerMaze, onLeft, '_', '.')
	growGroup(markerMaze, onRight, '_', ':')

	leftCount := 0
	rightCount := 0
	for _, line := range markerMaze {
		for _, r := range line {
			if r == '.' {
				leftCount++
			}
			if r == ':' {
				rightCount++
			}
		}
	}

	for _, line := range markerMaze {
		for _, r := range line {
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}

	innerCount := leftCount
	if leftIsOuter {
		innerCount = rightCount
	}
	return innerCount
}

func growGroup(markerMaze [][]rune, toProcess map[Pair]struct{}, emptyMarker rune, processedMarker rune) bool {
	offset := []Pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	outerGroupConfirmed := false
	for len(toProcess) != 0 {
		nextToProcess := map[Pair]struct{}{}
		for location := range toProcess {
			for _, offset := range offset {
				candidate := Pair{location.I + offset.I, location.J + offset.J}
				if candidate.I < 0 || candidate.I >= len(markerMaze) || candidate.J < 0 || candidate.J >= len(markerMaze[0]) {
					continue
				}
				if candidate.I == 0 || candidate.I == len(markerMaze)-1 || candidate.J == 0 || candidate.J == len(markerMaze[0])-1 {
					outerGroupConfirmed = true
				}
				if markerMaze[candidate.I][candidate.J] == emptyMarker {
					markerMaze[candidate.I][candidate.J] = processedMarker
					nextToProcess[candidate] = struct{}{}
				}
			}
		}
		toProcess = nextToProcess
	}
	return outerGroupConfirmed
}

func getSides(markerMaze [][]rune, location, direction Pair, pipeSymbol rune) (Pair, Pair) {
	leftOffset, rightOffset := GetLeftRightOffsets(direction)
	leftLocation := Pair{location.I + leftOffset.I, location.J + leftOffset.J}
	rightLocation := Pair{location.I + rightOffset.I, location.J + rightOffset.J}
	return leftLocation, rightLocation
}

func insertIfUnmarked(markerMaze [][]rune, tracker map[Pair]struct{}, location Pair) bool {
	inserted := false
	if location.I < 0 || location.I >= len(markerMaze) || location.J < 0 || location.J >= len(markerMaze[0]) {
		return inserted
	}
	if markerMaze[location.I][location.J] == '_' {
		tracker[location] = struct{}{}
		inserted = true
	}
	return inserted
}

func GetLeftRightOffsets(direction Pair) (Pair, Pair) {
	directions := []Pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	i := 0
	for ; i < len(directions); i++ {
		if directions[i] == direction {
			break
		}
	}
	left := (i + 3) % 4
	right := (i + 1) % 4
	return directions[left], directions[right]
}

func FindS(maze []string) (int, int) {
	for i := 0; i < len(maze); i++ {
		for j := 0; j < len(maze[0]); j++ {
			if rune(maze[i][j]) == 'S' {
				return i, j
			}
		}
	}
	return -1, -1
}

type Pair struct {
	I int
	J int
}

type PathPoint struct {
	Location  Pair
	Direction Pair
}

func TrackPipePath(maze []string, i, j int) []PathPoint {
	path := []PathPoint{}

	delta_i := 0
	delta_j := 0
	i_offsets := []int{0, 0, 1, -1}
	j_offsets := []int{1, -1, 0, 0}
	pipeSymbol := 'S'
	for n := range i_offsets {
		delta_i = i_offsets[n]
		delta_j = j_offsets[n]
		next_i := delta_i + i
		next_j := delta_j + j
		if next_i < 0 || next_i >= len(maze) || next_j < 0 || next_j >= len(maze[0]) {
			continue
		}
		pipeSymbol = rune(maze[next_i][next_j])
		if CanEnterPipe(pipeSymbol, delta_i, delta_j) {
			break
		}
	}

	for pipeSymbol != 'S' {
		path = append(path, PathPoint{Pair{i, j}, Pair{delta_i, delta_j}})
		i = delta_i + i
		j = delta_j + j
		pipeSymbol = rune(maze[i][j])
		if pipeSymbol != 'S' {
			delta_i, delta_j = ThroughPipe(pipeSymbol, delta_i, delta_j)
		}
	}

	return path
}

func CanEnterPipe(r rune, delta_i, delta_j int) bool {
	switch r {
	case '-':
		{
			if (delta_i == 0 && delta_j == 1) || (delta_i == 0 && delta_j == -1) {
				return true
			}
			break
		}
	case '|':
		{
			if (delta_i == 1 && delta_j == 0) || (delta_i == -1 && delta_j == 0) {
				return true
			}
			break
		}
	case 'J':
		{
			if (delta_i == 1 && delta_j == 0) || (delta_i == 0 && delta_j == 1) {
				return true
			}
			break
		}
	case 'F':
		{
			if (delta_i == 0 && delta_j == -1) || (delta_i == -1 && delta_j == 0) {
				return true
			}
			break
		}
	case 'L':
		{
			if (delta_i == 1 && delta_j == 0) || (delta_i == 0 && delta_j == -1) {
				return true
			}
			break
		}
	case '7':
		{
			if (delta_i == 0 && delta_j == 1) || (delta_i == -1 && delta_j == 0) {
				return true
			}
			break
		}
	case 'S':
		return true
	default:
		{
			return false
		}
	}
	return false
}

func ThroughPipe(r rune, delta_i, delta_j int) (int, int) {
	if r == 'J' || r == 'F' {
		return -delta_j, -delta_i
	}
	if r == 'L' || r == '7' {
		return delta_j, delta_i
	}
	if r == '-' || r == '|' {
		return delta_i, delta_j
	}
	log.Fatalf("Unexpected pipe symbol received: %c\n", r)
	return 0, 0
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	maze := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	fmt.Println("Part one")
	fmt.Printf("Number of steps till the farthest point in the loop: %v\n", PartOne(maze))
	fmt.Println("Part two")
	fmt.Printf("Number of internal points: %v\n", PartTwo(maze))
}
