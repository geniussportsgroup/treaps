// Package Treap exports a ordered set of arbitrary keys implemented through treaps.
// A treap is a kind of balanced binary Search tree where their operations are O(log n).
package treaps

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const notFound = -1

// Node The structure of every node
type Node struct {
	key      interface{} // generic key
	priority uint64      // priority value for heap order balancing
	count    int         // number of nodes that I, as tree, contain
	llink    *Node       // left child pointer
	rlink    *Node       // right child pointer
}

func (p *Node) swap(q *Node) {
	p.key, q.key = q.key, p.key
	p.priority, q.priority = q.priority, p.priority
	p.count, q.count = q.count, p.count
	p.llink, q.llink = q.llink, p.llink
	p.rlink, q.rlink = q.rlink, p.rlink
}

func (p *Node) reset() {
	p.llink = nullNodePtr
	p.rlink = nullNodePtr
	p.count = 1
}

// This node, supposed to be immutable, represents the empty tree, as well as an
// external node
var nullNodePtr *Node = &Node{
	key:      nil,
	priority: math.MaxUint64, // Empty tree always has maximum priority value
	count:    0,              // empty tree has zero nodes
	llink:    nil,
	rlink:    nil,
}

// The Treap object which represents a set of ordered keys whose operations exhibit
// O(log n) expected complexity
type Treap struct {
	seed          int64
	randGenerator *rand.Rand
	rootPtr       **Node
	head          Node // header node dummy parent of rootPtr
	headPtr       *Node
	Less          func(i1, i2 interface{}) bool
}

// helper for implementing == with < operation
func __equal(i1, i2 interface{}, less func(i1, i2 interface{}) bool) bool {
	return !less(i1, i2) && !less(i2, i1)
}

// helper for implementing <= with only < operation
func __lessOrEqual(i1, i2 interface{}, less func(i1, i2 interface{}) bool) bool {
	return less(i1, i2) || __equal(i1, i2, less)
}

func __greater(i1, i2 interface{}, less func(i1, i2 interface{}) bool) bool {
	return less(i2, i1)
}

func __greaterOrEqual(i1, i2 interface{}, less func(i1, i2 interface{}) bool) bool {
	return __greater(i1, i2, less) || __equal(i1, i2, less)
}

// Swap two treaps in O(1)
func (tree *Treap) Swap(other interface{}) interface{} {

	rhs := other.(*Treap)
	tree.seed, rhs.seed = rhs.seed, tree.seed
	tree.randGenerator, rhs.randGenerator = rhs.randGenerator, tree.randGenerator
	*tree.rootPtr, *rhs.rootPtr = *rhs.rootPtr, *tree.rootPtr
	tree.Less, rhs.Less = rhs.Less, tree.Less
	return tree
}

// New Create a new treap with a random generator set to seed and comparison function less
func New(seed int64, less func(i1, i2 interface{}) bool, items ...interface{}) *Treap {

	src := rand.NewSource(seed)
	tree := &Treap{
		seed:          seed,
		randGenerator: rand.New(src),
		Less:          less,
	}

	tree.head.llink = nullNodePtr
	tree.head.rlink = nullNodePtr
	tree.headPtr = &(tree.head)
	tree.rootPtr = &(tree.headPtr.rlink)

	for _, item := range items {
		tree.InsertDup(item)
	}

	return tree
}

// Clear Empty the set
func (tree *Treap) Clear() {
	*tree.rootPtr = nullNodePtr
}

// IsEmpty Return true is set is empty
func (tree *Treap) IsEmpty() bool { return *tree.rootPtr == nullNodePtr }

// NewTreap Create a new tree with random seed chosen from system clock
func NewTreap(less func(i1, i2 interface{}) bool, items ...interface{}) *Treap {
	return New(time.Now().UTC().UnixNano(), less, items...)
}

func (tree *Treap) Create(items ...interface{}) interface{} {
	return New(time.Now().UTC().UnixNano(), tree.Less, items...)
}

// Helper function that perform an exact topological Copy of tree rooted by p
func __copy(p *Node) *Node {

	if p == nullNodePtr {
		return nullNodePtr
	}

	return &Node{
		key:      p.key,
		priority: p.priority,
		count:    p.count,
		llink:    __copy(p.llink),
		rlink:    __copy(p.rlink),
	}
}

