package main

import "testing"

func Test_yellowsMatches(t *testing.T) {
	yellows := []string{
		"_i___",
	}
	input := []rune("lipas")

	if !yellowsMatches(yellows, input) {
		t.Fatalf("expect %v to match %v", yellows, string(input))
	}
}
