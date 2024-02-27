package main

import (
	"testing"
)

func SnailNumsEqual(lhs, rhs []SnailNum) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}

func TestSnailNumAdd(t *testing.T) {
	tests := []struct {
		lhs  []SnailNum
		rhs  []SnailNum
		want []SnailNum
	}{
		{
			[]SnailNum{{1, 0}, {2, 0}},
			[]SnailNum{{3, 1}, {4, 1}, {5, 0}},
			[]SnailNum{{1, 1}, {2, 1}, {3, 2}, {4, 2}, {5, 1}},
		},
		{
			[]SnailNum{{4, 3}, {3, 3}, {4, 2}, {4, 1}, {7, 1}, {8, 3}, {4, 3}, {9, 2}},
			[]SnailNum{{1, 0}, {1, 0}},
			[]SnailNum{{4, 4}, {3, 4}, {4, 3}, {4, 2}, {7, 2}, {8, 4}, {4, 4}, {9, 3}, {1, 1}, {1, 1}},
		},
	}
	for _, test := range tests {
		if got := Add(test.lhs, test.rhs); !SnailNumsEqual(got, test.want) {
			t.Errorf("TestSnailNumAdd(%v, %v) => %v; want %v", test.lhs, test.rhs, got, test.want)
		}
	}
}

func TestSnailExplode(t *testing.T) {
	tests := []struct {
		input     []SnailNum
		explodeAt int
		want      []SnailNum
	}{
		{
			[]SnailNum{{9, 4}, {8, 4}, {1, 3}, {2, 2}, {3, 1}, {4, 0}},
			0,
			[]SnailNum{{0, 3}, {9, 3}, {2, 2}, {3, 1}, {4, 0}},
		},
		{
			[]SnailNum{{7, 0}, {6, 1}, {5, 2}, {4, 3}, {3, 4}, {2, 4}},
			4,
			[]SnailNum{{7, 0}, {6, 1}, {5, 2}, {7, 3}, {0, 3}},
		},
		{
			[]SnailNum{{3, 1}, {2, 2}, {1, 3}, {7, 4}, {3, 4}, {6, 1}, {5, 2}, {4, 3}, {3, 4}, {2, 4}},
			3,
			[]SnailNum{{3, 1}, {2, 2}, {8, 3}, {0, 3}, {9, 1}, {5, 2}, {4, 3}, {3, 4}, {2, 4}},
		},
		{
			[]SnailNum{{3, 1}, {2, 2}, {8, 3}, {0, 3}, {9, 1}, {5, 2}, {4, 3}, {3, 4}, {2, 4}},
			7,
			[]SnailNum{{3, 1}, {2, 2}, {8, 3}, {0, 3}, {9, 1}, {5, 2}, {7, 3}, {0, 3}},
		},
	}
	for _, test := range tests {
		if got := explode(test.input, test.explodeAt); !SnailNumsEqual(got, test.want) {
			t.Errorf("explode(%v, %v) => %v; want %v", test.input, test.explodeAt, got, test.want)
		}
	}
}

func TestSnailSplit(t *testing.T) {
	tests := []struct {
		input     []SnailNum
		splitAt int
		want      []SnailNum
	}{
		{
			[]SnailNum{{1, 0}, {10, 0}},
			1,
			[]SnailNum{{1, 0}, {5, 1}, {5, 1}},
		},
		{
			[]SnailNum{{1, 0}, {11, 0}},
			1,
			[]SnailNum{{1, 0}, {5, 1}, {6, 1}},
		},
		{
			[]SnailNum{{10, 0}, {1, 0}},
			0,
			[]SnailNum{{5, 1}, {5, 1}, {1, 0}},
		},
		{
			[]SnailNum{{11, 0}, {1, 0}},
			0,
			[]SnailNum{{5, 1}, {6, 1}, {1, 0}},
		},
	}
	for _, test := range tests {
		if got := split(test.input, test.splitAt); !SnailNumsEqual(got, test.want) {
			t.Errorf("split(%v, %v) => \n%v; want %v", test.input, test.splitAt, got, test.want)
		}
	}
}

