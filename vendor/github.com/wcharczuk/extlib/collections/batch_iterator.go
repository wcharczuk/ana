/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

// BatchIterator is an iterator for interface{}
type BatchIterator[T any] struct {
	Items     []T
	BatchSize int
	Cursor    int
}

// HasNext returns if we should process another batch.
func (bi *BatchIterator[T]) HasNext() bool {
	return bi.Cursor < (len(bi.Items) - 1)
}

// Next yields the next batch of size `BatchSize` or smaller.
func (bi *BatchIterator[T]) Next() []T {
	if bi.BatchSize == 0 {
		return nil
	}
	if bi.Cursor >= len(bi.Items) {
		return nil
	}

	if (bi.Cursor + bi.BatchSize) < len(bi.Items) {
		output := bi.Items[bi.Cursor : bi.Cursor+bi.BatchSize]
		bi.Cursor += len(output)
		return output
	}

	output := bi.Items[bi.Cursor:]
	bi.Cursor += len(output)
	return output
}
