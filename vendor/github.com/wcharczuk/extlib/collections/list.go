/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

// List is a slice api.
type List[T any] []T

// Append adds elements to the list and returns a copy.
func (l List[T]) Append(v ...T) (output List[T]) {
	output = make([]T, len(l)+len(v))
	copy(output, l)
	copy(output[len(l):], v)
	return
}

// Filter returns a filtered version of the list.
func (l List[T]) Filter(predicate func(T) bool) (output List[T]) {
	output = make(List[T], 0, len(l))
	for _, i := range l {
		if predicate(i) {
			output = append(output, i)
		}
	}
	return
}

// InsertAt inserts an element at a given index.
func (l List[T]) InsertAt(v T, index int) (output List[T]) {
	output = make(List[T], len(l)+1)
	copy(output, l[:index])
	output[index] = v
	copy(output[index+1:], l[index:])
	return
}

// Offset returns an offset version of the list.
func (l List[T]) Offset(offset int) (output List[T]) {
	if len(l) == 0 {
		output = make([]T, 0)
		return
	}
	if len(l) < offset {
		output = make([]T, len(l))
		copy(output, l)
		return
	}
	output = make(List[T], 0, len(l)-offset)
	for x := offset; x < len(l); x++ {
		output = append(output, l[x])
	}
	return
}

// Limit returns an limited version of the list.
func (l List[T]) Limit(limit int) (output List[T]) {
	if len(l) == 0 {
		output = make([]T, 0)
		return
	}
	if len(l) <= limit {
		output = make([]T, len(l))
		copy(output, l)
		return
	}

	output = make(List[T], 0, limit)
	for x := 0; x < limit; x++ {
		output = append(output, l[x])
	}
	return
}

// Copy copies the list.
func (l List[T]) Copy() (output List[T]) {
	output = make(List[T], len(l))
	copy(output, l)
	return
}

// OrderBy orders the list by a given list of comparers.
func (l List[T]) OrderBy(comparers ...SorterComparer[T]) (output List[T]) {
	output = l.Copy()
	Sort(output, comparers...)
	return
}
