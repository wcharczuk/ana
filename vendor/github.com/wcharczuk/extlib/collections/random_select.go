/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

import "math/rand"

// RandomSelect returns a random selection from the values.
func RandomSelect[A any](r *rand.Rand, values []A, count int) (output []A) {
	maxIndex := len(values) - 1
	output = make([]A, count)
	for x := 0; x < count; x++ {
		output[x] = values[r.Intn(maxIndex)]
	}
	return
}
