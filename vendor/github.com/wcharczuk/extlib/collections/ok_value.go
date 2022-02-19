/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

// OkValue returns just the value from a (A,bool) return.
func OkValue[A any](v A, ok bool) A {
	return v
}

// Ok returns just the bool from a (A,bool) return.
func Ok[A any](_ A, ok bool) bool {
	return ok
}
