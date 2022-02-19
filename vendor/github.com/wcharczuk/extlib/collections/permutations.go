/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

func Permutations[A any](values ...A) [][]A {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return [][]A{values}
	}
	return _permutations(values, 0, nil)
}

func _permutations[A any](values []A, index int, working []A) (output [][]A) {
	if index == len(values) {
		return [][]A{working}
	}

	for x := 0; x <= len(working); x++ {
		output = append(output, _permutations(values, index+1, InsertAt(working, x, values[index]))...)
	}
	return
}
