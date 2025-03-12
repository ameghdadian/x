// Package iter provides utility functions for working with iterators.
package iter

import (
	"iter"
	"math/rand/v2"
	"slices"
)

// Concat receives an arbitrary number of slices and returns an iterator to
// iterate on them successively, in the order they are given.
func Concat[Slice []T, T any](seqs ...Slice) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for _, v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ConcatIter receives an arbitrary number of iterators and returns an iterator to
// iterate over the elements, in the order they are given.
func ConcatIter[T any](iters ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range iters {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Shuffle receives a slice and a swap function. It returns an iterator over the slice shuffled elements.
//
// NOTE: In order not to allocate extra memory, it shuffles given slice(array) in place.
// If you need the original slice(array) untouched, you can clone your slice using [slices.Clone]
// and then pass it as the argument.
func Shuffle[Slice []T, T any](seq Slice, swapFn func(i, j int)) iter.Seq[T] {
	return func(yield func(T) bool) {
		// swapFunc := func(i int, j int) {
		// 	seq[i], seq[j] = seq[j], seq[i]
		// }
		rand.Shuffle(len(seq), swapFn)

		for v := range slices.Values(seq) {
			if !yield(v) {
				return
			}
		}
	}
}

// Filter returns an iterator that contains the elements of `it` for which filterFn
// returns true.
func Filter[T comparable](it iter.Seq[T], filterFn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if filterFn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Map executes a user-supplied function on each element of the `it` and returns
// an iterator over modified elements.
func Map[T, R any](it iter.Seq[T], mapFn func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range it {
			if !yield(mapFn(v)) {
				return
			}
		}
	}
}

// ForEach executes a user-supplied function on each element of `it`.
func ForEach[T any](it iter.Seq[T], fn func(int, T)) {
	var indx int
	for v := range it {
		fn(indx, v)
		indx += 1
	}
}

// Reduce executes a user-supplied `reducer` function on each element of the iterator `it`,
// in order, passing in the return value from the calculation on the preceding element.
func Reduce[T, R any](it iter.Seq[T], reducer func(acc R, cur T) R, init R) R {
	for v := range it {
		init = reducer(init, v)
	}

	return init
}
