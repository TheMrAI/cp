package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Triplet struct {
	X int
	Y int
	Z int
}

type Hail struct {
	Pos Triplet
	Vel Triplet
}

func ParseInput(lines []string) []Hail {
	hail := []Hail{}

	numberRegexp := regexp.MustCompile("-?\\d+")

	for _, line := range lines {
		parts := strings.Split(line, "@")
		if len(parts) != 2 {
			log.Fatalf("Unexpected number of sections in input line: '%v'\n", line)
		}
		position := parseNumbers(parts[0], numberRegexp)
		velocity := parseNumbers(parts[1], numberRegexp)
		if len(position) != 3 && len(velocity) != 3 {
			log.Fatalf("Unexpected number of position/velocity data: %v, %v\n", position, velocity)
		}
		hail = append(hail, Hail{Triplet{position[0], position[1], position[2]}, Triplet{velocity[0], velocity[1], velocity[2]}})
	}

	return hail
}

func parseNumbers(input string, numberRegexp *regexp.Regexp) []int {
	numberSections := numberRegexp.FindAllString(input, -1)
	numbers := []int{}
	for _, text := range numberSections {
		number, err := strconv.Atoi(text)
		if err != nil {
			fmt.Printf("Integer conversion failed on input: '%v'\n", text)
		}
		numbers = append(numbers, number)
	}
	return numbers
}

type LineEq struct {
	M float64
	C float64
}

func hailReachPointInFuture(hail Hail, x float64, y float64) bool {
	xDir := (x - float64(hail.Pos.X)) / float64(hail.Vel.X)
	yDir := (y - float64(hail.Pos.Y)) / float64(hail.Vel.Y)
	return xDir > 0 && yDir > 0
}

func PartOne(hail []Hail, rangeStart, rangeLast float64) int {
	lineEquations := make([]LineEq, len(hail))
	for i := range hail {
		lineEquations[i] = hailToLineEquation(hail[i])
	}

	collisionCountInFuture := 0
	for i := 0; i < len(lineEquations)-1; i++ {
		for j := i + 1; j < len(lineEquations); j++ {
			x, y := intersectionPoint(lineEquations[i], lineEquations[j])
			// parallel lines
			if (x == math.Inf(-1) || x == math.Inf(1)) && (y == math.Inf(-1) || y == math.Inf(1)) {
				continue
			}
			if rangeStart <= x && x <= rangeLast && rangeStart <= y && y <= rangeLast && hailReachPointInFuture(hail[i], x, y) && hailReachPointInFuture(hail[j], x, y) {
				collisionCountInFuture++
			}
		}
	}
	return collisionCountInFuture
}

// only using X, Y coordinates
func hailToLineEquation(hail Hail) LineEq {
	m := float64(hail.Vel.Y) / float64(hail.Vel.X)
	c := float64(hail.Pos.Y) - m*float64(hail.Pos.X)
	return LineEq{m, c}
}

func intersectionPoint(a, b LineEq) (float64, float64) {
	x := (b.C - a.C) / (a.M - b.M)
	y := a.M*x + a.C
	return x, y
}

