package main

import "testing"

func Test_yellowsMatchesAll(t *testing.T) {
	yellows := []string{
		"_i___",
	}
	input := "lipas"

	if !yellowsMatchesAll(yellows, input) {
		t.Fatalf("expect %v to match %v", yellows, string(input))
	}
}

func Test_yellowsMatchesAny(t *testing.T) {
	yellows := []string{
		"_e___",
		"_ru_l",
	}
	inputs := []string{
		"urali",
		"lidar",
	}

	for _, input := range inputs {
		if !yellowsMatchesAny(yellows, input) {
			t.Fatalf("expect %v to not match %v", yellows, input)
		}
	}
}
