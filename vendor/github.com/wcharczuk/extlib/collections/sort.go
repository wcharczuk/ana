/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

import (
	"constraints"
	"sort"
)

// Sort sorts a list of elements by a given list of comparers.
func Sort[T any](elems []T, comparers ...SorterComparer[T]) {
	if len(comparers) == 0 {
		return
	}
	sort.Slice(elems, func(i, j int) bool {
		var compare int
		for _, sp := range comparers {
			compare = sp.Compare(elems[i], elems[j])
			if compare == 0 {
				continue
			}
			return compare < 0
		}
		return false
	})
}

// SortAsc is an identity sort ascending.
func SortAsc[T constraints.Ordered]() SorterComparer[T] {
	return SorterComparerFunc[T](func(i, j T) int {
		switch {
		case i == j:
			return 0
		case i < j:
			return -1
		default:
			return 1
		}
	})
}

// SortDesc is an identity sort descending.
func SortDesc[T constraints.Ordered]() SorterComparer[T] {
	return SorterComparerFunc[T](func(i, j T) int {
		switch {
		case i == j:
			return 0
		case i > j:
			return -1
		default:
			return 1
		}
	})
}

// SortKey is a sort comparer that extracts a key and sorts by it ascending.
func SortKey[T any, V constraints.Ordered](fn func(T) V) SorterComparer[T] {
	return SorterComparerFunc[T](func(i, j T) int {
		iv := fn(i)
		jv := fn(j)
		switch {
		case iv == jv:
			return 0
		case iv < jv:
			return -1
		default:
			return 1
		}
	})
}

// SortKeyDesc is a sort comparer that extracts a key and sorts by it descending.
func SortKeyDesc[T any, V constraints.Ordered](fn func(T) V) SorterComparer[T] {
	return SorterComparerFunc[T](func(i, j T) int {
		iv := fn(i)
		jv := fn(j)
		switch {
		case iv == jv:
			return 0
		case iv > jv:
			return -1
		default:
			return 1
		}
	})
}

// SorterComparer is a specific field or component of
// a multi-level sort.
type SorterComparer[T any] interface {
	Compare(i, j T) int
}

// SorterComparer is a predicate for comparing two elements.
type SorterComparerFunc[T any] func(T, T) int

// Compare implements SorterComparer.
func (scf SorterComparerFunc[T]) Compare(i, j T) int {
	return scf(i, j)
}
