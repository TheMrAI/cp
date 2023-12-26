package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Node struct {
	Name      string
	Neighbors map[int]struct{}
}

func ParseInput(lines []string) []Node {
	adjacencyList := []Node{}
	nodeIndex := map[string]int{}

	wordRegexp := regexp.MustCompile("\\w+")

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalf("Expected column in the input: '%s'\n", line)
		}
		node := wordRegexp.FindString(parts[0])
		if len(node) == 0 {
			log.Fatalf("Expected to find a string in input: '%s'\n", parts[0])
		}
		neighbors := wordRegexp.FindAllString(parts[1], -1)
		if len(neighbors) == 0 {
			log.Fatalf("Expected to find strings in input: '%s'\n", parts[1])
		}

		insertIfMissing(nodeIndex, node, &adjacencyList)
		for _, neighbor := range neighbors {
			insertIfMissing(nodeIndex, neighbor, &adjacencyList)
			adjacencyList[nodeIndex[node]].Neighbors[nodeIndex[neighbor]] = struct{}{}
		}
	}

	for id, node := range adjacencyList {
		for neighborId := range node.Neighbors {
			adjacencyList[id].Neighbors[neighborId] = struct{}{}
			adjacencyList[neighborId].Neighbors[id] = struct{}{}
		}
	}

	return adjacencyList
}

func insertIfMissing(nodeIndex map[string]int, nodeName string, adjacencyList *[]Node) {
	_, ok := nodeIndex[nodeName]
	if !ok {
		*adjacencyList = append(*adjacencyList, Node{nodeName, map[int]struct{}{}})
		nodeIndex[nodeName] = len(*adjacencyList) - 1
	}
}

type EdgeId struct {
	A int // we only need to ensure that A < B
	B int
}

// Well this is an annoying problem, like on day 24 part two.
// Had been thinking about for a day, then checked the solutions and half
// of them are simply incorrect (for example they don't work on the example)
// or they aren't solutions, just the employment of some external libraries.
// There was one that seems reasonable, and I had the same idea, but didn't want to
// do it, hoping there was something better. The idea is this. Since we are told
// there are 3 edges that hold the two components of the graphs together, we can be sure
// that all shortest paths from one group to another, must go through these edges.
// Unfortunately we don't know which node is in which group, so the best thing we can do
// is to just generate shortest paths and watch the edge usage (how many times a given edge)
// appears in a shortest path. Delete it, then repeat the whole process, two more times.
// We can't assume that the 3 most used edges will be the ones that need to be cut. It is entirely
// possible that shortest paths will look completely differently after we remove one such important edge.
// (Many online solutions just hand waive this away and their solution just happens to solve their input.)
// I don't want to find all shortest paths between all nodes, because that will take a long time, so I will
// just hand waive it away and say that after processing the third of the nodes we are convinced in the result.
// This will not be strictly true, but at least we aren't lying about the complexities we have glossed over either.
func PartOne(adjacencyList []Node) int {
	for removedEdgeCount := 0; removedEdgeCount < 3; removedEdgeCount++ {
		edgeCountOnShortestPaths := map[EdgeId]int{}
		toProcessNodesCount := len(adjacencyList) / 3

		if len(adjacencyList) < toProcessNodesCount {
			toProcessNodesCount = len(adjacencyList)
		}
		for i := 0; i < toProcessNodesCount; i++ {
			_, paths := DijkstraFrom(adjacencyList, i)
			// trace back the shortest paths from all other edges
			for j := 0; j < len(adjacencyList); j++ {
				TraceShortestPaths(paths, j, &edgeCountOnShortestPaths)
			}
		}
		// Find most common edge
		mostUsedEdge := EdgeId{-1, -1}
		for edge, count := range edgeCountOnShortestPaths {
			if mostUsedEdge.A == -1 && mostUsedEdge.B == -1 {
				mostUsedEdge = edge
				continue
			}
			if edgeCountOnShortestPaths[mostUsedEdge] < count {
				mostUsedEdge = edge
			}
		}
		// Remove edge
		a := mostUsedEdge.A
		b := mostUsedEdge.B
		fmt.Printf("Removing edge: %v - %v\n", adjacencyList[a].Name, adjacencyList[b].Name)
		for i := range adjacencyList {
			if i == a || i == b {
				delete(adjacencyList[i].Neighbors, a)
				delete(adjacencyList[i].Neighbors, b)
			}
		}
	}

	oneGroupCount := countBFT(adjacencyList, 0)
	fmt.Println(oneGroupCount)
	return oneGroupCount * (len(adjacencyList) - oneGroupCount)
}

func DijkstraFrom(adjacencyList []Node, fromId int) ([]int, []int) {
	distances := make([]int, len(adjacencyList))
	predecessors := make([]int, len(adjacencyList))
	marked := make([]bool, len(adjacencyList))
	for i := range distances {
		distances[i] = -1
		predecessors[i] = -1
	}
	distances[fromId] = 0
	predecessors[fromId] = fromId

	for i := 0; i < len(adjacencyList); i++ {
		closest := -1
		for mI := range marked {
			if marked[mI] == true {
				continue
			}
			if distances[mI] != -1 {
				if closest == -1 {
					closest = mI
				}
				if distances[mI] < distances[closest] {
					closest = mI
				}
			}
		}

		for neighborId := range adjacencyList[closest].Neighbors {
			if distances[neighborId] == -1 || distances[closest]+1 < distances[neighborId] {
				distances[neighborId] = distances[closest] + 1
				predecessors[neighborId] = closest
			}
		}
		marked[closest] = true
	}

	return distances, predecessors
}

func TraceShortestPaths(predecessors []int, from int, edgeCountOnShortestPaths *map[EdgeId]int) {
	for predecessors[from] != from {
		next := predecessors[from]
		(*edgeCountOnShortestPaths)[edgeId(from, next)] += 1
		from = next
	}
}

func countBFT(nodes []Node, nodeId int) int {
	count := 0
	seen := map[int]struct{}{}
	next := []int{nodeId}

	for len(next) != 0 {
		check := next[0]
		next = next[1:]
		_, ok := seen[check]
		if ok {
			continue
		}
		seen[check] = struct{}{}

		count++
		for neighbor := range nodes[check].Neighbors {
			next = append(next, neighbor)
		}
	}

	return count
}

func edgeId(a int, b int) EdgeId {
	if a > b {
		return EdgeId{b, a}
	}
	return EdgeId{a, b}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	adjacencyList := ParseInput(lines)
	for k, v := range adjacencyList {
		fmt.Printf("%v - %v\n", k, v)
	}

	fmt.Println("Part one")
	fmt.Printf("The multiple of the two group sizes: %v\n", PartOne(adjacencyList))
}