// Copy Get an exact Copy of tree
func (tree *Treap) Copy() *Treap {

	ret := New(tree.seed, tree.Less)
	*ret.rootPtr = __copy(*tree.rootPtr)

	return ret
}

// Helper for topological comparison of two trees
func __topologicalEqual(t1, t2 *Node, less func(i1, i2 interface{}) bool) bool {

	if t1 == nullNodePtr && t2 == nullNodePtr {
		return true
	}

	if (t1 == nullNodePtr && t2 != nullNodePtr) || (t1 != nullNodePtr && t2 == nullNodePtr) {
		return false
	}

	if !__equal(t1.key, t2.key, less) {
		return false // keys are different
	}

	return __topologicalEqual(t1.llink, t2.llink, less) &&
		__topologicalEqual(t1.rlink, t2.rlink, less)
}

// Return true if tree is topologically equivalent to rhs
func (tree *Treap) TopologicalEqual(rhs *Treap) bool {
	return __topologicalEqual(*tree.rootPtr, *rhs.rootPtr, tree.Less)
}

// Helper for inserting node p into the tree root. BST order is handled through less function
func __insertNode(root, p *Node, less func(i1, i2 interface{}) bool) *Node {

	if root == nullNodePtr {
		return p
	}

	resultNode := nullNodePtr
	if less(p.key, root.key) {
		resultNode = __insertNode(root.llink, p, less)
		if resultNode == nullNodePtr { // was p inserted?
			return nullNodePtr // key is already in tree ==> insertion fails
		}

		root.llink = resultNode
		root.count++
		if resultNode.priority < root.priority {
			root = rotateRight(root)
		}
		return root
	}

	if less(root.key, p.key) {
		resultNode = __insertNode(root.rlink, p, less)
		if resultNode == nullNodePtr { // was p inserted?
			return nullNodePtr // key is already in tree ==> insertion fails
		}

		root.rlink = resultNode
		root.count++
		if resultNode.priority < root.priority {
			root = rotateLeft(root)
		}
		return root
	}

	return nullNodePtr // key is already in tree ==> insertion fails
}

// Insert item into the tree. Return nil if key is already contained; otherwise
// returns the value of the just inserted item
func (tree *Treap) Insert(item interface{}) interface{} {

	p := &Node{
		key:      item,
		priority: tree.randGenerator.Uint64(),
		count:    1,
		llink:    nullNodePtr,
		rlink:    nullNodePtr,
	}

	result := __insertNode(*tree.rootPtr, p, tree.Less)
	if result == nullNodePtr {
		return nil
	}

	*tree.rootPtr = result
	return p.key
}

// Append equivalent to insert. Put for supporting functional operations
func (tree *Treap) Append(item interface{}, items ...interface{}) interface{} {
	tree.Insert(item)
	for _, i := range items {
		tree.Insert(i)
	}
	return tree
}

// Helper for inserting node p into the tree root. BST order is handled through less function.
// key stored in p can be already present in the tree,. In this case, The key will be duplicated
func __insertNodeDup(root, p *Node, less func(i1, i2 interface{}) bool) *Node {

	if root == nullNodePtr {
		return p
	}

	resultNode := nullNodePtr
	if less(p.key, root.key) {
		resultNode = __insertNodeDup(root.llink, p, less)
		root.llink = resultNode
		root.count++
		if resultNode.priority < root.priority {
			root = rotateRight(root)
		}
		return root
	}

	resultNode = __insertNodeDup(root.rlink, p, less)
	root.rlink = resultNode
	root.count++
	if resultNode.priority < root.priority {
		root = rotateLeft(root)
	}

	return root
}

// Insert item into the tree. Return nil if key is already contained; otherwise
// returns the value of the just inserted item
func (tree *Treap) InsertDup(item interface{}) interface{} {

	p := &Node{
		key:      item,
		priority: tree.randGenerator.Uint64(),
		count:    1,
		llink:    nullNodePtr,
		rlink:    nullNodePtr,
	}

	result := __insertNodeDup(*tree.rootPtr, p, tree.Less)

	*tree.rootPtr = result
	return p.key
}