func TestSnailMagnitude(t *testing.T) {
	tests := []struct {
		input []SnailNum
		want  int
	}{
		{
			[]SnailNum{{9, 1}, {1, 1}, {1, 1}, {9, 1}},
			129,
		},
		{
			[]SnailNum{{1, 1}, {2, 1}, {3, 2}, {4, 2}, {5, 1}},
			143,
		},
		{
			[]SnailNum{{0, 3}, {7, 3}, {4, 2}, {7, 3}, {8, 3}, {6, 3}, {0, 3}, {8, 1}, {1, 1}},
			1384,
		},
		{
			[]SnailNum{{5, 3}, {0, 3}, {7, 3}, {4, 3}, {5, 2}, {5, 2}, {6, 1}, {6, 1}},
			1137,
		},
		{
			[]SnailNum{{8, 3}, {7, 3}, {7, 3}, {7, 3}, {8, 3}, {6, 3}, {7, 3}, {7, 3}, {0, 3}, {7, 3}, {6, 3}, {6, 3}, {8, 2}, {7, 2}},
			3488,
		},
		{
			[]SnailNum{{6, 3}, {6, 3}, {7, 3}, {6, 3}, {7, 3}, {7, 3}, {7, 3}, {0, 3}, {7, 3}, {7, 3}, {7, 3}, {7, 3}, {7, 3}, {8, 3}, {9, 3}, {9, 3}},
			4140,
		},
	}
	for _, test := range tests {
		if got := Magnitude(test.input); got != test.want {
			t.Errorf("Magnitude(...) => %v; want %v", got, test.want)
		}
	}
}

func TestListAdditionReductionSequence(t *testing.T) {
	tests := []struct {
		input []string
		want  []SnailNum
	}{
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]"},
			[]SnailNum{{1, 3}, {1, 3}, {2, 3}, {2, 3}, {3, 2}, {3, 2}, {4, 1}, {4, 1}},
		},
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]"},
			[]SnailNum{{3, 3}, {0, 3}, {5, 3}, {3, 3}, {4, 2}, {4, 2}, {5, 1}, {5, 1}},
		},
		{
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
				"[6,6]"},
			[]SnailNum{{5, 3}, {0, 3}, {7, 3}, {4, 3}, {5, 2}, {5, 2}, {6, 1}, {6, 1}},
		},
		{
			[]string{
				"[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			},
			[]SnailNum{{3, 1}, {2, 2}, {8, 3}, {0, 3}, {9, 1}, {5, 2}, {7, 3}, {0, 3}},
		},
		{
			[]string{
				"[[[[4,3],4],4],[7,[[8,4],9]]]",
				"[1,1]",
			},
			[]SnailNum{{0, 3}, {7, 3}, {4, 2}, {7, 3}, {8, 3}, {6, 3}, {0, 3}, {8, 1}, {1, 1}},
		},
		{
			[]string{
				"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
				"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
			},
			[]SnailNum{{4, 3}, {0, 3}, {5, 3}, {4, 3}, {7, 3}, {7, 3}, {6, 3}, {0, 3}, {8, 2}, {7, 3}, {7, 3}, {7, 3}, {9, 3}, {5, 3}, {0, 3}},
		},
		{
			[]string{
				"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
				"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
				"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
				"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
				"[7,[5,[[3,8],[1,4]]]]",
				"[[2,[2,2]],[8,[8,1]]]",
				"[2,9]",
				"[1,[[[9,3],9],[[9,0],[0,7]]]]",
				"[[[5,[7,4]],7],1]",
				"[[[[4,2],2],6],[8,7]]",
			},
			[]SnailNum{{8, 3}, {7, 3}, {7, 3}, {7, 3}, {8, 3}, {6, 3}, {7, 3}, {7, 3}, {0, 3}, {7, 3}, {6, 3}, {6, 3}, {8, 2}, {7, 2}},
		},
	}
	for _, test := range tests {
		if got := AdditionReductionSequence(ParseInput(test.input)); !SnailNumsEqual(got, test.want) {
			t.Errorf("AdditionReductionSequence(...) => \n%v\n want \n%v", got, test.want)
		}
	}
}
