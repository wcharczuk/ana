package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var flagDictPath = flag.String("dict", "/usr/share/dict/american-english", "The dictionary path")
var flagVerbose = flag.Bool("verbose", false, "If we should show verbose output")

func init() {
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("ana [flags] <input>")
		oldUsage()
	}
	flag.Parse()
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

func main() {
	if len(flag.Args()) != 1 {
		usagef("please supply a <input>")
	}

	input := flag.Args()[0]
	inputPermutations := permutations(input)
	if *flagVerbose {
		fmt.Println("Permutations:")
		for p := range inputPermutations {
			fmt.Println(p)
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
			fmt.Println(dictWord)
		}
	}
}

type Set[A comparable] map[A]struct{}

func (s Set[A]) Add(v A)           { s[v] = struct{}{} }
func (s Set[A]) Has(v A) (ok bool) { _, ok = s[v]; return }

func permutations(input string) Set[string] {
	output := make(Set[string])

	inputRunes := []rune(input)

	for index, c := range inputRunes {
		_permutations(inputRunes, index, []rune{c}, output)
	}

	return output
}

func _permutations(inputRunes []rune, index int, working []rune, output Set[string]) {
	if index == len(inputRunes)-1 {
		if len(working) == len(inputRunes) {
			output.Add(string(working))
		}
		return
	}

	for subIndex := index + 1; subIndex < len(inputRunes); subIndex++ {
		_permutations(inputRunes, subIndex, append(working, inputRunes[subIndex]), output)
		_permutations(inputRunes, subIndex, append([]rune{inputRunes[subIndex]}, working...), output)
	}
}