// Search in tree key. If key is found, then the value contained in the set is returned.
// Otherwise, the key was not found, nil value is returned
func (tree *Treap) Search(key interface{}) interface{} {

	root := *tree.rootPtr
	for root != nullNodePtr {

		if tree.Less(key, root.key) {
			root = root.llink
		} else if tree.Less(root.key, key) {
			root = root.rlink
		} else {
			break // key found!
		}
	}

	if root == nullNodePtr {
		return nil
	}

	return root.key
}

// Return true if key is found in tree
func (tree *Treap) Has(key interface{}) bool {
	return tree.Search(key) != nil
}

// Helper function for searching a node and eventually Insert it into the tree if it is not found
func __searchOrInsertNode(root **Node, p *Node, less func(i1, i2 interface{}) bool) *Node {

	if *root == nullNodePtr {
		*root = p
		return p
	}

	if less(p.key, (*root).key) {
		ret := __searchOrInsertNode(&(*root).llink, p, less)
		if ret == p {
			(*root).count++
			if ret.priority < (*root).priority {
				*root = rotateRight(*root)
			}
		}
		return ret
	}

	if less((*root).key, p.key) {
		ret := __searchOrInsertNode(&(*root).rlink, p, less)
		if ret == p {
			(*root).count++
			if ret.priority < (*root).priority {
				*root = rotateLeft(*root)
			}
		}
		return ret
	}

	return *root // key is already in tree ==> insertion fails
}

// Search in tree item. If it is found, then the pair (false, item-value) is returned.
// Otherwise, the item is inserted into the tree and the pair (true, item) is returned
func (tree *Treap) SearchOrInsert(item interface{}) (bool, interface{}) {

	p := &Node{
		key:      item,
		priority: tree.randGenerator.Uint64(),
		count:    1,
		llink:    nullNodePtr,
		rlink:    nullNodePtr,
	}

	result := __searchOrInsertNode(tree.rootPtr, p, tree.Less)
	if result != p {
		return false, result.key
	}

	return true, p.key
}

// Helper for removing key from a tree. Returns the removed node if this one is found.
// Otherwise, nullNodePte is returned.
func __remove(rootPtr **Node, key interface{}, less func(i1, i2 interface{}) bool) *Node {

	if *rootPtr == nullNodePtr {
		return nullNodePtr
	}

	var retVal *Node
	if less(key, (*rootPtr).key) {
		retVal = __remove(&(*rootPtr).llink, key, less)
	} else if less((*rootPtr).key, key) {
		retVal = __remove(&(*rootPtr).rlink, key, less)
	} else { // key found
		retVal = *rootPtr // this node will be deleted
		*rootPtr = __joinExclusive(&(*rootPtr).llink, &(*rootPtr).rlink)
		retVal.reset()
		return retVal
	}

	if retVal == nullNodePtr {
		return nullNodePtr // key not found
	}

	(*rootPtr).count--

	return retVal
}

// Remove key from the tree. Return the removed value if the removal was successful.
// Otherwise, the item was not found and the value nil is returned as signal of the failure
func (tree *Treap) Remove(key interface{}) interface{} {

	retVal := __remove(tree.rootPtr, key, tree.Less)
	if retVal == nullNodePtr {
		return nil // key not found
	}

	return retVal.key
}

func __removePos(rootPtr **Node, i int) *Node {

	root := *rootPtr
	var retVal *Node
	if i == root.llink.count {
		retVal = root
		*rootPtr = __joinExclusive(&(*rootPtr).llink, &(*rootPtr).rlink)
		retVal.reset()
		return retVal
	} else if i < root.llink.count {
		retVal = __removePos(&(*rootPtr).llink, i)
	} else {
		retVal = __removePos(&(*rootPtr).rlink, i-(root.llink.count+1))
	}

	root.count--

	return retVal
}

func (tree *Treap) RemoveByPos(i int) interface{} {

	if i >= tree.Size() {
		panic(fmt.Sprintf("Invalid position %d", i))
	}

	retVal := __removePos(tree.rootPtr, i)
	return retVal.key
}

// Return the smallest item contained in the tree
func (tree *Treap) Min() interface{} {

	root := *tree.rootPtr
	if root == nullNodePtr {
		return nil
	}

	for root.llink != nullNodePtr {
		root = root.llink
	}

	return root.key
}

