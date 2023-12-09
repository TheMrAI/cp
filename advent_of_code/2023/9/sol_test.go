package main

import (
	"testing"
)

func TestExtrapolateNextValue(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{0, 3, 6, 9, 12, 15}, 18},
		{[]int{1, 3, 6, 10, 15, 21}, 28},
		{[]int{10, 13, 16, 21, 30, 45}, 68},
		{[]int{0}, 0},
		{[]int{1}, -1},
		{[]int{4, 1, 3}, -1},
		{[]int{0, 0, 0, 0, 0, 0}, 0},
		{[]int{-2, -4, -6}, -8},
		{[]int{8, 6, 4}, 2},
		{[]int{24, 32, 35, 33, 26, 14, -3, -25, -52, -84, -121, -163, -210, -262, -319, -381, -448, -520, -597, -679, -766}, -858},
		{[]int{-2, -4, -6, -8, -10, -12, -14, -16, -18, -20, -22, -24, -26, -28, -30, -32, -34, -36, -38, -40, -42}, -44},
	}
	for _, test := range tests {
		input := make([]int, len(test.input))
		copy(input, test.input)
		if got := ExtrapolateNextValue(input); got != test.want {
			t.Errorf("ExtrapolateNextValue(%v) => %v; want %v", test.input, got, test.want)
		}
	}
}
