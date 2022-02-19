/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

import "constraints"

// InsertSorted performs an insertion at the index that would satisfy
// that the resulting array would be sorted ascending.
func InsertSorted[A constraints.Ordered](working []A, v A) []A {
	insertAt := Search(working, v)
	working = append(working, v)
	copy(working[insertAt+1:], working[insertAt:])
	working[insertAt] = v
	return working
}
