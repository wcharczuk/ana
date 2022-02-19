package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/wcharczuk/extlib/cli"
	"github.com/wcharczuk/extlib/collections"
)

//go:embed dictionary.txt
var dictionary []byte

var alphabet = []rune("abcdefghijklmnopqrstuvwxyz")

var app = &cli.App{
	Name:  "wordle",
	Usage: "filter wordle dictionary words",
	Action: func(c *cli.Context) error {
		return action(c)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dict",
			Usage: "The dictionary path (optional, will use embedded dictionary by default)",
		},
		&cli.StringFlag{
			Name:  "mask",
			Usage: "The position mask to match (e.g. ?i?e??)",
		},
		&cli.StringFlag{
			Name:  "exclude",
			Usage: "The excluded letter set",
		},
		&cli.BoolFlag{
			Name:  "analyze",
			Usage: "If we should analyze results",
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "If we should limit results",
		},
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	dict, err := getDictionary(c.String("dict"))
	if err != nil {
		return err
	}

	exclude := collections.NewSet([]rune(c.String("exclude")))
	var maybe []rune
	for _, c := range alphabet {
		if !exclude.Has(c) {
			maybe = append(maybe, c)
		}
	}

	var inputPermutations collections.Set[string]
	if flagKnown := c.String("exclude"); flagKnown != "" {
		inputPermutations = permutations(flagKnown, string(maybe), c.String("mask"))
	}

	mask := []rune(c.String("mask"))
	analyzeDict := make(collections.Set[string])
	analyzeResults := &collections.Heap[WordStats]{
		LessFn: func(a, b WordStats) bool {
			return (a.Green + a.Yellow) > (b.Green + b.Yellow)
		},
	}

	var count int
	for dictWord := range dict {
		if inputPermutations != nil && !inputPermutations.Has(dictWord) {
			continue
		}
		if !matchesPositionMask(mask, []rune(dictWord)) {
			continue
		}
		if c.Bool("analyze") {
			analyzeDict.Add(dictWord)
		} else {
			if flagLimit := c.Int("limit"); flagLimit == 0 || (flagLimit > 0 && count < flagLimit) {
				fmt.Println(dictWord)
				count++
			}
		}
	}
	if c.Bool("analyze") {
		for word := range analyzeDict {
			analyzeResults.Push(analyze(analyzeDict, word))
		}
		for _, ws := range analyzeResults.Values {
			if flagLimit := c.Int("limit"); flagLimit == 0 || (flagLimit > 0 && count < flagLimit) {
				fmt.Printf("%s: %d/%d\n", ws.Word, ws.Green, ws.Yellow)
				count++
			}
		}
	}
	return nil
}

func getDictionary(dictPath string) (collections.Set[string], error) {
	r, err := getDictionaryReader(dictPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	output := make(collections.Set[string])
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		output.Add(scanner.Text())
	}
	return output, err
}

func getDictionaryReader(dictPath string) (io.ReadCloser, error) {
	if dictPath != "" {
		dictFile, err := os.Open(dictPath)
		if err != nil {
			return nil, err
		}
		return dictFile, nil
	}
	return io.NopCloser(bytes.NewReader(dictionary)), nil
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

func permutations(known, maybe, mask string) collections.Set[string] {
	knownRunes := []rune(known)
	maybeRunes := []rune(maybe)
	maskRunes := []rune(mask)
	if len(knownRunes) == len(maskRunes) {
		return collections.NewSet(_permutations(knownRunes, 0, maskRunes, nil))
	}

	output := make(collections.Set[string])
	missing := len(maskRunes) - len(knownRunes)

	maybeRunes = concat(maybeRunes, knownRunes...)
	for _, adds := range chooseAny(maybeRunes, missing) {
		results := _permutations(concat(knownRunes, adds...), 0, maskRunes, nil)
		for _, res := range results {
			output.Add(res)
		}
	}
	return output
}

func _permutations(input []rune, index int, mask, working []rune) (output []string) {
	if index == len(input) {
		if matchesPositionMask(mask, working) {
			output = []string{string(working)}
		}
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
	copy(output[:index], input[:index])
	output[index] = r
	copy(output[index+1:], input[index:])
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

		// start with all the bits of x
		// we only consider up to x
		// because left of x's value will be zeros
		for y := x; y > 0; y >>= 1 {
			// test if the _last_ bit is on
			// or off, if it's on, add the char
			// we do this with the value (1)
			// instead of all ones
			// if 1011 & 1 == 1, then the
			// last bit was one
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

func chooseAny(input []rune, count int) [][]rune {
	return _chooseAny(input, count, nil)
}

func _chooseAny(input []rune, count int, working []rune) (output [][]rune) {
	if len(working) == count {
		return [][]rune{working}
	}
	for _, r := range input {
		output = append(output,
			_chooseAny(input, count, concat(working, r))...,
		)
	}
	return
}

func analyze(dict collections.Set[string], word string) (output WordStats) {
	output.Word = word
	for w := range dict {
		if w == word {
			continue
		}
		green, yellow, miss := analyzeScore(w, word)
		output.Green += green
		output.Yellow += yellow
		output.Miss += miss
	}
	return
}

func analyzeScore(w0, w1 string) (green, yellow, miss int) {
	w0r := []rune(w0)
	w0s := collections.NewSet(w0r)
	w1r := []rune(w1)

	for x := 0; x < len(w0r); x++ {
		if w0r[x] == w1r[x] {
			green++
			continue
		}
		if w0s.Has(w1r[x]) {
			yellow++
			continue
		}
		miss++
		continue
	}
	return
}

type WordStats struct {
	Word   string
	Green  int
	Yellow int
	Miss   int
}
