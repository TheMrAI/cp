package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func PartOne(maze []string) int {
	i, j := FindS(maze)
	if i == -1 || j == -1 {
		log.Fatalf("No 'S' was found in the input.")
	}

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
			i = next_i
			j = next_j
			delta_i, delta_j = ThroughPipe(pipeSymbol, delta_i, delta_j)
			break
		}
	}

	loopLength := 1
	for pipeSymbol != 'S' {
		i = delta_i + i
		j = delta_j + j
		pipeSymbol = rune(maze[i][j])
		if pipeSymbol != 'S' {
			delta_i, delta_j = ThroughPipe(pipeSymbol, delta_i, delta_j)
		}
		loopLength++
	}

	return loopLength / 2
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
}