// Return the greatest item contained in the tree
func (tree *Treap) Max() interface{} {

	root := *tree.rootPtr
	if root == nullNodePtr {
		return nil
	}

	for root.rlink != nullNodePtr {
		root = root.rlink
	}

	return root.key
}

// Return in O(1) the number of keys contained in the tree
func (tree *Treap) Size() int { return (*tree.rootPtr).count }

// Helper function for splitting a tree according to key. The function returns two new trees.
// tsRoot contains all the keys less or equal than key and tgRoot contains the keys greater to
// key. The original tree in root remains in inconsistent state and it should not be used.
func __splitByKeyDup(root *Node, key interface{},
	less func(i1, i2 interface{}) bool) (tsRoot, tgRoot *Node) {

	if root == nullNodePtr {
		return nullNodePtr, nullNodePtr
	}

	if less(key, root.key) {
		tgRootAux := nullNodePtr
		tgRoot = root
		tsRoot, tgRootAux = __splitByKeyDup(root.llink, key, less)
		tgRoot.llink = tgRootAux
		tgRoot.count -= tsRoot.count
	} else {
		tsRootAux := nullNodePtr
		tsRoot = root
		tsRootAux, tgRoot = __splitByKeyDup(root.rlink, key, less)
		tsRoot.rlink = tsRootAux
		tsRoot.count -= tgRoot.count
	}
	return tsRoot, tgRoot
}

// SplitByKey tree in two trees tsTree and tgTres. tsTree contains all the keys of tree in
// [tree.Min(), key) and tgTree contains those ones in [key, tree.Max]. After completion,
// tree becomes empty.
func (tree *Treap) SplitByKey(key interface{}) (tsTree, tgTree *Treap) {

	tsTree = New(tree.seed, tree.Less)
	tgTree = New(tree.seed, tree.Less)

	*tsTree.rootPtr, *tgTree.rootPtr = __splitByKeyDup(*tree.rootPtr, key, tree.Less)

	*tree.rootPtr = nullNodePtr

	return
}

// Helper that joins two range-disjoint trees. By range-disjoint we mean that all the keys
// in tsRootPtr are less than any key in tgRootPtr. The helper returns the resulting join
// and the originals trees are emptied
func __joinExclusive(tsRootPtr, tgRootPtr **Node) *Node {

	if *tsRootPtr == nullNodePtr {
		return *tgRootPtr
	}

	if *tgRootPtr == nullNodePtr {
		return *tsRootPtr
	}

	if (*tsRootPtr).priority < (*tgRootPtr).priority {
		(*tsRootPtr).count += (*tgRootPtr).count
		(*tsRootPtr).rlink = __joinExclusive(&(*tsRootPtr).rlink, tgRootPtr)
		return *tsRootPtr
	}

	(*tgRootPtr).count += (*tsRootPtr).count
	(*tgRootPtr).llink = __joinExclusive(tsRootPtr, &(*tgRootPtr).llink)
	return *tgRootPtr
}

// join exclusive of tsTree with tgTree. Equivalent to append tgTree to tsTree.
// tgTree must be greater than tsTree. Panic is thrown if this condition is not met
func (tsTree *Treap) JoinExclusive(tgTree *Treap) {

	if tsTree.Size() != 0 && tgTree.Size() != 0 && !tsTree.Less(tsTree.Max(), tgTree.Min()) {
		panic("Trees are not range-disjoint")
	}

	*tsTree.rootPtr = __joinExclusive(tsTree.rootPtr, tgTree.rootPtr)
	*tgTree.rootPtr = nullNodePtr
}

func __joinDup(rootPtr **Node, root *Node, less func(k1, k2 interface{}) bool) {

	if root == nullNodePtr {
		return
	}

	l, r := root.llink, root.rlink
	root.llink, root.rlink, root.count = nullNodePtr, nullNodePtr, 1
	*rootPtr = __insertNodeDup(*rootPtr, root, less)
	__joinDup(rootPtr, l, less)
	__joinDup(rootPtr, r, less)
}

