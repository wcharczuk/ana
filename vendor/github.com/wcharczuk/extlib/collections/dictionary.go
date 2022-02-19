/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the 
LICENSE file at the root of the repository.

*/

package collections

// NewDictionary returns a new dictionary.
func NewDictionary[K comparable, V any]() Dictionary[K, V] {
	return make(Dictionary[K, V])
}

// Dictionary is a map with a better api.
type Dictionary[K comparable, V any] map[K]V

// Get gets a value from the dictionary.
func (d Dictionary[K, V]) Get(k K) (v V, ok bool) {
	v, ok = d[k]
	return
}

// Set sets a value in the dictionary.
func (d Dictionary[K, V]) Set(k K, v V) {
	d[k] = v
}

// Has returns if the dictionary contains an
// element with a given key.
func (d Dictionary[K, V]) Has(k K) (ok bool) {
	_, ok = d[k]
	return
}

// Remove removes an element.
func (d Dictionary[K, V]) Remove(k K) {
	delete(d, k)
}

// Keys returns the dictionary keys.
func (d Dictionary[K, V]) Keys() (output []K) {
	output = make([]K, 0, len(d))
	for k := range d {
		output = append(output, k)
	}
	return
}

// Values returns the dictionary values.
func (d Dictionary[K, V]) Values() (output []V) {
	output = make([]V, 0, len(d))
	for _, v := range d {
		output = append(output, v)
	}
	return
}
