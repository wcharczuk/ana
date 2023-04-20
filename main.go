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
	Name:  "ana",
	Usage: "filter dictionary words to solve anagrams",
	Action: func(c *cli.Context) error {
		return action(c)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dict-path",
			Usage: "The dictionary path (optional, will use embedded dictionary by default)",
		},
		&cli.StringFlag{
			Name:  "mask",
			Usage: "The known character position mask to match (e.g. ?i?e??)",
			Value: "?????",
		},
		&cli.StringFlag{
			Name:  "known",
			Usage: "The full known alphabet to search with",
			Value: string(alphabet),
		},
		&cli.StringFlag{
			Name:  "include",
			Usage: "The letter set that must be included in any match",
		},
		&cli.StringFlag{
			Name:  "exclude",
			Usage: "The letter set that if any appear in a word it will be disqualified",
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

func action(ctx *cli.Context) error {

	dictionaryPath := ctx.String("dict-path")
	known := ctx.String("known")
	include := ctx.String("include")
	exclude := ctx.String("exclude")
	mask := ctx.String("mask")

	dictionary, err := getDictionary(dictionaryPath)
	if err != nil {
		return err
	}

	excludeLookup := collections.NewSet([]rune(exclude))

	// build the alphabet
	var alphabetRunes []rune
	for _, c := range known {
		if !excludeLookup.Has(c) {
			alphabetRunes = append(alphabetRunes, c)
		}
	}

	fmt.Printf("using alphabet: %s\n", string(alphabetRunes))
	fmt.Printf("using dictionary: %s\n", dictionaryPath)
	fmt.Printf("using known: %s\n", known)
	fmt.Printf("using include: %s\n", include)
	fmt.Printf("using exclude: %s\n", exclude)
	fmt.Printf("using mask: %s\n", mask)

	knownPermutations := permutations(include, string(alphabetRunes), mask)

	maskRunes := []rune(mask)
	for word := range dictionary {
		if knownPermutations != nil && !knownPermutations.Has(word) {
			continue
		}
		if !matchesPositionMask(maskRunes, []rune(word)) {
			continue
		}
		if excludeMatches(excludeLookup, []rune(word)) {
			continue
		}
		fmt.Fprintln(os.Stdout, word)
	}
	return nil
}

func getDictionary(dictPath string) (collections.Set[string], error) {
	r, err := getDictionaryByPath(dictPath)
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

func getDictionaryByPath(dictPath string) (io.ReadCloser, error) {
	if dictPath != "" {
		dictFile, err := os.Open(dictPath)
		if err != nil {
			return nil, err
		}
		return dictFile, nil
	}
	return io.NopCloser(bytes.NewReader(dictionary)), nil
}

func excludeMatches(excludes collections.Set[rune], wordRunes []rune) bool {
	for _, r := range wordRunes {
		if excludes.Has(r) {
			return true
		}
	}
	return false
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
	if known == "" {
		return nil
	}

	knownRunes := []rune(known)
	maybeRunes := []rune(maybe)
	maskRunes := []rune(mask)
	if len(knownRunes) == len(maskRunes) {
		return collections.NewSet(_permutations(knownRunes, 0, maskRunes, nil))
	}

	output := make(collections.Set[string])
	missing := len(maskRunes) - len(knownRunes)

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
	if count <= 0 {
		return nil
	}
	return _chooseAny(input, count, 0, nil)
}

func _chooseAny(input []rune, count, index int, working []rune) (output [][]rune) {
	if len(working) == count {
		return [][]rune{working}
	}
	if index == len(input) {
		return nil
	}
	output = append(output,
		_chooseAny(input, count, index+1, concat(working, input[index]))...,
	)
	output = append(output,
		_chooseAny(input, count, index+1, working)...,
	)
	return
}
