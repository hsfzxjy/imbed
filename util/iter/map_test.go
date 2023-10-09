package iter_test

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/util/iter"
)

func Collect[T any](it core.Iterator[T]) []T {
	var res []T
	for it.HasNext() {
		res = append(res, it.Next())
	}
	return res
}

func Id[T any](x T) (T, bool) { return x, true }

func ExampleFlatFilterMap() {
	its := []core.Iterator[int]{
		iter.Slice(1, 2, 3),
		iter.Slice(4, 5, 6),
		iter.Slice[int](),
	}
	flatten := iter.FlatFilterMap(iter.Slice(its...), Id[int])
	fmt.Printf("%v\n", Collect(flatten))
	// Output:
	// [1 2 3 4 5 6]
}
