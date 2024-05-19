package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
)

// both end inclusive range
type Range struct {
	F int
	L int
}

func (r Range) Overlap(other Range) bool {
	return (r.F <= other.F && other.F <= r.L) || (r.F <= other.L && other.L <= r.L) || (r.F > other.F && other.L > r.L)
}

func (r Range) Contains(other Range) bool {
	return (r.F <= other.F && other.L <= r.L)
}

func (r Range) Length() int {
	return 1 + r.L - r.F
}

// Cut the ranges into pieces based on the boundary points, if they overlap in any way.
// It doesn't matter which piece do or do not overlap.
// May produce -1 length Ranges, which can, be ignored, but they help us identify the overlapping segment.
// Either there is exactly one Range result, when we have a perfect overlap, or there will always be 3 segments,
// in which case always the middle one has to be the overlap.
func (r Range) Chop(other Range) []Range {
	if !r.Overlap(other) {
		return []Range{}
	}
	points := []int{r.F, r.L, other.F, other.L}
	slices.Sort(points)
	// if we know two ranges overlap and if the points are put in order, we can say with utter surety
	// that the overlapping segment will be [point[1], point[2]], regardless from which range the actual point is
	// coming from
	return []Range{{points[0], points[1] - 1}, {points[1], points[2]}, {points[2] + 1, points[3]}}
}

type Instruction struct {
	SwitchOn bool
	X        Range
	Y        Range
	Z        Range
}

func convertToInt(input string) int {
	val, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalf("Integer parsing failed on input: '%s'", input)
	}
	return val
}

func ParseInput(lines []string) []Instruction {
	instructions := []Instruction{}

	instructionRegexp := regexp.MustCompile(`(on|off) x=(-?\d+)..(-?\d+),y=(-?\d+)..(-?\d+),z=(-?\d+)..(-?\d+)`)

	for _, line := range lines {
		matches := instructionRegexp.FindStringSubmatch(line)
		instruction := Instruction{}
		if matches[1] == "on" {
			instruction.SwitchOn = true
		} else {
			instruction.SwitchOn = false
		}
		instruction.X = Range{convertToInt(matches[2]), convertToInt(matches[3])}
		instruction.Y = Range{convertToInt(matches[4]), convertToInt(matches[5])}
		instruction.Z = Range{convertToInt(matches[6]), convertToInt(matches[7])}

		instructions = append(instructions, instruction)
	}

	return instructions
}

type Cuboid struct {
	X Range
	Y Range
	Z Range
}

func (v *Cuboid) Volume() int {
	return v.X.Length() * v.Y.Length() * v.Z.Length()
}

func (v Cuboid) Overlap(rhs Cuboid) bool {
	return v.X.Overlap(rhs.X) && v.Y.Overlap(rhs.Y) && v.Z.Overlap(rhs.Z)
}

func (v Cuboid) Intersect(rhs Cuboid) (bool, Cuboid) {
	x := v.X.Chop(rhs.X)
	y := v.Y.Chop(rhs.Y)
	z := v.Z.Chop(rhs.Z)

	if len(x) == 0 || len(y) == 0 || len(z) == 0 {
		return false, Cuboid{}
	}

	overlappingX := x[0]
	if len(x) == 3 {
		overlappingX = x[1]
	}
	overlappingY := y[0]
	if len(y) == 3 {
		overlappingY = y[1]
	}
	overlappingZ := z[0]
	if len(z) == 3 {
		overlappingZ = z[1]
	}
	return true, Cuboid{overlappingX, overlappingY, overlappingZ}
}

// It will chop up cuboid v using rhs as a guide to identify the pieces
// of v that are outside of rhs, and within.
func (v Cuboid) SubdivideWith(rhs Cuboid) (bool, []Cuboid, []Cuboid) {
	x := v.X.Chop(rhs.X)
	xPieces := []Range{}
	for i := range x {
		if x[i].Length() <= 0 || !v.X.Contains(x[i]) {
			continue
		}
		xPieces = append(xPieces, x[i])
	}
	y := v.Y.Chop(rhs.Y)
	yPieces := []Range{}
	for i := range y {
		if y[i].Length() <= 0 || !v.Y.Contains(y[i]) {
			continue
		}
		yPieces = append(yPieces, y[i])
	}
	z := v.Z.Chop(rhs.Z)
	zPieces := []Range{}
	for i := range z {
		if z[i].Length() <= 0 || !v.Z.Contains(z[i]) {
			continue
		}
		zPieces = append(zPieces, z[i])
	}

	if len(xPieces) == 0 || len(yPieces) == 0 || len(zPieces) == 0 {
		return false, []Cuboid{}, []Cuboid{}
	}

	intersection := []Cuboid{}
	difference := []Cuboid{}
	for xi := range xPieces {
		for yi := range yPieces {
			for zi := range zPieces {
				cuboid := Cuboid{xPieces[xi], yPieces[yi], zPieces[zi]}
				if rhs.Overlap(cuboid) {
					intersection = append(intersection, cuboid)
				} else {
					difference = append(difference, cuboid)
				}
			}
		}
	}
	return true, intersection, difference
}

type Octree struct {
	on       bool
	children []Octree
	cuboid    Cuboid
}

func (o Octree) OnCount() int {
	if o.leaf() {
		if o.on {
			return o.cuboid.X.Length() * o.cuboid.Y.Length() * o.cuboid.Z.Length()
		} else {
			return 0
		}
	}
	sum := 0
	for i := range o.children {
		sum += o.children[i].OnCount()
	}
	return sum
}

