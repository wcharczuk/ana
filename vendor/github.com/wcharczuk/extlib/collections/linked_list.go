/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

// LinkedList is an implementation of a fifo buffer using nodes and poitners.
// Remarks; it is not threadsafe. It is constant(ish) time in all ops.
type LinkedList[T any] struct {
	head   *LinkedListNode[T]
	tail   *LinkedListNode[T]
	length int
}

// Len returns the length of the list in constant time.
func (q *LinkedList[T]) Len() int {
	return q.length
}

// PushBack adds a new value to the front of the list.
func (q *LinkedList[T]) PushFront(value T) {
	node := &LinkedListNode[T]{Value: value}
	if q.head == nil { //the list is empty, that is to say head is nil
		q.head = node
		q.tail = node
	} else { //the list is not empty, we have a (valid) tail pointer
		q.head.Next = node
		node.Previous = q.head
		q.head = node
	}
	q.length++
}

// Push adds a new value to the linked list.
func (q *LinkedList[T]) Push(value T) {
	node := &LinkedListNode[T]{Value: value}
	if q.head == nil { //the list is empty, that is to say head is nil
		q.head = node
		q.tail = node
	} else { //the list is not empty, we have a (valid) tail pointer
		q.tail.Previous = node
		node.Next = q.tail
		q.tail = node
	}
	q.length++
}

// Pop removes the head element from the list.
func (q *LinkedList[T]) Pop() (out T, ok bool) {
	if q.head == nil {
		return
	}
	out = q.head.Value
	ok = true
	if q.length == 1 && q.head == q.tail {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.Previous
		if q.head != nil {
			q.head.Next = nil
		}
	}
	q.length--
	return
}

// PopBack removes the tail element from the list.
func (q *LinkedList[T]) PopBack() (out T, ok bool) {
	if q.tail == nil {
		return
	}
	out = q.tail.Value
	ok = true

	if q.length == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.tail = q.tail.Next
		if q.tail != nil {
			q.tail.Previous = nil
		}
	}
	q.length--
	return
}

// Peek returns the first element of the list but does not remove it.
func (q *LinkedList[T]) Peek() (out T, ok bool) {
	if q.head == nil {
		return
	}
	out = q.head.Value
	ok = true
	return
}

// PeekBack returns the last element of the list.
func (q *LinkedList[T]) PeekBack() (out T, ok bool) {
	if q.tail == nil {
		return
	}
	out = q.tail.Value
	ok = true
	return
}

// Clear clears the linked list.
func (q *LinkedList[T]) Clear() {
	q.tail = nil
	q.head = nil
	q.length = 0
}

// Each calls the consumer for each element of the linked list.
func (q *LinkedList[T]) Each(consumer func(value T)) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		consumer(nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
}

// LinkedListNode is a linked list node.
type LinkedListNode[T any] struct {
	// Next points towards the head.
	Next *LinkedListNode[T]
	// Previous points towards the tail.
	Previous *LinkedListNode[T]
	// Value holds the value of the node.
	Value T
}
