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

type Set[A comparable] map[A]int

func (s Set[A]) Add(v A)           { s[v] = s[v] + 1 }
func (s Set[A]) Has(v A) (ok bool) { _, ok = s[v]; return }
func (s Set[A]) Remove(v A) {
	count, ok := s[v]
	if !ok {
		return
	}
	if count == 1 {
		delete(s, v)
		return
	}
	s[v] = count - 1
}
func (s Set[A]) Pop() (v A) {
	for v = range s {
		break
	}
	s.Remove(v)
	return
}
func (s Set[A]) Copy() Set[A] {
	output := make(Set[A])
	for key, count := range s {
		output[key] = count
	}
	return output
}

func (s Set[A]) Count() (output int) {
	for _, c := range s {
		output = output + c
	}
	return
}

// a,b,c,d

// a
// b,a | a,b
// c,b,a | a,b,c
// d,c,b,a | a,b,c,d
// b
// a,b | b,a
// c,a,b |

func permutations(input string) Set[string] {
	output := make(Set[string])
	seen := make(Set[string])
	remaining := make(Set[rune])

	for _, r := range input {
		remaining.Add(r)
	}
	_permutations(remaining.Count(), nil, remaining, seen, output)
	return output
}

func _permutations(totalLen int, working []rune, remaining Set[rune], seen, output Set[string]) {
	if len(remaining) == 0 {
		if len(working) == totalLen {
			output.Add(string(working))
		}
		return
	}

	subRemaining := remaining.Copy()
	c := subRemaining.Pop()
	before := append([]rune{c}, working...)
	if !seen.Has(string(before)) {
		seen.Add(string(before))
		_permutations(totalLen, before, subRemaining, seen, output)
	}
	after := append(working, c)
	if !seen.Has(string(after)) {
		seen.Add(string(after))
		_permutations(totalLen, after, subRemaining, seen, output)
	}
}
