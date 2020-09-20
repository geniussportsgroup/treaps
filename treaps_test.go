package treaps

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func cmpInt(i1, i2 interface{}) bool {
	item1, ok := i1.(int)
	if !ok {
		panic("First parameter is not int")
	}
	item2, ok := i2.(int)
	if !ok {
		panic("Second parameter is not int")
	}
	return item1 < item2
}

func TestNewTreap(t *testing.T) {

	s := New(9, cmpInt)
	assert.NotNil(t, s)
	assert.Equal(t, 0, s.size())
}

func TestTreap_insert(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 500

	for i := 0; i < N; i++ {
		ret := tree.insert(i)
		assert.NotNil(t, ret)
	}
}

func insertNRandomItems(tree *Treap, n int) {
	for i := 0; i < n; i++ {
		val := rand.Intn(100 * n)
		for tree.insert(val) == nil {
			val = rand.Intn(100 * n)
		}
	}
}

func TestRandomInsertions(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 500
	for i := 0; i < N; i++ {
		val := rand.Intn(100 * N)
		for tree.insert(val) == nil {
			val = rand.Intn(100 * N)
		}
		assert.Equal(t, val, tree.search(val))
	}

	assert.Equal(t, N, int(tree.size()))
}

func TestTreap_remove(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100

	insertNRandomItems(tree, N)

	values := make([]int, 0, tree.size())
	for it := NewIterator(tree); it.hasCurr(); it.next() {
		values = append(values, it.getCurr().(int))
	}

	for _, val := range values {
		assert.Equal(t, val, tree.search(val))
		assert.Equal(t, val, tree.remove(val))
		assert.Equal(t, nil, tree.search(val))
	}
}

func TestTreap_split(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 100
	insertNRandomItems(tree, N)

	t1, t2 := tree.split(552)

	assert.Equal(t, N, t1.size()+t2.size())
	assert.Equal(t, 0, tree.size())

	for it := NewIterator(t1); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	tree.joinExclusive(t1)

	for it := NewIterator(tree); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	tree.joinExclusive(t2)
	assert.Equal(t, 0, t1.size())
	assert.Equal(t, 0, t2.size())

	for it := NewIterator(tree); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()
}

func TestTreap_searchOrInsert(t *testing.T) {

	const N = 1000
	tree := New(1, cmpInt)
	failures := New(2, cmpInt)

	for i := 0; i < N; i++ {
		val := rand.Intn(100 * N)
		ok, res := tree.searchOrInsert(val)
		assert.Equal(t, val, res)
		if !ok {
			failures.insert(val)
		}
	}

	fmt.Println("tree.size() = ", tree.size())
	fmt.Println("failures.size() = ", failures.size())

	for it := NewIterator(failures); it.hasCurr(); it.next() {
		assert.Equal(t, tree.search(it.getCurr()), it.getCurr())
	}
}

func TestTreap_choose(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 1000
	insertNRandomItems(tree, N)

	for i, it := 0, NewIterator(tree); it.hasCurr(); it.next() {
		item := tree.choose(i)
		assert.Equal(t, item, it.getCurr())
		i++
	}
}

func TestTreap_rank(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 1000
	insertNRandomItems(tree, N)

	for i, it := 0, NewIterator(tree); it.hasCurr(); it.next() {
		ok, pos := tree.rank(it.getCurr())
		assert.True(t, ok)
		assert.Equal(t, i, pos)
		i++
	}
}

func TestTreap_splitPos(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100
	insertNRandomItems(tree, N)

	min1, max1, min2, max2 := tree.min(), tree.choose(N/2-1), tree.choose(N/2), tree.max()

	t1, t2 := tree.splitPos(N / 2)

	for it := NewIterator(t1); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	assert.Equal(t, 0, tree.size())
	assert.NotNil(t, t1)
	assert.NotNil(t, t2)
	assert.Equal(t, N/2, t1.size())
	assert.Equal(t, N/2, t2.size())
	assert.Equal(t, min1, t1.min())
	assert.Equal(t, max1, t1.max())
	assert.Equal(t, min2, t2.min())
	assert.Equal(t, max2, t2.max())
}

func TestTreap_copy(t *testing.T) {
	t1 := New(2, cmpInt)
	const N = 100
	insertNRandomItems(t1, N)

	assert.True(t, checkBST(*t1.rootPtr, t1.less))
	assert.True(t, checkTreap(*t1.rootPtr))
	assert.True(t, checkCounter(*t1.rootPtr))

	t2 := t1.copy()

	for it := NewIterator(t1); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

	assert.True(t, t1.topologicalEqual(t2))
}

func TestTreap_removeRange(t *testing.T) {
	tree := New(2, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.insert(i)
	}

	midRange := tree.extractRange(40, 60)

	for key, it := 40, NewIterator(midRange); it.hasCurr(); it.next() {
		assert.Equal(t, key, it.getCurr())
		key++
	}
	fmt.Println()

	for it := NewIterator(tree); it.hasCurr(); it.next() {
		fmt.Print(it.getCurr(), " ")
	}
	fmt.Println()

}
