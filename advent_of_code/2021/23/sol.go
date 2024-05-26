package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Position struct {
	Occupied  bool
	Neighbors map[*Position]struct{}
	OnMap     Pair
}

type Pair struct {
	I int
	J int
}

const (
	ToMove = iota
	ToPark
	Parked
)

type Amphipod struct {
	Pos    *Position
	Marker byte
	State  int // man Go's enums are just garbage, pure garbage
}

type PositionGroups struct {
	Movable   map[*Position]struct{}
	Forbidden map[*Position]struct{}
	Parking   map[*Position]byte // position to parking type
}

// list of amphipods, then position classification
func ParseInput(lines []string) ([]Amphipod, PositionGroups) {
	amphipods := []Amphipod{}
	// find first '.' field
	i := 0
	j := 0
finder:
	for i = range lines {
		for j = range lines[i] {
			if lines[i][j] == '.' {
				break finder
			}
		}
	}

	// the order of the found arthropods is in pairs of parking types
	parkingSigns := []byte{'A', 'B', 'C', 'D'}

	positionGroups := PositionGroups{Movable: map[*Position]struct{}{},
		Forbidden: map[*Position]struct{}{},
		Parking:   map[*Position]byte{}}
	start := &Position{Occupied: false, Neighbors: map[*Position]struct{}{}, OnMap: Pair{i, j}}

	allPositions := map[Pair]*Position{}
	allPositions[Pair{i, j}] = start
	toCheckCoordinates := []Pair{start.OnMap}

	offsets := []Pair{{0, 1}, {-1, 0}, {1, 0}, {0, -1}}
	for len(toCheckCoordinates) > 0 {
		checkingCoordinate := toCheckCoordinates[len(toCheckCoordinates)-1]
		toCheckCoordinates = toCheckCoordinates[:len(toCheckCoordinates)-1]

		checkedPosition := allPositions[checkingCoordinate]
		for oi := range offsets {
			pos := Pair{checkedPosition.OnMap.I + offsets[oi].I, checkedPosition.OnMap.J + offsets[oi].J}
			r := lines[pos.I][pos.J]
			if r == '#' {
				continue
			}
			address, seen := allPositions[pos]
			if !seen {
				address = &Position{Occupied: false, Neighbors: map[*Position]struct{}{}, OnMap: pos}
				allPositions[pos] = address
				if r != '.' {
					address.Occupied = true
					amphipods = append(amphipods, Amphipod{Pos: address, Marker: r, State: ToMove})
					// register position as parking
					positionGroups.Parking[address] = parkingSigns[(len(amphipods)-1)/2] // placeholder
				}
				toCheckCoordinates = append(toCheckCoordinates, pos)
			}
			checkedPosition.Neighbors[address] = struct{}{}
		}
	}

	// classify the rest of the positions
	for ai := range allPositions {
		address := allPositions[ai]
		_, isParking := positionGroups.Parking[address]
		if isParking {
			continue
		}
		hasParkingNeighbor := false
		for neighbor := range address.Neighbors {
			_, park := positionGroups.Parking[neighbor]
			if park {
				hasParkingNeighbor = true
				break
			}
		}
		if hasParkingNeighbor {
			positionGroups.Forbidden[address] = struct{}{}
		} else {
			positionGroups.Movable[address] = struct{}{}
		}
	}

	return amphipods, positionGroups
}

func PartOne(amphipods []Amphipod, positionGroups PositionGroups, costs map[byte]int) int {
	leastCost := math.MaxInt
	// set the amphipods that are already in their final position immediately to parked
	for ai := range amphipods {
		if amphipods[ai].Marker == positionGroups.Parking[amphipods[ai].Pos] && len(amphipods[ai].Pos.Neighbors) == 1 {
			amphipods[ai].State = Parked
		}
	}
	optimalOrganizationCost(amphipods, positionGroups, 0, &leastCost, costs)
	return leastCost
}

