package treaps

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// This example shows how to get some percentiles of a sample of a 10 million of random heights

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
		return p1.height < p2.height
	})

	// we will generate 1e6 samples of random heights according to a normal dist with mean 1600 mm
	// and standard deviation of 400 mm/
	// Pay attention to the fact that heights can repeat. So we must use InsertDup
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
	// disadvantage that the set is modified, but with the advantage that the set is gotten in O(log n)
	p99 := set.ExtractRange(posOfPercentile99, set.Size()-1)

	assert.Equal(t, percentile99Size, p99.Size())
	assert.Equal(t, N-p99.Size(), set.Size())

	for i, it := 0, NewIterator(p99); i < len(p99Slice); i, it = i+1, it.Next().(*Iterator) {
		assert.Equal(t, p99Slice[i].id, it.GetCurr().(*Sample).id)
	}
}
