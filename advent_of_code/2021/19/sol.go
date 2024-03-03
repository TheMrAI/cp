package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Vec3 [3]int

func (lhs Vec3) Add(rhs Vec3) Vec3 {
	return Vec3{lhs[0] + rhs[0], lhs[1] + rhs[1], lhs[2] + rhs[2]}
}

func (lhs Vec3) Sub(rhs Vec3) Vec3 {
	return Vec3{lhs[0] - rhs[0], lhs[1] - rhs[1], lhs[2] - rhs[2]}
}

func (lhs Vec3) Dot(rhs Vec3) int {
	return lhs[0]*rhs[0] + lhs[1]*rhs[1] + lhs[2]*rhs[2]
}

func (lhs Vec3) Cross(rhs Vec3) Vec3 {
	return Vec3{lhs[1]*rhs[2] - lhs[2]*rhs[1], lhs[2]*rhs[0] - lhs[0]*rhs[2], lhs[0]*rhs[1] - lhs[1]*rhs[0]}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v[0], -v[1], -v[2]}
}

type Scanner struct {
	Id      int
	Beacons map[Vec3]struct{}
}

func ParseInput(lines []string) []Scanner {
	scanners := []Scanner{}

	for i := 0; i < len(lines); {
		scanner, nextI := parseScanner(lines, i)
		scanners = append(scanners, scanner)
		nextI++ // skip empty line
		i = nextI
	}

	return scanners
}

func parseScanner(lines []string, from int) (Scanner, int) {
	scannerRexexp := regexp.MustCompile(`--- scanner (\d+) ---`)
	coordinatesRegexp := regexp.MustCompile(`(-?\d+),(-?\d+),(-?\d+)`)

	scannerIdText := scannerRexexp.FindStringSubmatch(lines[from])
	fmt.Println(scannerIdText)
	if len(scannerIdText) != 2 {
		log.Fatalf("Failed to parse scanner id on input: `%v`", lines[from])
	}
	scannerId, err := strconv.Atoi(scannerIdText[1])
	if err != nil {
		log.Fatalf("Integer conversion failed on input: `%v`", scannerIdText[0])
	}
	from++

	beacons := map[Vec3]struct{}{}
	i := from
	for ; i < len(lines) && len(lines[i]) != 0; i++ {
		coordinatesText := coordinatesRegexp.FindStringSubmatch(lines[i])
		if len(coordinatesText) != 4 {
			log.Fatalf("Expected 3 coordinates, but got input: `%v`", lines[i])
		}
		coordinates := make([]int, 0, 3)
		for cI := 1; cI < len(coordinatesText); cI++ {
			coordinate, err := strconv.Atoi(coordinatesText[cI])
			if err != nil {
				log.Fatalf("Integer conversion failed on input: `%v`", coordinatesText[cI])
			}
			coordinates = append(coordinates, coordinate)
		}
		beacons[Vec3{coordinates[0], coordinates[1], coordinates[2]}] = struct{}{}
	}
	return Scanner{scannerId, beacons}, i
}

func PartOne(scanners []Scanner) int {
	unifiedScanners := unifyScanners(scanners)
	fmt.Printf("Unified scanners: %v\n", unifiedScanners)
	distinctBeacons := map[Vec3]struct{}{}
	for i := range unifiedScanners {
		for beacon := range unifiedScanners[i].Beacons {
			distinctBeacons[beacon] = struct{}{}
		}
	}
	return len(distinctBeacons)
}

func unifyScanners(scanners []Scanner) []Scanner {
	unifiedScanners := []Scanner{scanners[0]} // we start coordinate system from the zeroth scanner
	resolverQueue := []int{0}                 // all the scanners we have resolved, and now can be used to resolve others
	resolvedScanners := map[int]struct{}{}
	resolvedScanners[0] = struct{}{}

	for len(resolvedScanners) < len(scanners) {
		if len(resolverQueue) == 0 {
			fmt.Println("Oh-oh no more resolves available")
			break
		}
		resolveUsingScannerId := resolverQueue[0]
		resolverQueue = resolverQueue[1:]

		for toResolveScannerId := 0; toResolveScannerId < len(scanners); toResolveScannerId++ {
			_, resolved := resolvedScanners[toResolveScannerId]
			if resolved {
				continue
			}
			fmt.Printf("Trying to resolve %v using %v\n", toResolveScannerId, resolveUsingScannerId)
			resolved, resolvedBeaconPositions := tryResolveBeaconPositions(scanners[resolveUsingScannerId].Beacons, scanners[toResolveScannerId].Beacons)
			if resolved {
				fmt.Println("Resolved")
				// fmt.Printf("Resolved beacon positions: %v\n", resolvedBeaconPositions)
				resolverQueue = append(resolverQueue, toResolveScannerId)
				resolvedScanners[toResolveScannerId] = struct{}{}
				unifiedScanners = append(unifiedScanners, Scanner{scanners[toResolveScannerId].Id, resolvedBeaconPositions})
			}
		}
	}

	return unifiedScanners
}

