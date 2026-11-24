package adapter

import (
	"fmt"
	"gourd/internal/common"
	"maps"
	"slices"
)

type AnyAdapter interface {
	Apply(doc any) (*AnyResult, error)
}

type AnyResult struct {
	Result any
	Key    string
}
type Effect[D, R any] func(D) (R, error)

type resultWrapper[T any] struct {
	Result T
	Key    string
}

type adapterWrapper[D, R any] struct {
	effect Effect[D, *resultWrapper[R]]
}

func (w adapterWrapper[D, R]) Apply(doc any) (*AnyResult, error) {
	typed, ok := doc.(D)
	if !ok {
		return nil, fmt.Errorf("adapter expects %T, got %T", *new(D), doc)
	}
	res, err := w.effect(typed)
	if err != nil {
		return nil, err
	}
	return &AnyResult{
		Result: res.Result,
		Key:    res.Key,
	}, nil
}

func Wrap[D, R any](key string, fn Effect[D, R]) AnyAdapter {
	return adapterWrapper[D, R]{
		effect: func(d D) (*resultWrapper[R], error) {
			res, err := fn(d)
			if err != nil {
				return nil, err
			}
			return &resultWrapper[R]{
				Result: res,
				Key:    key,
			}, nil
		},
	}
}

type AdaptersMap map[common.FileType]AnyAdapter

func (m AdaptersMap) AvailableFileTypes() []common.FileType {
	return slices.Collect(maps.Keys(m))
}
