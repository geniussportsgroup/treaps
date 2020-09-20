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
	assert.Equal(t, 0, s.Size())
}

func TestTreap_insert(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 500

	for i := 0; i < N; i++ {
		ret := tree.Insert(i)
		assert.NotNil(t, ret)
	}
}

func insertNRandomItems(tree *Treap, n int) {
	for i := 0; i < n; i++ {
		val := rand.Intn(100 * n)
		for tree.Insert(val) == nil {
			val = rand.Intn(100 * n)
		}
	}
}

func TestRandomInsertions(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 500
	for i := 0; i < N; i++ {
		val := rand.Intn(100 * N)
		for tree.Insert(val) == nil {
			val = rand.Intn(100 * N)
		}
		assert.Equal(t, val, tree.Search(val))
	}

	assert.Equal(t, N, int(tree.Size()))
}

func TestTreap_remove(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100

	insertNRandomItems(tree, N)

	values := make([]int, 0, tree.Size())
	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		values = append(values, it.GetCurr().(int))
	}

	for _, val := range values {
		assert.Equal(t, val, tree.Search(val))
		assert.Equal(t, val, tree.Remove(val))
		assert.Equal(t, nil, tree.Search(val))
	}
}

func TestTreap_split(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 100
	insertNRandomItems(tree, N)

	t1, t2 := tree.SplitByKey(552)

	assert.Equal(t, N, t1.Size()+t2.Size())
	assert.Equal(t, 0, tree.Size())

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	tree.JoinExclusive(t1)

	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	tree.JoinExclusive(t2)
	assert.Equal(t, 0, t1.Size())
	assert.Equal(t, 0, t2.Size())

	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()
}

func TestTreap_searchOrInsert(t *testing.T) {

	const N = 1000
	tree := New(1, cmpInt)
	failures := New(2, cmpInt)

	for i := 0; i < N; i++ {
		val := rand.Intn(100 * N)
		ok, res := tree.SearchOrInsert(val)
		assert.Equal(t, val, res)
		if !ok {
			failures.Insert(val)
		}
	}

	fmt.Println("tree.Size() = ", tree.Size())
	fmt.Println("failures.Size() = ", failures.Size())

	for it := NewIterator(failures); it.HasCurr(); it.Next() {
		assert.Equal(t, tree.Search(it.GetCurr()), it.GetCurr())
	}
}

func TestTreap_choose(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 1000
	insertNRandomItems(tree, N)

	for i, it := 0, NewIterator(tree); it.HasCurr(); it.Next() {
		item := tree.Choose(i)
		assert.Equal(t, item, it.GetCurr())
		i++
	}
}

func TestTreap_rank(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 1000
	insertNRandomItems(tree, N)

	for i, it := 0, NewIterator(tree); it.HasCurr(); it.Next() {
		ok, pos := tree.RankInOrder(it.GetCurr())
		assert.True(t, ok)
		assert.Equal(t, i, pos)
		i++
	}
}

func TestTreap_splitPos(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100
	insertNRandomItems(tree, N)

	min1, max1, min2, max2 := tree.Min(), tree.Choose(N/2-1), tree.Choose(N/2), tree.Max()

	t1, t2 := tree.SplitByPosition(N / 2)

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	assert.Equal(t, 0, tree.Size())
	assert.NotNil(t, t1)
	assert.NotNil(t, t2)
	assert.Equal(t, N/2, t1.Size())
	assert.Equal(t, N/2, t2.Size())
	assert.Equal(t, min1, t1.Min())
	assert.Equal(t, max1, t1.Max())
	assert.Equal(t, min2, t2.Min())
	assert.Equal(t, max2, t2.Max())
}

func TestTreap_copy(t *testing.T) {
	t1 := New(2, cmpInt)
	const N = 100
	insertNRandomItems(t1, N)

	assert.True(t, checkBST(*t1.rootPtr, t1.less))
	assert.True(t, checkTreap(*t1.rootPtr))
	assert.True(t, checkCounter(*t1.rootPtr))

	t2 := t1.Copy()

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	assert.True(t, t1.TopologicalEqual(t2))
}

func TestTreap_removeRange(t *testing.T) {
	tree := New(2, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	midRange := tree.ExtractRange(40, 60)

	for key, it := 40, NewIterator(midRange); it.HasCurr(); it.Next() {
		assert.Equal(t, key, it.GetCurr())
		key++
	}
	fmt.Println()

	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

}
