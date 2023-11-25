package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Direction uint8

const (
	UP Direction = iota
	FORWARD
	DOWN
)

type Move struct {
	MoveDir   Direction
	Amplitude int
}

func ToDir(direction string) (Direction, error) {
	switch direction {
	case "forward":
		{
			return FORWARD, nil
		}
	case "up":
		{
			return UP, nil
		}
	case "down":
		{
			return DOWN, nil
		}
	default:
		{
			return UP, fmt.Errorf("Invalid direction: %v", direction)
		}
	}
}

func trackPosition(horizontal int, vertical int, moves []Move) (int, int) {
	for _, move := range moves {
		switch move.MoveDir {
		case FORWARD:
			{
				horizontal += move.Amplitude
				break
			}
		case UP:
			{
				vertical -= move.Amplitude
				break
			}
		case DOWN:
			{
				vertical += move.Amplitude
				break
			}
		}
	}
	return horizontal, vertical
}

func trackPositionWithAim(horizontal int, vertical int, aim int, moves []Move) (int, int) {
	for _, move := range moves {
		switch move.MoveDir {
		case FORWARD:
			{
				horizontal += move.Amplitude
				vertical += aim * move.Amplitude
				break
			}
		case UP:
			{
				aim -= move.Amplitude
				break
			}
		case DOWN:
			{
				aim += move.Amplitude
				break
			}
		}
	}
	return horizontal, vertical
}

func main() {
	moves := []Move{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lineSections := strings.Split(line, " ")
		if len(lineSections) != 2 {
			fmt.Printf("Invalid line received: %s", line)
			os.Exit(1)
		}
		direction, err := ToDir(lineSections[0])
		if err != nil {
			log.Fatal(err)
		}
		amplitude, err := strconv.Atoi(lineSections[1])
		if err != nil {
			log.Fatal(err)
		}
		moves = append(moves, Move{direction, amplitude})
	}
	horizontal, vertical := trackPosition(0, 0, moves)
	fmt.Printf("Final position, horizontal: %v vertical: %v, multiplied: %v\n", horizontal, vertical, horizontal*vertical)

	horizontal, vertical = trackPositionWithAim(0, 0, 0, moves)
	fmt.Printf("Final position, horizontal: %v vertical: %v, multiplied: %v\n", horizontal, vertical, horizontal*vertical)
}
