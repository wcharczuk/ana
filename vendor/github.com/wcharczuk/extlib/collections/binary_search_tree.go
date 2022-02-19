/*

Copyright (c) 2022 - Present. Will Charczuk. All rights reserved.
Use of this source code is governed by a MIT license that can be found in the
LICENSE file at the root of the repository.

*/

package collections

// BinarySearchTree is a AVL balanced tree which holds the properties
// that nodes are ordered left to right.
//
// The choice to use AVL to balance the tree means the use cases skew
// towards fast lookups at the expense of more costly mutations.
type BinarySearchTree[K Ordered, V any] struct {
	root *BinarySearchTreeNode[K, V]
}

// Insert adds a new value to the binary search tree.
func (bst *BinarySearchTree[K, V]) Insert(k K, v V) {
	bst.root = bst._insert(bst.root, k, v)
}

// Delete deletes a value from the tree, and returns if it existed.
func (bst *BinarySearchTree[K, V]) Delete(k K) {
	bst.root = bst._delete(bst.root, k)
	return
}

// Search searches for a node with a given key, returning the value
// and a boolean indicating the key was found.
func (bst *BinarySearchTree[K, V]) Search(k K) (v V, ok bool) {
	v, ok = bst._search(bst.root, k)
	return
}

// Min returns the minimum key and value.
func (bst *BinarySearchTree[K, V]) Min() (k K, v V, ok bool) {
	if bst.root == nil {
		return
	}
	k, v, ok = bst.root.Key, bst.root.Value, true
	current := bst.root
	for current.Left != nil {
		current = current.Left
		k, v = current.Key, current.Value
	}
	return
}

// Max returns the maximum key and value.
func (bst *BinarySearchTree[K, V]) Max() (k K, v V, ok bool) {
	if bst.root == nil {
		return
	}
	k, v, ok = bst.root.Key, bst.root.Value, true
	current := bst.root
	for current.Right != nil {
		current = current.Right
		k, v = current.Key, current.Value
	}
	return
}

// InOrder traversal returns the sorted values in the tree.
func (bst *BinarySearchTree[K, V]) InOrder(fn func(K, V)) {
	bst._inOrder(bst.root, fn)
}

// PreOrder traversal returns the values in the tree in pre-order.
func (bst *BinarySearchTree[K, V]) PreOrder(fn func(K, V)) {
	bst._preOrder(bst.root, fn)
}

// PostOrder traversal returns the values in the tree in post-order.
func (bst *BinarySearchTree[K, V]) PostOrder(fn func(K, V)) {
	bst._postOrder(bst.root, fn)
}

// KeysEqual is a function that can be used to deeply compare two trees based on their keys.
//
// Values are _not_ considered because values are not comparable by design.
func (bst *BinarySearchTree[K, V]) KeysEqual(other *BinarySearchTree[K, V]) bool {
	return bst._keysEqual(bst.root, other.root)
}

//
// internal methods
//

func (bst *BinarySearchTree[K, V]) _height(n *BinarySearchTreeNode[K, V]) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func (bst *BinarySearchTree[K, V]) _inOrder(n *BinarySearchTreeNode[K, V], fn func(K, V)) {
	if n == nil {
		return
	}
	bst._inOrder(n.Left, fn)
	fn(n.Key, n.Value)
	bst._inOrder(n.Right, fn)
}

func (bst *BinarySearchTree[K, V]) _preOrder(n *BinarySearchTreeNode[K, V], fn func(K, V)) {
	if n == nil {
		return
	}
	fn(n.Key, n.Value)
	bst._preOrder(n.Left, fn)
	bst._preOrder(n.Right, fn)
}

func (bst *BinarySearchTree[K, V]) _postOrder(n *BinarySearchTreeNode[K, V], fn func(K, V)) {
	if n == nil {
		return
	}
	bst._postOrder(n.Left, fn)
	bst._postOrder(n.Right, fn)
	fn(n.Key, n.Value)
}

func (bst *BinarySearchTree[K, V]) _insert(n *BinarySearchTreeNode[K, V], k K, v V) *BinarySearchTreeNode[K, V] {
	if n == nil {
		return &BinarySearchTreeNode[K, V]{
			Key:    k,
			Value:  v,
			Height: 1,
		}
	}

	if k < n.Key {
		n.Left = bst._insert(n.Left, k, v)
	} else if k > n.Key {
		n.Right = bst._insert(n.Right, k, v)
	} else {
		n.Value = v
		return n
	}

	n.Height = max(bst._height(n.Left), bst._height(n.Right)) + 1

	balanceFactor := bst._getBalanceFactor(n)
	if balanceFactor > 1 && k < n.Left.Key {
		return bst._rotateRight(n)
	}
	if balanceFactor < -1 && k > n.Right.Key {
		return bst._rotateLeft(n)
	}
	if balanceFactor > 1 && k > n.Left.Key {
		n.Left = bst._rotateLeft(n.Left)
		return bst._rotateRight(n)
	}
	if balanceFactor < -1 && k < n.Right.Key {
		n.Right = bst._rotateRight(n.Right)
		return bst._rotateLeft(n)
	}
	return n
}

