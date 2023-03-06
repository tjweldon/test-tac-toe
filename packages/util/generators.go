package util

import "fmt"

type Generator[T any] func() (next T)

type FiniteGen[T any] func() (next T, ok bool) 

func GenFrom[T any](items ...T) FiniteGen[T] {
	idx := 0
	gen := func() (next T, ok bool) {
		if idx == len(items) {
			return
		}
		item := items[idx]
		idx++
		return item, true
	}

	return gen
}

func LoopFrom[T any](items ...T) (g Generator[T], err error) {
	if len(items) == 0 {
		return g, fmt.Errorf("Cannot initialise a looping generator with an empty collection")
	}

	gen := GenFrom(items...)
	loop := func() (next T) {
		var ok bool
		if next, ok = gen(); !ok {
			gen = GenFrom(items...)
			next, _ = gen()
		}
		return next
	}

	return loop, nil
}
