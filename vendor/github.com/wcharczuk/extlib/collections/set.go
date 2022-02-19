/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

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
