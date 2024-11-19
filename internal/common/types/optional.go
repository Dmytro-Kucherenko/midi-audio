package types

type Optional[T any] struct {
	Value T
	Valid bool
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{value, true}
}

func Equal[T comparable](optional Optional[T], value T) bool {
	return optional.Valid && optional.Value == value
}

func EqualOptional[T comparable](optional Optional[T], value T) bool {
	return !optional.Valid || optional.Value == value
}