func optimalOrganizationCost(amphipods []Amphipod, positionGroups PositionGroups, currentCost int, bestCost *int, costs map[byte]int) {
	// Debugging
	printState(amphipods)
	// time.Sleep(1 * time.Second)
	if currentCost > *bestCost {
		return
	}
	// fmt.Println(amphipods)
	if allParked(amphipods) {
		if *bestCost > currentCost {
			*bestCost = currentCost
		}
		return
	}

	for ai := range amphipods {
		var choices []Target
		if amphipods[ai].State == Parked {
			continue
		} else if amphipods[ai].State == ToMove {
			choices = reacheableTragets(amphipods, ai, ToPark, positionGroups)
		} else { // ToPark
			choices = reacheableTragets(amphipods, ai, Parked, positionGroups)
		}

		if len(choices) == 0 {
			continue // in this configuration the amphipod may not take a move right now, but could do so later, perhaps
		}

		for ci := range choices {
			// move
			originalPos := amphipods[ai].Pos
			originalPos.Occupied = false
			originalState := amphipods[ai].State

			amphipods[ai].Pos = choices[ci].Pos
			amphipods[ai].Pos.Occupied = true
			if amphipods[ai].State == ToMove {
				amphipods[ai].State = ToPark
			} else {
				amphipods[ai].State = Parked
			}
			optimalOrganizationCost(amphipods, positionGroups, currentCost+choices[ci].Distance*costs[amphipods[ai].Marker], bestCost, costs)
			amphipods[ai].Pos.Occupied = false

			// restore amphipod on backtracking
			amphipods[ai].Pos = originalPos
			originalPos.Occupied = true
			amphipods[ai].State = originalState
		}
	}
}

type Target struct {
	Pos      *Position
	Distance int
}

type DFS struct {
	Prev     *Position
	Pos      *Position
	Distance int
}

func reacheableTragets(amphipods []Amphipod, ai int, targetState int, positionGroups PositionGroups) []Target {
	targets := []Target{}

	toProcess := []DFS{{nil, amphipods[ai].Pos, 0}}
	for len(toProcess) > 0 {
		processing := toProcess[len(toProcess)-1]
		toProcess = toProcess[:len(toProcess)-1]

		// check if valid target
		if targetState == ToPark {
			_, ok := positionGroups.Movable[processing.Pos]
			if ok {
				targets = append(targets, Target{processing.Pos, processing.Distance})
			}
		} else if targetState == Parked {
			parkingSpotType, ok := positionGroups.Parking[processing.Pos]
			if ok && parkingSpotType == amphipods[ai].Marker {
				// now the spot either has to have only one neighbor in which case it is valid
				if len(processing.Pos.Neighbors) == 1 {
					targets = append(targets, Target{processing.Pos, processing.Distance})
				} else if len(processing.Pos.Neighbors) == 2 {
					// or it has to have two neighbors and the other parking spot must be occupied by the same type
					var otherParkingSpot *Position
					for pos := range processing.Pos.Neighbors {
						_, isParking := positionGroups.Parking[pos]
						if isParking {
							otherParkingSpot = pos
						}
					}

					if !otherParkingSpot.Occupied {
						continue
					}
					mismatched := false
					for check := range amphipods {
						if amphipods[check].Pos == otherParkingSpot {
							if amphipods[check].Marker != amphipods[ai].Marker {
								mismatched = true
							}
							break
						}
					}
					if !mismatched {
						targets = append(targets, Target{processing.Pos, processing.Distance})
					}
				}
			}
		}

		for neighbor := range processing.Pos.Neighbors {
			if neighbor.Occupied || neighbor == processing.Prev {
				continue
			}
			toProcess = append(toProcess, DFS{processing.Pos, neighbor, processing.Distance + 1})
		}
	}

	return targets
}

// doesn't validate if it is correctly parked
func allParked(amphipods []Amphipod) bool {
	for ai := range amphipods {
		if amphipods[ai].State != Parked {
			return false
		}
	}
	return true
}

func printState(amphipods []Amphipod) {
	situation := [][]byte{[]byte("#############"),
		[]byte("#...........#"),
		[]byte("###.#.#.#.###"),
		[]byte("  #.#.#.#.#  "),
		[]byte("  #########  "),
	}
	for ai := range amphipods {
		onMap := amphipods[ai].Pos.OnMap
		situation[onMap.I][onMap.J] = amphipods[ai].Marker
	}
	for i := range situation {
		for j := range situation[i] {
			fmt.Printf("%c", situation[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	amphipods, positionGroupings := ParseInput(lines)
	for i := range amphipods {
		fmt.Printf("%c - %p - %v\n", amphipods[i].Marker, amphipods[i].Pos, amphipods[i].Pos.OnMap)
		for n := range amphipods[i].Pos.Neighbors {
			fmt.Printf("\t%p - %v\n", n, n.OnMap)
		}
	}

	for a, t := range positionGroupings.Parking {
		fmt.Printf("%c - %p\n", t, a)
	}

	costs := map[byte]int{
		'A': 1,
		'B': 10,
		'C': 100,
		'D': 1000,
	}
	fmt.Printf("The least energy to organize the amphipods is: %v\n", PartOne(amphipods, positionGroupings, costs))
}
