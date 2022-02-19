/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

import "math/rand"

// RandomShuffle shuffles the given values in place.
func RandomShuffle[A any](r *rand.Rand, values []A) {
	r.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
}
