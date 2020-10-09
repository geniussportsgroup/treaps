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

	assert.True(t, tree.check())

	// test that insert fails for duplicated key
	for i := 0; i < N; i++ {
		ret := tree.Insert(i)
		assert.Nil(t, ret)
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

	assert.Equal(t, N, tree.Size())
	assert.True(t, tree.check())
}

func TestTreap_remove(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100

	insertNRandomItems(tree, N)

	values := make([]int, 0, tree.Size())
	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		values = append(values, it.GetCurr().(int))
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for _, val := range values {
		assert.Equal(t, val, tree.Search(val), "key must be in tree")
		assert.Equal(t, val, tree.Remove(val))
		assert.True(t, tree.check())
		assert.Equal(t, nil, tree.Search(val))
		assert.Nil(t, tree.Remove(val), "key already removed should fail")
	}

	assert.True(t, tree.check())
}

func TestTreap_split(t *testing.T) {

	tree := New(1, cmpInt)
	const N = 100
	insertNRandomItems(tree, N)

	t1, t2 := tree.SplitByKey(552)

	assert.Equal(t, N, t1.Size()+t2.Size())
	assert.Equal(t, 0, tree.Size())

	assert.True(t, t1.check())
	assert.True(t, t2.check())

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	tree.JoinExclusive(t1)

	assert.True(t, tree.check())

	for it := NewIterator(tree); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	tree.JoinExclusive(t2)
	assert.True(t, tree.check())
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

	assert.True(t, tree.check())
	assert.True(t, failures.check())

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

	assert.True(t, tree.check())
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

	assert.True(t, tree.check())
}

func TestTreap_splitPos(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	min1, max1, min2, max2 := tree.Min(), tree.Choose(N/2), tree.Choose(N/2+1), tree.Max()

	t1, t2 := tree.SplitByPosition(N / 2)

	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, 0, tree.Size())
	assert.NotNil(t, t1)
	assert.NotNil(t, t2)
	assert.Equal(t, N/2+1, t1.Size())
	assert.Equal(t, N/2-1, t2.Size())
	assert.Equal(t, min1, t1.Min())
	assert.Equal(t, max1, t1.Max())
	assert.Equal(t, min2, t2.Min())
	assert.Equal(t, max2, t2.Max())

	for i, it := 0, NewIterator(t1); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr().(int))
	}

	for i, it := N/2+1, NewIterator(t2); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr().(int))
	}

	t0, t1 := t1.SplitByPosition(0)

	assert.True(t, t0.check())
	assert.True(t, t1.check())

	for i, it := 0, NewIterator(t0); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}

	for i, it := 1, NewIterator(t1); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}
}

func TestTreap_SplitByPositionCorners(t *testing.T) {
	tree := New(1, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	t1, t2 := tree.SplitByPosition(tree.Size() - 1)
	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, N, t1.Size())
	assert.Equal(t, 0, t2.Size())
	for i, it := 0, NewIterator(t1); i < N; i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}

	t1, t2 = t1.SplitByPosition(0)
	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, 1, t1.Size())
	assert.Equal(t, N-1, t2.Size())
	assert.Equal(t, 0, t1.Min())
	for i, it := 1, NewIterator(t2); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}

	t1, t2 = t2.SplitByPosition(0)
	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, 1, t1.Size())
	assert.Equal(t, 1, t1.Min())
	assert.Equal(t, N-2, t2.Size())
	for i, it := 2, NewIterator(t2); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}

	t1, t2 = t2.SplitByPosition(t2.Size() - 2)
	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, N-3, t1.Size())
	assert.Equal(t, 1, t2.Size())
	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

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

	assert.True(t, tree.check())
	assert.True(t, midRange.check())

	for key, it := 40, NewIterator(midRange); it.HasCurr(); it.Next() {
		assert.Equal(t, key, it.GetCurr())
		key++
	}
}

func TestTreap_ExtractRangeCorners(t *testing.T) {
	tree := New(2, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	res := tree.ExtractRange(0, tree.Size()-1)
	for it := NewIterator(res); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	assert.Equal(t, N-1, res.Size())
	assert.Equal(t, 1, tree.Size())
}

func TestTreap_IteratorNext(t *testing.T) {
	tree := New(3, cmpInt)
	const N = 100
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	i, it := 0, NewIterator(tree)
	for ; it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}
	assert.Equal(t, i, N)

	for i, it = N-1, it.Prev(); it.HasCurr(); i, it = i-1, it.Prev() {
		assert.Equal(t, i, it.GetCurr())
	}
	assert.Equal(t, i, -1)
}

func TestNewReverseIterator(t *testing.T) {
	tree := New(3, cmpInt)
	const N = 10000
	for i := 0; i < N; i++ {
		tree.Insert(i)
	}

	i, it := N-1, NewReverseIterator(tree)
	for ; it.HasCurr(); i, it = i-1, it.Prev() {
		assert.Equal(t, i, it.GetCurr())
	}
	assert.Equal(t, i, -1)

	for i, it = 0, it.Next(); it.HasCurr(); i, it = i+1, it.Next() {
		assert.Equal(t, i, it.GetCurr())
	}
	assert.Equal(t, i, N)
}

func TestTreap_joinDup(t *testing.T) {

	const N = 1000
	t1, t2 := NewTreap(cmpInt), NewTreap(cmpInt)

	insertNRandomItems(t1, N)
	insertNRandomItems(t2, N)

	n1, n2 := t1.Size(), t2.Size()

	t1.JoinDup(t2)

	assert.True(t, checkAll(*t1.rootPtr, t1.less))
	assert.Equal(t, n1+n2, t1.Size())
	assert.True(t, t1.check())

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()
}

