package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

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

func main() {
	var flagDictPath = flag.String("dict", "/usr/share/dict/american-english", "The dictionary path")
	var flagVerbose = flag.Bool("verbose", false, "If we should show verbose output")

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
	if *flagVerbose {
		fmt.Println("Permutations:")
		for p := range inputPermutations {
			fmt.Println("\t" + p)
		}
		fmt.Println("Found:")
	}

	dictFile, err := os.Open(*flagDictPath)
	fatal(err)
	defer dictFile.Close()

	dictScanner := bufio.NewScanner(dictFile)
	var dictWord string
	for dictScanner.Scan() {
		dictWord = dictScanner.Text()
		if inputPermutations.Has(dictWord) {
			if *flagVerbose {
				fmt.Println("\t" + dictWord)
			} else {
				fmt.Println(dictWord)
			}
		}
	}
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
				string(sliceRemove([]rune(input), index)),
				before,
				seen,
				output,
			)
		}
		after := working + string(c)
		if !seen.Has(after) {
			seen.Add(after)
			_permutations(
				string(sliceRemove([]rune(input), index)),
				after,
				seen,
				output,
			)
		}
	}
}

func sliceRemove[A any](values []A, index int) []A {
	if len(values) == 0 {
		return nil
	}
	return append(values[:index], values[index+1:]...)
}

type Set[A comparable] map[A]struct{}

func (s Set[A]) Add(v A)           { s[v] = struct{}{} }
func (s Set[A]) Has(v A) (ok bool) { _, ok = s[v]; return }