func tryResolveBeaconPositions(resolveUsingBeacons map[Vec3]struct{}, toResolveBeacons map[Vec3]struct{}) (bool, map[Vec3]struct{}) {
	transformers := getTransformers()
	for i := 0; i < len(transformers); i++ {
		remappedBeacons := remapCoordinates(toResolveBeacons, transformers[i])
		for pointA := range resolveUsingBeacons {
			for pointB := range remappedBeacons {
				delta := pointA.Sub(pointB)
				adjustedBeacons := adjustCoordinates(remappedBeacons, delta)
				// fmt.Printf("Adjusted beacon candidates: %v\n", adjustedBeacons)
				scannerResolved := atLeastTwelveOverlap(resolveUsingBeacons, adjustedBeacons)
				if scannerResolved {
					return true, adjustedBeacons
				}
			}
		}
	}
	return false, map[Vec3]struct{}{}
}

type Transformer [3][3]int

func (transformer Transformer) Mul(rhs Vec3) Vec3 {
	rowOne := Vec3{transformer[0][0], transformer[0][1], transformer[0][2]}
	rowTwo := Vec3{transformer[1][0], transformer[1][1], transformer[1][2]}
	rowThree := Vec3{transformer[2][0], transformer[2][1], transformer[2][2]}
	return Vec3{rowOne.Dot(rhs), rowTwo.Dot(rhs), rowThree.Dot(rhs)}
}

func getTransformers() []Transformer {
	transformers := make([]Transformer, 0, 24)

	// generate all transformation matrices, by setting the first base, then picking one orthogonal
	// direction for the second base, then calculating the third using cross product
	cardinalDirections := []Vec3{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}
	for _, baseOne := range cardinalDirections {
		transformers = append(transformers, genTransformMatrices(baseOne, cardinalDirections)...)
	}

	return transformers
}

func genTransformMatrices(firstBase Vec3, secondBasesCandidates []Vec3) []Transformer {
	transformers := make([]Transformer, 0, 4)
	for _, secondBase := range secondBasesCandidates {
		// check if firstBase and secondBase isn't on the same axis
		onSameAxis := false
		for i := 0; i < 3; i++ {
			if (firstBase[i] == 1 || firstBase[i] == -1) && (firstBase[i] == secondBase[i] || firstBase[i] == -secondBase[i]) {
				onSameAxis = true
				break
			}
		}
		if onSameAxis {
			continue
		}

		thirdBase := firstBase.Cross(secondBase)
		transformer := Transformer{
			{firstBase[0], secondBase[0], thirdBase[0]},
			{firstBase[1], secondBase[1], thirdBase[1]},
			{firstBase[2], secondBase[2], thirdBase[2]}}
		transformers = append(transformers, transformer)
	}
	return transformers
}

func remapCoordinates(toRemapBeacons map[Vec3]struct{}, transformer Transformer) map[Vec3]struct{} {
	remapped := map[Vec3]struct{}{}
	for k := range toRemapBeacons {
		remapped[transformer.Mul(k)] = struct{}{}
	}
	return remapped
}

func adjustCoordinates(toRemap map[Vec3]struct{}, delta Vec3) map[Vec3]struct{} {
	adjusted := map[Vec3]struct{}{}
	for v := range toRemap {
		adjusted[v.Add(delta)] = struct{}{}
	}

	return adjusted
}

func atLeastTwelveOverlap(beaconSetLHS, beaconSetRHS map[Vec3]struct{}) bool {
	overlapCount := 0
	for k := range beaconSetLHS {
		_, contains := beaconSetRHS[k]
		if contains {
			overlapCount++
		}
	}
	return overlapCount >= 12
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	scanners := ParseInput(lines)
	fmt.Printf("The number of distinct beacons: %v\n", PartOne(scanners))
}
