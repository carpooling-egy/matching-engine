package collections

type Tuple2[A any, B any] struct {
	First  A
	Second B
}

func NewTuple2[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{First: a, Second: b}
}
