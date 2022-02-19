/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

// LRU is a cache that evicts items based on
// which were used last.
type LRU[K comparable, V any] struct {
	Capacity int
	OnEvict  func(K, V)

	head   *lruItem[K, V]
	tail   *lruItem[K, V]
	lookup map[K]*lruItem[K, V]
}

// Get returns an item with a given key.
func (lru *LRU[K, V]) Get(k K) (v V, ok bool) {
	if lru.lookup == nil {
		return
	}
	var i *lruItem[K, V]
	if i, ok = lru.lookup[k]; !ok {
		return
	}
	lru.moveToTail(i)
	v = i.value
	return
}

// Touch moves a given key to the end of the lru queue.
func (lru *LRU[K, V]) Touch(k K) (ok bool) {
	if lru.lookup == nil {
		return
	}

	item, ok := lru.lookup[k]
	if !ok {
		return
	}
	lru.moveToTail(item)
	return
}

// Get returns an item with a given key.
func (lru *LRU[K, V]) Set(k K, v V) {
	if lru.lookup == nil {
		lru.lookup = make(map[K]*lruItem[K, V])
	}

	if item, ok := lru.lookup[k]; ok {
		item.value = v
		lru.moveToTail(item)
		return
	}

	newItem := &lruItem[K, V]{
		key:   k,
		value: v,
	}
	lru.lookup[k] = newItem
	if lru.head == nil {
		lru.head = newItem
		lru.tail = newItem
		return
	}
	lru.moveToTail(newItem)

	if lru.Capacity > 0 && len(lru.lookup) > lru.Capacity {
		delete(lru.lookup, lru.head.key)
		if lru.OnEvict != nil {
			lru.OnEvict(lru.head.key, lru.head.value)
		}
		lru.removeHead()
	}
}

// Remove removes an element.
func (lru *LRU[K, V]) Remove(k K) (ok bool) {
	if lru.lookup == nil {
		return
	}

	var i *lruItem[K, V]
	if i, ok = lru.lookup[k]; !ok {
		return
	}
	delete(lru.lookup, k)

	if lru.head == i {
		lru.removeHead()
		return
	}
	lru.removeItem(i)
	return
}

// Head returns the head, or oldest, key and value.
func (lru *LRU[K, V]) Head() (k K, v V, ok bool) {
	if lru.head == nil {
		return
	}
	k = lru.head.key
	v = lru.head.value
	ok = true
	return
}

// Tail returns the tail, or most recently used, key and value.
func (lru *LRU[K, V]) Tail() (k K, v V, ok bool) {
	if lru.tail == nil {
		return
	}
	k = lru.tail.key
	v = lru.tail.value
	ok = true
	return
}

// Len returns the number of items in the lru cache.
func (lru *LRU[K, V]) Len() int { return len(lru.lookup) }

// Each calls a given function for each element in the lru cache.
func (lru *LRU[K, V]) Each(fn func(K, V)) {
	current := lru.head
	for current != nil {
		fn(current.key, current.value)
		current = current.previous
	}
}

//
// internal helpers
//

func (lru *LRU[K, V]) moveToTail(i *lruItem[K, V]) {
	if lru.tail == i {
		return
	}

	// remove item from existing place in list
	if lru.head == i {
		lru.head = i.previous
		lru.head.next = nil
	} else {
		after := i.previous
		before := i.next
		if after != nil {
			after.next = before
		}
		if before != nil {
			before.previous = after
		}
	}

	// append to tail
	i.next = lru.tail
	i.previous = nil
	lru.tail.previous = i
	lru.tail = i
	return
}

func (lru *LRU[K, V]) removeHead() {
	if lru.head == nil {
		return
	}

	// if we have a single element,
	// we will need to change the tail
	// pointer as well
	if lru.head == lru.tail {
		lru.head = nil
		lru.tail = nil
		return
	}

	// remove from head
	after := lru.head.previous
	if after != nil {
		after.next = nil
	}
	lru.head = after
}

func (lru *LRU[K, V]) removeItem(i *lruItem[K, V]) {
	after := i.previous
	before := i.next
	if after != nil {
		after.next = before
	}
	if before != nil {
		before.previous = after
	}
	if lru.tail == i {
		lru.tail = i.next
		if lru.tail != nil {
			lru.tail.previous = nil
		}
	}
}

type lruItem[K comparable, V any] struct {
	key   K
	value V

	// next points towards the head
	next *lruItem[K, V]
	// previous points towards the tail
	previous *lruItem[K, V]
}
