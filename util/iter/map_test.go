package iter_test

import (
	"cmp"
	"fmt"

	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

func Collect[T any](it iter.Ator[T]) []T {
	var res []T
	for it.HasNext() {
		res = append(res, it.Next().Unwrap())
	}
	return res
}

func ExampleFlatFilterMap() {
	its := []iter.Nexter[int]{
		iter.Slice(1, 2, 3),
		iter.Slice(4, 5, 6),
		iter.Slice[int](),
	}
	flatten := iter.FlatFilterMap(iter.Slice(its...), tipe.Ok[int])
	fmt.Printf("%v\n", Collect(flatten))
	chain := iter.Chain(iter.Slice(1, 2, 3), iter.Slice(4, 5, 6))
	fmt.Printf("%v\n", Collect(chain))
	// Output:
	// [1 2 3 4 5 6]
	// [1 2 3 4 5 6]
}

func ExampleSorted() {
	it := iter.Sorted(iter.Slice(2, 4, 6, 5, 3, 1), cmp.Compare[int], false)
	fmt.Printf("%v\n", Collect(it))
	// Output:
	// [1 2 3 4 5 6]
}
