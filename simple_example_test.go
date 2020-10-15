package treaps

import (
	"fmt"
	Set "github.com/geniussportsgroup/treaps"
	"testing"
)

func Test_Simple(t *testing.T) {

	// we create a simple tree of 15 integer keys
	tree := Set.NewTreap(func(i1, i2 interface{}) bool {
		return i1.(int) < i2.(int)
	}, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)

	// Now we traverse the tree with the iterator
	for it := Set.NewIterator(tree); it.HasCurr(); it.Next() {
		fmt.Print("[key = ", it.GetCurr())
		_, pos := tree.RankInOrder(it.GetCurr())
		fmt.Print(" pos = ", pos, "]")
	}
	fmt.Println()

	// Now we access the keys by their inorder position, search in the tree and compute their ordinal
	for i := 0; i < 16; i++ {
		key := tree.Choose(i) // access for position inorder
		fmt.Print("[key = ", key, " ")

		foundKey := tree.Search(key)
		fmt.Print("found key = ", foundKey, " ")

		_, pos := tree.RankInOrder(key)
		fmt.Print(" pos = ", pos, "]")
	}
	fmt.Println()
}