func (bst *BinarySearchTree[K, V]) _delete(n *BinarySearchTreeNode[K, V], k K) *BinarySearchTreeNode[K, V] {
	if n == nil {
		return nil
	}

	if k < n.Key {
		n.Left = bst._delete(n.Left, k)
	} else if k > n.Key {
		n.Right = bst._delete(n.Right, k)
	} else {
		if n.Left == nil || n.Right == nil {
			var temp *BinarySearchTreeNode[K, V]
			if n.Left == nil {
				temp = n.Right
			} else {
				temp = n.Left
			}
			if temp == nil {
				temp = n
				n = nil
			} else {
				n = temp
			}
		} else {
			temp := bst._searchMin(n.Right)
			n.Key, n.Value = temp.Key, temp.Value
			n.Right = bst._delete(n.Right, temp.Key)
		}
	}

	if n == nil {
		return nil
	}

	n.Height = max(bst._height(n.Left), bst._height(n.Right)) + 1

	balanceFactor := bst._getBalanceFactor(n)
	if balanceFactor > 1 && bst._getBalanceFactor(n.Left) >= 0 {
		return bst._rotateRight(n)
	}
	if balanceFactor > 1 && bst._getBalanceFactor(n.Left) < 0 {
		n.Left = bst._rotateLeft(n.Left)
		return bst._rotateRight(n)
	}

	if balanceFactor < -1 && bst._getBalanceFactor(n.Right) <= 0 {
		return bst._rotateLeft(n)
	}

	if balanceFactor < -1 && bst._getBalanceFactor(n.Right) > 0 {
		n.Right = bst._rotateRight(n.Right)
		return bst._rotateLeft(n)
	}
	return n
}

func (bst *BinarySearchTree[K, V]) _searchMin(n *BinarySearchTreeNode[K, V]) (min *BinarySearchTreeNode[K, V]) {
	min = n
	for min.Left != nil {
		min = min.Left
	}
	return
}

func (bst *BinarySearchTree[K, V]) _search(n *BinarySearchTreeNode[K, V], k K) (v V, ok bool) {
	if n == nil {
		return
	}
	if n.Key == k {
		v = n.Value
		ok = true
		return
	}
	if k < n.Key {
		v, ok = bst._search(n.Left, k)
		return
	}
	v, ok = bst._search(n.Right, k)
	return
}

func (bst *BinarySearchTree[K, V]) _rotateRight(y *BinarySearchTreeNode[K, V]) *BinarySearchTreeNode[K, V] {
	if y.Left == nil {
		return y
	}
	x := y.Left
	t2 := x.Right
	x.Right = y
	y.Left = t2
	y.Height = max(bst._height(y.Left), bst._height(y.Right)) + 1
	x.Height = max(bst._height(x.Left), bst._height(x.Right)) + 1
	return x
}

func (bst *BinarySearchTree[K, V]) _rotateLeft(x *BinarySearchTreeNode[K, V]) *BinarySearchTreeNode[K, V] {
	if x.Right == nil {
		return x
	}

	y := x.Right
	t2 := y.Left
	y.Left = x
	x.Right = t2
	x.Height = max(bst._height(x.Left), bst._height(x.Right)) + 1
	y.Height = max(bst._height(y.Left), bst._height(y.Right)) + 1
	return y
}

func (bst *BinarySearchTree[K, V]) _getBalanceFactor(n *BinarySearchTreeNode[K, V]) int {
	if n == nil {
		return 0
	}
	return bst._height(n.Left) - bst._height(n.Right)
}

func (bst *BinarySearchTree[K, V]) _keysEqual(a, b *BinarySearchTreeNode[K, V]) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}
	if a.Key != b.Key {
		return false
	}
	if a.Height != b.Height {
		return false
	}
	return bst._keysEqual(a.Left, b.Left) && bst._keysEqual(a.Right, b.Right)
}

// BinarySearchTreeNode is a node in a BinarySearchTree.
type BinarySearchTreeNode[K Ordered, V any] struct {
	Key    K
	Value  V
	Left   *BinarySearchTreeNode[K, V]
	Right  *BinarySearchTreeNode[K, V]
	Height int
}

func max[K Ordered](keys ...K) (k K) {
	if len(keys) == 0 {
		return
	}
	k = keys[0]
	for x := 1; x < len(keys); x++ {
		if keys[x] > k {
			k = keys[x]
		}
	}
	return
}
