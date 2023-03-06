package util

type Optional[T any] struct {
	isSome bool
	some T
}

func (opt Optional[T]) IsNone() bool {
	return !opt.isSome
}

func (opt Optional[T]) IsSome() bool {
	return opt.isSome
}

func (opt Optional[T]) Some() T {
	var some T
	if opt.isSome {
		some = opt.some
	}

	return some
}

func Some[T any](t T) Optional[T] {
	return Optional[T]{
		isSome: true,
		some: t,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

