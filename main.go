package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

//go:embed dictionary.txt
var dictionary []byte

// MASK_CHAR is the character we use as a wildcard in masks.
const MASK_CHAR = '_'

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
			Name:  "green",
			Usage: "The position of the matched letters in mask form (e.g. 'WO__L')",
		},
		&cli.StringSliceFlag{
			Name:  "yellow",
			Usage: "The yellows in position mask form (can be multiple!)",
		},
		&cli.StringFlag{
			Name:  "gray",
			Usage: "The excluded letter set as a string (e.g. 'ergv')",
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "If we should limit the number of results shown.",
		},
		&cli.BoolFlag{
			Name:    "match",
			Aliases: []string{"m"},
			Usage:   "If we should show match results.",
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
	dict, err := getDictionary(ctx.String("dict"))
	if err != nil {
		return err
	}

	flagLimit := ctx.Int("limit")
	yellows := ctx.StringSlice("yellow")
	green := []rune(ctx.String("green"))
	gray := []rune(ctx.String("gray"))

	var isDebug = os.Getenv("DEBUG") != ""
	debugf := func(format string, args ...any) {
		if isDebug {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		}
	}

	if ctx.Bool("match") {
		var matched []wordWithScore
		for dictWord := range dict {
			dictWordRunes := []rune(dictWord)
			if !greenMatches(green, dictWordRunes) {
				debugf("skipping %q; doesn't match greens %q", dictWord, string(green))
				continue
			}
			if !yellowsMatchesAll(yellows, dictWord) {
				debugf("skipping %q; doesn't match yellows %q", dictWord, strings.Join(yellows, ", "))
				continue
			}
			if !grayMatches(gray, dictWordRunes) {
				debugf("skipping %q; doesn't match grays %q", dictWord, string(gray))
				continue
			}
			matched = append(matched, wordWithScore{
				Word:  dictWord,
				Score: scoreWordMatch(dictWord),
			})
		}
		sort.SliceStable(matched, func(i, j int) bool {
			return matched[i].Score > matched[j].Score
		})
		for index, ws := range matched {
			fmt.Printf("%s (%d)\n", ws.Word, ws.Score)
			if flagLimit > 0 && index > flagLimit {
				break
			}
		}
	} else {
		var discover []wordWithScore
		for dictWord := range dict {
			dictWordRunes := []rune(dictWord)
			if len(green) > 0 {
				if greenMatches(green, dictWordRunes) {
					debugf("skipping %q for discovery; matches greens %q", dictWord, green)
					continue
				}
			}
			if yellowsMatchesAny(yellows, dictWord) {
				debugf("skipping %q for discovery; matches yellows %q", dictWord, strings.Join(yellows, ", "))
				continue
			}
			if !grayMatches(gray, dictWordRunes) {
				debugf("skipping %q for discovery; doesn't match grays %q", dictWord, string(gray))
				continue
			}
			discover = append(discover, wordWithScore{
				Word:  dictWord,
				Score: scoreWordMatch(dictWord),
			})
		}
		sort.SliceStable(discover, func(i, j int) bool {
			return discover[i].Score > discover[j].Score
		})

		if len(discover) > 0 {
			for index, ws := range discover {
				fmt.Printf("%s (%d)\n", ws.Word, ws.Score)
				if flagLimit > 0 && index > flagLimit {
					break
				}
			}
		}
	}

	return nil
}

type wordWithScore struct {
	Word  string
	Score int
}

func getDictionary(dictPath string) (Set[string], error) {
	r, err := getDictionaryReader(dictPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	output := make(Set[string])
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

func greenMatches(greens, input []rune) bool {
	if len(greens) == 0 && len(input) == 0 {
		return false
	}
	if len(greens) == 0 && len(input) > 0 {
		return true
	}
	if len(greens) != len(input) {
		return false
	}
	for index, r := range input {
		if greens[index] == MASK_CHAR {
			continue
		}
		if greens[index] != r {
			return false
		}
	}
	return true
}

// yellowsMatchesAll returns true if _all_ of the given yellow masks
// match the give input.
//
// a yellow match is given as the yellow letter counts being a
// strict subset of the input letter counts.
func yellowsMatchesAll(yellows []string, input string) bool {
	inputCounts := runeCounts(input)
	for _, y := range yellows {
		yCounts := runeCounts(y)
		if !runeCountsWithin(yCounts, inputCounts) {
			return false
		}
	}
	return true
}

// yellowsMatchesNone returns true if _any_ of the given yellow masks
// match the give input.
//
// a yellow match is given as the yellow letter counts being a
// strict subset of the input letter counts.
func yellowsMatchesAny(yellows []string, input string) bool {
	inputCounts := runeCounts(input)
	yCounts := runeCounts(yellows...)
	for key := range yCounts {
		if _, ok := inputCounts[key]; ok {
			return true
		}
	}
	return false
}

// grayMatches returns if _none_ of the runes in the gray list
// appear in the input.
func grayMatches(grays, input []rune) bool {
	for _, gc := range grays {
		for _, gi := range input {
			if gc == gi {
				return false
			}
		}
	}
	return true
}

// runeCounts returns a map of each rune in a given input
// mapped to the count or number of times that rune appears
// in the input list.
func runeCounts(inputs ...string) map[rune]int {
	output := make(map[rune]int)
	for _, input := range inputs {
		for _, c := range input {
			if c != MASK_CHAR {
				output[c] += 1
			}
		}
	}
	return output
}

// runeCountsWithin returns if a is a strict subset of b.
//
// strictly "subset" means every key in a exists in b, and
// the counts for each key in a is less than or equal to the count in b.
func runeCountsWithin(a, b map[rune]int) bool {
	for key, aCount := range a {
		bCount, ok := b[key]
		if !ok {
			return false
		}
		if aCount > bCount {
			return false
		}
	}
	return true
}

// NewSet creates a new set.
func NewSet[A comparable](values []A) Set[A] {
	s := make(Set[A])
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Set is a generic set.
type Set[A comparable] map[A]struct{}

// Add adds a given element.
func (s *Set[A]) Add(v A) {
	(*s)[v] = struct{}{}
}

// Has returns if a given element exists.
func (s *Set[A]) Has(v A) bool {
	_, ok := (*s)[v]
	return ok
}

/*
scrabble score is given as:
1 point – A   E   I   O   U   L   N   S   T   R
2 points – D   G
3 points – B   C   M   P
4 points – F   H   V   W   Y
5 points – K
8 points – J  X
10 points – Q  Z
*/

func scoreWordInvertedScrabble(c rune) int {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'l', 'n', 's', 't', 'r':
		return 10
	case 'd', 'g':
		return 8
	case 'b', 'c', 'm', 'p':
		return 5
	case 'f', 'h', 'v', 'w', 'y':
		return 4
	case 'k':
		return 3
	case 'j', 'x':
		return 2
	case 'q', 'z':
		return 1
	default:
		panic("invalid letter for inverted scrabble scoring")
	}
}

func scoreWordUnique(word string) int {
	s := make(Set[rune])
	for _, c := range word {
		s.Add(c)
	}
	return len(s)
}

func scoreWordMatch(word string) int {
	var invertedScrabbleScore int
	for _, c := range word {
		invertedScrabbleScore += scoreWordInvertedScrabble(c)
	}
	uniqueScore := scoreWordUnique(word)
	return (uniqueScore * 20) + invertedScrabbleScore
}