// join rhs with tree. The result is equivalent to the union of tree and rhs
// Notice that keys could be repeated. At the end of operation rhs becomes empty
func (tree *Treap) JoinDup(rhs *Treap) {

	__joinDup(tree.rootPtr, *rhs.rootPtr, tree.Less)
	*rhs.rootPtr = nullNodePtr
}

// Union of root tree on tree pointer by rootPtr. Keys of root that are not in rootPtr are
// copied without mutating root
func __union(rootPtr **Node, root *Node, less func(k1, k2 interface{}) bool) {

	if root == nullNodePtr {
		return
	}

	p := &Node{
		key:      root.key,
		priority: root.priority,
		count:    1,
		llink:    nullNodePtr,
		rlink:    nullNodePtr,
	}

	result := __insertNode(*rootPtr, p, less)
	if result != nullNodePtr {
		*rootPtr = result
	}
	__union(rootPtr, root.llink, less)
	__union(rootPtr, root.rlink, less)
}

// Do the union of keys of rhs with tree. The result is equivalent to the union of tree and rhs
// Notice that keys should not be repeated.
// At the end of operation the original sets become emtpy. If the keys are no repeated, then
// diff1 and diff2 contain the exact differences
func (tree *Treap) Union(rhs *Treap) {

	__union(tree.rootPtr, *rhs.rootPtr, tree.Less)
}

// helper for intersecting. root tree is traversed in preorder and its nodes inserted into
// the intersection result or in diff1. nodes of rhs belonging to the intersection are deleted.
func __intersectionPrefix(root *Node, rhsPtr, result, diff1, diff2 **Node,
	less func(k1, k2 interface{}) bool) {

	if root == nullNodePtr {
		return
	}

	key := root.key
	l, r := root.llink, root.rlink
	p1 := root
	p1.reset() // children saved in l and r
	p2 := __remove(rhsPtr, key, less)
	if p2 != nullNodePtr { // is the key in both sets?
		q := __insertNode(*result, p1, less)
		if q != nil { // p1.key could be duplicated in rootPtr. In this case we delete
			*result = q
		}
	} else {
		*diff1 = __insertNodeDup(*diff1, p1, less)
	}

	__intersectionPrefix(l, rhsPtr, result, diff1, diff2, less)
	__intersectionPrefix(r, rhsPtr, result, diff1, diff2, less)
}

// Compute the intersection of tree with rhs. Intersection is put on result and remaining keys
// are put on diff1 and diff2 respectively
func (tree *Treap) Intersection(rhs *Treap) (result, diff1, diff2 *Treap) {

	result = NewTreap(tree.Less)
	diff1 = NewTreap(tree.Less)
	diff2 = NewTreap(tree.Less)

	__intersectionPrefix(*tree.rootPtr, rhs.rootPtr, result.rootPtr,
		diff1.rootPtr, diff2.rootPtr, tree.Less)

	*tree.rootPtr = nullNodePtr
	diff2.JoinDup(rhs)

	return
}

// Return the pos-th node
func __choose(root *Node, pos int) *Node {

	for i := pos; i != root.llink.count; {
		if i < root.llink.count {
			root = root.llink
		} else {
			i -= root.llink.count + 1
			root = root.rlink
		}
	}
	return root
}

// Return the key located in the position pos respect to the order of the keys.
// The item is retrieved in O(log n) expected time.
// Panic if pos is greater or equal to the number of elements stored into the tree
func (tree *Treap) Choose(pos int) interface{} {

	root := *tree.rootPtr
	if pos >= root.count {
		panic(fmt.Sprintf("Position %d out of range", pos))
	}

	return __choose(*tree.rootPtr, pos).key
}

// Helper that computes the position of key respect to the ordered kes stored in the tree
// root. It returns nullNodePtr if key is not contained in the tree.
func __rank(root *Node, key interface{}, less func(i1, i2 interface{}) bool) int {

	if root == nullNodePtr {
		return notFound
	}

	if less(key, root.key) {
		return __rank(root.llink, key, less)
	}

	if less(root.key, key) {
		ret := __rank(root.rlink, key, less)
		if ret != notFound {
			return ret + root.llink.count + 1
		}
		return notFound
	}

	return root.llink.count // key found
}

