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
	var flagLimit = flag.Int("limit", 0, "If we should limit results")
	var flagVerbose = flag.Bool("verbose", false, "If we should show verbose output")

	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("ana [flags]")
		oldUsage()
	}
	flag.Parse()

	dict := getDictionary(*flagDictPath)

	var inputPermutations Set[string]
	if *flagKnown != "" {
		inputPermutations = permutations(*flagKnown, *flagMaybe, *flagMask)
	}

	if *flagVerbose {
		fmt.Println("permutations")
		for w := range inputPermutations {
			fmt.Println(w)
		}
		fmt.Println("---")
		fmt.Println("dictionary:", len(dict), "words")
	}

	mask := []rune(*flagMask)

	analyzeResults := &Heap[WordStats]{
		Less: func(a, b WordStats) bool {
			return a.Green > b.Green
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
		if *flagAnalyze {
			analyzeResults.Push(analyze(dict, dictWord))
		} else {
			if *flagLimit == 0 || (*flagLimit > 0 && count < *flagLimit) {
				fmt.Println(dictWord)
				count++
			}
		}
	}
	if *flagAnalyze {
		for _, ws := range analyzeResults.Values {
			if *flagLimit == 0 || (*flagLimit > 0 && count < *flagLimit) {
				fmt.Printf("%s: %d/%d\n", ws.Word, ws.Green, ws.Yellow)
				count++
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

func getDictionary(dictPath string) Set[string] {
	r := getDictionaryReader(dictPath)
	defer r.Close()

	output := make(Set[string])
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		output.Add(scanner.Text())
	}
	return output
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

func analyze(dict Set[string], word string) (output WordStats) {
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
	w0s := NewSet(w0r)
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

// Heap is a generic priority queue.
//
// You must provide a `Less(...) bool` function, but values can be omitted.
type Heap[A any] struct {
	Values []A
	Less   func(A, A) bool
}

// OkValue returns just the value from a (A,bool) return.
func OkValue[A any](v A, ok bool) A {
	return v
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func (h *Heap[A]) Init() {
	n := len(h.Values)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Len returns the length, or number of items in the heap.
func (h *Heap[A]) Len() int {
	return len(h.Values)
}

// Push pushes values onto the heap.
func (h *Heap[A]) Push(v A) {
	h.Values = append(h.Values, v)
	h.up(len(h.Values) - 1)
}

// Peek returns the first (smallest) element in the heap.
func (h *Heap[A]) Peek() (output A, ok bool) {
	if len(h.Values) == 0 {
		return
	}
	output = h.Values[0]
	ok = true
	return
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *Heap[A]) Pop() (output A, ok bool) {
	if len(h.Values) == 0 {
		return
	}

	// heap pop
	n := len(h.Values) - 1
	h.swap(0, n)
	h.down(0, n)

	// intheap pop
	old := h.Values
	n = len(old)
	output = old[n-1]
	ok = true
	h.Values = old[0 : n-1]
	return
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[A]) Fix(i int) {
	if !h.down(i, len(h.Values)) {
		h.up(i)
	}
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[A]) Remove(i int) (output A, ok bool) {
	n := len(h.Values) - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.Pop()
}

//
// internal helpers
//

func (h *Heap[A]) swap(i, j int) {
	h.Values[i], h.Values[j] = h.Values[j], h.Values[i]
}

func (h *Heap[A]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(h.Values[j], h.Values[i]) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[A]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(h.Values[j2], h.Values[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(h.Values[j], h.Values[i]) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}
