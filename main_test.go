package main

import (
	"testing"
)

func Test_permutations(t *testing.T) {
	p := permutations("abc", "")

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

func Test_matchesPositionMask(t *testing.T) {
	testCases := [...]struct {
		Mask     string
		Input    string
		Expected bool
	}{
		{"", "", false},
		{"", "abc", true},
		{"?bc?", "abc", false},
		{"?b?", "abc", true},
		{"?c?", "abc", false},
		{"?b?e", "abce", true},
		{"????", "abce", true},
		{"abce", "abce", true},
		{"?b?d", "abce", false},
	}

	for _, tc := range testCases {
		actual := matchesPositionMask(tc.Mask, tc.Input)
		if tc.Expected != actual {
			t.Fatalf("expected: %v for mask: %s and input: %s", tc.Expected, tc.Mask, tc.Input)
		}
	}
}

func Test_removeAt(t *testing.T) {
	testCases := [...]struct {
		Input    string
		Index    int
		Expected string
	}{
		{
			Input:    "",
			Index:    0,
			Expected: "",
		},
		{
			Input:    "1",
			Index:    0,
			Expected: "",
		},
		{
			Input:    "12345",
			Index:    0,
			Expected: "2345",
		},
		{
			Input:    "12345",
			Index:    4,
			Expected: "1234",
		},
		{
			Input:    "12345",
			Index:    1,
			Expected: "1345",
		},
		{
			Input:    "12345",
			Index:    2,
			Expected: "1245",
		},
		{
			Input:    "12345",
			Index:    3,
			Expected: "1235",
		},
	}
	for _, tc := range testCases {
		actual := removeAt(tc.Input, tc.Index)
		if tc.Expected != actual {
			t.Fatalf("expected: %v to equal actual: %v", tc.Expected, actual)
		}
	}
}