// My first idea was to get two planes, from a 4 initial velocity vectors,
// intersecting them to get the line all hail must
// cross at one point for the possibility of having a solution,
// then find an integer value on it, a velocity vector and done.
// This sounds nice and all, the only problem is that I can't make the planes as there
// is simply not enough information to make them.
//
// The next option I have seen others do. Observe that if our stone
// starts from Ps with velocity Vs then any other stone i with Pi, Vi
// will collide with it at time t then Ps+Vs*t = Pi+Vi*t. This is a bit
// nasty having 3 unknowns, but if we rearrange it:
// Ps = Pi + (Vi-Vs)*t, then we have modified the problem and are now stating:
// If our stone was made stationary and instead we removed its velocity from all other
// hail, under what Vs and t will all produce the same Ps? It is easier to see
// that know we only have 2 unknowns. Better yet, we don't care about t at all
// only that the two lines all collide at the same point for a given w.
// So in essence we have to search for such a w, for which all modified hail stone
// trajectories will collide at the same point. That point is our Ps.
// This still means that we have three dimensions (the components of the Vs) vector that
// can vary. This is still very annoying. But lastly if we realize there can be exactly one
// solution (as the exercise told us) we can realize that we can throw away one of the dimensions
// , for example z and still get the same result. (To aid understanding, just try thinking of the 3D vectors, as if
// instead of the original vectors, we were working with their projections from a top-down perspective.)
// This we already solved in part 1.
func PartTwo(hails []Hail) int {
	xVelocities := identifyVelocities(hails, 0)
	yVelocities := identifyVelocities(hails, 1)
	zVelocities := identifyVelocities(hails, 2)

	fmt.Printf("Matching X: %v\n", xVelocities)
	fmt.Printf("Matching Y: %v\n", yVelocities)
	fmt.Printf("Matching Z: %v\n", zVelocities)

	xVelocity := 0
	for v := range xVelocities {
		xVelocity = v
		break
	}
	yVelocity := 0
	for v := range yVelocities {
		yVelocity = v
		break
	}
	zVelocity := 0
	for v := range zVelocities {
		zVelocity = v
		break
	}
	// Now intersect two hail lines
	hailOne := hails[0]
	hailTwo := hails[1]

	ma := float64(hailOne.Vel.Y-yVelocity) / float64(hailOne.Vel.X-xVelocity)
	mb := float64(hailTwo.Vel.Y-yVelocity) / float64(hailTwo.Vel.X-xVelocity)

	ca := float64(hailOne.Pos.Y) - (ma * float64(hailOne.Pos.X))
	cb := float64(hailTwo.Pos.Y) - (mb * float64(hailTwo.Pos.X))

	posX := (cb - ca) / (ma - mb)
	posY := ma*posX + ca
	t := (posX - float64(hailOne.Pos.X)) / float64(hailOne.Vel.X-xVelocity)
	posZ := float64(hailOne.Pos.Z) + float64(hailOne.Vel.Z-zVelocity)*t

	fmt.Println(posX, posY, posZ)
	return int(posX + posY + posZ)
}

func identifyVelocities(hails []Hail, dimension int) map[int]struct{} {
	velocityMatcher := map[int][]Hail{}
	for _, hail := range hails {
		velocityComponent := velByDim(hail, dimension)
		_, ok := velocityMatcher[velocityComponent]
		if !ok {
			velocityMatcher[velocityComponent] = []Hail{}
		}
		velocityMatcher[velocityComponent] = append(velocityMatcher[velocityComponent], hail)
	}

	velocityCandidates := map[int]struct{}{}
	for velocity, matchedHails := range velocityMatcher {
		if len(matchedHails) > 1 {
			for i := 0; i < len(matchedHails)-1; i++ {
				for j := 0; j < len(matchedHails); j++ {
					validVelocities := map[int]struct{}{}
					distance := posByDim(matchedHails[i], dimension) - posByDim(matchedHails[j], dimension)
					for v := -1000; v <= 1000; v++ {
						deltaV := v - velocity
						if deltaV == 0 {
							continue
						}
						if distance%deltaV == 0 {
							validVelocities[v] = struct{}{}
						}
					}
					if len(velocityCandidates) == 0 {
						velocityCandidates = validVelocities
						continue
					}
					reducedValids := map[int]struct{}{}
					for v := range velocityCandidates {
						_, ok := validVelocities[v]
						if ok {
							reducedValids[v] = struct{}{}
						}
					}
					velocityCandidates = reducedValids
				}
			}
		}
	}
	return velocityCandidates
}

func posByDim(h Hail, dim int) int {
	switch dim {
	case 0:
		return h.Pos.X
	case 1:
		return h.Pos.Y
	case 2:
		return h.Pos.Z
	}
	return 0
}

func velByDim(h Hail, dim int) int {
	switch dim {
	case 0:
		return h.Vel.X
	case 1:
		return h.Vel.Y
	case 2:
		return h.Vel.Z
	}
	return 0
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	hail := ParseInput(lines)

	fmt.Println("Part one")
	fmt.Printf("The number of hailstones that will collide between 200000000000000 - 400000000000000 based on their X,Y coordinates: %v\n", PartOne(hail, 200000000000000, 400000000000000))
	fmt.Println("Part two")
	fmt.Printf("The sum of x, y, z starting coordinate of the stone that will smash all hail is: %v\n", PartTwo(hail))
}
