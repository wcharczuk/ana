package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var flagDictPath = flag.String("dict", "/usr/share/dict/american-english", "The dictionary path")
	var flagMask = flag.String("m", "?????", "If we should match a position mask (e.g. ?i?e??)")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("ana [flags] <input>")
		oldUsage()
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		usagef("please supply a <input>")
	}

	input := flag.Args()[0]
	inputPermutations := permutations(input)

	dictFile, err := os.Open(*flagDictPath)
	fatal(err)
	defer dictFile.Close()

	dictScanner := bufio.NewScanner(dictFile)
	var dictWord string
	for dictScanner.Scan() {
		dictWord = dictScanner.Text()
		if inputPermutations.Has(dictWord) {
			if matchesPositionMask(*flagMask, dictWord) {
				fmt.Println(dictWord)
			}
		}
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

func permutations(input string) Set[string] {
	output := make(Set[string])
	seen := make(Set[string])
	var working string
	_permutations(input, working, seen, output)
	return output
}

func _permutations(input, working string, seen, output Set[string]) {
	if len(input) == 0 {
		output.Add(string(working))
		return
	}
	for index, c := range input {
		before := string(c) + working
		if !seen.Has(before) {
			seen.Add(before)
			_permutations(
				removeAt(input, index),
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