// Compute the position of key respect to the order of the full set. If the key is found,
// then the pair (true, pos) is returned, where pos is the position of key respect to the
// order of all keys stored in the tree. Otherwise, the method returns (false, Undetermined)
// for indicating that the key is not in the tree.
// The computation spends O(log n) expected time
func (tree *Treap) RankInOrder(key interface{}) (ok bool, pos int) {

	pos = __rank(*tree.rootPtr, key, tree.Less)
	ok = pos != notFound
	return
}

// Helper that SplitByKey tree root by position i. l = [0, i] r = [i + 1, N - 1]
func __splitPos(root *Node, i int) (l, r *Node) {

	if i == root.llink.count {
		l = root
		r = root.rlink
		l.rlink = nullNodePtr
		l.count -= r.count
		return
	}

	if i < root.llink.count {
		lp, rp := __splitPos(root.llink, i)
		l = lp
		r = root
		r.llink = rp
		r.count -= l.count
	} else {
		lp, rp := __splitPos(root.rlink, i-(root.llink.count+1))
		r = rp
		l = root
		l.rlink = lp
		l.count -= r.count
	}
	return
}

// SplitByKey tree in ts = [Min, i] and tg = (i, Max). After operation tree becomes empty
func (tree *Treap) SplitByPosition(i int) (ts, tg *Treap) {

	root := *tree.rootPtr
	if i < 0 || i >= root.count {
		panic(fmt.Sprintf("Position %d out of range", i))
	}

	ts = New(tree.seed, tree.Less)
	tg = New(tree.seed, tree.Less)

	if i == root.count-1 {
		*ts.rootPtr = *tree.rootPtr
		*tree.rootPtr = nullNodePtr
		*tg.rootPtr = nullNodePtr
		return
	}

	*ts.rootPtr, *tg.rootPtr = __splitPos(*tree.rootPtr, i)
	*tree.rootPtr = nullNodePtr

	return
}

// Extract from tree all the keys in [beginPos, endPos]. tree looses the extracted range
func (tree *Treap) ExtractRange(beginPos, endPos int) *Treap {

	if beginPos > endPos || endPos > (*tree.rootPtr).count-1 {
		panic(fmt.Sprintf("Invalid positions %d %d respect to number of keys %d",
			beginPos, endPos, (*tree.rootPtr).count))
	}

	begPos := beginPos - 1
	if beginPos == 0 {
		begPos = 0
	}
	treeAux, endTree := tree.SplitByPosition(endPos)
	beginTree, result := treeAux.SplitByPosition(begPos)

	beginTree.JoinExclusive(endTree)

	tree.Swap(beginTree)

	return result
}

func (tree *Treap) lexicographicCmp(rhs *Treap) int {

	it1, it2 := NewIterator(tree), NewIterator(rhs)
	for it1.HasCurr() && it2.HasCurr() {
		item1 := it1.GetCurr()
		item2 := it2.GetCurr()
		if tree.Less(item1, item2) {
			return -1
		} else if tree.Less(item2, item1) {
			return 1
		}
		it1.Next()
		it2.Next()
	}

	if !it1.HasCurr() && !it2.HasCurr() {
		return 0
	}

	if it1.HasCurr() {
		return 1
	}

	return -1
}

// Rotate p to the right. Left child becomes root
func rotateRight(p *Node) *Node {
	q := p.llink
	p.llink = q.rlink
	q.rlink = p
	p.count -= 1 + q.llink.count
	q.count += 1 + p.rlink.count
	return q
}

// Rotate p to the left. Right child becomes root
func rotateLeft(p *Node) *Node {
	q := p.rlink
	p.rlink = q.llink
	q.llink = p
	p.count -= 1 + q.rlink.count
	q.count += 1 + p.llink.count
	return q
}

// Iterator on Treap. Traversal is ordered
type Iterator struct {
	root *Node
	curr *Node
	pos  int
	N    int
}

// Initialize a treap iterator
func initialize(it *Iterator) {
	if it.N <= 0 {
		return
	}
	it.curr = __choose(it.root, 0)
	it.pos = 0
}

func (tree *Treap) CreateIterator() interface{} {

	return NewIterator(tree)
}

// Return a iterator on the treap tree
func NewIterator(tree *Treap) *Iterator {
	it := &Iterator{
		root: *tree.rootPtr,
		curr: nil,
		pos:  -1,
		N:    tree.Size(),
	}
	initialize(it)
	return it
}

