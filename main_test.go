package main

import (
	"reflect"
	"testing"
)

func Test_permutations(t *testing.T) {
	p := permutations("abc")

	expectHas := func(v string) {
		if !p.Has(v) {
			t.Fatalf("expected to have %q", v)
		}
	}

	expectHas("abc")
	expectHas("bca")
	expectHas("cab")
	expectHas("bac")
	expectHas("cba")
	expectHas("acb")
}

func Test_sliceRemove(t *testing.T) {
	testCases := [...]struct {
		Input    []int
		Index    int
		Expected []int
	}{
		{
			Input:    nil,
			Index:    0,
			Expected: nil,
		},
		{
			Input:    []int{1},
			Index:    0,
			Expected: nil,
		},
		{
			Input:    []int{1, 2, 3, 4, 5},
			Index:    0,
			Expected: []int{2, 3, 4, 5},
		},
		{
			Input:    []int{1, 2, 3, 4, 5},
			Index:    4,
			Expected: []int{1, 2, 3, 4},
		},
		{
			Input:    []int{1, 2, 3, 4, 5},
			Index:    1,
			Expected: []int{1, 3, 4, 5},
		},
		{
			Input:    []int{1, 2, 3, 4, 5},
			Index:    2,
			Expected: []int{1, 2, 4, 5},
		},
		{
			Input:    []int{1, 2, 3, 4, 5},
			Index:    3,
			Expected: []int{1, 2, 3, 5},
		},
	}
	for _, tc := range testCases {
		actual := sliceRemove(tc.Input, tc.Index)
		if !reflect.DeepEqual(tc.Expected, actual) {
			t.Fatalf("expected: %v to equal actual: %v", tc.Expected, actual)
		}
	}
}

func Test_sliceRemoveRepeat(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	output := sliceRemove(
		sliceRemove(
			sliceRemove(input, 4),
			3,
		),
		2,
	)
	if !reflect.DeepEqual([]int{1, 2}, output) {
		t.Fatalf("expected: %v to equal actual: %v", []int{1, 2}, output)
	}
}
