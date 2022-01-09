package main

import "testing"

func Test_insertAt(t *testing.T) {
	cases := [...]struct {
		Input    []rune
		Index    int
		Rune     rune
		Expected string
	}{
		{nil, 0, 'a', "a"},
		{[]rune("abc"), 0, 'd', "dabc"},
		{[]rune("abc"), 1, 'd', "adbc"},
		{[]rune("abc"), 2, 'd', "abdc"},
		{[]rune("abc"), 3, 'd', "abcd"},
	}

	for _, c := range cases {
		actual := insertAt(c.Input, c.Rune, c.Index)
		if string(actual) != c.Expected {
			t.Fatalf("expected %q, got %q", c.Expected, actual)
		}
	}
}

func Test_chose(t *testing.T) {
	results := choose([]rune("abcde"), 3)
	if len(results) == 0 {
		t.Fatal("results are empty")
	}
	if len(results) != 10 {
		for _, r := range results {
			t.Log(string(r))
		}
		t.Fatalf("expected 10 results, got %d", len(results))
	}
}

func Test_chose_one(t *testing.T) {
	results := choose([]rune("abcde"), 1)
	if len(results) == 0 {
		t.Fatal("results are empty")
	}
	if len(results) != 5 {
		for _, r := range results {
			t.Log(string(r))
		}
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}

func Test_choseAny(t *testing.T) {
	results := chooseAny([]rune("abcde"), 3)
	if len(results) == 0 {
		t.Fatal("results are empty")
	}
	if len(results) != 125 {
		for _, r := range results {
			t.Log(string(r))
		}
		t.Fatalf("expected 125 results, got %d", len(results))
	}
}

func Test_choseAny_one(t *testing.T) {
	results := chooseAny([]rune("abcde"), 1)
	if len(results) == 0 {
		t.Fatal("results are empty")
	}
	if len(results) != 5 {
		for _, r := range results {
			t.Log(string(r))
		}
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}
