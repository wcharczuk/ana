/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

// InsertAt inserts an element into a given array at a given index.
func InsertAt[A any](working []A, index int, v A) (output []A) {
	output = make([]A, len(working)+1)
	copy(output, working[:index])
	output[index] = v
	copy(output[index+1:], working[index:])
	return output
}
