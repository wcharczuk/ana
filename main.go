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

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("ana [flags]")
		oldUsage()
	}
	flag.Parse()

	var inputPermutations Set[string]
	if *flagKnown != "" {
		inputPermutations = permutations(*flagKnown, *flagMaybe, *flagMask)
	}

	dictFile := getDictionaryReader(*flagDictPath)
	defer dictFile.Close()

	dictScanner := bufio.NewScanner(dictFile)
	var dictWord string
	for dictScanner.Scan() {
		dictWord = dictScanner.Text()
		if inputPermutations != nil && !inputPermutations.Has(dictWord) {
			continue
		}
		if *flagMask != "" && !matchesPositionMask(*flagMask, dictWord) {
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

func matchesPositionMask(mask, input string) bool {
	if mask == "" && input == "" {
		return false
	}
	if mask == "" && input != "" {
		return true
	}
	if len(mask) != len(input) {
		return false
	}
	maskRunes := []rune(mask)
	for index, r := range input {
		if maskRunes[index] == '?' {
			continue
		}
		if maskRunes[index] != r {
			return false
		}
	}
	return true
}

func permutations(known, maybe, mask string) Set[string] {
	output := make(Set[string])
	var working string

	if len(known) == len(mask) {
		seen := make(Set[string])
		_permutations(known, mask, working, seen, output)
		return output
	}

	knownRunes := []rune(known)
	maybeRunes := []rune(maybe)
	if len(known) == len(mask)-1 {
		for _, r := range knownRunes {
			seen := make(Set[string])
			_permutations(string(append(knownRunes, r)), mask, working, seen, output)
		}
		for _, r := range maybeRunes {
			seen := make(Set[string])
			_permutations(string(append(knownRunes, r)), mask, working, seen, output)
		}
	}
	return output
}

func _permutations(input, mask, working string, seen, output Set[string]) {
	if len(input) == 0 {
		if mask != "" && matchesPositionMask(mask, string(working)) {
			output.Add(string(working))
		}
		return
	}
	for index, c := range input {
		before := string(c) + working
		if !seen.Has(before) {
			seen.Add(before)
			_permutations(
				removeAt(input, index),
				mask,
				before,
				seen,
				output,
			)
		}
		after := working + string(c)
		if !seen.Has(after) {
			seen.Add(after)
			_permutations(
				removeAt(input, index),
				mask,
				after,
				seen,
				output,
			)
		}
	}
}

func removeAt(input string, index int) string {
	if input == "" {
		return ""
	}
	inputRunes := []rune(input)
	return string(append(inputRunes[:index], inputRunes[index+1:]...))
}

type Set[A comparable] map[A]struct{}

func (s Set[A]) Add(v A)           { s[v] = struct{}{} }
func (s Set[A]) Has(v A) (ok bool) { _, ok = s[v]; return }
