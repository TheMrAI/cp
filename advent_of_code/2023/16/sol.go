package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// func ParseInput(lines []string) []int {
// 	data := []int{}

// 	numberRegexp := regexp.MustCompile("\\d+")

// 	for _, line := range lines {

// 	}

// 	return data
// }

type Pair struct {
	I int
	J int
}

type Direction int

// there can only be 4 or RotateDirection90Degrees breaks
const (
	Up Direction = iota
	Right
	Down
	Left
)

func StepBeam(location Pair, dir Direction) Pair {
	switch dir {
	case Up:
		return Pair{location.I - 1, location.J}
	case Right:
		return Pair{location.I, location.J + 1}
	case Down:
		return Pair{location.I + 1, location.J}
	default:
		return Pair{location.I, location.J - 1}
	}
}

func RotateDirection90Degrees(dir Direction, clockwise bool) Direction {
	if clockwise {
		return (dir + 1) % 4
	}
	return (dir + 3) % 4
}

type BeamData struct {
	Location Pair
	Dir      Direction
}

func BounceBeam(dir Direction, tiltedMirror byte) Direction {
	if tiltedMirror != '/' && tiltedMirror != '\\' {
		log.Fatalf("Non supported mirror symnbol received: '%v'\n", tiltedMirror)
	}

	if tiltedMirror == '/' {
		if dir == Up || dir == Down {
			return RotateDirection90Degrees(dir, true)
		}
		return RotateDirection90Degrees(dir, false)
	}
	if dir == Up || dir == Down {
		return RotateDirection90Degrees(dir, false)
	}
	return RotateDirection90Degrees(dir, true)
}

func PartOne(tiles []string) int {
	return CountEnergizedTiles(tiles, BeamData{Pair{0, 0}, Right})
}

func PartTwo(tiles []string) int {
	maximum := 0

	for i := 0; i < len(tiles); i++ {
		leftMax := CountEnergizedTiles(tiles, BeamData{Pair{i, 0}, Right})
		if leftMax > maximum {
			maximum = leftMax
		}
		rightMax := CountEnergizedTiles(tiles, BeamData{Pair{i, len(tiles[0]) - 1}, Left})
		if rightMax > maximum {
			maximum = rightMax
		}
	}

	for j := 0; j < len(tiles[0]); j++ {
		upMax := CountEnergizedTiles(tiles, BeamData{Pair{0, j}, Down})
		if upMax > maximum {
			maximum = upMax
		}
		downMax := CountEnergizedTiles(tiles, BeamData{Pair{len(tiles) - 1, j}, Up})
		if downMax > maximum {
			maximum = downMax
		}
	}

	return maximum
}

func CountEnergizedTiles(tiles []string, startingBeam BeamData) int {
	visited := make([][][]bool, len(tiles))
	// this is not the most optimal, we could use a bitmask, but good enough for this
	for i := range visited {
		visited[i] = make([][]bool, len(tiles[0]))
		for j := range visited[i] {
			visited[i][j] = make([]bool, 4)
		}
	}

	// BFT beam propagation until there are no more paths to check
	toProcess := []BeamData{}
	toProcess = append(toProcess, startingBeam)
	for len(toProcess) > 0 {
		beam := toProcess[0]
		toProcess = toProcess[1:]
		if beam.Location.I < 0 || beam.Location.I >= len(tiles) || beam.Location.J < 0 || beam.Location.J >= len(tiles[0]) {
			continue
		}
		if visited[beam.Location.I][beam.Location.J][beam.Dir] {
			continue
		}
		visited[beam.Location.I][beam.Location.J][beam.Dir] = true

		tile := tiles[beam.Location.I][beam.Location.J]
		switch tile {
		case '.':
			toProcess = append(toProcess, BeamData{StepBeam(beam.Location, beam.Dir), beam.Dir})
		case '\\':
			rotatedDir := BounceBeam(beam.Dir, tile)
			toProcess = append(toProcess, BeamData{StepBeam(beam.Location, rotatedDir), rotatedDir})
		case '/':
			rotatedDir := BounceBeam(beam.Dir, tile)
			toProcess = append(toProcess, BeamData{StepBeam(beam.Location, rotatedDir), rotatedDir})
		case '|':
			if beam.Dir == Up || beam.Dir == Down {
				toProcess = append(toProcess, BeamData{StepBeam(beam.Location, beam.Dir), beam.Dir})
				continue
			}
			beamDirAnticlockwise := RotateDirection90Degrees(beam.Dir, false)
			beamDirClockwise := RotateDirection90Degrees(beam.Dir, true)
			toProcess = append(toProcess, BeamData{StepBeam(beam.Location, beamDirAnticlockwise), beamDirAnticlockwise}, BeamData{StepBeam(beam.Location, beamDirClockwise), beamDirClockwise})
		case '-':
			if beam.Dir == Left || beam.Dir == Right {
				toProcess = append(toProcess, BeamData{StepBeam(beam.Location, beam.Dir), beam.Dir})
				continue
			}
			beamDirAnticlockwise := RotateDirection90Degrees(beam.Dir, false)
			beamDirClockwise := RotateDirection90Degrees(beam.Dir, true)
			toProcess = append(toProcess, BeamData{StepBeam(beam.Location, beamDirAnticlockwise), beamDirAnticlockwise}, BeamData{StepBeam(beam.Location, beamDirClockwise), beamDirClockwise})
		}
	}

	energizedTileCount := 0
	for i := range visited {
		for _, tile := range visited[i] {
			for _, energized := range tile {
				if energized {
					energizedTileCount++
					break
				}
			}
		}
	}

	return energizedTileCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	fmt.Println("Part one")
	fmt.Printf("Number of energized tiles: %v\n", PartOne(lines))
	fmt.Println("Part two")
	fmt.Printf("Maximum number of energized tiles: %v\n", PartTwo(lines))
}
