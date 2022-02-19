/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

// Powerset returns all possible selections from a given
// set of values with their order preserved.
func Powerset[A any](values ...A) (output [][]A) {
	max := 1 << len(values)
	for x := 1; x < max; x++ {
		var index int
		var working []A
		for y := x; y > 0; y >>= 1 {
			if y&1 == 1 {
				working = append(working, values[index])
			}
			index++
		}
		output = append(output, working)
	}
	return
}
