/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

const (
	queueMinimumGrow     = 4
	queueGrowFactor      = 200
	queueDefaultCapacity = 4
)

// Queue is a fifo buffer that is backed by a pre-allocated array, as opposed to a linked-list
// which would allocate a whole new struct for each element, which saves GC churn.
// Push can be O(n), Dequeue can be O(1).
type Queue[A any] struct {
	array []A
	head  int
	tail  int
	size  int
}

// Len returns the length of the queue (as it is currently populated).
//
// Actual memory footprint may be different, use `Cap()` to return total memory.
func (q *Queue[A]) Len() int {
	return q.size
}

// Cap returns the total capacity of the queue, including empty elements.
func (q *Queue[A]) Cap() int {
	return len(q.array)
}

// Clear removes all objects from the Queue.
func (q *Queue[A]) Clear() {
	if q.head < q.tail {
		arrayClear(q.array, q.head, q.size)
	} else {
		arrayClear(q.array, q.head, len(q.array)-q.head)
		arrayClear(q.array, 0, q.tail)
	}
	q.head = 0
	q.tail = 0
	q.size = 0
}

// Push adds an element to the "back" of the Queue.
func (q *Queue[A]) Push(v A) {
	if len(q.array) == 0 {
		q.array = make([]A, queueDefaultCapacity)
	} else if q.size == len(q.array) {
		newCapacity := int(len(q.array) * int(queueGrowFactor/100))
		if newCapacity < (len(q.array) + queueMinimumGrow) {
			newCapacity = len(q.array) + queueMinimumGrow
		}
		q.setCapacity(newCapacity)
	}
	q.array[q.tail] = v
	q.tail = (q.tail + 1) % len(q.array)
	q.size++
}

// Pop removes the first (oldest) element from the Queue.
func (q *Queue[A]) Pop() (output A, ok bool) {
	if q.size == 0 {
		return
	}
	output = q.array[q.head]
	ok = true
	q.head = (q.head + 1) % len(q.array)
	q.size--
	return
}

// Pop removes the last (newest) element from the Queue.
func (q *Queue[A]) PopBack() (output A, ok bool) {
	if q.size == 0 {
		return
	}

	if q.tail == 0 {
		output = q.array[len(q.array)-1]
		q.tail = len(q.array) - 1
	} else {
		output = q.array[q.tail-1]
		q.tail = q.tail - 1
	}
	ok = true
	q.size--
	return
}

// Peek returns but does not remove the first element.
func (q *Queue[A]) Peek() (output A, ok bool) {
	if q.size == 0 {
		return
	}
	output = q.array[q.head]
	ok = true
	return
}

// PeekBack returns but does not remove the last element.
func (q *Queue[A]) PeekBack() (output A, ok bool) {
	if q.size == 0 {
		return
	}
	if q.tail == 0 {
		output = q.array[len(q.array)-1]
		ok = true
		return
	}
	output = q.array[q.tail-1]
	ok = true
	return
}

// Each calls the fn for each element in the buffer.
func (q *Queue[A]) Each(fn func(A)) {
	if q.size == 0 {
		return
	}
	if q.head < q.tail {
		for cursor := q.head; cursor < q.tail; cursor++ {
			fn(q.array[cursor])
		}
	} else {
		for cursor := q.head; cursor < len(q.array); cursor++ {
			fn(q.array[cursor])
		}
		for cursor := 0; cursor < q.tail; cursor++ {
			fn(q.array[cursor])
		}
	}
}

// ReverseEach calls fn in reverse order (tail to head).
func (q *Queue[A]) ReverseEach(fn func(A)) {
	if q.size == 0 {
		return
	}
	if q.head < q.tail {
		for cursor := q.tail - 1; cursor >= q.head; cursor-- {
			fn(q.array[cursor])
		}
	} else {
		for cursor := q.tail; cursor > 0; cursor-- {
			fn(q.array[cursor])
		}
		for cursor := len(q.array) - 1; cursor >= q.head; cursor-- {
			fn(q.array[cursor])
		}
	}
}

func (q *Queue[A]) setCapacity(capacity int) {
	newArray := make([]A, capacity)
	if q.size > 0 {
		if q.head < q.tail {
			arrayCopy(q.array, q.head, newArray, 0, q.size)
		} else {
			arrayCopy(q.array, q.head, newArray, 0, len(q.array)-q.head)
			arrayCopy(q.array, 0, newArray, len(q.array)-q.head, q.tail)
		}
	}
	q.array = newArray
	q.head = 0
	if q.size == capacity {
		q.tail = 0
	} else {
		q.tail = q.size
	}
}

// trimExcess resizes the buffer to better fit the contents.
func (q *Queue[A]) trimExcess() {
	threshold := float64(len(q.array)) * 0.9
	if q.size < int(threshold) {
		q.setCapacity(q.size)
	}
}

func arrayClear[A any](source []A, index, length int) {
	var zero A
	for x := 0; x < length; x++ {
		absoluteIndex := x + index
		source[absoluteIndex] = zero
	}
}

func arrayCopy[A any](source []A, sourceIndex int, destination []A, destinationIndex, length int) {
	for x := 0; x < length; x++ {
		from := sourceIndex + x
		to := destinationIndex + x
		destination[to] = source[from]
	}
}