func TestTreap_Union(t *testing.T) {

	const N = 1000
	t1, t2 := NewTreap(cmpInt), NewTreap(cmpInt)

	insertNRandomItems(t1, N)
	insertNRandomItems(t2, N)

	n1, n2 := t1.Size(), t2.Size()

	t1.Union(t2)

	assert.True(t, t1.check())
	assert.Equal(t, n2, t2.Size())
	assert.Less(t, n1, t1.Size())

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		assert.Equal(t, it.GetCurr(), t1.Search(it.GetCurr()))
	}
}

func Test_checkBST(t *testing.T) {

	root := &Node{
		key:      10,
		priority: 0,
		count:    0,
		llink: &Node{
			key:      5,
			priority: 0,
			count:    0,
			llink:    nullNodePtr,
			rlink:    nullNodePtr,
		},
		rlink: &Node{
			key:      15,
			priority: 0,
			count:    0,
			llink:    nullNodePtr,
			rlink:    nullNodePtr,
		},
	}

	assert.True(t, checkBST(root, cmpInt))

	root.llink.key = 11
	assert.False(t, checkBST(root, cmpInt))
}

func TestTreap_SimpleIntersection(t *testing.T) {

	t1 := New(1, cmpInt, 1, 3, 5, 7, 9, 10, 11, 13, 15, 17, 19)
	t2 := New(1, cmpInt, 2, 4, 6, 8, 9, 10, 12, 14, 16, 18, 20)

	for it := NewIterator(t1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(t2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	result, d1, d2 := t1.Intersection(t2)

	assert.True(t, t1.check())
	assert.True(t, t2.check())
	assert.Equal(t, 0, t1.Size())
	assert.Equal(t, 0, t2.Size())

	for it := NewIterator(result); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(d1); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()

	for it := NewIterator(d2); it.HasCurr(); it.Next() {
		fmt.Print(it.GetCurr(), " ")
	}
	fmt.Println()
}

func TestTreap_Swap(t *testing.T) {

	t1 := createSamples(10)
	t2 := createSamples(15)

	assert.True(t, t1.check())
	assert.True(t, t2.check())

	t1.Swap(t2)

	assert.True(t, t1.check())
	assert.True(t, t2.check())
}

func TestTreap_Intersection(t *testing.T) {
	const N = 100000
	t1, t2 := New(1, cmpInt), New(1, cmpInt)
	insertNRandomItems(t1, N)
	insertNRandomItems(t2, N)

	c1 := t1.Copy()
	c2 := t2.Copy()

	inter, diff1, diff2 := t1.Intersection(t2)

	assert.True(t, t1.check())
	assert.True(t, t2.check())

	assert.Equal(t, 0, t1.Size())
	assert.Equal(t, 0, t2.Size())

	for it := NewIterator(inter); it.HasCurr(); it.Next() {
		curr := it.GetCurr()
		assert.Equal(t, curr, c1.Search(curr))
		assert.Equal(t, curr, c2.Search(curr))
		assert.Equal(t, nil, diff1.Search(curr))
		assert.Equal(t, nil, diff2.Search(curr))
	}
}

func TestTreap_IntersectionCorners(t *testing.T) {

	inter, d1, d2 := NewTreap(cmpInt).Intersection(NewTreap(cmpInt))
	assert.Equal(t, 0, inter.Size())
	assert.Equal(t, 0, d1.Size())
	assert.Equal(t, 0, d2.Size())
	assert.True(t, inter.check())
	assert.True(t, d1.check())
	assert.True(t, d2.check())

	t1 := NewTreap(cmpInt, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19)
	t2 := NewTreap(cmpInt, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20)
	c1, c2 := t1.Copy(), t2.Copy()
	n1, n2 := t1.Size(), t2.Size()

	inter, d1, d2 = t1.Intersection(t2)
	assert.Equal(t, 0, inter.Size())
	assert.True(t, inter.check())
	assert.True(t, d1.check())
	assert.True(t, d2.check())
	assert.Equal(t, 0, t1.Size())
	assert.Equal(t, 0, t2.Size())
	assert.Equal(t, n1, d1.Size())
	assert.Equal(t, n2, d2.Size())
	assert.Equal(t, 0, c1.lexicographicCmp(d1))
	assert.Equal(t, 0, c2.lexicographicCmp(d2))
}

func TestTreap_lexicographicCmp(t *testing.T) {

	t1 := NewTreap(cmpInt, 1, 2, 3)
	t2 := NewTreap(cmpInt, 1, 2)
	t3 := NewTreap(cmpInt, 1)
	t4 := NewTreap(cmpInt)
	t5 := NewTreap(cmpInt, 2, 3, 4)

	assert.Equal(t, 1, t1.lexicographicCmp(t2))
	assert.Equal(t, -1, t2.lexicographicCmp(t1))

	assert.Equal(t, -1, t4.lexicographicCmp(t3))
	assert.Equal(t, -1, t4.lexicographicCmp(t2))
	assert.Equal(t, -1, t4.lexicographicCmp(t1))

	assert.Equal(t, 0, t1.lexicographicCmp(t1.Copy()))

	assert.Equal(t, -1, t3.lexicographicCmp(t2))

	assert.Equal(t, -1, t1.lexicographicCmp(t5))

	assert.Equal(t, 1, t5.lexicographicCmp(t1))
}

func TestTreap_RemoveByPos(t *testing.T) {

	tree := NewTreap(cmpInt, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17)

	assert.Equal(t, 0, tree.RemoveByPos(0))
	assert.True(t, tree.check())

	assert.Equal(t, 17, tree.RemoveByPos(16))
	assert.True(t, tree.check())
}
