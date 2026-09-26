package adapter

import (
	"fmt"
)

type AnyAdapter interface {
	Apply(doc any) error
}
type Effect[T any] func(T) error

type adapterWrapper[T any] struct {
	effect Effect[T]
}

func (w adapterWrapper[T]) Apply(doc any) error {
	typed, ok := doc.(T)
	if !ok {
		return fmt.Errorf("adapter expects %T, got %T", *new(T), doc)
	}
	return w.effect(typed)
}

func Wrap[T any](fn Effect[T]) AnyAdapter {
	return adapterWrapper[T]{effect: fn}
}