func NewReverseIterator(tree *Treap) *Iterator {
	it := &Iterator{
		root: *tree.rootPtr,
		curr: nil,
		pos:  -1,
		N:    tree.Size(),
	}

	return it.ResetLast()
}

// Reset the iterator to the first item of the set
func (it *Iterator) ResetFirst() interface{} {
	initialize(it)
	return it
}

// Reset the iterator to the last item of the set
func (it *Iterator) ResetLast() *Iterator {
	if it.N == 0 {
		panic("Tree is empty")
	}
	it.pos = it.N - 1
	it.curr = __choose(it.root, it.pos)

	return it
}

func (it *Iterator) getPos() int { return it.pos }

// Return true if iterator is positioned on an item. Otherwise it return false
func (it *Iterator) HasCurr() bool {
	return it.pos >= 0 && it.pos < it.N
}

// Return the current item on which the iterator is positioned. Panic if there is not current item
func (it *Iterator) GetCurr() interface{} {
	if !it.HasCurr() {
		panic("Iterator has not current item")
	}
	return it.curr.key
}

// Advance iterator to the next item in the ordered sequence
func (it *Iterator) Next() interface{} {
	if it.pos == it.N {
		panic("Iterator overflow")
	}

	it.pos++
	if it.pos == it.N {
		it.curr = nullNodePtr
		return it
	}

	it.curr = __choose(it.root, it.pos)
	return it
}

// Advance iterator to the previous item in the ordered sequence
func (it *Iterator) Prev() *Iterator {
	if it.pos == -1 {
		panic("Iterator underflow")
	}

	it.pos--
	if it.pos == -1 {
		it.curr = nullNodePtr
		return it
	}

	it.curr = __choose(it.root, it.pos)
	return it
}

// Traverse inorder the whole set and execute operation on each key.
// The function stops if operation return false. Otherwise the function continues toward the
// following key.
// Return true if all the set was traversed, false otherwise.
// WARNING: it is not supposed that operation might modify the key
func (tree *Treap) Traverse(operation func(key interface{}) bool) bool {

	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		if !operation(it.GetCurr()) {
			return false
		}
	}

	return true
}

// Simple BST checker; Not completely correct
func checkBST(node *Node, less func(i1, i2 interface{}) bool) bool {

	if node == nullNodePtr {
		return true
	}

	if node.llink != nullNodePtr {
		if !less(node.llink.key, node.key) && !__equal(node.llink.key, node.key, less) {
			return false
		}
		if !checkBST(node.llink, less) {
			return false
		}
	}

	if node.rlink != nullNodePtr {
		if !less(node.key, node.rlink.key) && !__equal(node.key, node.rlink.key, less) {
			return false
		}
		if !checkBST(node.rlink, less) {
			return false
		}
	}

	return true
}

// Simple priority checker
func checkTreap(node *Node) bool {

	if node == nullNodePtr {
		return true
	}

	if node.priority > node.llink.priority || node.priority > node.rlink.priority {
		return false
	}

	return checkTreap(node.llink) && checkTreap(node.rlink)
}

func checkCounter(p *Node) bool {

	if p == nullNodePtr {
		return true
	}

	if p.llink.count+1+p.rlink.count != p.count {
		return false
	}

	return checkCounter(p.llink) && checkCounter(p.rlink)
}

func checkAll(p *Node, less func(i1, i2 interface{}) bool) bool {

	// put thus for making debugging easier
	if !checkBST(p, less) {
		return false
	}
	if !checkTreap(p) {
		return false
	}
	if !checkCounter(p) {
		return false
	}
	return true
}

func (tree *Treap) check() bool {

	if !checkSentinel() {
		return false
	}

	if tree.head.llink != nullNodePtr {
		return false
	}

	if tree.headPtr != &tree.head {
		return false
	}

	if tree.rootPtr != &(tree.headPtr.rlink) {
		return false
	}

	return checkAll(*tree.rootPtr, tree.Less)
}

func checkSentinel() bool {
	return nullNodePtr.key == nil && nullNodePtr.llink == nil && nullNodePtr.rlink == nil &&
		nullNodePtr.count == 0
}
