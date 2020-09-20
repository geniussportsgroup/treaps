// Package Treap exports a ordered set of arbitrary keys implemented through treaps.
// A treap is a kind of balanced binary Search tree where their operations are O(log n).
package treaps

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/golang-collections/collections/stack"
)

const notFound = -1

// The structure of every node
type Node struct {
	key      interface{} // generic key
	priority uint64      // priority value for heap order balancing
	count    int         // number of nodes that I, as tree, contain
	llink    *Node       // left child pointer
	rlink    *Node       // right child pointer
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
	less          func(i1, i2 interface{}) bool
}

// Swap two treaps in O(1)
func (tree *Treap) Swap(rhs *Treap) {
	tree.seed, rhs.seed = rhs.seed, tree.seed
	tree.randGenerator, rhs.randGenerator = rhs.randGenerator, tree.randGenerator
	tree.rootPtr, rhs.rootPtr = rhs.rootPtr, tree.rootPtr
	tree.headPtr, rhs.headPtr = rhs.headPtr, tree.headPtr
}

// Create a new treap with a random generator set to seed and comparison function less
func New(seed int64, less func(i1, i2 interface{}) bool) *Treap {

	src := rand.NewSource(seed)
	tree := &Treap{
		seed:          seed,
		randGenerator: rand.New(src),
		less:          less,
	}

	tree.head.llink = nullNodePtr
	tree.head.rlink = nullNodePtr
	tree.headPtr = &(tree.head)
	tree.rootPtr = &(tree.headPtr.rlink)

	return tree
}

