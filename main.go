package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
)

//go:embed dictionary.txt
var dictionary []byte

func main() {
	var flagDictPath = flag.String("dict", "", "The dictionary path (optional, will use embedded dictionary by default)")
	var flagMask = flag.String("mask", "?????", "If we should match a position mask (e.g. ?i?e??)")
	var flagKnown = flag.String("known", "", "The known letter set")
	var flagMaybe = flag.String("maybe", "", "The maybe letter set")

	var flagAnalyze = flag.Bool("analyze", false, "If we should print analysis results")
	var flagPermutations = flag.Bool("permutations", false, "If we should show permutations")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("ana [flags]")
		oldUsage()
	}
	flag.Parse()

	if *flagAnalyze {
		fmt.Println("analyze mode")
		return
	}

	var inputPermutations Set[string]
	if *flagKnown != "" {
		inputPermutations = permutations(*flagKnown, *flagMaybe, *flagMask)
	}

	if *flagPermutations {
		fmt.Printf("yielded %d permutations:\n", len(inputPermutations))
		for w := range inputPermutations {
			fmt.Println(w)
		}
		fmt.Println("---")
	}

	dictFile := getDictionaryReader(*flagDictPath)
	defer dictFile.Close()

	mask := []rune(*flagMask)

	dictScanner := bufio.NewScanner(dictFile)
	var dictWord string
	for dictScanner.Scan() {
		dictWord = dictScanner.Text()
		if inputPermutations != nil && !inputPermutations.Has(dictWord) {
			continue
		}
		if !matchesPositionMask(mask, []rune(dictWord)) {
			continue
		}
		fmt.Println(dictWord)
	}
}

func usagef(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	flag.Usage()
	os.Exit(1)
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %+v\n", err)
		os.Exit(1)
	}
}

func getDictionaryReader(dictPath string) io.ReadCloser {
	if dictPath != "" {
		dictFile, err := os.Open(dictPath)
		fatal(err)
		return dictFile
	}
	return io.NopCloser(bytes.NewReader(dictionary))
}

func matchesPositionMask(mask, input []rune) bool {
	if len(mask) == 0 && len(input) == 0 {
		return false
	}
	if len(mask) == 0 && len(input) > 0 {
		return true
	}
	if len(mask) != len(input) {
		return false
	}
	for index, r := range input {
		if mask[index] == '?' {
			continue
		}
		if mask[index] != r {
			return false
		}
	}
	return true
}

func permutations(known, maybe, mask string) Set[string] {
	knownRunes := []rune(known)
	maybeRunes := []rune(maybe)
	maskRunes := []rune(mask)
	if len(knownRunes) == len(maskRunes) {
		return NewSet(_permutations(knownRunes, 0, maskRunes, nil))
	}

	output := make(Set[string])
	missing := len(maskRunes) - len(knownRunes)

	maybeRunes = concat(maybeRunes, knownRunes...)
	for _, adds := range choose(maybeRunes, missing) {
		results := _permutations(concat(knownRunes, adds...), 0, maskRunes, nil)
		for _, res := range results {
			output.Add(res)
		}
	}
	return output
}

func _permutations(input []rune, index int, mask, working []rune) (output []string) {
	if index == len(input) && matchesPositionMask(mask, working) {
		output = []string{string(working)}
		return
	}

	c := input[index]
	for x := 0; x <= len(working); x++ {
		output = append(output,
			_permutations(input, index+1, mask, insertAt(working, c, x))...,
		)
	}
	return
}

func insertAt(input []rune, r rune, index int) []rune {
	output := make([]rune, len(input)+1)
	if index > 0 {
		copy(output[:index], input[:index])
	}
	output[index] = r
	if index < len(input) {
		copy(output[index+1:], input[index:])
	}
	return output
}

func concat(a []rune, b ...rune) []rune {
	output := make([]rune, 0, len(a)+len(b))
	for _, r := range a {
		output = append(output, r)
	}
	for _, r := range b {
		output = append(output, r)
	}
	return output
}

func choose(input []rune, count int) (output [][]rune) {
	max := 1 << len(input)

	for x := 0; x < max; x++ {
		var index int
		var w []rune
		for y := x; y > 0; y >>= 1 {
			if (y & 1) == 1 {
				w = append(w, input[index])
			}
			index++
		}
		if len(w) == count {
			output = append(output, w)
		}
	}
	return
}

// NewSet creates a new set.
func NewSet[A comparable](values []A) Set[A] {
	s := make(Set[A])
	for _, v := range values {
		s.Add(v)
	}
	return s
}

type Set[A comparable] map[A]struct{}

func (s Set[A]) Add(v A)           { s[v] = struct{}{} }
func (s Set[A]) Has(v A) (ok bool) { _, ok = s[v]; return }