func (o *Octree) Modify(v Cuboid, turnOn bool) {
	if o.cuboid == v {
		o.on = turnOn
		o.children = []Octree{}
		return
	}
	if o.leaf() {
		o.subdivide()
	}
	for i := range o.children {
		ok, subVoxel := o.children[i].cuboid.Intersect(v)
		if !ok {
			continue
		}
		o.children[i].Modify(subVoxel, turnOn)
	}
}

func (o *Octree) subdivide() {
	xF := o.cuboid.X.F
	xL := o.cuboid.X.L
	yF := o.cuboid.Y.F
	yL := o.cuboid.Y.L
	zF := o.cuboid.Z.F
	zL := o.cuboid.Z.L

	xMid := xF + (o.cuboid.X.Length() / 2) - 1
	yMid := yF + (o.cuboid.Y.Length() / 2) - 1
	zMid := zF + (o.cuboid.Z.Length() / 2) - 1

	bottomCloseLeft := Octree{o.on, []Octree{}, Cuboid{Range{xF, xMid}, Range{yF, yMid}, Range{zF, zMid}}}
	bottomCloseRight := Octree{o.on, []Octree{}, Cuboid{Range{xMid + 1, xL}, Range{yF, yMid}, Range{zF, zMid}}}
	bottomFarLeft := Octree{o.on, []Octree{}, Cuboid{Range{xF, xMid}, Range{yMid + 1, yL}, Range{zF, zMid}}}
	bottomFarRight := Octree{o.on, []Octree{}, Cuboid{Range{xMid + 1, xL}, Range{yMid + 1, yL}, Range{zF, zMid}}}
	topCloseLeft := Octree{o.on, []Octree{}, Cuboid{Range{xF, xMid}, Range{yF, yMid}, Range{zMid + 1, zL}}}
	topCloseRight := Octree{o.on, []Octree{}, Cuboid{Range{xMid + 1, xL}, Range{yF, yMid}, Range{zMid + 1, zL}}}
	topFarLeft := Octree{o.on, []Octree{}, Cuboid{Range{xF, xMid}, Range{yMid + 1, yL}, Range{zMid + 1, zL}}}
	topFarRight := Octree{o.on, []Octree{}, Cuboid{Range{xMid + 1, xL}, Range{yMid + 1, yL}, Range{zMid + 1, zL}}}

	o.children = []Octree{
		bottomCloseLeft,
		bottomCloseRight,
		bottomFarLeft,
		bottomFarRight,
		topCloseLeft,
		topCloseRight,
		topFarLeft,
		topFarRight,
	}
}

func (o Octree) leaf() bool {
	return len(o.children) == 0
}

type Point struct {
	X int
	Y int
	Z int
}

// The Octree based solution is more than sufficient for part one, but then again a brute force solution worked as well.
// For part two it has no chance of working though, because big cubes would require so much space partitioning that
// we easily run out of memory even for small tests.
// Leaving this part using the Octree implementation, because this was the first time I have used the technique.
func PartOne(instructions []Instruction) int {
	xMin := math.MaxInt
	xMax := math.MinInt
	yMin := math.MaxInt
	yMax := math.MinInt
	zMin := math.MaxInt
	zMax := math.MinInt
	for _, instruction := range instructions {
		if instruction.X.F < -50 || instruction.Y.F < -50 || instruction.Z.F < -50 ||
			instruction.X.F > 50 || instruction.Y.F > 50 || instruction.Z.F > 50 {
			continue
		}
		xMin = min(xMin, instruction.X.F)
		xMax = max(xMax, instruction.X.L)
		yMin = min(yMin, instruction.Y.F)
		yMax = max(yMax, instruction.Y.L)
		zMin = min(zMin, instruction.Z.F)
		zMax = max(zMax, instruction.Z.L)
	}
	octree := Octree{false, []Octree{}, Cuboid{Range{xMin, xMax}, Range{yMin, yMax}, Range{zMin, zMax}}}
	for _, instruction := range instructions {
		if instruction.X.F < -50 || instruction.Y.F < -50 || instruction.Z.F < -50 ||
			instruction.X.F > 50 || instruction.Y.F > 50 || instruction.Z.F > 50 {
			continue
		}

		cuboid := Cuboid{instruction.X, instruction.Y, instruction.Z}
		octree.Modify(cuboid, instruction.SwitchOn)
	}
	return octree.OnCount()
}

// A simple algorithm where we just cut turned on cubes into smaller cubes if the instruction
// would introduce a cube that would collide. Add the differences and discard the intersection.
// Then if the instruction was to turn on the piece, we add it back.
func PartTwo(instructions []Instruction) int {
	cuboids := map[Cuboid]struct{}{}
	for i := range instructions {
		tmp := map[Cuboid]struct{}{}
		candidate := Cuboid{instructions[i].X, instructions[i].Y, instructions[i].Z}
		for v := range cuboids {
			ok, _, differences := v.SubdivideWith(candidate)
			if !ok {
				tmp[v] = struct{}{}
				continue
			}
			for di := range differences {
				tmp[differences[di]] = struct{}{}
			}
		}
		if instructions[i].SwitchOn {
			tmp[candidate] = struct{}{}
		}
		cuboids = tmp
	}

	onCount := 0
	for cuboid := range cuboids {
		onCount += cuboid.Volume()
	}
	return onCount
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	instructions := ParseInput(lines)

	fmt.Printf("The number of ON cubes after the initialization sequence: %v\n", PartOne(instructions))
	fmt.Printf("The number of ON cubes after the reboot sequence: %v\n", PartTwo(instructions))
}