// Create a new tree with random seed chosen from system clock
func NewTreap(less func(i1, i2 interface{}) bool) *Treap {
	return New(time.Now().UTC().UnixNano(), less)
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

// Get an exact Copy of tree
func (tree *Treap) Copy() *Treap {

	ret := New(tree.seed, tree.less)
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

	if !(!less(t1.key, t2.key) && !less(t2.key, t1.key)) {
		fmt.Println(t1.key, t2.key)
		return false // keys are different
	}

	return __topologicalEqual(t1.llink, t2.llink, less) &&
		__topologicalEqual(t1.rlink, t2.rlink, less)
}

// Return true if tree is topologically equivalent to rhs
func (tree *Treap) TopologicalEqual(rhs *Treap) bool {
	return __topologicalEqual(*tree.rootPtr, *rhs.rootPtr, tree.less)
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

	result := __insertNode(*tree.rootPtr, p, tree.less)
	if result == nullNodePtr {
		return nil
	}

	*tree.rootPtr = result
	return p.key
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

	resultNode = __insertNodeDup(root.rlink, p, less)
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

	result := __insertNodeDup(*tree.rootPtr, p, tree.less)

	*tree.rootPtr = result
	return p.key
}

// Search in tree key. If key is found, then the value contained in the set is returned.
// Otherwise, the key was not found, nil value is returned
func (tree *Treap) Search(key interface{}) interface{} {

	root := *tree.rootPtr
	for root != nullNodePtr {

		if tree.less(key, root.key) {
			root = root.llink
		} else if tree.less(root.key, key) {
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

	result := __searchOrInsertNode(tree.rootPtr, p, tree.less)
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
		return retVal
	}

	if retVal == nil {
		return nullNodePtr // key not found
	}

	retVal.llink = nullNodePtr
	retVal.rlink = nullNodePtr
	retVal.count = 1
	(*rootPtr).count--

	return retVal
}

// Remove key from the tree. Return the removed value if the removal was successful.
// Otherwise, the item was not found and the value nil is returned as signal of the failure
func (tree *Treap) Remove(key interface{}) interface{} {

	retVal := __remove(tree.rootPtr, key, tree.less)
	if retVal == nullNodePtr {
		return nil // key not found
	}

	return retVal.key
}

// Return the smallest item contained in the tree
func (tree *Treap) Min() interface{} {

	root := *tree.rootPtr
	if root == nullNodePtr {
		panic("Treap is empty")
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
		panic("Treap is empty")
	}

	for root.rlink != nullNodePtr {
		root = root.rlink
	}

	return root.key
}

// Return in O(1) the number of keys contained in the tree
func (tree *Treap) Size() int { return (*tree.rootPtr).count }

// Helper function for splitting a tree according to key. The function returns two new trees.
// tsRoot contains all the keys less than key and tgRoot contains the keys greater or equal to
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
// [tree.Min(), key> and tgTree contains those ones in [key, tree.Max]. After completion,
// tree becomes empty.
func (tree *Treap) SplitByKey(key interface{}) (tsTree, tgTree *Treap) {

	tsTree = New(tree.seed, tree.less)
	tgTree = New(tree.seed, tree.less)

	*tsTree.rootPtr, *tgTree.rootPtr = __splitByKeyDup(*tree.rootPtr, key, tree.less)

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

	if tsTree.Size() != 0 && tgTree.Size() != 0 && !tsTree.less(tsTree.Max(), tgTree.Min()) {
		panic("Trees are not range-disjoint")
	}

	*tsTree.rootPtr = __joinExclusive(tsTree.rootPtr, tgTree.rootPtr)
	*tgTree.rootPtr = nullNodePtr
}

// Return the key located in the position pos respect to the order of the keys.
// The item is retrieved in O(log n) expected time.
// Panic if pos is greater or equal to the number of elements stored into the tree
func (tree *Treap) Choose(pos int) interface{} {

	root := *tree.rootPtr
	if pos >= root.count {
		panic(fmt.Sprintf("Position %d out of range", pos))
	}

	for i := pos; i != root.llink.count; {
		if i < root.llink.count {
			root = root.llink
		} else {
			i -= root.llink.count + 1
			root = root.rlink
		}
	}

	return root.key
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

	pos = __rank(*tree.rootPtr, key, tree.less)
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

	ts = New(tree.seed, tree.less)
	tg = New(tree.seed, tree.less)

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
	s    stack.Stack
}

func advanceToMin(it *Iterator, p *Node) *Node {
	for p.llink != nullNodePtr {
		it.s.Push(p)
		p = p.llink
	}
	return p
}

func advanceToMax(it *Iterator, p *Node) *Node {
	for p != nullNodePtr {
		p = p.rlink
	}
	return p
}

func emptyStack(it *Iterator) {
	for it.s.Len() > 0 {
		it.s.Pop()
	}
}

func initialize(it *Iterator) {
	if it.root != nullNodePtr {
		it.curr = advanceToMin(it, it.root)
	} else {
		it.curr = nullNodePtr
	}
}

// Return a iterator on the treap tree
func NewIterator(tree *Treap) *Iterator {
	it := &Iterator{
		root: *tree.rootPtr,
		curr: nil,
		s:    stack.Stack{},
	}
	initialize(it)
	return it
}

// Reset the iterator to the first item of the set
func (it *Iterator) ResetFirst() {
	emptyStack(it)
	initialize(it)
}

// Reset the iterator to the last item of the set
func (it *Iterator) ResetLast() {
	emptyStack(it)
	advanceToMax(it, it.root)
}

// Return true if iterator is positioned on an item. Otherwise it return false
func (it *Iterator) HasCurr() bool {
	return it.curr != nullNodePtr
}

// Return the current item on which the iterator is positioned. Panic if there is not current item
func (it *Iterator) GetCurr() interface{} {
	if !it.HasCurr() {
		panic("Iterator has not current item")
	}
	return it.curr.key
}

// Advance iterator to the next item in the ordered sequence
func (it *Iterator) Next() *Iterator {
	if !it.HasCurr() {
		panic("Iterator has not current item")
	}
	it.curr = it.curr.rlink
	if it.curr != nullNodePtr {
		it.curr = advanceToMin(it, it.curr)
		return it
	}

	if it.s.Len() == 0 {
		it.curr = nullNodePtr
	} else {
		it.curr = it.s.Pop().(*Node)
	}

	return it
}

// TODO Implement Prev method

// Simple BST checker; Not completely correct
func checkBST(node *Node, less func(i1, i2 interface{}) bool) bool {

	if node == nullNodePtr {
		return true
	}

	if node.llink != nullNodePtr {
		if !less(node.llink.key, node.key) {
			return false
		}
		if !checkBST(node.llink, less) {
			return false
		}
	}

	if node.rlink != nullNodePtr {
		if !less(node.key, node.rlink.key) {
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
	return checkBST(p, less) && checkTreap(p) && checkCounter(p)
}

func (tree *Treap) check() bool {

	if tree.head.llink != nullNodePtr {
		return false
	}

	if tree.headPtr != &tree.head {
		return false
	}

	if tree.rootPtr != &(tree.headPtr.rlink) {
		return false
	}

	return true
}
