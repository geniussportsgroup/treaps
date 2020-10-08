# Golang Ordered Sets

This package implements ordered sets implemented through Treaps. It accepts any type of key. 

The user must specify a function `less(k1, k2)` which receives a pair of keys and determine whether `k1` is or not less than `k2`. In the `less()` function, the user can specify access to the key; this is the way as the package can handle arbitrary types of keys.

Given a set, the following operations are supported:
- Insertion, search, and deletion in O(log n) expected case. The package can handle repeated keys.
- To know the i-th key inside the order in O(log n) expected case.
- To know what is the inorder position of any key in O(log n) expected case.
- To split the set by a key's value or at any specific position in O(log n) expected case.
- To extract any rank [i, j] of the set in O(log n) where i and j are position respect to the order.
- Union and interception of sets in O(m log n).
- Join of rank disjoint sets in O(max(log n, log m)).

## Installation

    go get -u github.com/geniussportsgroup/treaps
    
## Example

This example shows how to get the 99th percentile of a set

    package treaps

    import (
        "github.com/stretchr/testify/assert"
        "math/rand"
        "testing"
    )

    // This example shows how to get the some percentiles of a sample of a ten millions of random heights

    // The following type simulate the representation of the height of a person
    type Sample struct {
        id     int // id of person
        height int // height in millimeters
    }

    const N = int(1e7)

    func createSamples(n int) *Treap {

        set := NewTreap(func(i1, i2 interface{}) bool {
            p1, ok := i1.(*Sample)
            if !ok {
                panic("First parameter is not of type *Sample")
            }
            p2, ok := i2.(*Sample)
            if !ok {
                panic("Second parameter is not of type *Sample")
            }
            return p1.height < p2.height // sort by height
        })

        // we will generate 1e6 samples of random heights according to a normal dist with mean 1600 mm
        // and standard deviation of 400 mm/
        // Pay attention to the fact that heights can be repeat. So we must use InsertDup
        for id := 0; id < n; id++ {
            set.InsertDup(&Sample{
                id:     id,
                height: int(rand.NormFloat64()*400 + 1600),
            })
        }

        return set
    }

    func TestExample_99Percentiles(t *testing.T) {

        set := createSamples(N)

        posOfPercentile99 := int((set.Size() * 99) / 100) // index of first item belonging to 99 percentile

        // we have two ways for getting the samples included in the 99 percentile

        percentile99Size := set.Size() - posOfPercentile99

        // First we can use the method choose(i) for getting all the samples. This way has complexity
        // O(n log(N)) where n is the n ~ 0.01*N and N is the number of samples, which definitely is faster
        // than O(N log N) in time and O(N) in space if you take the elements and sort them
        p99Slice := make([]Sample, 0)
        for i := posOfPercentile99; i < set.Size(); i++ {
            samplePtr := set.Choose(i)
            p99Slice = append(p99Slice, *samplePtr.(*Sample))
        }

        // The second method is just to extract the whole percentile from the set, with the eventual
        // disadvantage that the set is modified, but with the advantage that that the set is gotten in O(log n)
        p99 := set.ExtractRange(posOfPercentile99, set.Size()-1)

        assert.Equal(t, percentile99Size, p99.Size())
        assert.Equal(t, N - p99.Size(), set.Size())

        for i, it := 0, NewIterator(p99); i < len(p99Slice); i, it = i+1, it.Next() {
            assert.Equal(t, p99Slice[i].id, it.GetCurr().(*Sample).id)
        }
    }
