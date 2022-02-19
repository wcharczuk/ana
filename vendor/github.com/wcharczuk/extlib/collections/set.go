/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

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
